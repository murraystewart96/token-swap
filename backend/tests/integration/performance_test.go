package integration

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/murraystewart96/token-swap/internal/worker"
	"github.com/murraystewart96/token-swap/pkg/tracing"
	"github.com/murraystewart96/token-swap/tests/integration/testutils"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/trace"
)

// TraceAnalyzer helps analyze performance from trace data
type TraceAnalyzer struct {
	producerExporter *tracing.InMemoryExporter
	consumerExporter *tracing.InMemoryExporter
}

func (ta *TraceAnalyzer) AnalyzeAllTraces() []*tracing.TraceMetrics {
	// Combine spans from both producer and consumer
	allSpans := append(ta.producerExporter.GetSpans(), ta.consumerExporter.GetSpans()...)

	// Group spans by trace ID
	traceSpans := make(map[string][]trace.ReadOnlySpan)
	for _, span := range allSpans {
		traceID := span.SpanContext().TraceID().String()
		traceSpans[traceID] = append(traceSpans[traceID], span)
	}

	var results []*tracing.TraceMetrics
	for traceID, spans := range traceSpans {
		// Create temporary exporter for analysis
		tempExporter := tracing.NewInMemoryExporter()
		tempExporter.ExportSpans(context.Background(), spans)

		if metrics := tempExporter.AnalyzeTrace(traceID); metrics != nil {
			results = append(results, metrics)
		}
	}

	return results
}

// TestPerformance_EndToEndLatency measures true end-to-end latency using traces
func TestPerformance_EndToEndLatency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	// Setup infrastructure
	infra, err := testutils.SetupTestInfrastructure(context.Background())
	require.NoError(t, err)
	defer infra.Cleanup()
	require.NoError(t, infra.Reset())

	// Initialize trace exporters
	producerExporter, producerShutdown, err := tracing.InitTestTracer("test-producer")
	require.NoError(t, err)
	defer producerShutdown()

	consumerExporter, consumerShutdown, err := tracing.InitTestTracer("test-consumer")
	require.NoError(t, err)
	defer consumerShutdown()

	analyzer := &TraceAnalyzer{producerExporter, consumerExporter}

	// Start worker
	workerService, err := worker.New(
		infra.KafkaConsumer,
		[]string{config.TradeHistoryTopic, config.ReserveHistoryTopic},
		infra.PoolCache,
		infra.DB,
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	go workerService.Start(ctx)

	t.Log("Testing end-to-end latency with 100 events...")

	// Warmup phase
	t.Log("Warming up services...")

	for i := 0; i < 5; i++ {
		sendTradeEventWithoutTracing(t, ctx, infra, i)
	}
	time.Sleep(2 * time.Second)

	t.Log("Starting performance measurement...")

	// Send 100 events and measure latency

	for i := 0; i < 100; i++ {
		sendTradeEventWithTracing(t, ctx, infra, i)

		// Small delay to avoid overwhelming the system
		if i%10 == 0 && i > 0 { // every multiple of 10
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Wait for processing
	time.Sleep(3 * time.Second)

	// Analyze results
	metrics := analyzer.AnalyzeAllTraces()
	require.Greater(t, len(metrics), 0, "Should have captured trace metrics")

	var latencies []time.Duration
	var processingTimes []time.Duration

	for _, m := range metrics {
		latencies = append(latencies, m.TotalDuration)

		// Calculate actual processing time (excludes Kafka transport delays)
		actualProcessing := m.KafkaProduceTime + m.WorkerProcessingTime + m.DatabaseOperationTime + m.CacheOperationTime
		processingTimes = append(processingTimes, actualProcessing)

		t.Logf("Trace %s: End-to-End=%v, Processing=%v (Kafka: %v, Worker: %v, DB: %v)",
			m.TraceID[:8], m.TotalDuration, actualProcessing, m.KafkaProduceTime,
			m.WorkerProcessingTime, m.DatabaseOperationTime)
	}

	// Calculate statistics for end-to-end latency (includes Kafka transport)
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	e2eP50 := latencies[len(latencies)/2]
	e2eP95 := latencies[int(0.95*float64(len(latencies)))]

	// Calculate statistics for processing time (excludes transport delays)
	sort.Slice(processingTimes, func(i, j int) bool { return processingTimes[i] < processingTimes[j] })
	procP50 := processingTimes[len(processingTimes)/2]
	procP95 := processingTimes[int(0.95*float64(len(processingTimes)))]

	t.Logf("📊 Performance Results:")
	t.Logf("   End-to-End P50: %v, P95: %v (includes Kafka transport)", e2eP50, e2eP95)
	t.Logf("   Processing P50: %v, P95: %v (actual work time)", procP50, procP95)

	// Realistic assertions for test environment
	require.Less(t, e2eP50, 150*time.Millisecond, "P50 end-to-end latency")
	require.Less(t, e2eP95, 300*time.Millisecond, "P95 end-to-end latency")

	require.Less(t, procP50, 2*time.Millisecond, "Application processing should be fast")
	require.Less(t, procP95, 10*time.Millisecond, "Even P95 processing should be reasonable")
}

// Helper functions
func generateUniqueTxHash() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return "0x" + hex.EncodeToString(bytes)
}

func sendTradeEventWithTracing(t *testing.T, ctx context.Context, infra *testutils.TestInfrastructure, i int) {
	// Create event listener span
	ctx, span := tracing.StartSpan(ctx, "events.handleSwapEvent")
	defer span.End()

	txHash := generateUniqueTxHash()
	tradeEvent := &models.TradeEvent{
		TxHash:           txHash,
		BlockNumber:      uint64(1000000 + i),
		TransactionIndex: uint(i % 100),
		TokenIn:          "MET",
		TokenOut:         "YOU",
		AmountIn:         "100.0",
		AmountOut:        "150.0",
		Timestamp:        time.Now().Unix(),
	}

	span.SetAttributes(
		tracing.BlockchainAttributes(tradeEvent.BlockNumber, tradeEvent.TxHash)...,
	)

	eventJSON, err := json.Marshal(tradeEvent)
	require.NoError(t, err)

	err = infra.KafkaProducer.Produce(ctx, config.TradeHistoryTopic, []byte(tradeEvent.TxHash), eventJSON)
	require.NoError(t, err)
}

func sendTradeEventWithoutTracing(t *testing.T, ctx context.Context, infra *testutils.TestInfrastructure, i int) {
	txHash := generateUniqueTxHash()
	tradeEvent := &models.TradeEvent{
		TxHash:           txHash,
		BlockNumber:      uint64(1000000 + i),
		TransactionIndex: uint(i % 100),
		TokenIn:          "MET",
		TokenOut:         "YOU",
		AmountIn:         "100.0",
		AmountOut:        "150.0",
		Timestamp:        time.Now().Unix(),
	}

	eventJSON, err := json.Marshal(tradeEvent)
	require.NoError(t, err)

	err = infra.KafkaProducer.Produce(ctx, config.TradeHistoryTopic, []byte(tradeEvent.TxHash), eventJSON)
	require.NoError(t, err)
}

func sendReserveEventWithoutTracing(t *testing.T, ctx context.Context, infra *testutils.TestInfrastructure, i int) {
	txHash := generateUniqueTxHash()
	reserveEvent := &models.ReserveEvent{
		TxHash:      txHash,
		BlockNumber: uint64(1000000 + i),
		PoolAddress: "0x1234567890123456789012345678901234567890",
		METReserve:  fmt.Sprintf("%.2f", float64(1000000-i*100)),
		YOUReserve:  fmt.Sprintf("%.2f", float64(1500000+i*150)),
		Timestamp:   time.Now().Unix(),
	}

	eventJSON, err := json.Marshal(reserveEvent)
	require.NoError(t, err)

	err = infra.KafkaProducer.Produce(ctx, config.ReserveHistoryTopic, []byte(reserveEvent.TxHash), eventJSON)
	require.NoError(t, err)
}

func sendReserveEventWithTracing(t *testing.T, ctx context.Context, infra *testutils.TestInfrastructure, i int) {
	// Create event listener span
	ctx, span := tracing.StartSpan(ctx, "events.handleSyncEvent")
	defer span.End()

	txHash := generateUniqueTxHash()
	reserveEvent := &models.ReserveEvent{
		TxHash:      txHash,
		BlockNumber: uint64(1000000 + i),
		PoolAddress: "0x1234567890123456789012345678901234567890",
		METReserve:  fmt.Sprintf("%.2f", float64(1000000-i*100)),
		YOUReserve:  fmt.Sprintf("%.2f", float64(1500000+i*150)),
		Timestamp:   time.Now().Unix(),
	}

	span.SetAttributes(
		tracing.BlockchainAttributes(reserveEvent.BlockNumber, reserveEvent.TxHash)...,
	)

	eventJSON, err := json.Marshal(reserveEvent)
	require.NoError(t, err)

	err = infra.KafkaProducer.Produce(ctx, config.ReserveHistoryTopic, []byte(reserveEvent.TxHash), eventJSON)
	require.NoError(t, err)
}

func calculateAverage(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}
