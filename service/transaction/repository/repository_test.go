package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) (*sqlx.DB, *Repositories, *sqlx.DB, func()) {
	t.Helper()

	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := sqlx.Connect("pgx", connStr)
	require.NoError(t, err)

	runMigrations(t, ctx, pool, "../../../pkg/database/migrations")

	repos := NewRepositories(pool, nil, nil, nil, nil)

	cleanup := func() {
		pool.Close()
		pgContainer.Terminate(ctx)
	}

	return pool, repos, pool, cleanup
}

func runMigrations(t *testing.T, ctx context.Context, pool *sqlx.DB, migrationDir string) {
	t.Helper()

	migrations := []string{
		"20250204052452_add_pgcrypto_and_uuid_columns.sql",
		"20250204052518_create_users.sql",
		"20250204052555_create_role.sql",
		"20250204052605_create_user_role.sql",
		"20250204052613_create_refreshtoken.sql",
		"20250204052614_create_reset_token.sql",
		"20250206050402_create_merchants.sql",
		"20250206050403_create_merchant_documents.sql",
		"20250206050448_create_cashiers.sql",
		"20250206050455_create_categories.sql",
		"20250206050505_create_products.sql",
		"20250206050519_create_orders.sql",
		"20250206050525_create_order_items.sql",
		"20250206050554_create_transactions.sql",
	}

	for _, file := range migrations {
		path := migrationDir + "/" + file
		sql, err := os.ReadFile(path)
		require.NoError(t, err)

		_, err = pool.ExecContext(ctx, string(sql))
		if err != nil {
			t.Logf("Warning executing %s: %v", file, err)
		}
	}
}

func ptr[T any](v T) *T {
	return &v
}

func seedTransactionDeps(t *testing.T, ctx context.Context, pool *sqlx.DB) (int, int, int) {
	t.Helper()

	var userID int
	err := pool.QueryRowContext(ctx, `INSERT INTO users (firstname, lastname, email, password, verification_code)
		VALUES ('Trans', 'User', 'trans@example.com', 'password', 'CODE') RETURNING user_id`).Scan(&userID)
	require.NoError(t, err)

	var merchantID int
	err = pool.QueryRowContext(ctx, `INSERT INTO merchants (user_id, name, status) VALUES ($1, 'Trans Merchant', 'active') RETURNING merchant_id`, userID).Scan(&merchantID)
	require.NoError(t, err)

	var cashierID int
	err = pool.QueryRowContext(ctx, `INSERT INTO cashiers (merchant_id, user_id, name) VALUES ($1, $2, 'Trans Cashier') RETURNING cashier_id`, merchantID, userID).Scan(&cashierID)
	require.NoError(t, err)

	var orderID int
	err = pool.QueryRowContext(ctx, `INSERT INTO orders (merchant_id, cashier_id, total_price) VALUES ($1, $2, 100000) RETURNING order_id`, merchantID, cashierID).Scan(&orderID)
	require.NoError(t, err)

	return merchantID, cashierID, orderID
}

func TestTransactionCommand_Create(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, _, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	txn, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID:       orderID,
		CashierID:     0, // Not used in create
		MerchantID:    merchantID,
		PaymentMethod: "cash",
		Amount:        100000,
		ChangeAmount:  ptr(0),
		PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)
	require.NotNil(t, txn)
	assert.Equal(t, int32(orderID), txn.OrderID)
	assert.Equal(t, int32(merchantID), txn.MerchantID)
	assert.Equal(t, "cash", txn.PaymentMethod)
	assert.Equal(t, int32(100000), txn.Amount)
	assert.Equal(t, "completed", txn.PaymentStatus)
	assert.NotZero(t, txn.TransactionID)
}

func TestTransactionQuery_FindByID(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, _, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	created, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "card",
		Amount: 50000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)

	found, err := repos.TransactionQueryRepository.FindById(ctx, int(created.TransactionID))
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.TransactionID, found.TransactionID)
	assert.Equal(t, "card", found.PaymentMethod)
}

func TestTransactionQuery_FindByOrderID(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, _, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	_, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "qris",
		Amount: 75000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)

	found, err := repos.TransactionQueryRepository.FindByOrderId(ctx, orderID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, int32(orderID), found.OrderID)
}

func TestTransactionCommand_Update(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	created, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "cash",
		Amount: 50000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)

	newStatus := "refunded"
	updated, err := repos.TransactionCommandRepository.UpdateTransaction(ctx, &requests.UpdateTransactionRequest{
		TransactionID: ptr(int(created.TransactionID)),
		OrderID:       orderID,
		CashierID:     cashierID,
		MerchantID:    merchantID,
		PaymentMethod: "transfer",
		Amount:        60000,
		PaymentStatus: &newStatus,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "transfer", updated.PaymentMethod)
	assert.Equal(t, "refunded", updated.PaymentStatus)
	assert.Equal(t, int32(60000), updated.Amount)
}

func TestTransactionQuery_FindAll(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, _, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	_, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "cash",
		Amount: 10000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)
	_, err = repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "card",
		Amount: 20000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)

	results, total, err := repos.TransactionQueryRepository.FindAllTransactions(ctx, &requests.FindAllTransaction{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 2)
	assert.GreaterOrEqual(t, len(results), 2)
}

func TestTransactionCommand_Trash(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, _, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	created, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "cash",
		Amount: 30000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)

	trashed, err := repos.TransactionCommandRepository.TrashTransaction(ctx, int(created.TransactionID))
	require.NoError(t, err)
	require.NotNil(t, trashed)
	assert.True(t, trashed.DeletedAt.Valid)
}

func TestTransactionQuery_FindByTrashed(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, _, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	created, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "cash",
		Amount: 40000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)

	_, err = repos.TransactionCommandRepository.TrashTransaction(ctx, int(created.TransactionID))
	require.NoError(t, err)

	results, total, err := repos.TransactionQueryRepository.FindByTrashed(ctx, &requests.FindAllTransaction{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestTransactionCommand_Restore(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, _, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	created, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "cash",
		Amount: 50000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)

	_, err = repos.TransactionCommandRepository.TrashTransaction(ctx, int(created.TransactionID))
	require.NoError(t, err)

	restored, err := repos.TransactionCommandRepository.RestoreTransaction(ctx, int(created.TransactionID))
	require.NoError(t, err)
	require.NotNil(t, restored)
	assert.False(t, restored.DeletedAt.Valid)
}

func TestTransactionQuery_FindByActive(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, _, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	created, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "cash",
		Amount: 60000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)

	results, total, err := repos.TransactionQueryRepository.FindByActive(ctx, &requests.FindAllTransaction{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)

	found := false
	for _, r := range results {
		if int(r.TransactionID) == int(created.TransactionID) {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestTransactionCommand_DeletePermanent(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, _, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	created, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "cash",
		Amount: 70000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)

	_, err = repos.TransactionCommandRepository.TrashTransaction(ctx, int(created.TransactionID))
	require.NoError(t, err)

	deleted, err := repos.TransactionCommandRepository.DeleteTransactionPermanently(ctx, int(created.TransactionID))
	require.NoError(t, err)
	assert.True(t, deleted)

	_, err = repos.TransactionQueryRepository.FindById(ctx, int(created.TransactionID))
	assert.Error(t, err)
}

func TestTransactionQuery_FindByMerchant(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, _, orderID := seedTransactionDeps(t, ctx, pool)
	paymentStatus := "completed"

	_, err := repos.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		OrderID: orderID, MerchantID: merchantID, PaymentMethod: "cash",
		Amount: 80000, PaymentStatus: &paymentStatus,
	})
	require.NoError(t, err)

	results, total, err := repos.TransactionQueryRepository.FindByMerchant(ctx, &requests.FindAllTransactionByMerchant{
		MerchantID: merchantID,
		Search:     "",
		Page:       1,
		PageSize:   10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

// ----- Transaction Stats Tests -----

func TestTransactionStats_GetMonthlyAmountSuccess(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	results, err := repos.TransactionStatsRepository.GetMonthlyAmountSuccess(ctx, &requests.MonthAmountTransaction{
		Year:  2026,
		Month: 1,
	})
	require.NoError(t, err)
	require.NotNil(t, results)
}

func TestTransactionStats_GetYearlyAmountSuccess(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	results, err := repos.TransactionStatsRepository.GetYearlyAmountSuccess(ctx, 2026)
	require.NoError(t, err)
	require.NotNil(t, results)
}

func TestTransactionStats_GetMonthlyAmountFailed(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	results, err := repos.TransactionStatsRepository.GetMonthlyAmountFailed(ctx, &requests.MonthAmountTransaction{
		Year:  2026,
		Month: 1,
	})
	require.NoError(t, err)
	require.NotNil(t, results)
}

func TestTransactionStats_GetYearlyAmountFailed(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	results, err := repos.TransactionStatsRepository.GetYearlyAmountFailed(ctx, 2026)
	require.NoError(t, err)
	require.NotNil(t, results)
}
