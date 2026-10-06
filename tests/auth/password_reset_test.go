package auth_test

import (
	"context"
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	"net"
	"testing"

	auth_cache "github.com/MamangRust/monolith-pointofsale-auth/cache"
	"github.com/MamangRust/monolith-pointofsale-auth/handler"
	"github.com/MamangRust/monolith-pointofsale-auth/repository"
	"github.com/MamangRust/monolith-pointofsale-auth/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/auth"
	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// AuthPasswordResetTestSuite covers the three password-flow gRPC resolvers
// (VerifyCode, ForgotPassword, ResetPassword) end-to-end against a real
// PostgreSQL + Redis stack, and verifies the ResetPassword hashing fix:
// the persisted password must be a bcrypt hash (never the plaintext) and
// login must succeed with the new password afterwards.
type AuthPasswordResetTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *sqlx.DB
	redisClient *redis.Client
	client      pb.AuthServiceClient
	conn        *grpc.ClientConn
	grpcServer  *grpc.Server
	hasher      hash.HashPassword
	roleConn    *grpc.ClientConn
	userConn    *grpc.ClientConn
	depsServers []*grpc.Server
}

func (s *AuthPasswordResetTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	queries := pool

	log, _ := logger.NewLogger("test", nil)
	s.hasher = hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	obs, _ := observability.NewObservability("test", log)

	// Auth reaches user and role over gRPC, so both services must be running.
	roleConn, userConn, servers, err := tests.StartRoleUserServices(queries, log, cacheStore, obs)
	s.Require().NoError(err)
	s.roleConn = roleConn
	s.userConn = userConn
	s.depsServers = servers

	repos := repository.NewRepositories(&repository.Deps{
		Db:                queries,
		UserQueryClient:   pbuser.NewUserQueryServiceClient(userConn),
		UserCommandClient: pbuser.NewUserCommandServiceClient(userConn),
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
		Hash:          s.hasher,
		Kafka:         nil,
		Observability: obs,
	})

	h := handler.NewAuthHandleGrpc(svc, log)
	s.grpcServer = grpc.NewServer()
	pb.RegisterAuthServiceServer(s.grpcServer, h)

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)
	go func() {
		_ = s.grpcServer.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewAuthServiceClient(conn)
}

func (s *AuthPasswordResetTestSuite) TearDownSuite() {
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
	if s.roleConn != nil {
		s.roleConn.Close()
	}
	if s.userConn != nil {
		s.userConn.Close()
	}
	for _, srv := range s.depsServers {
		srv.Stop()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

// registerUser creates a fresh user through the register resolver and returns
// its user_id.
func (s *AuthPasswordResetTestSuite) registerUser(email, password string) int32 {
	ctx := context.Background()
	res, err := s.client.RegisterUser(ctx, &pb.RegisterRequest{
		Firstname:       "Reset",
		Lastname:        "Tester",
		Email:           email,
		Password:        password,
		ConfirmPassword: password,
	})
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotZero(res.Data.Id)
	return res.Data.Id
}

func (s *AuthPasswordResetTestSuite) markUserVerified(userID int32) {
	_, err := s.dbPool.ExecContext(context.Background(),
		"UPDATE users SET is_verified = true WHERE user_id = $1", userID)
	s.Require().NoError(err)
}

func (s *AuthPasswordResetTestSuite) fetchVerificationCode(userID int32) string {
	var code string
	err := s.dbPool.QueryRowContext(context.Background(),
		"SELECT verification_code FROM users WHERE user_id = $1", userID).Scan(&code)
	s.Require().NoError(err)
	return code
}

func (s *AuthPasswordResetTestSuite) fetchResetToken(userID int32) string {
	var token string
	err := s.dbPool.QueryRowContext(context.Background(),
		"SELECT token FROM reset_tokens WHERE user_id = $1", userID).Scan(&token)
	s.Require().NoError(err)
	return token
}

func (s *AuthPasswordResetTestSuite) storedPassword(userID int32) string {
	var pw string
	err := s.dbPool.QueryRowContext(context.Background(),
		"SELECT password FROM users WHERE user_id = $1", userID).Scan(&pw)
	s.Require().NoError(err)
	return pw
}

// ---------------------------------------------------------------------------
// VerifyCode resolver
// ---------------------------------------------------------------------------

func (s *AuthPasswordResetTestSuite) Test1_VerifyCodeResolver_Success() {
	userID := s.registerUser("verifycode.success@example.com", "password123")
	code := s.fetchVerificationCode(userID)
	s.Require().NotEmpty(code)

	res, err := s.client.VerifyCode(context.Background(), &pb.VerifyCodeRequest{Code: code})
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Equal("success", res.Status)

	var isVerified bool
	err = s.dbPool.QueryRowContext(context.Background(),
		"SELECT is_verified FROM users WHERE user_id = $1", userID).Scan(&isVerified)
	s.Require().NoError(err)
	s.True(isVerified)
}

func (s *AuthPasswordResetTestSuite) Test2_VerifyCodeResolver_InvalidCode() {
	_, err := s.client.VerifyCode(context.Background(), &pb.VerifyCodeRequest{Code: "does-not-exist"})
	s.Require().Error(err)
	s.Equal(codes.NotFound, status.Code(err))
}

// ---------------------------------------------------------------------------
// ForgotPassword resolver
// ---------------------------------------------------------------------------

func (s *AuthPasswordResetTestSuite) Test3_ForgotPasswordResolver_Success() {
	userID := s.registerUser("forgot.success@example.com", "password123")

	res, err := s.client.ForgotPassword(context.Background(), &pb.ForgotPasswordRequest{Email: "forgot.success@example.com"})
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Equal("success", res.Status)

	// A reset token must be persisted for the user.
	s.Require().NotEmpty(s.fetchResetToken(userID))
}

func (s *AuthPasswordResetTestSuite) Test4_ForgotPasswordResolver_UnknownEmail_DoesNotRevealAccount() {
	res, err := s.client.ForgotPassword(context.Background(), &pb.ForgotPasswordRequest{Email: "ghost@example.com"})
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Equal("success", res.Status)
}

// ---------------------------------------------------------------------------
// ResetPassword resolver + hashing fix
// ---------------------------------------------------------------------------

func (s *AuthPasswordResetTestSuite) Test5_ResetPasswordResolver_Success_AndHashingFix() {
	const newPassword = "freshpassword456"

	userID := s.registerUser("reset.success@example.com", "password123")

	_, err := s.client.ForgotPassword(context.Background(), &pb.ForgotPasswordRequest{Email: "reset.success@example.com"})
	s.Require().NoError(err)

	token := s.fetchResetToken(userID)
	s.Require().NotEmpty(token)

	res, err := s.client.ResetPassword(context.Background(), &pb.ResetPasswordRequest{
		ResetToken:      token,
		Password:        newPassword,
		ConfirmPassword: newPassword,
	})
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Equal("success", res.Status)

	// --- Hashing fix: the DB must hold a bcrypt hash, never the plaintext ---
	stored := s.storedPassword(userID)
	s.NotEqual(newPassword, stored, "plaintext password must never be persisted")
	s.NoError(s.hasher.ComparePassword(stored, newPassword),
		"stored password should verify against the new password via bcrypt")

	// The reset token must be consumed after a successful reset.
	var count int
	err = s.dbPool.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM reset_tokens WHERE user_id = $1", userID).Scan(&count)
	s.Require().NoError(err)
	s.Equal(0, count)

	// --- Login with the new password must succeed after the reset ---
	s.markUserVerified(userID)
	loginRes, err := s.client.LoginUser(context.Background(), &pb.LoginRequest{
		Email:    "reset.success@example.com",
		Password: newPassword,
	})
	s.Require().NoError(err)
	s.Require().NotNil(loginRes)
	s.NotEmpty(loginRes.Data.AccessToken)
}

func (s *AuthPasswordResetTestSuite) Test6_ResetPasswordResolver_InvalidToken() {
	_, err := s.client.ResetPassword(context.Background(), &pb.ResetPasswordRequest{
		ResetToken:      "bogus-token",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	})
	s.Require().Error(err)
	s.Equal(codes.NotFound, status.Code(err))
}

func (s *AuthPasswordResetTestSuite) Test7_ResetPasswordResolver_PasswordMismatch() {
	userID := s.registerUser("reset.mismatch@example.com", "password123")

	_, err := s.client.ForgotPassword(context.Background(), &pb.ForgotPasswordRequest{Email: "reset.mismatch@example.com"})
	s.Require().NoError(err)

	token := s.fetchResetToken(userID)
	s.Require().NotEmpty(token)

	_, err = s.client.ResetPassword(context.Background(), &pb.ResetPasswordRequest{
		ResetToken:      token,
		Password:        "newpassword123",
		ConfirmPassword: "differentpassword",
	})
	s.Require().Error(err)
}

func TestAuthPasswordResetSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthPasswordResetTestSuite))
}
