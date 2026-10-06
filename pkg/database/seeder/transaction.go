package seeder

import (
	"context"
	"math/rand"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

const seedGetOrders = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM orders
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR total_price::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const seedCreateTransaction = `INSERT INTO transactions (
    merchant_id,
    payment_method,
    amount,
    change_amount,
    payment_status,
    order_id,
    created_at,
    updated_at,
    deleted_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    NULL
)
RETURNING *`

type transactionSeeder struct {
	db     *sqlx.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTransactionSeeder(db *sqlx.DB, ctx context.Context, logger logger.LoggerInterface) *transactionSeeder {
	return &transactionSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *transactionSeeder) Seed() error {
	var orders []*db.GetOrdersRow
	if err := r.db.SelectContext(r.ctx, &orders, seedGetOrders, "", int32(20), int32(0)); err != nil {
		r.logger.Error("Failed to get transactions:", zap.Any("error", err))
		return err
	}

	var merchants []*db.GetMerchantsRow
	if err := r.db.SelectContext(r.ctx, &merchants, seedGetMerchants, "", int32(20), int32(0)); err != nil {
		r.logger.Error("Failed to get transactions:", zap.Any("error", err))
		return err
	}

	for i := 0; i < 10; i++ {
		selectedMerchantId := merchants[rand.Intn(len(merchants))]
		selectedOrderId := orders[rand.Intn(len(orders))]

		var paymentMethod string
		var amount, changeAmount float64
		var paymentStatus string

		paymentMethod = "Credit Card"
		amount = float64(100 + i)
		changeAmount = float64(5 + i)
		paymentStatus = "Completed"

		var transaction db.Transaction
		if err := r.db.GetContext(r.ctx, &transaction, seedCreateTransaction,
			selectedMerchantId.MerchantID,
			paymentMethod,
			int32(amount),
			ptrInt32(int32(changeAmount)),
			paymentStatus,
			selectedOrderId.OrderID,
		); err != nil {
			r.logger.Error("Failed to create transaction:", zap.Any("error", err))
			return err
		}
	}

	r.logger.Info("Successfully seeded 10 transactions.", zap.Int("count", 10))
	return nil
}
