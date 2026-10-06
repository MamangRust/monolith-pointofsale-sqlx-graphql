package transaction_test

import (
	"context"
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"

	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	trans_cache "github.com/MamangRust/monolith-pointofsale-transacton/cache"
	"github.com/MamangRust/monolith-pointofsale-transacton/repository"
	"github.com/MamangRust/monolith-pointofsale-transacton/service"
	"github.com/stretchr/testify/suite"
)

type TransactionServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *TransactionServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// gRPC service nyata: cashier(→user/merchant), merchant, order(→product/order-item), category.
	s.SetupTransactionService()
	s.SetupCategoryService()

	queries := s.SQLxDB()

	// Transaction repositories with real gRPC clients
	mencache := trans_cache.NewMencache(s.GetCacheStore())
	repos := repository.NewRepositories(
		queries,
		pbcashier.NewCashierQueryServiceClient(s.Conns["cashier"]),
		pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		pborder.NewOrderQueryServiceClient(s.Conns["order"]),
		pb.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
	)

	s.svc = service.NewService(&service.Deps{
		Kafka:         nil,
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
}

func (s *TransactionServiceTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *TransactionServiceTestSuite) TestTransactionLifecycle() {
	ctx := context.Background()

	userID := s.SeedUser(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	categoryID := s.SeedCategory(ctx)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)
	s.SeedOrderItem(ctx, orderID, productID)

	// cashier_id as seeded by SeedOrder (cashiers table)
	var cashierID int
	err := s.DBPool().QueryRow(ctx,
		`SELECT cashier_id FROM cashiers WHERE user_id = $1 AND merchant_id = $2 AND deleted_at IS NULL LIMIT 1`,
		userID, merchantID,
	).Scan(&cashierID)
	s.Require().NoError(err)

	// 2. Create Transaction
	req := &requests.CreateTransactionRequest{
		OrderID:       orderID,
		CashierID:     cashierID,
		MerchantID:    merchantID,
		PaymentMethod: "Transfer Bank",
		Amount:        1000000,
	}
	created, err := s.svc.TransactionCommand.CreateTransaction(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	transactionID := int(created.TransactionID)

	// 3. FindByID
	found, err := s.svc.TransactionQuery.FindById(ctx, transactionID)
	s.Require().NoError(err)
	s.Equal("success", found.PaymentStatus)

	// 4. Update
	newPaymentMethod := "GOPAY"
	updateReq := &requests.UpdateTransactionRequest{
		TransactionID: &transactionID,
		OrderID:       orderID,
		CashierID:     cashierID,
		PaymentMethod: newPaymentMethod,
		Amount:        req.Amount,
	}
	updated, err := s.svc.TransactionCommand.UpdateTransaction(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(newPaymentMethod, updated.PaymentMethod)

	// 5. FindAll
	_, total, err := s.svc.TransactionQuery.FindAllTransactions(ctx, &requests.FindAllTransaction{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 6. Trash
	_, err = s.svc.TransactionCommand.TrashedTransaction(ctx, transactionID)
	s.Require().NoError(err)

	// 7. FindTrashed
	_, totalTrashed, err := s.svc.TransactionQuery.FindByTrashed(ctx, &requests.FindAllTransaction{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 8. FindActive
	active, _, err := s.svc.TransactionQuery.FindByActive(ctx, &requests.FindAllTransaction{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, tx := range active {
		s.NotEqual(transactionID, int(tx.TransactionID))
	}

	// 9. Restore
	_, err = s.svc.TransactionCommand.RestoreTransaction(ctx, transactionID)
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, err = s.svc.TransactionCommand.TrashedTransaction(ctx, transactionID)
	s.Require().NoError(err)
	success, err := s.svc.TransactionCommand.DeleteTransactionPermanently(ctx, transactionID)
	s.Require().NoError(err)
	s.True(success)

	// 11. RestoreAll & DeleteAll
	o1 := s.SeedOrder(ctx, userID, merchantID, productID)
	s.SeedOrderItem(ctx, o1, productID)

	o2 := s.SeedOrder(ctx, userID, merchantID, productID)
	s.SeedOrderItem(ctx, o2, productID)

	s.svc.TransactionCommand.TrashedTransaction(ctx, o1)
	s.svc.TransactionCommand.TrashedTransaction(ctx, o2)

	resRestoreAll, err := s.svc.TransactionCommand.RestoreAllTransactions(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.TransactionCommand.TrashedTransaction(ctx, o1)
	s.svc.TransactionCommand.TrashedTransaction(ctx, o2)

	resDeleteAll, err := s.svc.TransactionCommand.DeleteAllTransactionPermanent(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestTransactionServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionServiceTestSuite))
}
