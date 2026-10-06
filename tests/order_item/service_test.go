package order_item_test

import (
	"context"
	"testing"

	item_cache "github.com/MamangRust/monolith-pointofsale-order-item/cache"
	"github.com/MamangRust/monolith-pointofsale-order-item/repository"
	"github.com/MamangRust/monolith-pointofsale-order-item/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	tests "github.com/MamangRust/monolith-pointofsale-test"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type OrderItemServiceTestSuite struct {
	suite.Suite
	ts  *tests.TestSuite
	svc *service.Service
}

func (s *OrderItemServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(opts)

	queries := pool

	log, _ := logger.NewLogger("test", nil)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)
	mencache := item_cache.NewMencache(cacheStore)

	obs, _ := observability.NewObservability("test", log)

	repos := repository.NewRepositories(queries)

	s.svc = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache,
		Observability: obs,
	})
}

func (s *OrderItemServiceTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *OrderItemServiceTestSuite) TestOrderItemLifecycle() {
	ctx := context.Background()

	// Seed data directly (FK-valid: user → merchant → cashier → category → product → order)
	var userID, merchantID, categoryID, productID, orderID int
	err := s.ts.DBPool().QueryRow(ctx,
		`INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ($1, $2, $3, $4, 'test-verify', true) RETURNING user_id`,
		"OItem", "Svc", "oitem.svc@example.com", "password123",
	).Scan(&userID)
	s.Require().NoError(err)

	err = s.ts.DBPool().QueryRow(ctx,
		`INSERT INTO merchants (user_id, name, description, address, contact_email, contact_phone, status)
		 VALUES ($1, 'OI Merchant', 'Desc', 'Addr', 'oi@example.com', '123', 'active') RETURNING merchant_id`,
		userID,
	).Scan(&merchantID)
	s.Require().NoError(err)

	_, err = s.ts.DBPool().Exec(ctx,
		`INSERT INTO cashiers (merchant_id, user_id, name) VALUES ($1, $2, 'OI Cashier')`,
		merchantID, userID)
	s.Require().NoError(err)

	err = s.ts.DBPool().QueryRow(ctx,
		`INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING category_id`,
		"Category", "Desc",
	).Scan(&categoryID)
	s.Require().NoError(err)

	err = s.ts.DBPool().QueryRow(ctx,
		`INSERT INTO products (merchant_id, category_id, name, description, price, count_in_stock, brand, weight, image_product) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING product_id`,
		merchantID, categoryID, "Product", "Desc", 1000, 10, "Brand", 1, "img.jpg",
	).Scan(&productID)
	s.Require().NoError(err)

	err = s.ts.DBPool().QueryRow(ctx,
		`INSERT INTO orders (merchant_id, cashier_id, total_price) VALUES ($1, $2, $3) RETURNING order_id`,
		merchantID, 1, 5000,
	).Scan(&orderID)
	s.Require().NoError(err)

	// 1. Create OrderItem (direct DB since service is query-only and doesn't have command)
	// The order_item service is query-only (OrderItemQuery), items are created via order service
	// Test the query interface

	// Create an order item directly for query testing
	var oItemID int
	err = s.ts.DBPool().QueryRow(ctx,
		`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING order_item_id`,
		orderID, productID, 2, 5000,
	).Scan(&oItemID)
	s.Require().NoError(err)

	// 2. FindByOrder
	items, err := s.svc.OrderItemQuery.FindOrderItemByOrder(ctx, orderID)
	s.Require().NoError(err)
	s.NotEmpty(items)

	// 3. FindAll
	_, total, err := s.svc.OrderItemQuery.FindAllOrderItems(ctx, &requests.FindAllOrderItems{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 4. FindByActive
	active, _, err := s.svc.OrderItemQuery.FindByActive(ctx, &requests.FindAllOrderItems{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(active)

	// 5. Trash the order item directly
	_, err = s.ts.DBPool().Exec(ctx, `UPDATE order_items SET deleted_at = NOW() WHERE order_item_id = $1`, oItemID)
	s.Require().NoError(err)

	// 6. FindByTrashed
	_, totalTrashed, err := s.svc.OrderItemQuery.FindByTrashed(ctx, &requests.FindAllOrderItems{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)
}

func TestOrderItemServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemServiceTestSuite))
}
