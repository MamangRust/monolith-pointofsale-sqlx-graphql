package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

const createRole = `INSERT INTO roles (
    role_name,
    created_at,
    updated_at
) VALUES (
    $1,
    current_timestamp,
    current_timestamp
) RETURNING
    role_id,
    role_name,
    created_at,
    updated_at,
    deleted_at`

type roleSeeder struct {
	db     *sqlx.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewRoleSeeder(db *sqlx.DB, ctx context.Context, logger logger.LoggerInterface) *roleSeeder {
	return &roleSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *roleSeeder) Seed() error {
	randomRoles := []string{
		"ROLE_ADMIN",
		"Admin Access 1",
		"Super Admin",
		"Admin",
		"Store Manager",
		"Cashier",
		"Inventory Staff",
		"Support",
		"Auditor",
		"Viewer",
	}

	totalRoles := len(randomRoles)

	for i, roleName := range randomRoles {
		var role db.Role
		if err := r.db.GetContext(r.ctx, &role, createRole, roleName); err != nil {
			r.logger.Error("failed to seed role", zap.Int("role", i+1), zap.String("roleName", roleName), zap.Error(err))
			return fmt.Errorf("failed to seed role %d (%s): %w", i+1, roleName, err)
		}
	}

	r.logger.Info("role seeded successfully", zap.Int("totalRoles", totalRoles))
	return nil
}
