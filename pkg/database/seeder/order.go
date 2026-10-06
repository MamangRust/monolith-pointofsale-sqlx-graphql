package seeder

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

const seedGetCashiers = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM cashiers
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const seedCreateOrder = `INSERT INTO orders (merchant_id, cashier_id, total_price)
VALUES ($1, $2, $3)
RETURNING *`

const seedGetProductsByMerchant = `WITH filtered_products AS (
    SELECT
        p.product_id,
        p.name,
        p.description,
        p.price,
        p.count_in_stock,
        p.brand,
        p.image_product,
        p.created_at,
        c.name AS category_name
    FROM
        products p
    JOIN
        categories c ON p.category_id = c.category_id
    WHERE
        p.deleted_at IS NULL
        AND p.merchant_id = $1
        AND (
            p.name ILIKE '%' || COALESCE($2, '') || '%'
            OR p.description ILIKE '%' || COALESCE($2, '') || '%'
            OR $2 IS NULL
        )
        AND (
            c.category_id = NULLIF($3, 0)
            OR NULLIF($3, 0) IS NULL
        )
        AND (
            p.price >= COALESCE(NULLIF($4, 0), 0)
            AND p.price <= COALESCE(NULLIF($5, 0), 999999999)
        )
)
SELECT
    (SELECT COUNT(*) FROM filtered_products) AS total_count,
    fp.*
FROM
    filtered_products fp
ORDER BY
    fp.created_at DESC
LIMIT $6 OFFSET $7`

const seedCreateOrderItem = `INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING *`

type orderSeeder struct {
	db     *sqlx.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewOrderSeeder(db *sqlx.DB, ctx context.Context, logger logger.LoggerInterface) *orderSeeder {
	return &orderSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *orderSeeder) Seed() error {
	var merchants []*db.GetMerchantsRow
	if err := r.db.SelectContext(r.ctx, &merchants, seedGetMerchants, "", int32(20), int32(0)); err != nil {
		r.logger.Error("Failed to get merchants", zap.Error(err))
		return err
	}

	var cashiers []*db.GetCashiersRow
	if err := r.db.SelectContext(r.ctx, &cashiers, seedGetCashiers, "", int32(20), int32(0)); err != nil {
		r.logger.Error("Failed to get cashiers", zap.Error(err))
		return err
	}

	if len(merchants) == 0 || len(cashiers) == 0 {
		r.logger.Error("No merchants or cashiers found, skipping order seeding")
		return nil
	}

	for i := 0; i < 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		cashier := cashiers[rand.Intn(len(cashiers))]
		totalPrice := int32(rand.Intn(500000) + 50000)

		var order db.Order
		if err := r.db.GetContext(r.ctx, &order, seedCreateOrder,
			merchant.MerchantID,
			cashier.CashierID,
			int64(totalPrice),
		); err != nil {
			r.logger.Error("Failed to create order", zap.Error(err))
			return err
		}

		orderID := order.OrderID

		var products []*db.GetProductsByMerchantRow
		if err := r.db.SelectContext(r.ctx, &products, seedGetProductsByMerchant,
			merchant.MerchantID,
			nil,
			0,
			0,
			0,
			int32(10),
			int32(0),
		); err != nil {
			r.logger.Error("Failed to get products", zap.Error(err))
			return err
		}

		if len(products) == 0 {
			r.logger.Debug("No products found for merchant", zap.Int32("merchant_id", merchant.MerchantID))
			continue
		}

		for j := 0; j < rand.Intn(5)+1; j++ {
			product := products[rand.Intn(len(products))]
			quantity := int32(rand.Intn(5) + 1)
			price := product.Price * quantity

			var orderItem db.OrderItem
			if err := r.db.GetContext(r.ctx, &orderItem, seedCreateOrderItem,
				orderID,
				product.ProductID,
				quantity,
				price,
			); err != nil {
				r.logger.Error("Failed to create order item", zap.Error(err))
				return err
			}
		}
	}

	r.logger.Info("Order seeding completed successfully.", zap.Int("count", 10))
	return nil
}
