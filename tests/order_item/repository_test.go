package order_item_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-order-item/repository"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-pointofsale-test"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type OrderItemRepositoryTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	repo        *repository.Repositories
	orderItemID int
	orderID     int
	productID   int
}

func (s *OrderItemRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)

	queries := pool
	s.repo = repository.NewRepositories(queries)

	// Seed data
	var userID int
	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ($1, $2, $3, $4, 'test-verify', true) RETURNING user_id`,
		"OItem", "Repo", "oitem.repo@example.com", "password123",
	).Scan(&userID)
	s.Require().NoError(err)

	var merchantID int
	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO merchants (user_id, name, status) VALUES ($1, $2, $3) RETURNING merchant_id`,
		userID, "OItem Merchant", "active",
	).Scan(&merchantID)
	s.Require().NoError(err)

	var categoryID int
	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO categories (name) VALUES ($1) RETURNING category_id`,
		"OItem Category",
	).Scan(&categoryID)
	s.Require().NoError(err)

	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO products (merchant_id, category_id, name, price, count_in_stock) VALUES ($1, $2, $3, $4, $5) RETURNING product_id`,
		merchantID, categoryID, "OItem Product", 1000, 10,
	).Scan(&s.productID)
	s.Require().NoError(err)

	var cashierID int
	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO cashiers (merchant_id, user_id, name) VALUES ($1, $2, $3) RETURNING cashier_id`,
		merchantID, userID, "Test Cashier",
	).Scan(&cashierID)
	s.Require().NoError(err)

	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO orders (merchant_id, cashier_id, total_price) VALUES ($1, $2, $3) RETURNING order_id`,
		merchantID, cashierID, 5000,
	).Scan(&s.orderID)
	s.Require().NoError(err)

	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING order_item_id`,
		s.orderID, s.productID, 2, 5000,
	).Scan(&s.orderItemID)
	s.Require().NoError(err)
}

func (s *OrderItemRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *OrderItemRepositoryTestSuite) Test1_FindAll() {
	ctx := context.Background()

	results, total, err := s.repo.OrderItemQuery.FindAllOrderItems(ctx, &requests.FindAllOrderItems{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(total)
	s.GreaterOrEqual(*total, 1)
	s.GreaterOrEqual(len(results), 1)
}

func (s *OrderItemRepositoryTestSuite) Test2_FindOrderItemByOrder() {
	s.Require().NotZero(s.orderID)
	ctx := context.Background()

	found, err := s.repo.OrderItemQuery.FindOrderItemByOrder(ctx, s.orderID)
	s.NoError(err)
	s.NotEmpty(found)
	s.Equal(int32(s.orderID), found[0].OrderID)
}

func (s *OrderItemRepositoryTestSuite) Test3_FindByActive() {
	ctx := context.Background()

	results, total, err := s.repo.OrderItemQuery.FindByActive(ctx, &requests.FindAllOrderItems{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(total)
	s.GreaterOrEqual(*total, 1)
	s.GreaterOrEqual(len(results), 1)
}

func (s *OrderItemRepositoryTestSuite) Test4_FindByTrashed() {
	ctx := context.Background()

	// First trash the order item
	_, err := s.ts.DBPool().Exec(ctx, `UPDATE order_items SET deleted_at = NOW() WHERE order_item_id = $1`, s.orderItemID)
	s.Require().NoError(err)

	results, total, err := s.repo.OrderItemQuery.FindByTrashed(ctx, &requests.FindAllOrderItems{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(total)
	s.GreaterOrEqual(*total, 1)
	s.GreaterOrEqual(len(results), 1)

	// Restore the order item for cleanup
	_, err = s.ts.DBPool().Exec(ctx, `UPDATE order_items SET deleted_at = NULL WHERE order_item_id = $1`, s.orderItemID)
	s.Require().NoError(err)
}

func TestOrderItemRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemRepositoryTestSuite))
}
