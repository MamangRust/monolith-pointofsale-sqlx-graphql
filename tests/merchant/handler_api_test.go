package merchant_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"

	graphtest "github.com/MamangRust/monolith-pointofsale-apigateway/graphtest"
)

type MerchantApiTestSuite struct {
	tests.BaseTestSuite
	redisClient *redis.Client
	handler     http.Handler
	tokenMgr    auth.TokenManager
	token       string
	merchantID  int
	userID      int
}

func (s *MerchantApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()

	// Seed user
	s.userID = s.SeedUser(context.Background())

	s.redisClient = s.RedisClient()
	s.redisClient.FlushAll(context.Background())

	log, _ := logger.NewLogger("test", nil)

	tokenMgr, err := auth.NewManager("mysecret")
	s.Require().NoError(err)
	s.tokenMgr = tokenMgr
	s.token, err = tokenMgr.GenerateToken(s.userID, "api")
	s.Require().NoError(err)

	conns := &graphtest.ServiceConnections{
		UserClient:     s.Conns["user"],
		MerchantClient: s.Conns["merchant"],
	}
	resolver := graphtest.NewResolver(conns, log, s.redisClient)
	s.handler = graphtest.NewAuthHandler(resolver, tokenMgr, log)
}

func (s *MerchantApiTestSuite) TearDownSuite() {
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	s.BaseTestSuite.TearDownSuite()
}

func (s *MerchantApiTestSuite) TestMerchantApiLifecycle() {
	// 1. Create
	createQ := `mutation {
  createMerchant(input: {
    user_id: ` + strconv.Itoa(s.userID) + `,
    name: "Test Merchant",
    description: "Test Description",
    address: "Test Address",
    contact_email: "merchant@example.com",
    contact_phone: "123456789",
    status: "active"
  }) {
    status
    message
    data {
      id
      name
    }
  }
}`
	resp, err := graphtest.ExecuteGraphQL(s.handler, createQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "createMerchant errors: %+v", resp.Errors)
	create, ok := resp.Data["createMerchant"].(map[string]interface{})
	s.True(ok, "expected createMerchant in data, got %+v", resp.Data)
	createData := create["data"].(map[string]interface{})
	s.Equal("Test Merchant", createData["name"])
	s.merchantID = int(createData["id"].(float64))

	// 2. FindById
	findQ := `query {
  findByIdMerchant(input: {id: ` + strconv.Itoa(s.merchantID) + `}) {
    status
    message
    data {
      id
      name
    }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, findQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByIdMerchant errors: %+v", resp.Errors)
	found, ok := resp.Data["findByIdMerchant"].(map[string]interface{})
	s.True(ok, "expected findByIdMerchant in data, got %+v", resp.Data)
	foundData := found["data"].(map[string]interface{})
	s.Equal(float64(s.merchantID), foundData["id"])

	// 3. FindAll
	allQ := `query {
  findAllMerchant(input: {page: 1, page_size: 10}) {
    status
    message
    data { id name }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, allQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findAllMerchant errors: %+v", resp.Errors)

	// 4. FindByActive
	activeQ := `query {
  findByActiveMerchant(input: {page: 1, page_size: 10}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, activeQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByActiveMerchant errors: %+v", resp.Errors)

	// 5. Update
	updateQ := `mutation {
  updateMerchant(input: {
    merchant_id: ` + strconv.Itoa(s.merchantID) + `,
    user_id: ` + strconv.Itoa(s.userID) + `,
    name: "Updated Merchant",
    description: "Updated Description",
    address: "Updated Address",
    contact_email: "updated@example.com",
    contact_phone: "987654321",
    status: "active"
  }) {
    status
    message
    data { id name }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, updateQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "updateMerchant errors: %+v", resp.Errors)
	update, ok := resp.Data["updateMerchant"].(map[string]interface{})
	s.True(ok, "expected updateMerchant in data, got %+v", resp.Data)
	updateData := update["data"].(map[string]interface{})
	s.Equal("Updated Merchant", updateData["name"])

	// 6. Trash
	trashQ := `mutation {
  trashedMerchant(input: {id: ` + strconv.Itoa(s.merchantID) + `}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "trashedMerchant errors: %+v", resp.Errors)

	// 7. FindByTrashed
	trashedQ := `query {
  findByTrashedMerchant(input: {page: 1, page_size: 10}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashedQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByTrashedMerchant errors: %+v", resp.Errors)

	// 8. Restore
	restoreQ := `mutation {
  restoreMerchant(input: {id: ` + strconv.Itoa(s.merchantID) + `}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "restoreMerchant errors: %+v", resp.Errors)

	// 9. DeletePermanent
	// Trash first then permanent
	trashAgainQ := `mutation {
  trashedMerchant(input: {id: ` + strconv.Itoa(s.merchantID) + `}) {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashAgainQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors)

	deleteQ := `mutation {
  deleteMerchantPermanent(input: {id: ` + strconv.Itoa(s.merchantID) + `}) {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "deleteMerchantPermanent errors: %+v", resp.Errors)

	// 10. RestoreAll
	restoreAllQ := `mutation {
  restoreAllMerchant {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreAllQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "restoreAllMerchant errors: %+v", resp.Errors)

	// 11. DeleteAll
	deleteAllQ := `mutation {
  deleteAllMerchantPermanent {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteAllQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "deleteAllMerchantPermanent errors: %+v", resp.Errors)
}

func TestMerchantApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantApiTestSuite))
}
