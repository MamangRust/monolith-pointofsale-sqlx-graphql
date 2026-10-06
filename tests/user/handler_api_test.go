package user_test

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

type UserHandlerTestSuite struct {
	tests.BaseTestSuite
	redisClient *redis.Client
	handler     http.Handler
	tokenMgr    auth.TokenManager
	token       string
	userID      int
	userEmail   string
}

func (s *UserHandlerTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	s.SetupRoleService()

	// Seed default role required by user service CreateUser
	_, err := s.DBPool().Exec(s.Ctx,
		`INSERT INTO roles (role_name, created_at, updated_at)
		 VALUES ('Admin Access 1', current_timestamp, current_timestamp)
		 ON CONFLICT (role_name) DO NOTHING`)
	s.Require().NoError(err)

	s.SetupUserService()

	s.redisClient = s.RedisClient()
	s.redisClient.FlushAll(context.Background())

	log, _ := logger.NewLogger("test", nil)

	tokenMgr, err := auth.NewManager("mysecret")
	s.Require().NoError(err)
	s.tokenMgr = tokenMgr
	s.token, err = tokenMgr.GenerateToken(1, "api")
	s.Require().NoError(err)

	conns := &graphtest.ServiceConnections{
		RoleClient: s.Conns["role"],
		UserClient: s.Conns["user"],
	}
	resolver := graphtest.NewResolver(conns, log, s.redisClient)
	s.handler = graphtest.NewAuthHandler(resolver, tokenMgr, log)
}

func (s *UserHandlerTestSuite) TearDownSuite() {
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	s.BaseTestSuite.TearDownSuite()
}

func (s *UserHandlerTestSuite) TestUserApiLifecycle() {
	// 1. Create
	s.userEmail = "handler.user@example.com"
	createQ := `mutation {
  createUser(input: {
    firstname: "Handler",
    lastname: "User",
    email: "` + s.userEmail + `",
    password: "password123",
    confirm_password: "password123"
  }) {
    status
    message
    data {
      id
      firstname
      lastname
      email
    }
  }
}`
	resp, err := graphtest.ExecuteGraphQL(s.handler, createQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "createUser errors: %+v", resp.Errors)
	create, ok := resp.Data["createUser"].(map[string]interface{})
	s.True(ok, "expected createUser in data, got %+v", resp.Data)
	createData := create["data"].(map[string]interface{})
	s.userID = int(createData["id"].(float64))

	// 2. FindAll
	allQ := `query {
  findAllUsers(input: {page: 1, page_size: 10}) {
    status
    message
    data { id email }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, allQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findAllUsers errors: %+v", resp.Errors)

	// 3. FindById
	s.Require().NotZero(s.userID)
	findQ := `query {
  findByIdUser(input: {id: ` + strconv.Itoa(s.userID) + `}) {
    status
    message
    data { id email }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, findQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByIdUser errors: %+v", resp.Errors)
	found, ok := resp.Data["findByIdUser"].(map[string]interface{})
	s.True(ok, "expected findByIdUser in data, got %+v", resp.Data)
	foundData := found["data"].(map[string]interface{})
	s.Equal(float64(s.userID), foundData["id"])

	// 4. FindByActive
	activeQ := `query {
  findByActiveUsers(input: {page: 1, page_size: 10}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, activeQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByActiveUsers errors: %+v", resp.Errors)

	// 5. Update
	updateQ := `mutation {
  updateUser(input: {
    id: ` + strconv.Itoa(s.userID) + `,
    firstname: "Updated",
    lastname: "User",
    email: "` + s.userEmail + `",
    password: "password123",
    confirm_password: "password123"
  }) {
    status
    message
    data { id firstname }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, updateQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "updateUser errors: %+v", resp.Errors)
	update, ok := resp.Data["updateUser"].(map[string]interface{})
	s.True(ok, "expected updateUser in data, got %+v", resp.Data)
	updateData := update["data"].(map[string]interface{})
	s.Equal("Updated", updateData["firstname"])

	// 6. Trash
	trashQ := `mutation {
  trashedUser(input: {id: ` + strconv.Itoa(s.userID) + `}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "trashedUser errors: %+v", resp.Errors)

	// 7. FindByTrashed
	trashedQ := `query {
  findByTrashedUsers(input: {page: 1, page_size: 10}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashedQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByTrashedUsers errors: %+v", resp.Errors)

	// 8. Restore
	restoreQ := `mutation {
  restoreUser(input: {id: ` + strconv.Itoa(s.userID) + `}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "restoreUser errors: %+v", resp.Errors)

	// 9. DeletePermanent
	// Trash first then permanent
	trashAgainQ := `mutation {
  trashedUser(input: {id: ` + strconv.Itoa(s.userID) + `}) {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashAgainQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors)

	deleteQ := `mutation {
  deleteUserPermanent(input: {id: ` + strconv.Itoa(s.userID) + `}) {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "deleteUserPermanent errors: %+v", resp.Errors)

	// 10. RestoreAll
	restoreAllQ := `mutation {
  restoreAllUser {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreAllQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "restoreAllUser errors: %+v", resp.Errors)

	// 11. DeleteAll
	deleteAllQ := `mutation {
  deleteAllUserPermanent {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteAllQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "deleteAllUserPermanent errors: %+v", resp.Errors)
}

func TestUserHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserHandlerTestSuite))
}
