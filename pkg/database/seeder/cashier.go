package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

const seedGetMerchants = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM merchants
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR contact_email ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const seedCreateCashier = `INSERT INTO cashiers (merchant_id, user_id, name)
VALUES ($1, $2, $3) RETURNING *`

type cashierSeeder struct {
	db     *sqlx.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCashierSeeder(db *sqlx.DB, ctx context.Context, logger logger.LoggerInterface) *cashierSeeder {
	return &cashierSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *cashierSeeder) Seed() error {
	var merchants []*db.GetMerchantsRow
	if err := r.db.SelectContext(r.ctx, &merchants, seedGetMerchants, "", int32(20), int32(0)); err != nil {
		r.logger.Error("Failed to fetch merchants:", zap.Any("err", err))
		return err
	}

	var users []*db.GetUsersRow
	if err := r.db.SelectContext(r.ctx, &users, seedGetUsers, "", int32(20), int32(0)); err != nil {
		r.logger.Error("Failed to fetch users:", zap.Any("error", err))
		return err
	}

	if len(merchants) == 0 || len(users) == 0 {
		r.logger.Error("Merchants or Users not found. Seed operation aborted.")
		return fmt.Errorf("no merchants or users found")
	}

	for i := 1; i <= 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		user := users[rand.Intn(len(users))]

		cashierName := fmt.Sprintf("Cashier %d", i)
		var cashier db.Cashier
		if err := r.db.GetContext(r.ctx, &cashier, seedCreateCashier,
			merchant.MerchantID,
			user.UserID,
			cashierName,
		); err != nil {
			r.logger.Error("Failed to create cashier:", zap.Any("error", err))
			return err
		}

	}

	r.logger.Info("Cashier seeding completed successfully.", zap.Int("count", 10))
	return nil
}
