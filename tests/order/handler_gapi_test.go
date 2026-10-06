package order_test

import (
	"context"
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	"testing"

	order_cache "github.com/MamangRust/monolith-pointofsale-order/cache"
	order_handler "github.com/MamangRust/monolith-pointofsale-order/handler"
	order_repo "github.com/MamangRust/monolith-pointofsale-order/repository"
	order_service "github.com/MamangRust/monolith-pointofsale-order/service"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderGapiTestSuite struct {
	tests.BaseTestSuite
	clientQuery           pborder.OrderQueryServiceClient
	clientCommand         pborder.OrderCommandServiceClient
	clientStats           pborder.OrderStatsServiceClient
	clientStatsByMerchant pborder.OrderStatsByMerchantServiceClient
}

func (s *OrderGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupCashierService()
	s.SetupTransactionService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := s.SQLxDB()

	// Order dependencies
	mencache := order_cache.NewMencache(cacheStore)
	repos := order_repo.NewRepositories(&order_repo.Deps{
		Db:                     queries,
		CashierQueryClient:     pbcashier.NewCashierQueryServiceClient(s.Conns["cashier"]),
		MerchantQueryClient:    pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		ProductQueryClient:     pbproduct.NewProductQueryServiceClient(s.Conns["product"]),
		ProductCommandClient:   pbproduct.NewProductCommandServiceClient(s.Conns["product"]),
		OrderItemQueryClient:   pb.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		OrderItemCommandClient: pb.NewOrderItemCommandServiceClient(s.Conns["order-item"]),
	})
	svc := order_service.NewService(&order_service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := order_handler.NewHandler(&order_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pborder.RegisterOrderQueryServiceServer(server, handler.Order)
	pborder.RegisterOrderCommandServiceServer(server, handler.OrderCommand)
	pborder.RegisterOrderStatsServiceServer(server, handler.OrderStats)
	pborder.RegisterOrderStatsByMerchantServiceServer(server, handler.OrderStatsByMerchant)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.clientQuery = pborder.NewOrderQueryServiceClient(conn)
	s.clientCommand = pborder.NewOrderCommandServiceClient(conn)
	s.clientStats = pborder.NewOrderStatsServiceClient(conn)
	s.clientStatsByMerchant = pborder.NewOrderStatsByMerchantServiceClient(conn)
}

func (s *OrderGapiTestSuite) TestOrderGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	// Seed a cashier (orders.cashier_id references cashiers.cashier_id)
	var cashierID int
	err := s.DBPool().QueryRow(ctx,
		`INSERT INTO cashiers (merchant_id, user_id, name) VALUES ($1, $2, 'Order Gapi Cashier') RETURNING cashier_id`,
		merchID, userID,
	).Scan(&cashierID)
	s.Require().NoError(err)

	// 2. Create
	createRes, err := s.clientCommand.Create(ctx, &pborder.CreateOrderRequest{
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
	s.Require().NotNil(createRes)
	orderID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.clientQuery.FindById(ctx, &pborder.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)
	s.Equal(int32(userID), getRes.Data.CashierId)

	// 4. FindAll
	allRes, err := s.clientQuery.FindAll(ctx, &pborder.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.clientQuery.FindByActive(ctx, &pborder.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	// Fetch order items first
	itemClient := pb.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	itemsRes, err := itemClient.FindOrderItemByOrder(ctx, &pb.FindByIdOrderItemRequest{Id: orderID})
	s.Require().NoError(err)
	s.NotEmpty(itemsRes.Data)
	orderItemID := itemsRes.Data[0].Id

	_, err = s.clientCommand.Update(ctx, &pborder.UpdateOrderRequest{
		OrderId: orderID,
		Items: []*pborder.UpdateOrderItemRequest{
			{
				OrderItemId: orderItemID,
				ProductId:   int32(prodID),
				Quantity:    1,
			},
		},
	})
	s.Require().NoError(err)

	// 7. Trash
	_, err = s.clientCommand.TrashedOrder(ctx, &pborder.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.clientCommand.FindByTrashed(ctx, &pborder.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.clientCommand.RestoreOrder(ctx, &pborder.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.clientCommand.TrashedOrder(ctx, &pborder.FindByIdOrderRequest{Id: orderID})
	_, err = s.clientCommand.DeleteOrderPermanent(ctx, &pborder.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.clientCommand.RestoreAllOrder(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.clientCommand.DeleteAllOrderPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestOrderGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderGapiTestSuite))
}
