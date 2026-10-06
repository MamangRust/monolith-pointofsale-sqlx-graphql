package order_test

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

type OrderApiTestSuite struct {
	tests.BaseTestSuite
	handler   http.Handler
	orderID   int
	userID    int
	cashierID int
	merchID   int
	prodID    int
	authToken string
}

func (s *OrderApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupCategoryService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()

	// Seed dependencies
	ctx := context.Background()
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	// Seed a cashier (orders.cashier_id references cashiers.cashier_id)
	var cashierID int
	err := s.DBPool().QueryRow(ctx,
		`INSERT INTO cashiers (merchant_id, user_id, name) VALUES ($1, $2, 'Order Api Cashier') RETURNING cashier_id`,
		merchID, userID,
	).Scan(&cashierID)
	s.Require().NoError(err)

	s.userID = userID
	s.cashierID = cashierID
	s.merchID = merchID
	s.prodID = prodID

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

func (s *OrderApiTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *OrderApiTestSuite) TestOrderApiLifecycle() {
	// 1. Create
	createQuery := fmt.Sprintf(`mutation {
  createOrder(input: {
    merchant_id: %d
    cashier_id: %d
    items: [{product_id: %d, quantity: 1}]
  }) {
    status
    message
    data {
      id
    }
  }
}`, s.merchID, s.cashierID, s.prodID)

	resp, err := graphtest.ExecuteGraphQL(s.handler, createQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	create, ok := resp.Data["createOrder"].(map[string]interface{})
	s.True(ok, "expected createOrder in data, got %+v", resp.Data)
	s.Equal("success", create["status"])

	data, ok := create["data"].(map[string]interface{})
	s.True(ok, "expected createOrder.data, got %+v", create)
	s.orderID = int(data["id"].(float64))

	// 2. FindById
	findByIDQuery := fmt.Sprintf(`query {
  findByIdOrder(input: {id: %d}) {
    status
    message
    data {
      id
    }
  }
}`, s.orderID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, findByIDQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	found, ok := resp.Data["findByIdOrder"].(map[string]interface{})
	s.True(ok, "expected findByIdOrder in data, got %+v", resp.Data)
	data, ok = found["data"].(map[string]interface{})
	s.True(ok, "expected findByIdOrder.data, got %+v", found)
	s.Equal(float64(s.orderID), data["id"])

	// 3. FindAll
	findAllQuery := `query {
  findAllOrder(input: {page: 1, page_size: 10}) {
    status
    message
    data {
      id
    }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, findAllQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 4. FindByActive
	findActiveQuery := `query {
  findByActiveOrder(input: {page: 1, page_size: 10}) {
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

	// 5. Update — fetch order_item_id from DB first
	var orderItemID int
	err = s.DBPool().QueryRow(context.Background(),
		`SELECT order_item_id FROM order_items WHERE order_id = $1 LIMIT 1`, s.orderID,
	).Scan(&orderItemID)
	s.Require().NoError(err)

	updateQuery := fmt.Sprintf(`mutation {
  updateOrder(input: {
    order_id: %d
    items: [{order_item_id: %d, product_id: %d, quantity: 1}]
  }) {
    status
    message
    data {
      id
    }
  }
}`, s.orderID, orderItemID, s.prodID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, updateQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 6. Trash
	trashQuery := fmt.Sprintf(`mutation {
  trashedOrder(input: {id: %d}) {
    status
    message
    data {
      id
    }
  }
}`, s.orderID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, trashQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 7. FindByTrashed
	findTrashedQuery := `query {
  findByTrashedOrder(input: {page: 1, page_size: 10}) {
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

	// 8. Restore
	restoreQuery := fmt.Sprintf(`mutation {
  restoreOrder(input: {id: %d}) {
    status
    message
    data {
      id
    }
  }
}`, s.orderID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 9. DeletePermanent (trash first, then delete permanent)
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashQuery, nil, s.authToken)
	s.NoError(err)

	deletePermQuery := fmt.Sprintf(`mutation {
  deleteOrderPermanent(input: {id: %d}) {
    status
    message
  }
}`, s.orderID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, deletePermQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 10. RestoreAll
	restoreAllQuery := `mutation {
  restoreAllOrder {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreAllQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "restoreAllOrder with no trashed orders is a no-op success, got: %+v", resp)

	// 11. DeleteAll — no trashed orders remain, still a no-op success.
	deleteAllQuery := `mutation {
  deleteAllOrderPermanent {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteAllQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "deleteAllOrderPermanent with no trashed orders is a no-op success, got: %+v", resp)
}

func (s *OrderApiTestSuite) Test12_GetMonthlyTotalRevenue() {
	query := `query {
  findMonthlyTotalRevenue(input: {year: 2024, month: 1}) {
    status
    message
    data {
      year
      month
    }
  }
}`
	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)
}

func (s *OrderApiTestSuite) Test13_GetYearlyTotalRevenue() {
	query := `query {
  findYearlyTotalRevenue(input: {year: 2024}) {
    status
    message
    data {
      year
    }
  }
}`
	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)
}

func TestOrderApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderApiTestSuite))
}
