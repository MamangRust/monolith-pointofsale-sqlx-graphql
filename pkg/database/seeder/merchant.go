package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

const seedCreateMerchant = `INSERT INTO merchants (user_id, name, description, address, contact_email, contact_phone, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *`

type merchantSeeder struct {
	db     *sqlx.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantSeeder(db *sqlx.DB, ctx context.Context, logger logger.LoggerInterface) *merchantSeeder {
	return &merchantSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantSeeder) Seed() error {
	var users []*db.GetUsersRow
	if err := r.db.SelectContext(r.ctx, &users, seedGetUsers, "", int32(20), int32(0)); err != nil {
		r.logger.Error("Failed to fetch merchants:", zap.Any("error", err))
		return err
	}

	for i := 1; i <= 10; i++ {
		userID := users[i%len(users)].UserID

		var merchant db.Merchant
		if err := r.db.GetContext(r.ctx, &merchant, seedCreateMerchant,
			userID,
			fmt.Sprintf("Toko %d", i),
			ptrString(fmt.Sprintf("Deskripsi untuk Toko %d", i)),
			ptrString(fmt.Sprintf("Jl. Toko %d", i)),
			ptrString(fmt.Sprintf("toko%d@example.com", i)),
			ptrString(fmt.Sprintf("0812345678%d", i)),
			"active",
		); err != nil {
			r.logger.Error("Failed to create merchant:", zap.Error(err))
			return err
		}
	}

	r.logger.Info("Merchant seeding completed successfully.", zap.Int("count", 10))
	return nil
}
