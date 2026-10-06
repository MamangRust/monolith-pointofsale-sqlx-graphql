package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-pointofsale-auth/repository"
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
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type AuthRepositoryTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	repo        *repository.Repositories
	userID      int
	email       string
	roleConn    *grpc.ClientConn
	userConn    *grpc.ClientConn
	grpcServers []*grpc.Server
}

func (s *AuthRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)

	log, _ := logger.NewLogger("test", nil)
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.ts.RedisClient(), log, cacheMetrics)

	queries := pool
	roleConn, userConn, servers, err := tests.StartRoleUserServices(queries, log, cacheStore, obs)
	s.Require().NoError(err)
	s.roleConn = roleConn
	s.userConn = userConn
	s.grpcServers = servers

	s.repo = repository.NewRepositories(&repository.Deps{
		Db:                queries,
		UserQueryClient:   pbuser.NewUserQueryServiceClient(userConn),
		UserCommandClient: pbuser.NewUserCommandServiceClient(userConn),
		RoleQueryClient:   pbrole.NewRoleQueryServiceClient(roleConn),
		RoleCommandClient: pbrole.NewRoleCommandServiceClient(roleConn),
		UserRoleClient:    pbuserrole.NewUserRoleServiceClient(roleConn),
	})
	s.email = "auth.repo.test@example.com"
}

func (s *AuthRepositoryTestSuite) TearDownSuite() {
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

func (s *AuthRepositoryTestSuite) Test1_CreateUser() {
	ctx := context.Background()

	req := &requests.RegisterRequest{
		FirstName:       "Auth",
		LastName:        "Repo",
		Email:           s.email,
		Password:        "password123",
		ConfirmPassword: "password123",
		VerifiedCode:    "123456",
		IsVerified:      false,
	}

	res, err := s.repo.User.CreateUser(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(s.email, res.Email)
	s.userID = int(res.UserID)
}

func (s *AuthRepositoryTestSuite) Test2_FindByEmail() {
	s.Require().NotEmpty(s.email)
	ctx := context.Background()

	found, err := s.repo.User.FindByEmail(ctx, s.email)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.userID), found.UserID)
}

func (s *AuthRepositoryTestSuite) Test3_FindById() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	found, err := s.repo.User.FindById(ctx, s.userID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.userID), found.UserID)
}

func (s *AuthRepositoryTestSuite) Test4_UpdateVerification() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	updated, err := s.repo.User.UpdateUserIsVerified(ctx, s.userID, true)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(s.userID), updated.UserID)
}

func (s *AuthRepositoryTestSuite) Test5_UpdatePassword() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	updated, err := s.repo.User.UpdateUserPassword(ctx, s.userID, "newpassword123")
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(s.userID), updated.UserID)
	// The password hash is intentionally not returned over gRPC, so the stored
	// value is verified by reading it back through FindByEmail (auth-only RPC).
	found, err := s.repo.User.FindByEmail(ctx, s.email)
	s.NoError(err)
	s.Require().NotNil(found)
	s.Equal("newpassword123", found.Password)
}

func (s *AuthRepositoryTestSuite) Test6_FindByVerificationCode() {
	ctx := context.Background()

	found, err := s.repo.User.FindByVerificationCode(ctx, "123456")
	s.NoError(err)
	s.NotNil(found)
}

func (s *AuthRepositoryTestSuite) Test7_RefreshToken() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	token := "test-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")

	req := &requests.CreateRefreshToken{
		UserId:    s.userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	res, err := s.repo.RefreshToken.CreateRefreshToken(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(token, res.Token)

	found, err := s.repo.RefreshToken.FindByToken(ctx, token)
	s.NoError(err)
	s.NotNil(found)

	foundByUser, err := s.repo.RefreshToken.FindByUserId(ctx, s.userID)
	s.NoError(err)
	s.NotNil(foundByUser)

	err = s.repo.RefreshToken.DeleteRefreshToken(ctx, token)
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) Test8_ResetToken() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	token := "reset-token-123"
	expiresAt := time.Now().Add(1 * time.Hour).Format("2006-01-02 15:04:05")

	req := &requests.CreateResetTokenRequest{
		UserID:     s.userID,
		ResetToken: token,
		ExpiredAt:  expiresAt,
	}

	res, err := s.repo.ResetToken.CreateResetToken(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(token, res.Token)

	found, err := s.repo.ResetToken.FindByToken(ctx, token)
	s.NoError(err)
	s.NotNil(found)

	err = s.repo.ResetToken.DeleteResetToken(ctx, s.userID)
	s.NoError(err)
}

func TestAuthRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthRepositoryTestSuite))
}
