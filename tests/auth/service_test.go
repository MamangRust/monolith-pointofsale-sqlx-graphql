package auth_test

import (
	"context"
	"testing"

	mencache "github.com/MamangRust/monolith-pointofsale-auth/cache"
	"github.com/MamangRust/monolith-pointofsale-auth/repository"
	"github.com/MamangRust/monolith-pointofsale-auth/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	tests "github.com/MamangRust/monolith-pointofsale-test"

	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type AuthServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	pool        *sqlx.DB
	authService *service.Service
	email       string
	password    string
	roleConn    *grpc.ClientConn
	userConn    *grpc.ClientConn
	grpcServers []*grpc.Server
}

func (s *AuthServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)
	s.pool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(opts)

	queries := pool

	log, _ := logger.NewLogger("test", nil)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)
	mencacheService := mencache.NewMencache(cacheStore)

	tokenManager, _ := auth.NewManager("mysecretkey")

	obs, _ := observability.NewObservability("test", log)

	// Auth reaches user and role over gRPC, so both services must be running.
	roleConn, userConn, servers, err := tests.StartRoleUserServices(queries, log, cacheStore, obs)
	s.Require().NoError(err)
	s.roleConn = roleConn
	s.userConn = userConn
	s.grpcServers = servers

	repos := repository.NewRepositories(&repository.Deps{
		Db:                queries,
		UserQueryClient:   pbuser.NewUserQueryServiceClient(userConn),
		UserCommandClient: pbuser.NewUserCommandServiceClient(userConn),
		RoleQueryClient:   pbrole.NewRoleQueryServiceClient(roleConn),
		RoleCommandClient: pbrole.NewRoleCommandServiceClient(roleConn),
		UserRoleClient:    pbuserrole.NewUserRoleServiceClient(roleConn),
	})

	s.authService = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencacheService,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: obs,
	})

	s.email = "auth.service.test@example.com"
	s.password = "password123"
}

func (s *AuthServiceTestSuite) TearDownSuite() {
	if s.roleConn != nil {
		s.roleConn.Close()
	}
	if s.userConn != nil {
		s.userConn.Close()
	}
	for _, srv := range s.grpcServers {
		srv.Stop()
	}
	s.ts.Teardown()
}

func (s *AuthServiceTestSuite) TestAuthLifecycle() {
	ctx := context.Background()

	// 1. Register
	regReq := &requests.RegisterRequest{
		FirstName:       "Auth",
		LastName:        "Service",
		Email:           s.email,
		Password:        s.password,
		ConfirmPassword: s.password,
	}

	created, err := s.authService.Register.Register(ctx, regReq)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	s.Equal(s.email, created.Email)

	// 1b. Verify email (login hanya menerima user is_verified = true)
	var verifyCode string
	err = s.pool.QueryRowContext(ctx,
		"SELECT verification_code FROM users WHERE email = $1", s.email).Scan(&verifyCode)
	s.Require().NoError(err)

	verified, err := s.authService.PasswordReset.VerifyCode(ctx, verifyCode)
	s.Require().NoError(err)
	s.True(verified)

	// 2. Login
	loginReq := &requests.AuthRequest{
		Email:    s.email,
		Password: s.password,
	}

	tokenRes, err := s.authService.Login.Login(ctx, loginReq)
	s.Require().NoError(err)
	s.Require().NotNil(tokenRes)
	s.NotEmpty(tokenRes.AccessToken)
	s.NotEmpty(tokenRes.RefreshToken)

	// 3. ForgotPassword
	success, err := s.authService.PasswordReset.ForgotPassword(ctx, s.email)
	s.Require().NoError(err)
	s.True(success)
}

func TestAuthServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthServiceTestSuite))
}
