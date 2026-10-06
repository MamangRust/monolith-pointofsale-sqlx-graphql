package transaction_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	"github.com/stretchr/testify/suite"

	graphtest "github.com/MamangRust/monolith-pointofsale-apigateway/graphtest"
)

type TransactionStatsApiTestSuite struct {
	tests.BaseTestSuite
	handler    http.Handler
	merchantID int
	userID     int
	authToken  string
}

func (s *TransactionStatsApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Seed a successful transaction
	_, err := s.DBPool().Exec(ctx, `
		INSERT INTO transactions (merchant_id, order_id, amount, payment_method, payment_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		s.merchantID, orderID, 100000, "credit_card", "success", time.Now())
	s.Require().NoError(err)

	log, _ := logger.NewLogger("test", nil)
	tokenManager, _ := auth.NewManager("mysecret")
	token, err := tokenManager.GenerateToken(s.userID, "apigateway")
	s.Require().NoError(err)
	s.authToken = token

	conns := &graphtest.ServiceConnections{
		TransactionClient: s.Conns["transaction"],
		OrderClient:       s.Conns["order"],
		OrderItemClient:   s.Conns["order-item"],
		CashierClient:     s.Conns["cashier"],
		MerchantClient:    s.Conns["merchant"],
		ProductClient:     s.Conns["product"],
		UserClient:        s.Conns["user"],
	}
	resolver := graphtest.NewResolver(conns, log, s.RedisClient())
	s.handler = graphtest.NewAuthHandler(resolver, tokenManager, log)
}

func (s *TransactionStatsApiTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *TransactionStatsApiTestSuite) TestFindMonthStatusSuccess() {
	now := time.Now()
	query := fmt.Sprintf(`query {
  findMonthStatusSuccess(input: {year: %d, month: %d}) {
    status
    message
    data {
      year
      month
    }
  }
}`, now.Year(), int(now.Month()))

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findMonth, ok := resp.Data["findMonthStatusSuccess"].(map[string]interface{})
	s.True(ok, "expected findMonthStatusSuccess in data, got %+v", resp.Data)
	s.Equal("success", findMonth["status"])
	s.NotEmpty(findMonth["data"])
}

func (s *TransactionStatsApiTestSuite) TestFindYearStatusSuccess() {
	year := time.Now().Year()
	query := fmt.Sprintf(`query {
  findYearStatusSuccess(input: {year: %d}) {
    status
    message
    data {
      year
    }
  }
}`, year)

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findYear, ok := resp.Data["findYearStatusSuccess"].(map[string]interface{})
	s.True(ok, "expected findYearStatusSuccess in data, got %+v", resp.Data)
	s.Equal("success", findYear["status"])
	s.NotEmpty(findYear["data"])
}

func (s *TransactionStatsApiTestSuite) TestFindMonthMethodSuccess() {
	now := time.Now()
	query := fmt.Sprintf(`query {
  findMonthMethodSuccess(input: {year: %d, month: %d}) {
    status
    message
    data {
      month
    }
  }
}`, now.Year(), int(now.Month()))

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findMethod, ok := resp.Data["findMonthMethodSuccess"].(map[string]interface{})
	s.True(ok, "expected findMonthMethodSuccess in data, got %+v", resp.Data)
	s.Equal("success", findMethod["status"])
	s.NotEmpty(findMethod["data"])
}

func (s *TransactionStatsApiTestSuite) TestFindMonthStatusSuccessByMerchant() {
	now := time.Now()
	query := fmt.Sprintf(`query {
  findMonthStatusSuccessByMerchant(input: {year: %d, month: %d, merchantId: %d}) {
    status
    message
    data {
      year
      month
    }
  }
}`, now.Year(), int(now.Month()), s.merchantID)

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findByMerchant, ok := resp.Data["findMonthStatusSuccessByMerchant"].(map[string]interface{})
	s.True(ok, "expected findMonthStatusSuccessByMerchant in data, got %+v", resp.Data)
	s.Equal("success", findByMerchant["status"])
	s.NotEmpty(findByMerchant["data"])
}

func TestTransactionStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsApiTestSuite))
}
