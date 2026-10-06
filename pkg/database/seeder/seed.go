package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/jmoiron/sqlx"
)

type Deps struct {
	Db     *sqlx.DB
	Ctx    context.Context
	Logger logger.LoggerInterface
	Hash   hash.HashPassword
}

type Seeder struct {
	User        *userSeeder
	Role        *roleSeeder
	UserRole    *userRoleSeeder
	Cashier     *cashierSeeder
	Category    *categorySeeder
	Product     *productSeeder
	Merchant    *merchantSeeder
	Order       *orderSeeder
	Transaction *transactionSeeder
}

func NewSeeder(deps Deps) *Seeder {
	return &Seeder{
		User:        NewUserSeeder(deps.Db, deps.Hash, deps.Ctx, deps.Logger),
		Role:        NewRoleSeeder(deps.Db, deps.Ctx, deps.Logger),
		UserRole:    NewUserRoleSeeder(deps.Db, deps.Ctx, deps.Logger),
		Merchant:    NewMerchantSeeder(deps.Db, deps.Ctx, deps.Logger),
		Cashier:     NewCashierSeeder(deps.Db, deps.Ctx, deps.Logger),
		Category:    NewCategorySeeder(deps.Db, deps.Ctx, deps.Logger),
		Product:     NewProductSeeder(deps.Db, deps.Ctx, deps.Logger),
		Order:       NewOrderSeeder(deps.Db, deps.Ctx, deps.Logger),
		Transaction: NewTransactionSeeder(deps.Db, deps.Ctx, deps.Logger),
	}
}

func (s *Seeder) Run() error {
	if err := s.seedWithDelay("users", s.User.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("roles", s.Role.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("user_roles", s.UserRole.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("merchant", s.Merchant.Seed); err != nil {
		return nil
	}

	if err := s.seedWithDelay("cashier", s.Cashier.Seed); err != nil {
		return nil
	}

	if err := s.seedWithDelay("category", s.Category.Seed); err != nil {
		return nil
	}

	if err := s.seedWithDelay("product", s.Product.Seed); err != nil {
		return nil
	}

	if err := s.seedWithDelay("order", s.Order.Seed); err != nil {
		return nil
	}

	if err := s.seedWithDelay("transaction", s.Transaction.Seed); err != nil {
		return nil
	}

	return nil
}

func (s *Seeder) seedWithDelay(entityName string, seedFunc func() error) error {
	if err := seedFunc(); err != nil {
		return fmt.Errorf("failed to seed %s: %w", entityName, err)
	}

	time.Sleep(30 * time.Second)
	return nil
}

func ptrString(s string) *string {
	return &s
}

func ptrInt32(i int32) *int32 {
	return &i
}
