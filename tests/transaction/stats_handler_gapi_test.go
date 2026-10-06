package transaction_test

import (
	"context"
	"testing"
	"time"

	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	"github.com/stretchr/testify/suite"
)

type TransactionStatsGapiTestSuite struct {
	tests.BaseTestSuite
	client     pbtransaction.TransactionStatsServiceClient
	userID     int
	merchantID int
}

func (s *TransactionStatsGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	s.client = pbtransaction.NewTransactionStatsServiceClient(s.Conns["transaction"])

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Create a successful transaction
	_, err := s.DBPool().Exec(ctx, `
		INSERT INTO transactions (order_id, merchant_id, amount, payment_method, payment_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		orderID, s.merchantID, 10000, "credit_card", "success", time.Now())
	s.Require().NoError(err)

	_, err = s.DBPool().Exec(ctx, `
		INSERT INTO transactions (order_id, merchant_id, amount, payment_method, payment_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		orderID, s.merchantID, 5000, "bank_transfer", "failed", time.Now())
	s.Require().NoError(err)
}

func (s *TransactionStatsGapiTestSuite) TestFindMonthStatusSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransaction.FindMonthlyTransactionStatus{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.FindMonthStatusSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func (s *TransactionStatsGapiTestSuite) TestFindYearStatusSuccess() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransaction.FindYearlyTransactionStatus{
		Year: int32(year),
	}

	res, err := s.client.FindYearStatusSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func (s *TransactionStatsGapiTestSuite) TestFindMonthStatusFailed() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransaction.FindMonthlyTransactionStatus{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.FindMonthStatusFailed(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func (s *TransactionStatsGapiTestSuite) TestFindYearStatusFailed() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransaction.FindYearlyTransactionStatus{
		Year: int32(year),
	}

	res, err := s.client.FindYearStatusFailed(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func TestTransactionStatsGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsGapiTestSuite))
}
