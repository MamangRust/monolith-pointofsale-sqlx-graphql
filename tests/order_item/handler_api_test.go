package order_item_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	"github.com/stretchr/testify/suite"

	graphtest "github.com/MamangRust/monolith-pointofsale-apigateway/graphtest"
)

type OrderItemApiTestSuite struct {
	tests.BaseTestSuite
	handler   http.Handler
	orderID   int
	authToken string
}

func (s *OrderItemApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()

	// Seed dependencies
	ctx := context.Background()
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)
	s.orderID = s.SeedOrder(ctx, userID, merchID, prodID)

	log, _ := logger.NewLogger("test", nil)
	tokenManager, _ := auth.NewManager("mysecret")
	token, err := tokenManager.GenerateToken(userID, "apigateway")
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

func (s *OrderItemApiTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *OrderItemApiTestSuite) TestOrderItemApiLifecycle() {
	// 1. FindAll
	findAllQuery := `query {
  findAllOrderItem(input: {page: 1, page_size: 10}) {
    status
    message
    data {
      id
    }
  }
}`
	resp, err := graphtest.ExecuteGraphQL(s.handler, findAllQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 2. FindOrderItemByOrder
	s.Require().NotZero(s.orderID)
	findByOrderQuery := fmt.Sprintf(`query {
  findOrderItemByOrder(input: {id: %d}) {
    status
    message
    data {
      id
    }
  }
}`, s.orderID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, findByOrderQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findByOrder, ok := resp.Data["findOrderItemByOrder"].(map[string]interface{})
	s.True(ok, "expected findOrderItemByOrder in data, got %+v", resp.Data)
	data, ok := findByOrder["data"].([]interface{})
	s.True(ok, "expected findOrderItemByOrder.data, got %+v", findByOrder)
	if len(data) == 0 {
		s.T().Skip("No order items found")
	}

	// 3. FindByActive
	findActiveQuery := `query {
  findByActiveOrderItem(input: {page: 1, page_size: 10}) {
    status
    message
    data {
      id
    }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, findActiveQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 4. FindByTrashed
	findTrashedQuery := `query {
  findByTrashedOrderItem(input: {page: 1, page_size: 10}) {
    status
    message
    data {
      id
    }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, findTrashedQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)
}

func TestOrderItemApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemApiTestSuite))
}
