package transaction_test

import (
	"context"
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	trans_cache "github.com/MamangRust/monolith-pointofsale-transacton/cache"
	trans_handler "github.com/MamangRust/monolith-pointofsale-transacton/handler"
	trans_repo "github.com/MamangRust/monolith-pointofsale-transacton/repository"
	trans_service "github.com/MamangRust/monolith-pointofsale-transacton/service"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionGapiTestSuite struct {
	tests.BaseTestSuite
	clientQuery           pbtransaction.TransactionQueryServiceClient
	clientCommand         pbtransaction.TransactionCommandServiceClient
	clientStats           pbtransaction.TransactionStatsServiceClient
	clientStatsByMerchant pbtransaction.TransactionStatsByMerchantServiceClient
}

func (s *TransactionGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := s.SQLxDB()

	cashierClient := pbcashier.NewCashierQueryServiceClient(s.Conns["cashier"])
	merchantClient := pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"])
	orderClient := pborder.NewOrderQueryServiceClient(s.Conns["order"])
	orderItemClient := pb.NewOrderItemQueryServiceClient(s.Conns["order-item"])

	// Transaction dependencies
	mencache := trans_cache.NewMencache(cacheStore)
	repos := trans_repo.NewRepositories(queries, cashierClient, merchantClient, orderClient, orderItemClient)
	svc := trans_service.NewService(&trans_service.Deps{
		Kafka:         nil,
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := trans_handler.NewHandler(&trans_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbtransaction.RegisterTransactionQueryServiceServer(server, handler.Transaction)
	pbtransaction.RegisterTransactionCommandServiceServer(server, handler.TransactionCommand)
	pbtransaction.RegisterTransactionStatsServiceServer(server, handler.TransactionStats)
	pbtransaction.RegisterTransactionStatsByMerchantServiceServer(server, handler.TransactionStatsByMerchant)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.clientQuery = pbtransaction.NewTransactionQueryServiceClient(conn)
	s.clientCommand = pbtransaction.NewTransactionCommandServiceClient(conn)
	s.clientStats = pbtransaction.NewTransactionStatsServiceClient(conn)
	s.clientStatsByMerchant = pbtransaction.NewTransactionStatsByMerchantServiceClient(conn)
}

func (s *TransactionGapiTestSuite) TestTransactionGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, merchID, catID)
	orderID := s.SeedOrder(ctx, userID, merchID, prodID)

	// cashier_id as seeded by SeedOrder (cashiers table)
	var cashierID int
	err := s.DBPool().QueryRow(ctx,
		`SELECT cashier_id FROM cashiers WHERE user_id = $1 AND merchant_id = $2 AND deleted_at IS NULL LIMIT 1`,
		userID, merchID,
	).Scan(&cashierID)
	s.Require().NoError(err)

	// 2. Create
	createRes, err := s.clientCommand.Create(ctx, &pbtransaction.CreateTransactionRequest{
		OrderId:       int32(orderID),
		CashierId:     int32(cashierID),
		PaymentMethod: "E-Wallet",
		PaymentStatus: "pending",
		Amount:        100000,
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	transID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.clientQuery.FindById(ctx, &pbtransaction.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)
	s.Equal("E-Wallet", getRes.Data.PaymentMethod)

	// 4. FindAll
	allRes, err := s.clientQuery.FindAll(ctx, &pbtransaction.FindAllTransactionRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.clientQuery.FindByActive(ctx, &pbtransaction.FindAllTransactionRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.clientCommand.Update(ctx, &pbtransaction.UpdateTransactionRequest{
		TransactionId: transID,
		OrderId:       int32(orderID),
		CashierId:     int32(cashierID),
		PaymentMethod: "Credit Card",
		Amount:        150000,
	})
	s.Require().NoError(err)
	s.Equal("Credit Card", updateRes.Data.PaymentMethod)

	// 7. Trash
	_, err = s.clientCommand.TrashedTransaction(ctx, &pbtransaction.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.clientCommand.FindByTrashed(ctx, &pbtransaction.FindAllTransactionRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.clientCommand.RestoreTransaction(ctx, &pbtransaction.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.clientCommand.TrashedTransaction(ctx, &pbtransaction.FindByIdTransactionRequest{Id: transID})
	_, err = s.clientCommand.DeleteTransactionPermanent(ctx, &pbtransaction.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.clientCommand.RestoreAllTransaction(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.clientCommand.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestTransactionGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionGapiTestSuite))
}
