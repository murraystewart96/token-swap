package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/murraystewart96/token-swap/internal/storage"
	"github.com/rs/zerolog/log"
)

const (
	readTimeout     = time.Second * 60
	writeTimeout    = time.Second * 60
	shutdownTimeout = 5 * time.Second

	defaultTradesLimit = 30
	maxTradesLimit     = 100
)

type Handler struct {
	db    storage.DB
	cache storage.PoolCache
}

func NewHandler(db storage.DB, cache storage.PoolCache) *Handler {
	return &Handler{
		db:    db,
		cache: cache,
	}
}

func Serve(ctx context.Context, addr string, handler *Handler) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(
		gin.Recovery(),
	)

	router.GET("/api/pool/current-price", handler.GetCurrentPrice)
	router.GET("/api/pool/reserves", handler.GetReserves)
	router.GET("/api/trades", handler.GetTrades)

	router.GET("/api/analytics/volume", handler.GetVolumeAnalytics)
	router.GET("/api/analytics/price-history", handler.GetPriceHistory)
	router.GET("/api/analytics/activity", handler.GetActivityAnalytics)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	log.Info().Msgf("[HTTP] Starting server on %s", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Info().Msg("[HTTP] Shutting down server...")

	err := srv.Shutdown(shutdownCtx) //nolint:contextcheck
	if err != nil {
		panic(err)
	}

	<-shutdownCtx.Done()

	log.Info().Msg("[HTTP] Shutdown complete")
}

func (h *Handler) GetReserves(ctx *gin.Context) {
	poolAddr, found := ctx.GetQuery("pool_address")
	if !found {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "pool_address parameter required"})
		return
	}

	reserves, err := h.cache.GetReserves(ctx, poolAddr)
	if err != nil {
		log.Error().Err(err).Str("pool_address", poolAddr).Msg("failed getting reserves from cache")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve reserves"})
		return
	}

	resp := &models.ReservesResponse{
		METAmount: reserves.METAmount,
		YOUAmount: reserves.YOUAmount,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCurrentPrice(ctx *gin.Context) {
	tradingPair, found := ctx.GetQuery("trading_pair")
	if !found {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "trading_pair parameter required"})
		return
	}

	price, err := h.cache.GetPrice(ctx, tradingPair)
	if err != nil {
		log.Error().Err(err).Str("trading_pair", tradingPair).Msg("failed getting current price from cache")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve trading pair"})
		return
	}

	resp := models.CurrentPriceResponse{
		Price: price,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetTrades(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defaultTradesLimit
	}

	if limit > maxTradesLimit {
		limit = maxTradesLimit
	}

	cursor, _ := ctx.GetQuery("cursor")

	cursorBlock, cursorTx, err := parseCursor(cursor)
	if err != nil {
		log.Error().Err(err).Msgf("failed to parse cursor: %s", cursor)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse cursor param"})

		return
	}

	trades, err := h.db.GetTradesByCursor(cursorBlock, cursorTx, limit)
	if err != nil {
		log.Error().Err(err).Msg("failed getting trades")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve trading pair"})

		return
	}

	// Create response with next cursor
	var nextCursor *string
	if len(trades) > 0 {
		lastTrade := trades[len(trades)-1]
		cursor := createCursor(lastTrade.BlockNumber, lastTrade.TransactionIndex)
		nextCursor = &cursor
	}

	response := models.TradesResponse{
		Trades:     trades, // No conversion needed
		NextCursor: nextCursor,
		HasMore:    len(trades) == limit,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) GetVolumeAnalytics(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "24h")
	token := ctx.DefaultQuery("token", "all") // "all", "MET", "YOU"

	if token != "all" && token != "MET" && token != "YOU" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token parameter"})
		return
	}

	// Parse period into time range
	start, end, err := parsePeriod(period)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid period"})
		return
	}

	volumeData, err := h.db.GetVolumeAnalytics(start, end, token)
	if err != nil {
		log.Error().Err(err).Msg("failed getting volume analytics")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve volume data"})
		return
	}

	ctx.JSON(http.StatusOK, volumeData)
}

func (h *Handler) GetPriceHistory(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "24h")
	interval := ctx.DefaultQuery("interval", "1h") // "1m", "5m", "1h", "1d"

	start, end, err := parsePeriod(period)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid period"})
		return
	}

	// Parse interval for data point frequency
	intervalDuration, err := parseInterval(interval)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid interval"})
		return
	}

	priceData, err := h.db.GetPriceHistory(start, end, intervalDuration)
	if err != nil {
		log.Error().Err(err).Msg("failed getting price history")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve price history"})
		return
	}

	ctx.JSON(http.StatusOK, priceData)
}

func (h *Handler) GetActivityAnalytics(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "24h")

	start, end, err := parsePeriod(period)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid period"})
		return
	}

	activityData, err := h.db.GetActivityAnalytics(start, end)
	if err != nil {
		log.Error().Err(err).Msg("failed getting activity analytics")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve activity data"})
		return
	}

	ctx.JSON(http.StatusOK, activityData)
}

// *** HELPER ***

// Cursor format: "block_number:transaction_index"
func parseCursor(cursor string) (uint64, uint, error) {
	if cursor == "" {
		return 0, 0, nil // No cursor means start from beginning
	}

	parts := strings.Split(cursor, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid cursor format")
	}

	blockNumber, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid block number in cursor")
	}

	txIndex, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid transaction index in cursor")
	}

	return blockNumber, uint(txIndex), nil
}

func createCursor(blockNumber uint64, txIndex uint) string {
	return fmt.Sprintf("%d:%d", blockNumber, txIndex)
}

func parsePeriod(period string) (time.Time, time.Time, error) {
	now := time.Now()

	switch period {
	case "1h":
		return now.Add(-1 * time.Hour), now, nil
	case "24h":
		return now.Add(-24 * time.Hour), now, nil
	case "7d":
		return now.Add(-7 * 24 * time.Hour), now, nil
	case "30d":
		return now.Add(-30 * 24 * time.Hour), now, nil
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("unsupported period: %s", period)
	}
}

func parseInterval(interval string) (time.Duration, error) {
	switch interval {
	case "1m":
		return 1 * time.Minute, nil
	case "5m":
		return 5 * time.Minute, nil
	case "1h":
		return 1 * time.Hour, nil
	case "1d":
		return 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported interval: %s", interval)
	}
}
