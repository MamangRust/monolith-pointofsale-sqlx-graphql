package category_test

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

type CategoryStatsApiTestSuite struct {
	tests.BaseTestSuite
	handler    http.Handler
	categoryID int
	merchantID int
	userID     int
}

func (s *CategoryStatsApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	log, _ := logger.NewLogger("test", nil)

	// Build the GraphQL apigateway handler wired to the in-process gRPC servers.
	conns := &graphtest.ServiceConnections{
		CategoryClient: s.Conns["category"],
	}
	resolver := graphtest.NewResolver(conns, log, s.RedisClient())
	tokenManager, _ := auth.NewManager("mysecret")
	s.handler = graphtest.NewAuthHandler(resolver, tokenManager, log)

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	s.categoryID = s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, s.categoryID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	_, err := s.DBPool().Exec(ctx, "UPDATE orders SET created_at = $1 WHERE order_id = $2",
		time.Now(), orderID)
	s.Require().NoError(err)
}

func (s *CategoryStatsApiTestSuite) TestFindMonthTotalPrice() {
	tokenManager, _ := auth.NewManager("mysecret")
	token, _ := tokenManager.GenerateToken(1, "test")
	now := time.Now()

	query := fmt.Sprintf(`query {
  findMonthlyTotalPrices(input: {year: %d, month: %d}) {
    status
    message
    data {
      year
      month
      total_revenue
    }
  }
}`, now.Year(), int(now.Month()))

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findMonthly, ok := resp.Data["findMonthlyTotalPrices"].(map[string]interface{})
	s.True(ok, "expected findMonthlyTotalPrices in data, got %+v", resp.Data)
	s.Equal("success", findMonthly["status"])
	s.NotEmpty(findMonthly["data"])
}

func (s *CategoryStatsApiTestSuite) TestFindYearTotalPrice() {
	tokenManager, _ := auth.NewManager("mysecret")
	token, _ := tokenManager.GenerateToken(1, "test")
	year := time.Now().Year()

	query := fmt.Sprintf(`query {
  findYearlyTotalPrices(input: {year: %d}) {
    status
    message
    data {
      year
      total_revenue
    }
  }
}`, year)

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findYearly, ok := resp.Data["findYearlyTotalPrices"].(map[string]interface{})
	s.True(ok, "expected findYearlyTotalPrices in data, got %+v", resp.Data)
	s.Equal("success", findYearly["status"])
	s.NotEmpty(findYearly["data"])
}

func (s *CategoryStatsApiTestSuite) TestFindMonthPrice() {
	tokenManager, _ := auth.NewManager("mysecret")
	token, _ := tokenManager.GenerateToken(1, "test")
	year := time.Now().Year()

	query := fmt.Sprintf(`query {
  findMonthPrice(input: {year: %d}) {
    status
    message
    data {
      month
      category_id
      category_name
      order_count
      items_sold
      total_revenue
    }
  }
}`, year)

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findMonthPrice, ok := resp.Data["findMonthPrice"].(map[string]interface{})
	s.True(ok, "expected findMonthPrice in data, got %+v", resp.Data)
	s.Equal("success", findMonthPrice["status"])
	s.NotEmpty(findMonthPrice["data"])
}

func (s *CategoryStatsApiTestSuite) TestFindYearPrice() {
	tokenManager, _ := auth.NewManager("mysecret")
	token, _ := tokenManager.GenerateToken(1, "test")
	year := time.Now().Year()

	query := fmt.Sprintf(`query {
  findYearPrice(input: {year: %d}) {
    status
    message
    data {
      year
      category_id
      category_name
      order_count
      items_sold
      total_revenue
      unique_products_sold
    }
  }
}`, year)

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findYearPrice, ok := resp.Data["findYearPrice"].(map[string]interface{})
	s.True(ok, "expected findYearPrice in data, got %+v", resp.Data)
	s.Equal("success", findYearPrice["status"])
	s.NotEmpty(findYearPrice["data"])
}

func (s *CategoryStatsApiTestSuite) TestFindMonthTotalPriceById() {
	tokenManager, _ := auth.NewManager("mysecret")
	token, _ := tokenManager.GenerateToken(1, "test")
	now := time.Now()

	query := fmt.Sprintf(`query {
  findMonthlyTotalPricesById(input: {year: %d, month: %d, category_id: %d}) {
    status
    message
    data {
      year
      month
      total_revenue
    }
  }
}`, now.Year(), int(now.Month()), s.categoryID)

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findMonthlyByID, ok := resp.Data["findMonthlyTotalPricesById"].(map[string]interface{})
	s.True(ok, "expected findMonthlyTotalPricesById in data, got %+v", resp.Data)
	s.Equal("success", findMonthlyByID["status"])
	s.NotEmpty(findMonthlyByID["data"])
}

func (s *CategoryStatsApiTestSuite) TestFindMonthTotalPriceByMerchant() {
	tokenManager, _ := auth.NewManager("mysecret")
	token, _ := tokenManager.GenerateToken(1, "test")
	now := time.Now()

	query := fmt.Sprintf(`query {
  findMonthlyTotalPricesByMerchant(input: {year: %d, month: %d, merchant_id: %d}) {
    status
    message
    data {
      year
      month
      total_revenue
    }
  }
}`, now.Year(), int(now.Month()), s.merchantID)

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findMonthlyByMerchant, ok := resp.Data["findMonthlyTotalPricesByMerchant"].(map[string]interface{})
	s.True(ok, "expected findMonthlyTotalPricesByMerchant in data, got %+v", resp.Data)
	s.Equal("success", findMonthlyByMerchant["status"])
	s.NotEmpty(findMonthlyByMerchant["data"])
}

func TestCategoryStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryStatsApiTestSuite))
}
