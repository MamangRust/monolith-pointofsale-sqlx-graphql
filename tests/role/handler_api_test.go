package role_test

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	role_cache "github.com/MamangRust/monolith-pointofsale-role/cache"
	role_handler "github.com/MamangRust/monolith-pointofsale-role/handler"
	"github.com/MamangRust/monolith-pointofsale-role/repository"
	"github.com/MamangRust/monolith-pointofsale-role/service"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	tests "github.com/MamangRust/monolith-pointofsale-test"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	graphtest "github.com/MamangRust/monolith-pointofsale-apigateway/graphtest"

	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
)

type RoleApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	handler     http.Handler
	tokenMgr    auth.TokenManager
	token       string
	grpcServer  *grpc.Server
	conn        *grpc.ClientConn
	roleID      int
}

func (s *RoleApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)
	s.redisClient.FlushAll(context.Background())

	queries, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)
	repos := repository.NewRepositories(queries)

	log, _ := logger.NewLogger("test", nil)
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	mencache := role_cache.NewMencache(cacheStore)

	roleService := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache,
		Observability: obs,
	})

	// Start internal gRPC Server for Role module
	roleHandlerGrpc := role_handler.NewHandler(&role_handler.Deps{
		Service: roleService,
		Logger:  log,
	})
	server := grpc.NewServer()
	pbrole.RegisterRoleQueryServiceServer(server, roleHandlerGrpc.Role)
	pbrole.RegisterRoleCommandServiceServer(server, roleHandlerGrpc.RoleCommand)
	s.grpcServer = server

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn

	tokenMgr, err := auth.NewManager("mysecret")
	s.Require().NoError(err)
	s.tokenMgr = tokenMgr
	s.token, err = tokenMgr.GenerateToken(1, "api")
	s.Require().NoError(err)

	// Build the GraphQL apigateway handler wired to the in-process gRPC servers.
	conns := &graphtest.ServiceConnections{
		RoleClient: conn,
	}
	resolver := graphtest.NewResolver(conns, log, s.redisClient)
	s.handler = graphtest.NewAuthHandler(resolver, tokenMgr, log)
}

func (s *RoleApiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	s.ts.Teardown()
}

func (s *RoleApiTestSuite) TestRoleApiLifecycle() {
	// 1. Create
	createQ := `mutation {
  createRole(input: {name: "API Role"}) {
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
	s.Empty(resp.Errors, "createRole errors: %+v", resp.Errors)
	create, ok := resp.Data["createRole"].(map[string]interface{})
	s.True(ok, "expected createRole in data, got %+v", resp.Data)
	createData := create["data"].(map[string]interface{})
	s.Equal("API Role", createData["name"])
	s.roleID = int(createData["id"].(float64))

	// 2. FindAll
	allQ := `query {
  findAllRole(input: {page: 1, page_size: 10}) {
    status
    message
    data { id name }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, allQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findAllRole errors: %+v", resp.Errors)

	// 3. FindById
	s.Require().NotZero(s.roleID)
	findQ := `query {
  findByIdRole(input: {role_id: ` + strconv.Itoa(s.roleID) + `}) {
    status
    message
    data { id name }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, findQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByIdRole errors: %+v", resp.Errors)
	found, ok := resp.Data["findByIdRole"].(map[string]interface{})
	s.True(ok, "expected findByIdRole in data, got %+v", resp.Data)
	foundData := found["data"].(map[string]interface{})
	s.Equal(float64(s.roleID), foundData["id"])

	// 4. FindByActive
	activeQ := `query {
  findByActiveRole(input: {page: 1, page_size: 10}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, activeQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByActiveRole errors: %+v", resp.Errors)

	// 5. FindByTrashed
	trashedQ := `query {
  findByTrashedRole(input: {page: 1, page_size: 10}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashedQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByTrashedRole errors: %+v", resp.Errors)

	// 6. FindByUserId
	byUserQ := `query {
  findByUserIdRole(input: {user_id: 1}) {
    status
    message
    data { id name }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, byUserQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "findByUserIdRole errors: %+v", resp.Errors)

	// 7. Update
	updateQ := `mutation {
  updateRole(input: {id: ` + strconv.Itoa(s.roleID) + `, name: "Updated API Role"}) {
    status
    message
    data { id name }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, updateQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "updateRole errors: %+v", resp.Errors)

	// 8. Restore
	restoreQ := `mutation {
  restoreRole(input: {role_id: ` + strconv.Itoa(s.roleID) + `}) {
    status
    message
    data { id }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "restoreRole errors: %+v", resp.Errors)

	// 9. DeletePermanent
	// Trash first then permanent
	trashQ := `mutation {
  trashedRole(input: {role_id: ` + strconv.Itoa(s.roleID) + `}) {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors)

	deleteQ := `mutation {
  deleteRolePermanent(input: {role_id: ` + strconv.Itoa(s.roleID) + `}) {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "deleteRolePermanent errors: %+v", resp.Errors)

	// 10. RestoreAll
	restoreAllQ := `mutation {
  restoreAllRole {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreAllQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "restoreAllRole errors: %+v", resp.Errors)

	// 11. DeleteAll
	deleteAllQ := `mutation {
  deleteAllRolePermanent {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteAllQ, nil, s.token)
	s.NoError(err)
	s.Empty(resp.Errors, "deleteAllRolePermanent errors: %+v", resp.Errors)
}

func TestRoleApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleApiTestSuite))
}
