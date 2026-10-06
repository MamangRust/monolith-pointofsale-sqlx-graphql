package seeder

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

const seedGetUsers = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM users
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR firstname ILIKE '%' || $1 || '%' OR lastname ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const seedGetRoles = `SELECT
    role_id,
    role_name,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    roles
WHERE
    $1::TEXT IS NULL OR role_name ILIKE '%' || $1 || '%'
ORDER BY
    created_at ASC
LIMIT $2 OFFSET $3`

const seedAssignRoleToUser = `INSERT INTO user_roles (
    user_id,
    role_id,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    current_timestamp,
    current_timestamp
) RETURNING
    user_role_id,
    user_id,
    role_id,
    created_at,
    updated_at,
    deleted_at`

type userRoleSeeder struct {
	db     *sqlx.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserRoleSeeder(db *sqlx.DB, ctx context.Context, logger logger.LoggerInterface) *userRoleSeeder {
	return &userRoleSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *userRoleSeeder) Seed() error {
	var users []*db.GetUsersRow
	if err := r.db.SelectContext(r.ctx, &users, seedGetUsers, "", int32(20), int32(0)); err != nil {
		r.logger.Error("failed to fetch users", zap.Error(err))
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	var roles []*db.GetRolesRow
	if err := r.db.SelectContext(r.ctx, &roles, seedGetRoles, "", int32(4), int32(0)); err != nil {
		r.logger.Error("failed to fetch roles", zap.Error(err))
		return fmt.Errorf("failed to fetch roles: %w", err)
	}

	if len(users) == 0 || len(roles) == 0 {
		r.logger.Debug("no users or roles available for seeding")
		return nil
	}

	rand.Seed(uint64(time.Now().UnixNano()))

	for _, user := range users {
		role := roles[rand.Intn(len(roles))]

		var userRole db.UserRole
		if err := r.db.GetContext(r.ctx, &userRole, seedAssignRoleToUser, user.UserID, role.RoleID); err != nil {
			r.logger.Error("failed to assign role to user", zap.String("user", user.Email), zap.String("role", role.RoleName), zap.Error(err))
			return fmt.Errorf("failed to assign role %s to user %s: %w", role.RoleName, user.Email, err)
		}
	}

	r.logger.Info("user roles assigned successfully", zap.Int("users", len(users)), zap.Int("roles", len(roles)))
	return nil
}
