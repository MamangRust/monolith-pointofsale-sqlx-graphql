package order_test

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

type OrderStatsApiTestSuite struct {
	tests.BaseTestSuite
	handler    http.Handler
	merchantID int
	userID     int
	authToken  string
}

func (s *OrderStatsApiTestSuite) SetupSuite() {
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

	// Ensure created_at is set to current time to be picked up by stats
	_, err := s.DBPool().Exec(ctx, "UPDATE orders SET created_at = $1 WHERE order_id = $2",
		time.Now(), orderID)
	s.Require().NoError(err)

	log, _ := logger.NewLogger("test", nil)
	tokenManager, _ := auth.NewManager("mysecret")
	token, err := tokenManager.GenerateToken(s.userID, "apigateway")
	s.Require().NoError(err)
	s.authToken = token

	conns := &graphtest.ServiceConnections{
		OrderClient:     s.Conns["order"],
		OrderItemClient: s.Conns["order-item"],
		CashierClient:   s.Conns["cashier"],
		MerchantClient:  s.Conns["merchant"],
		ProductClient:   s.Conns["product"],
		UserClient:      s.Conns["user"],
	}
	resolver := graphtest.NewResolver(conns, log, s.RedisClient())
	s.handler = graphtest.NewAuthHandler(resolver, tokenManager, log)
}

func (s *OrderStatsApiTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyTotalRevenue() {
	now := time.Now()
	query := fmt.Sprintf(`query {
  findMonthlyTotalRevenue(input: {year: %d, month: %d}) {
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

	findMonthly, ok := resp.Data["findMonthlyTotalRevenue"].(map[string]interface{})
	s.True(ok, "expected findMonthlyTotalRevenue in data, got %+v", resp.Data)
	s.Equal("success", findMonthly["status"])
	s.NotEmpty(findMonthly["data"])
}

func (s *OrderStatsApiTestSuite) TestFindYearlyTotalRevenue() {
	year := time.Now().Year()
	query := fmt.Sprintf(`query {
  findYearlyTotalRevenue(input: {year: %d}) {
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

	findYearly, ok := resp.Data["findYearlyTotalRevenue"].(map[string]interface{})
	s.True(ok, "expected findYearlyTotalRevenue in data, got %+v", resp.Data)
	s.Equal("success", findYearly["status"])
	s.NotEmpty(findYearly["data"])
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyOrder() {
	year := time.Now().Year()
	query := fmt.Sprintf(`query {
  findMonthlyRevenue(input: {year: %d}) {
    status
    message
    data {
      month
    }
  }
}`, year)

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findMonthly, ok := resp.Data["findMonthlyRevenue"].(map[string]interface{})
	s.True(ok, "expected findMonthlyRevenue in data, got %+v", resp.Data)
	s.Equal("success", findMonthly["status"])
	s.NotEmpty(findMonthly["data"])
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyTotalRevenueByMerchant() {
	now := time.Now()
	query := fmt.Sprintf(`query {
  findMonthlyTotalRevenueByMerchant(input: {year: %d, month: %d, merchant_id: %d}) {
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

	findByMerchant, ok := resp.Data["findMonthlyTotalRevenueByMerchant"].(map[string]interface{})
	s.True(ok, "expected findMonthlyTotalRevenueByMerchant in data, got %+v", resp.Data)
	s.Equal("success", findByMerchant["status"])
	s.NotEmpty(findByMerchant["data"])
}

func TestOrderStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsApiTestSuite))
}
