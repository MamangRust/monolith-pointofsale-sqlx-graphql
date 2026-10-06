package transaction_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	"github.com/MamangRust/monolith-pointofsale-transacton/repository"
	"github.com/stretchr/testify/suite"
)

type TransactionRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *TransactionRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()

	queries := s.SQLxDB()
	s.repo = repository.NewRepositories(queries, nil, nil, nil, nil)
}

func (s *TransactionRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *TransactionRepositoryTestSuite) TestTransactionLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, merchID, catID)
	orderID := s.SeedOrder(ctx, userID, merchID, prodID)

	pageReq := &requests.FindAllTransaction{Search: "", Page: 1, PageSize: 10}

	// 1. Create Transaction
	status := "paid"
	req := &requests.CreateTransactionRequest{
		OrderID:       orderID,
		CashierID:     userID,
		MerchantID:    merchID,
		PaymentMethod: "Credit Card",
		PaymentStatus: &status,
		Amount:        100000,
	}
	created, err := s.repo.TransactionCommandRepository.CreateTransaction(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(orderID), created.OrderID)
	txnID := int(created.TransactionID)

	// 2. FindByID
	found, err := s.repo.TransactionQueryRepository.FindById(ctx, txnID)
	s.NoError(err)
	s.Equal(created.Amount, found.Amount)

	// 3. FindByOrderID
	foundByOrder, err := s.repo.TransactionQueryRepository.FindByOrderId(ctx, orderID)
	s.NoError(err)
	s.NotNil(foundByOrder)

	// 4. FindAll
	all, _, err := s.repo.TransactionQueryRepository.FindAllTransactions(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 5. FindActive (before trash)
	active, _, err := s.repo.TransactionQueryRepository.FindByActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 6. Update
	updatedStatus := "settled"
	updateReq := &requests.UpdateTransactionRequest{
		TransactionID: &txnID,
		OrderID:       orderID,
		MerchantID:    merchID,
		PaymentMethod: "Debit Card",
		Amount:        120000,
		PaymentStatus: &updatedStatus,
	}
	updated, err := s.repo.TransactionCommandRepository.UpdateTransaction(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(120000), updated.Amount)

	// 7. Trash
	trashed, err := s.repo.TransactionCommandRepository.TrashTransaction(ctx, txnID)
	s.NoError(err)
	s.NotNil(trashed)

	// 8. FindTrashed
	trashedItems, _, err := s.repo.TransactionQueryRepository.FindByTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 9. FindActive after trash — should NOT include
	activeAfterTrash, _, err := s.repo.TransactionQueryRepository.FindByActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(txnID, int(item.TransactionID))
	}

	// 10. Restore
	restored, err := s.repo.TransactionCommandRepository.RestoreTransaction(ctx, txnID)
	s.NoError(err)
	s.NotNil(restored)

	// 11. Trash again then DeletePermanent
	_, err = s.repo.TransactionCommandRepository.TrashTransaction(ctx, txnID)
	s.Require().NoError(err)

	success, err := s.repo.TransactionCommandRepository.DeleteTransactionPermanently(ctx, txnID)
	s.NoError(err)
	s.True(success)

	// 12. RestoreAll
	status2 := "pending"
	second, _ := s.repo.TransactionCommandRepository.CreateTransaction(ctx, &requests.CreateTransactionRequest{
		CashierID: userID, OrderID: orderID, MerchantID: merchID,
		PaymentMethod: "Bank Transfer", PaymentStatus: &status2, Amount: 50000,
	})
	s.repo.TransactionCommandRepository.TrashTransaction(ctx, int(second.TransactionID))

	resRestore, err := s.repo.TransactionCommandRepository.RestoreAllTransactions(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 13. DeleteAll
	s.repo.TransactionCommandRepository.TrashTransaction(ctx, int(second.TransactionID))
	resDelete, err := s.repo.TransactionCommandRepository.DeleteAllTransactionPermanent(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestTransactionRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
