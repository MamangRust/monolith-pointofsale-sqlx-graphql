package transaction_test

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

type TransactionApiTestSuite struct {
	tests.BaseTestSuite
	handler       http.Handler
	transactionID int
	userID        int
	cashierID     int
	merchID       int
	orderID       int
	authToken     string
}

func (s *TransactionApiTestSuite) SetupSuite() {
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
	merchID := s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, merchID, catID)
	orderID := s.SeedOrder(ctx, userID, merchID, prodID)
	s.SeedOrderItem(ctx, orderID, prodID)

	// cashier_id as seeded by SeedOrder (cashiers table)
	var cashierID int
	err := s.DBPool().QueryRow(ctx,
		`SELECT cashier_id FROM cashiers WHERE user_id = $1 AND merchant_id = $2 AND deleted_at IS NULL LIMIT 1`,
		userID, merchID,
	).Scan(&cashierID)
	s.Require().NoError(err)
	s.cashierID = cashierID

	s.userID = userID
	s.merchID = merchID
	s.orderID = orderID

	log, _ := logger.NewLogger("test", nil)
	tokenManager, _ := auth.NewManager("mysecret")
	token, err := tokenManager.GenerateToken(userID, "apigateway")
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

func (s *TransactionApiTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *TransactionApiTestSuite) TestTransactionApiLifecycle() {
	// 1. Create
	createQuery := fmt.Sprintf(`mutation {
  createTransaction(input: {
    orderId: %d
    cashierId: %d
    paymentMethod: "Transfer Bank"
    amount: 100000
    paymentStatus: "pending"
  }) {
    status
    message
    data {
      id
    }
  }
}`, s.orderID, s.cashierID)

	resp, err := graphtest.ExecuteGraphQL(s.handler, createQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	create, ok := resp.Data["createTransaction"].(map[string]interface{})
	s.True(ok, "expected createTransaction in data, got %+v", resp.Data)
	s.Equal("success", create["status"])

	data, ok := create["data"].(map[string]interface{})
	s.True(ok, "expected createTransaction.data, got %+v", create)
	s.transactionID = int(data["id"].(float64))

	// 2. FindAll
	findAllQuery := `query {
  findAllTransaction(input: {page: 1, pageSize: 10}) {
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

	// 3. FindById
	s.Require().NotZero(s.transactionID)
	findByIDQuery := fmt.Sprintf(`query {
  findByIdTransaction(input: {id: %d}) {
    status
    message
    data {
      id
    }
  }
}`, s.transactionID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, findByIDQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	found, ok := resp.Data["findByIdTransaction"].(map[string]interface{})
	s.True(ok, "expected findByIdTransaction in data, got %+v", resp.Data)
	data, ok = found["data"].(map[string]interface{})
	s.True(ok, "expected findByIdTransaction.data, got %+v", found)
	s.Equal(float64(s.transactionID), data["id"])

	// 4. FindByMerchant
	s.Require().NotZero(s.merchID)
	findByMerchantQuery := fmt.Sprintf(`query {
  findByMerchantTransaction(input: {merchantId: %d, page: 1, pageSize: 10}) {
    status
    message
    data {
      id
    }
  }
}`, s.merchID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, findByMerchantQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 5. FindByActive
	findActiveQuery := `query {
  findByActiveTransaction(input: {page: 1, pageSize: 10}) {
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

	// 6. Update
	updateQuery := fmt.Sprintf(`mutation {
  updateTransaction(input: {
    transactionId: %d
    orderId: %d
    cashierId: %d
    paymentMethod: "GOPAY"
    amount: 100000
  }) {
    status
    message
    data {
      id
    }
  }
}`, s.transactionID, s.orderID, s.cashierID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, updateQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 7. Trash
	trashQuery := fmt.Sprintf(`mutation {
  trashedTransaction(input: {id: %d}) {
    status
    message
    data {
      id
    }
  }
}`, s.transactionID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, trashQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 8. FindByTrashed
	findTrashedQuery := `query {
  findByTrashedTransaction(input: {page: 1, pageSize: 10}) {
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

	// 9. Restore
	restoreQuery := fmt.Sprintf(`mutation {
  restoreTransaction(input: {id: %d}) {
    status
    message
    data {
      id
    }
  }
}`, s.transactionID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 10. DeletePermanent (trash first, then delete permanent)
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashQuery, nil, s.authToken)
	s.NoError(err)

	deletePermQuery := fmt.Sprintf(`mutation {
  deleteTransactionPermanent(input: {id: %d}) {
    status
    message
  }
}`, s.transactionID)

	resp, err = graphtest.ExecuteGraphQL(s.handler, deletePermQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 11. RestoreAll
	restoreAllQuery := `mutation {
  restoreAllTransaction {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreAllQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	// 12. DeleteAll
	deleteAllQuery := `mutation {
  deleteAllTransactionPermanent {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteAllQuery, nil, s.authToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)
}

func (s *TransactionApiTestSuite) Test13_MonthlySuccessStats() {
	query := `query {
  findMonthStatusSuccess(input: {year: 2026, month: 4}) {
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

func TestTransactionApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionApiTestSuite))
}
