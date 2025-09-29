package integration

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/murraystewart96/token-swap/internal/config"
	"github.com/murraystewart96/token-swap/internal/contracts"
	"github.com/murraystewart96/token-swap/internal/events"
	"github.com/murraystewart96/token-swap/internal/models"
	"github.com/murraystewart96/token-swap/internal/storage"
	"github.com/murraystewart96/token-swap/internal/worker"
	"github.com/murraystewart96/token-swap/tests/integration/testutils"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

const (
	anvilDefaultPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	typeMeToken  = 0
	typeYouToken = 1
)

type TestContracts struct {
	MeToken  *contracts.MeToken
	YouToken *contracts.YouToken
	Pool     *contracts.Pool

	MeTokenAddr  common.Address
	YouTokenAddr common.Address
	PoolAddr     common.Address
}

// TestEventListener_EndToEndFlow tests the complete pipeline:
// Contract Event → Event Listener → Kafka → Worker → Database/Cache
func TestEventListener_EndToEndFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup test infrastructure
	infra, err := testutils.SetupTestInfrastructure(t.Context())
	require.NoError(t, err, "Failed to setup test infrastructure")
	defer infra.Cleanup()

	// Reset services to clean state
	require.NoError(t, infra.Reset(), "Failed to reset test infrastructure")

	// Step 1: Deploy test contract to Anvil
	t.Log("Step 1: Deploying test contract...")
	auth := createTransactor(t, infra.ETHClient)
	testContracts, err := deployTestContracts(t, infra.ETHClient, auth)
	require.NoError(t, err, "Failed to deploy test contracts")

	// Step 2: Start event listener
	t.Log("Step 2: Starting event listener...")
	listenerConfig := &config.Listener{
		RPCUrl:       "localhost:8546",
		ContractAddr: testContracts.PoolAddr.Hex(),
	}

	eventClient, err := events.NewClient(listenerConfig, infra.KafkaProducer, infra.DB)
	require.NoError(t, err)

	workerService, err := worker.New(infra.KafkaConsumer, []string{config.TradeHistoryTopic, config.ReserveHistoryTopic}, infra.PoolCache, infra.DB)
	require.NoError(t, err)

	// Start services
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	go eventClient.Listen(ctx)
	go workerService.Start(ctx)

	// Wait for both services to initialize
	time.Sleep(2 * time.Second)
	t.Log("Both event listener and worker started")

	// Step 3: Execute swap transaction to generate events
	t.Log("Step 3: Executing swap transaction...")
	swapTxHash, err := executeSwapTransaction(t, infra.ETHClient, testContracts, auth, 100, typeMeToken)
	require.NoError(t, err, "Failed to execute swap transaction")
	t.Logf("Swap transaction hash: %s", swapTxHash.Hex())

	// Step 4: Wait for event to be captured and processed
	t.Log("Step 4: Waiting for event processing...")
	require.Eventually(t, func() bool {
		// Verify swap created trade event in database
		tradeOk := verifyTradeInDatabase(t, infra.DB, swapTxHash.Hex())

		// Verify swap created reserve event in database
		reserveOk := verifyReserveInDatabase(t, infra.DB, swapTxHash.Hex())

		// Verify price in cache matches blockchain
		priceOk := verifyPriceInCache(t, testContracts.Pool, infra.PoolCache)

		return tradeOk && reserveOk && priceOk
	}, 20*time.Second, 1*time.Second)

	t.Log("Step 5: Event processing pipeline test completed!")

	t.Log("✅ End-to-end test completed successfully!")
}

func createTransactor(t *testing.T, client *ethclient.Client) *bind.TransactOpts {
	// Create test account from Anvil's default private key
	privateKey, err := crypto.HexToECDSA(anvilDefaultPrivateKey)
	require.NoError(t, err, "Failed to parse private key")

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok, "Failed to cast public key to ECDSA")

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	t.Logf("Deploying contract from account: %s", fromAddress.Hex())

	// Get nonce
	nonce, err := client.PendingNonceAt(t.Context(), fromAddress)
	require.NoError(t, err, "Failed to get nonce")

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(t.Context())
	require.NoError(t, err, "Failed to get gas price")

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(31337))
	require.NoError(t, err, "Failed to create transactor")

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	return auth
}

// deployTestContract deploys the token and pool contract to Anvil
// returns the address, auth and pool contract interface
func deployTestContracts(t *testing.T, client *ethclient.Client, auth *bind.TransactOpts) (*TestContracts, error) {
	// Deploy METoken
	meTokenAddr, meTokenTx, meToken, err := contracts.DeployMeToken(auth, client, big.NewInt(1000000))
	require.NoError(t, err)
	auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)
	t.Logf("METoken contract deployed at: %s", meTokenAddr.Hex())

	// Deploy YOU Token
	youTokenAddr, youTokenTx, youToken, err := contracts.DeployYouToken(auth, client, big.NewInt(100000))
	require.NoError(t, err)
	auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)
	t.Logf("YOUToken contract deployed at: %s", youTokenAddr.Hex())

	// Wait for token deployments
	meReceipt, err := bind.WaitMined(t.Context(), client, meTokenTx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), meReceipt.Status)

	youReceipt, err := bind.WaitMined(t.Context(), client, youTokenTx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), youReceipt.Status)

	// Deploy Pool
	poolAddr, poolTx, pool, err := contracts.DeployPool(auth, client, meTokenAddr, youTokenAddr)
	require.NoError(t, err)
	auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)
	t.Logf("Pool contract deployed at: %s", poolAddr.Hex())

	// Wait for pool deployment
	poolReceipt, err := bind.WaitMined(t.Context(), client, poolTx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), poolReceipt.Status)

	t.Logf("Pool deployed at: %s", poolAddr.Hex())

	metApproveTx, err := meToken.Approve(auth, poolAddr, big.NewInt(500000))
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)

	approveReceipt, err := bind.WaitMined(context.Background(), client, metApproveTx)
	if err != nil {
		return nil, err
	}
	require.Equal(t, uint64(1), approveReceipt.Status)

	youApproveTx, err := youToken.Approve(auth, poolAddr, big.NewInt(50000))
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)

	approveReceipt, err = bind.WaitMined(context.Background(), client, youApproveTx)
	if err != nil {
		return nil, err
	}
	require.Equal(t, uint64(1), approveReceipt.Status)

	// Add liquidity to pool
	liquidityTx, err := pool.AddLiquidity(auth, big.NewInt(500000), big.NewInt(50000))
	require.NoError(t, err)
	auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)

	// Wait for pool liquidity
	liquidityReceipt, err := bind.WaitMined(t.Context(), client, liquidityTx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), liquidityReceipt.Status)

	return &TestContracts{
		MeToken:      meToken,
		MeTokenAddr:  meTokenAddr,
		YouToken:     youToken,
		YouTokenAddr: youTokenAddr,
		Pool:         pool,
		PoolAddr:     poolAddr,
	}, nil
}

// executeSwapTransaction performs a swap on the deployed contract
func executeSwapTransaction(t *testing.T, client *ethclient.Client, testContracts *TestContracts, auth *bind.TransactOpts, amountIn int64, tokenType uint8) (common.Hash, error) {
	expectedOut, err := testContracts.Pool.GetAmountOut(nil, big.NewInt(amountIn), tokenType)
	if err != nil {
		return common.Hash{}, err
	}

	var approveTx, swapTx *types.Transaction

	switch tokenType {
	case typeMeToken:
		// Step 1: Submit approve transaction
		approveTx, err = testContracts.MeToken.Approve(auth, testContracts.PoolAddr, big.NewInt(amountIn))
		if err != nil {
			return common.Hash{}, err
		}
		auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)

		// Step 2: WAIT for approve to be mined before proceeding
		approveReceipt, err := bind.WaitMined(context.Background(), client, approveTx)
		if err != nil {
			return common.Hash{}, err
		}
		require.Equal(t, uint64(1), approveReceipt.Status)

		// Step 3: Now submit swap transaction
		swapTx, err = testContracts.Pool.SwapMeTokenForYouToken(auth, big.NewInt(amountIn), expectedOut)
		if err != nil {
			return common.Hash{}, err
		}
		auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)

	case typeYouToken:
		// Step 1: Submit approve transaction
		approveTx, err = testContracts.YouToken.Approve(auth, testContracts.PoolAddr, big.NewInt(amountIn))
		if err != nil {
			return common.Hash{}, err
		}
		auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)

		// Step 2: WAIT for approve to be mined before proceeding
		approveReceipt, err := bind.WaitMined(context.Background(), client, approveTx)
		if err != nil {
			return common.Hash{}, err
		}
		require.Equal(t, uint64(1), approveReceipt.Status)

		// Step 3: Now submit swap transaction
		auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)
		swapTx, err = testContracts.Pool.SwapYouTokenForMeToken(auth, big.NewInt(amountIn), expectedOut)
		if err != nil {
			return common.Hash{}, err
		}
		auth.Nonce = big.NewInt(auth.Nonce.Int64() + 1)
	}

	// Wait for swap to be mined
	swapReceipt, err := bind.WaitMined(context.Background(), client, swapTx)
	if err != nil {
		return common.Hash{}, err
	}
	require.Equal(t, uint64(1), swapReceipt.Status)

	return swapTx.Hash(), nil
}

func verifyTradeInDatabase(t *testing.T, db storage.DB, expectedTxHash string) bool {
	// Query trades table for the transaction hash
	trades, err := db.GetTradesByTimeRange(time.Now().Add(-10*time.Minute), time.Now())
	require.NoError(t, err, "Failed to query trades from database")

	// Look for our specific transaction
	var foundTrade *models.TradeEvent
	for _, trade := range trades {
		if trade.TxHash == expectedTxHash {
			foundTrade = trade
			break
		}
	}

	if foundTrade == nil {
		return false
	}

	return expectedTxHash == foundTrade.TxHash
}

func verifyReserveInDatabase(t *testing.T, db storage.DB, expectedTxHash string) bool {
	// Query reserves table for the transaction hash
	reserves, err := db.GetReservesByTimeRange(time.Now().Add(-10*time.Minute), time.Now())
	require.NoError(t, err, "Failed to query reserves from database")

	// Look for our specific transaction
	var foundReserve *models.ReserveEvent
	for _, reserve := range reserves {
		if reserve.TxHash == expectedTxHash {
			foundReserve = reserve
			break
		}
	}

	if foundReserve == nil {
		return false
	}

	return expectedTxHash == foundReserve.TxHash
}

func verifyPriceInCache(t *testing.T, pool *contracts.Pool, poolCache storage.PoolCache) bool {
	// Get cached price as decimal
	cachedPrice, err := poolCache.GetPrice(t.Context(), worker.MET_YOU_PAIR)
	if errors.Is(err, redis.Nil) {
		return false // entry doesn't exist yet
	} else {
		require.NoError(t, err, "Failed to get price from cache")
	}

	// Get latest price from the blockchain as decimal
	reserves, err := pool.GetReserves(nil)
	require.NoError(t, err, "Failed to get reserves from smart contract")

	// Calculate price from reserves (same as AMM formula)
	metReserve := decimal.NewFromBigInt(reserves.MeTokenReserve, 0)
	youReserve := decimal.NewFromBigInt(reserves.YouTokenReserve, 0)

	require.NoError(t, err, "Divide by zero")

	latestPrice := youReserve.Div(metReserve).StringFixed(6) // YOU per MET

	return cachedPrice == latestPrice
}

// TODO: BONUS TASKS (Optional)

//
// 3. Add error scenarios:
//    - Test with invalid transaction data
//    - Test with network interruptions
//    - Test with duplicate events
//
// 4. Add performance testing:
//    - Generate multiple swaps rapidly
//    - Measure event processing latency
//    - Verify no events are lost
