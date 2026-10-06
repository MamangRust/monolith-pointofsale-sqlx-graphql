package merchant_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-merchant/repository"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"

	tests "github.com/MamangRust/monolith-pointofsale-test"

	"net"

	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MerchantRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	repo       *repository.Repositories
	merchantID int
	userID     int
}

func (s *MerchantRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)

	// Create a placeholder gRPC listener to satisfy NewRepositories
	// (UserQuery gRPC methods are not called during DB-only tests)
	userLis, _ := net.Listen("tcp", "localhost:0")
	userConn, _ := grpc.NewClient(userLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	queries := pool
	s.repo = repository.NewRepositories(queries, pbuser.NewUserQueryServiceClient(userConn))

	// Seed a user ID directly for merchant tests
	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ($1, $2, $3, $4, 'test-verify', true) RETURNING user_id`,
		"Merchant", "RepoTest", "merchant.repo@example.com", "password123",
	).Scan(&s.userID)
	s.Require().NoError(err)
}

func (s *MerchantRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *MerchantRepositoryTestSuite) Test1_CreateMerchant() {
	ctx := context.Background()

	req := &requests.CreateMerchantRequest{
		UserID:       s.userID,
		Name:         "Test Merchant",
		Description:  "Detailed description",
		Address:      "Merchant Street No. 1",
		ContactEmail: "merchant@example.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}

	res, err := s.repo.MerchantCommand.CreateMerchant(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Name, res.Name)
	s.Equal(int32(s.userID), res.UserID)
	s.merchantID = int(res.MerchantID)
}

func (s *MerchantRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	found, err := s.repo.MerchantQuery.FindById(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.merchantID, int(found.MerchantID))
}

func (s *MerchantRepositoryTestSuite) Test3_UpdateMerchant() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	req := &requests.UpdateMerchantRequest{
		MerchantID:   &s.merchantID,
		UserID:       s.userID,
		Name:         "Updated Merchant",
		Description:  "Updated description",
		Address:      "Updated Street",
		ContactEmail: "updated@example.com",
		ContactPhone: "08987654321",
		Status:       "active",
	}

	res, err := s.repo.MerchantCommand.UpdateMerchant(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("Updated Merchant", res.Name)
}

func (s *MerchantRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	// Trash
	trashed, err := s.repo.MerchantCommand.TrashedMerchant(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(trashed)

	// Restore
	restored, err := s.repo.MerchantCommand.RestoreMerchant(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(restored)

	// Verify restored
	found, err := s.repo.MerchantQuery.FindById(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(found)
}

func (s *MerchantRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	// Must be trashed first for permanent delete
	_, err := s.repo.MerchantCommand.TrashedMerchant(ctx, s.merchantID)
	s.NoError(err)

	success, err := s.repo.MerchantCommand.DeleteMerchantPermanent(ctx, s.merchantID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantRepositoryTestSuite))
}
