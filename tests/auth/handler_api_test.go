package auth_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	auth_cache "github.com/MamangRust/monolith-pointofsale-auth/cache"
	"github.com/MamangRust/monolith-pointofsale-auth/handler"
	"github.com/MamangRust/monolith-pointofsale-auth/repository"
	"github.com/MamangRust/monolith-pointofsale-auth/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	role_cache "github.com/MamangRust/monolith-pointofsale-role/cache"
	role_handler "github.com/MamangRust/monolith-pointofsale-role/handler"
	role_repo "github.com/MamangRust/monolith-pointofsale-role/repository"
	role_service "github.com/MamangRust/monolith-pointofsale-role/service"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	user_cache "github.com/MamangRust/monolith-pointofsale-user/cache"
	user_handler "github.com/MamangRust/monolith-pointofsale-user/handler"
	user_repo "github.com/MamangRust/monolith-pointofsale-user/repository"
	user_service "github.com/MamangRust/monolith-pointofsale-user/service"
	"github.com/jmoiron/sqlx"

	graphtest "github.com/MamangRust/monolith-pointofsale-apigateway/graphtest"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/MamangRust/monolith-pointofsale-pb"
	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	userpb "github.com/MamangRust/monolith-pointofsale-pb/user"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
)

type AuthHandlerApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	handler     http.Handler
	email       string
	password    string
	accessToken string
	userID      int
}

func (s *AuthHandlerApiTestSuite) SetupSuite() {
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

	log, _ := logger.NewLogger("test", nil)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	obs, _ := observability.NewObservability("test", log)

	// 1. Setup Role Service & gRPC Server
	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(queries)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repositories:  roleRepos,
		Logger:        log,
		Mencache:      roleMencache,
		Observability: obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  log,
	})
	roleServer := grpc.NewServer()
	pbrole.RegisterRoleQueryServiceServer(roleServer, roleGapi.Role)
	pbrole.RegisterRoleCommandServiceServer(roleServer, roleGapi.RoleCommand)
	pbuserrole.RegisterUserRoleServiceServer(roleServer, roleGapi.UserRole)
	roleLis, _ := net.Listen("tcp", "localhost:0")
	go roleServer.Serve(roleLis)
	roleConn, _ := grpc.NewClient(roleLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 2. Setup User Service & gRPC Server
	userMencache := user_cache.NewMencache(cacheStore)
	userRepos := user_repo.NewRepositories(&user_repo.Deps{
		Db:              queries,
		RoleQueryClient: pbrole.NewRoleQueryServiceClient(roleConn),
		UserRoleClient:  pbuserrole.NewUserRoleServiceClient(roleConn),
	})
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        log,
		Hash:          hasher,
		Mencache:      userMencache,
		Observability: obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  log,
	})
	userServer := grpc.NewServer()
	userpb.RegisterUserQueryServiceServer(userServer, userGapi.User)
	userpb.RegisterUserCommandServiceServer(userServer, userGapi.UserCommand)
	userLis, _ := net.Listen("tcp", "localhost:0")
	go userServer.Serve(userLis)
	userConn, _ := grpc.NewClient(userLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 3. Setup Auth Service & gRPC Server
	repos := repository.NewRepositories(&repository.Deps{
		Db:                queries,
		UserQueryClient:   userpb.NewUserQueryServiceClient(userConn),
		UserCommandClient: userpb.NewUserCommandServiceClient(userConn),
		RoleQueryClient:   pbrole.NewRoleQueryServiceClient(roleConn),
		RoleCommandClient: pbrole.NewRoleCommandServiceClient(roleConn),
		UserRoleClient:    pbuserrole.NewUserRoleServiceClient(roleConn),
	})

	tokenManager, _ := auth.NewManager("mysecret")
	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      auth_cache.NewMencache(cacheStore),
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: obs,
	})

	h := handler.NewAuthHandleGrpc(svc, log)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, h)

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	// Build the GraphQL apigateway handler wired to the in-process gRPC servers.
	conns := &graphtest.ServiceConnections{
		AuthClient: conn,
		RoleClient: roleConn,
		UserClient: userConn,
	}
	resolver := graphtest.NewResolver(conns, log, s.redisClient)
	s.handler = graphtest.NewAuthHandler(resolver, tokenManager, log)

	s.email = "auth.handler.api.test@example.com"
	s.password = "password123"

	// ROLE_ADMIN sudah di-seed oleh tests.SetupTestSuite (SeedMinimalRoles).
	// CreateRole di bawah hanya memastikan wiring role service aktif;
	// konflik duplikat (sudah ada) diabaikan.
	roleClient := tests.NewRoleClient(roleConn)
	if _, err := roleClient.CreateRole(context.Background(), &pbrole.CreateRoleRequest{
		Name: "ROLE_ADMIN",
	}); err != nil {
		fmt.Printf("DEBUG: Seed ROLE_ADMIN gRPC (expected conflict): %v\n", err)
	}
}

func (s *AuthHandlerApiTestSuite) TearDownSuite() {
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	s.ts.Teardown()
}

// verifyUser simulates the email-verification step (login only accepts
// is_verified = true users).
func (s *AuthHandlerApiTestSuite) verifyUser(email string) {
	_, err := s.dbPool.Exec(context.Background(),
		"UPDATE users SET is_verified = true WHERE email = $1", email)
	s.Require().NoError(err)
}

func (s *AuthHandlerApiTestSuite) Test1_Register() {
	query := `mutation {
  registerUser(input: {
    firstname: "Auth"
    lastname: "API"
    email: "` + s.email + `"
    password: "` + s.password + `"
    confirm_password: "` + s.password + `"
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

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, "")
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	register, ok := resp.Data["registerUser"].(map[string]interface{})
	s.True(ok, "expected registerUser in data, got %+v", resp.Data)
	s.Equal("success", register["status"])

	data, ok := register["data"].(map[string]interface{})
	s.True(ok, "expected registerUser.data, got %+v", register)
	s.userID = int(data["id"].(float64))
}

func (s *AuthHandlerApiTestSuite) Test2_Login() {
	s.verifyUser(s.email)

	query := `mutation {
  loginUser(input: {
    email: "` + s.email + `"
    password: "` + s.password + `"
  }) {
    status
    message
    data {
      access_token
      refresh_token
    }
  }
}`

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, "")
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	login, ok := resp.Data["loginUser"].(map[string]interface{})
	s.True(ok, "expected loginUser in data, got %+v", resp.Data)
	s.Equal("success", login["status"])

	data, ok := login["data"].(map[string]interface{})
	s.True(ok, "expected loginUser.data, got %+v", login)
	s.accessToken = data["access_token"].(string)
	s.NotEmpty(s.accessToken)
}

func (s *AuthHandlerApiTestSuite) Test4_LoginLockout() {
	email := "locked.api@example.com"

	regQuery := `mutation {
  registerUser(input: {
    firstname: "Locked"
    lastname: "API"
    email: "` + email + `"
    password: "correctpassword"
    confirm_password: "correctpassword"
  }) {
    status
    message
  }
}`
	regResp, err := graphtest.ExecuteGraphQL(s.handler, regQuery, nil, "")
	s.NoError(err)
	s.Empty(regResp.Errors, "expected no GraphQL errors, got: %+v", regResp.Errors)
	s.verifyUser(email)

	loginQuery := `mutation {
  loginUser(input: {
    email: "` + email + `"
    password: "wrongpassword"
  }) {
    status
    message
  }
}`

	// Fail login 5 times (total 5) — password mismatch is an error response.
	for i := 0; i < 5; i++ {
		loginResp, err := graphtest.ExecuteGraphQL(s.handler, loginQuery, nil, "")
		s.NoError(err)
		s.NotEmpty(loginResp.Errors, "attempt %d should have failed", i+1)
	}

	// 6th attempt should also fail (ErrAccountLocked).
	lockedResp, err := graphtest.ExecuteGraphQL(s.handler, loginQuery, nil, "")
	s.NoError(err)
	s.NotEmpty(lockedResp.Errors, "locked account should still reject login")
}

func (s *AuthHandlerApiTestSuite) Test3_GetMe() {
	s.Require().NotZero(s.userID)
	s.Require().NotEmpty(s.accessToken)

	query := `query {
  getMe(input: {access_token: "` + s.accessToken + `"}) {
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

	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, s.accessToken)
	s.NoError(err)
	s.Empty(resp.Errors, "expected no GraphQL errors, got: %+v", resp.Errors)

	getMe, ok := resp.Data["getMe"].(map[string]interface{})
	s.True(ok, "expected getMe in data, got %+v", resp.Data)
	s.Equal("success", getMe["status"])

	data, ok := getMe["data"].(map[string]interface{})
	s.True(ok, "expected getMe.data, got %+v", getMe)
	s.Equal(s.email, data["email"])
}

func TestAuthHandlerApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthHandlerApiTestSuite))
}
