package tests

import (
	"context"
	"reflect"

	"github.com/jmoiron/sqlx"

	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
)

type BaseTestSuite struct {
	suite.Suite
	ts      *TestSuite
	Log     logger.LoggerInterface
	Obs     observability.TraceLoggerObservability
	Conns   map[string]*grpc.ClientConn
	Servers []*grpc.Server
	Ctx     context.Context
	Cancel  context.CancelFunc
}

func (s *BaseTestSuite) SetupSuite() {
	s.Ctx, s.Cancel = context.WithCancel(context.Background())

	ts, err := SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	s.Log, _ = logger.NewLogger("test", nil)

	if s.Log == nil || (reflect.ValueOf(s.Log).Kind() == reflect.Ptr && reflect.ValueOf(s.Log).IsNil()) {
		z, _ := zap.NewDevelopment()
		s.Log = &logger.Logger{Log: z}
	}

	s.Obs, err = observability.NewObservability("test", s.Log)
	s.Require().NoError(err)
	s.Require().NotNil(s.Obs)
	s.Conns = make(map[string]*grpc.ClientConn)
}

func (s *BaseTestSuite) TearDownSuite() {
	for _, conn := range s.Conns {
		conn.Close()
	}
	for _, server := range s.Servers {
		server.GracefulStop()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
	if s.Cancel != nil {
		s.Cancel()
	}
}

func (s *BaseTestSuite) DBPool() *pgxpool.Pool {
	return s.ts.DBPool()
}

// SQLxDB opens a sqlx handle over the test database. Repositories in this repo
// are written against sqlx (pgx driver), so service wiring under test needs
// this rather than the raw pgx pool.
func (s *BaseTestSuite) SQLxDB() *sqlx.DB {
	db, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)
	return db
}

// DBURL exposes the test database connection string for suites that need to
// open their own handles.
func (s *BaseTestSuite) DBURL() string {
	return s.ts.DBURL
}

func (s *BaseTestSuite) RedisClient() *goredis.Client {
	return s.ts.RedisClient()
}

func (s *BaseTestSuite) RegisterServer(server *grpc.Server) string {
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Servers = append(s.Servers, server)
	return addr
}

func (s *BaseTestSuite) GetConnection(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) SeedUser(ctx context.Context) int {
	_, err := s.DBPool().Exec(ctx, `
		INSERT INTO roles (role_name, created_at, updated_at) 
		VALUES ('Admin Access 1', current_timestamp, current_timestamp),
		       ('ROLE_ADMIN', current_timestamp, current_timestamp)
		ON CONFLICT (role_name) DO NOTHING
	`)
	s.Require().NoError(err)

	res, err := NewUserClient(s.Conns["user"]).Create(ctx, &pbuser.CreateUserRequest{
		Firstname:       "Seed",
		Lastname:        "User",
		Email:           "seed.user@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedCategory(ctx context.Context) int {
	res, err := NewCategoryClient(s.Conns["category"]).Create(ctx, &pbcategory.CreateCategoryRequest{
		Name:        "Seed Category",
		Description: "Seed Description",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedMerchant(ctx context.Context, userID int) int {
	res, err := NewMerchantClient(s.Conns["merchant"]).Create(ctx, &pbmerchant.CreateMerchantRequest{
		UserId:       int32(userID),
		Name:         "Seed Merchant",
		Description:  "Seed Description",
		Address:      "Seed Address",
		ContactEmail: "merchant@example.com",
		ContactPhone: "08123456789",
		Status:       "active",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedProduct(ctx context.Context, merchantID int, categoryID int) int {
	res, err := NewProductClient(s.Conns["product"]).Create(ctx, &pbproduct.CreateProductRequest{
		MerchantId:   int32(merchantID),
		CategoryId:   int32(categoryID),
		Name:         "Seed Product",
		Description:  "Seed Description",
		Price:        10000,
		CountInStock: 100,
		Brand:        "Seed Brand",
		Weight:       1000,
		ImageProduct: "seed.jpg",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedOrder(ctx context.Context, userID int, merchID int, prodID int) int {
	var cashierID int
	err := s.DBPool().QueryRow(ctx, `
		SELECT cashier_id FROM cashiers WHERE user_id = $1 AND merchant_id = $2 AND deleted_at IS NULL LIMIT 1
	`, userID, merchID).Scan(&cashierID)
	if err != nil {
		err = s.DBPool().QueryRow(ctx, `
			INSERT INTO cashiers (merchant_id, user_id, name, created_at, updated_at)
			VALUES ($1, $2, 'Seed Cashier', current_timestamp, current_timestamp)
			RETURNING cashier_id
		`, merchID, userID).Scan(&cashierID)
		s.Require().NoError(err)
	}

	res, err := NewOrderClient(s.Conns["order"]).Create(ctx, &pborder.CreateOrderRequest{
		MerchantId: int32(merchID),
		CashierId:  int32(cashierID),
		Items: []*pborder.CreateOrderItemRequest{
			{
				ProductId: int32(prodID),
				Quantity:  1,
			},
		},
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedOrderItem(ctx context.Context, orderID int, productID int) int {
	var id int
	err := s.DBPool().QueryRow(ctx,
		`INSERT INTO "order_items" (order_id, product_id, quantity, price) VALUES ($1, $2, 1, 1000) RETURNING order_item_id`,
		orderID, productID,
	).Scan(&id)
	s.Require().NoError(err)
	return id
}
