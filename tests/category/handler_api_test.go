package category_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
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

	cat_cache "github.com/MamangRust/monolith-pointofsale-category/cache"
	cat_handler "github.com/MamangRust/monolith-pointofsale-category/handler"
	cat_repo "github.com/MamangRust/monolith-pointofsale-category/repository"
	cat_service "github.com/MamangRust/monolith-pointofsale-category/service"

	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
)

type CategoryApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	handler     http.Handler
	categoryID  int
}

func (s *CategoryApiTestSuite) SetupSuite() {
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

	log, _ := logger.NewLogger("test", nil)
	queries, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	obs, _ := observability.NewObservability("test", log)

	// Setup Category Service & gRPC Server
	catMencache := cat_cache.NewMencache(cacheStore)
	catRepos := cat_repo.NewRepositories(queries)
	catSvc := cat_service.NewService(&cat_service.Deps{
		Mencache:      catMencache,
		Repositories:  catRepos,
		Logger:        log,
		Observability: obs,
	})
	catGapi := cat_handler.NewHandler(&cat_handler.Deps{
		Service: catSvc,
		Logger:  log,
	})
	catServer := grpc.NewServer()
	pbcategory.RegisterCategoryQueryServiceServer(catServer, catGapi.Category)
	pbcategory.RegisterCategoryCommandServiceServer(catServer, catGapi.CategoryCommand)
	pbcategory.RegisterCategoryStatsServiceServer(catServer, catGapi.CategoryStats)
	catLis, _ := net.Listen("tcp", "localhost:0")
	go catServer.Serve(catLis)
	catConn, _ := grpc.NewClient(catLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Build the GraphQL apigateway handler wired to the in-process gRPC servers.
	conns := &graphtest.ServiceConnections{
		CategoryClient: catConn,
	}
	resolver := graphtest.NewResolver(conns, log, s.redisClient)
	tokenManager, _ := auth.NewManager("mysecret")
	s.handler = graphtest.NewAuthHandler(resolver, tokenManager, log)
}

func (s *CategoryApiTestSuite) TearDownSuite() {
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	s.ts.Teardown()
}

func (s *CategoryApiTestSuite) TestCategoryApiLifecycle() {
	tokenManager, _ := auth.NewManager("mysecret")
	token, _ := tokenManager.GenerateToken(1, "test")

	// 1. Create
	createQuery := `mutation {
  createCategory(input: {
    name: "Test Category"
    description: "Test Description"
  }) {
    status
    message
    data {
      id
      name
      description
    }
  }
}`
	resp, err := graphtest.ExecuteGraphQL(s.handler, createQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	create, ok := resp.Data["createCategory"].(map[string]interface{})
	s.True(ok, "expected createCategory in data, got %+v", resp.Data)
	s.Equal("success", create["status"])

	catData, ok := create["data"].(map[string]interface{})
	s.True(ok, "expected createCategory.data, got %+v", create)
	s.categoryID = int(catData["id"].(float64))
	s.Equal("Test Category", catData["name"])

	// 2. FindById
	findByIDQuery := fmt.Sprintf(`query {
  findByIdCategory(input: {id: %d}) {
    status
    message
    data {
      id
      name
      description
    }
  }
}`, s.categoryID)
	resp, err = graphtest.ExecuteGraphQL(s.handler, findByIDQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findByID, ok := resp.Data["findByIdCategory"].(map[string]interface{})
	s.True(ok, "expected findByIdCategory in data, got %+v", resp.Data)
	s.Equal("success", findByID["status"])

	foundData, ok := findByID["data"].(map[string]interface{})
	s.True(ok, "expected findByIdCategory.data, got %+v", findByID)
	s.Equal(float64(s.categoryID), foundData["id"])

	// 3. FindAll
	findAllQuery := `query {
  findAllCategory(input: {page: 1, page_size: 10}) {
    status
    message
    data {
      id
      name
    }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, findAllQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findAll, ok := resp.Data["findAllCategory"].(map[string]interface{})
	s.True(ok, "expected findAllCategory in data, got %+v", resp.Data)
	s.Equal("success", findAll["status"])

	// 4. FindByActive
	findByActiveQuery := `query {
  findByActiveCategory(input: {page: 1, page_size: 10}) {
    status
    message
    data {
      id
      name
    }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, findByActiveQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findByActive, ok := resp.Data["findByActiveCategory"].(map[string]interface{})
	s.True(ok, "expected findByActiveCategory in data, got %+v", resp.Data)
	s.Equal("success", findByActive["status"])

	// 5. Update
	updateQuery := fmt.Sprintf(`mutation {
  updateCategory(input: {
    category_id: %d
    name: "Updated Category"
    description: "Updated Description"
  }) {
    status
    message
    data {
      id
      name
      description
    }
  }
}`, s.categoryID)
	resp, err = graphtest.ExecuteGraphQL(s.handler, updateQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	update, ok := resp.Data["updateCategory"].(map[string]interface{})
	s.True(ok, "expected updateCategory in data, got %+v", resp.Data)
	s.Equal("success", update["status"])

	updData, ok := update["data"].(map[string]interface{})
	s.True(ok, "expected updateCategory.data, got %+v", update)
	s.Equal("Updated Category", updData["name"])

	// 6. Trash
	trashQuery := fmt.Sprintf(`mutation {
  trashedCategory(input: {id: %d}) {
    status
    message
    data {
      id
      name
    }
  }
}`, s.categoryID)
	resp, err = graphtest.ExecuteGraphQL(s.handler, trashQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	trash, ok := resp.Data["trashedCategory"].(map[string]interface{})
	s.True(ok, "expected trashedCategory in data, got %+v", resp.Data)
	s.Equal("success", trash["status"])

	// 7. FindByTrashed
	findByTrashedQuery := `query {
  findByTrashedCategory(input: {page: 1, page_size: 10}) {
    status
    message
    data {
      id
      name
    }
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, findByTrashedQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	findByTrashed, ok := resp.Data["findByTrashedCategory"].(map[string]interface{})
	s.True(ok, "expected findByTrashedCategory in data, got %+v", resp.Data)
	s.Equal("success", findByTrashed["status"])

	// 8. Restore
	restoreQuery := fmt.Sprintf(`mutation {
  restoreCategory(input: {id: %d}) {
    status
    message
    data {
      id
      name
    }
  }
}`, s.categoryID)
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	restore, ok := resp.Data["restoreCategory"].(map[string]interface{})
	s.True(ok, "expected restoreCategory in data, got %+v", resp.Data)
	s.Equal("success", restore["status"])

	// 9. DeletePermanent (trash first then delete)
	trashQuery2 := fmt.Sprintf(`mutation {
  trashedCategory(input: {id: %d}) {
    status
    message
  }
}`, s.categoryID)
	_, _ = graphtest.ExecuteGraphQL(s.handler, trashQuery2, nil, token)

	deleteQuery := fmt.Sprintf(`mutation {
  deleteCategoryPermanent(input: {id: %d}) {
    status
    message
  }
}`, s.categoryID)
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	deleteRes, ok := resp.Data["deleteCategoryPermanent"].(map[string]interface{})
	s.True(ok, "expected deleteCategoryPermanent in data, got %+v", resp.Data)
	s.Equal("success", deleteRes["status"])

	// 10. RestoreAll
	restoreAllQuery := `mutation {
  restoreAllCategory {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, restoreAllQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	restoreAll, ok := resp.Data["restoreAllCategory"].(map[string]interface{})
	s.True(ok, "expected restoreAllCategory in data, got %+v", resp.Data)
	s.Equal("success", restoreAll["status"])

	// 11. DeleteAll
	deleteAllQuery := `mutation {
  deleteAllCategoryPermanent {
    status
    message
  }
}`
	resp, err = graphtest.ExecuteGraphQL(s.handler, deleteAllQuery, nil, token)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	deleteAll, ok := resp.Data["deleteAllCategoryPermanent"].(map[string]interface{})
	s.True(ok, "expected deleteAllCategoryPermanent in data, got %+v", resp.Data)
	s.Equal("success", deleteAll["status"])
}

func TestCategoryApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryApiTestSuite))
}
