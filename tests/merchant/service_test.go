package merchant_test

import (
	"context"
	"testing"

	merchant_cache "github.com/MamangRust/monolith-pointofsale-merchant/cache"
	"github.com/MamangRust/monolith-pointofsale-merchant/repository"
	"github.com/MamangRust/monolith-pointofsale-merchant/service"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"

	tests "github.com/MamangRust/monolith-pointofsale-test"

	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	"github.com/stretchr/testify/suite"
)

type MerchantServiceTestSuite struct {
	tests.BaseTestSuite
	merchantService *service.Service
	merchantID      int
	userID          int
}

func (s *MerchantServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	queries := s.SQLxDB()

	// User service nyata (gRPC) — CreateMerchant memvalidasi user via gRPC.
	s.SetupUserService()

	// Seed a user directly
	err := s.DBPool().QueryRow(s.Ctx,
		`INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ($1, $2, $3, $4, 'test-verify', true) RETURNING user_id`,
		"Merchant", "ServiceTest", "merchant.svc@example.com", "password123",
	).Scan(&s.userID)
	s.Require().NoError(err)

	mencache := merchant_cache.NewMencache(s.GetCacheStore())
	userClient := pbuser.NewUserQueryServiceClient(s.Conns["user"])
	repos := repository.NewRepositories(queries, userClient)

	s.merchantService = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        s.Log,
		Mencache:      mencache,
		Kafka:         nil,
		Observability: s.Obs,
	})
}

func (s *MerchantServiceTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *MerchantServiceTestSuite) TestMerchantLifecycle() {
	ctx := context.Background()

	// 1. Create
	req := &requests.CreateMerchantRequest{
		UserID:       s.userID,
		Name:         "Service Merchant",
		Description:  "Merchant created via service layer",
		Address:      "Service Street No. 1",
		ContactEmail: "service.merchant@example.com",
		ContactPhone: "08123456781",
		Status:       "active",
	}
	created, err := s.merchantService.MerchantCommand.CreateMerchant(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	s.Equal(req.Name, created.Name)
	merchantID := int(created.MerchantID)

	// 2. FindByID
	found, err := s.merchantService.MerchantQuery.FindById(ctx, merchantID)
	s.Require().NoError(err)
	s.Equal(merchantID, int(found.MerchantID))

	// 3. Update
	updateReq := &requests.UpdateMerchantRequest{
		MerchantID:   &merchantID,
		UserID:       s.userID,
		Name:         "Updated Service Merchant",
		Description:  "Updated description",
		Address:      "Updated Street",
		ContactEmail: "updated.merchant@example.com",
		ContactPhone: "08987654321",
		Status:       "active",
	}
	updated, err := s.merchantService.MerchantCommand.UpdateMerchant(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(updateReq.Name, updated.Name)

	// 4. FindAll
	_, total, err := s.merchantService.MerchantQuery.FindAll(ctx, &requests.FindAllMerchants{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.merchantService.MerchantCommand.TrashedMerchant(ctx, merchantID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.merchantService.MerchantQuery.FindByTrashed(ctx, &requests.FindAllMerchants{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.merchantService.MerchantQuery.FindByActive(ctx, &requests.FindAllMerchants{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, m := range active {
		s.NotEqual(merchantID, int(m.MerchantID))
	}

	// 8. Restore
	_, err = s.merchantService.MerchantCommand.RestoreMerchant(ctx, merchantID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.merchantService.MerchantCommand.TrashedMerchant(ctx, merchantID)
	s.Require().NoError(err)
	success, err := s.merchantService.MerchantCommand.DeleteMerchantPermanent(ctx, merchantID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	m1, _ := s.merchantService.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: s.userID, Name: "M1", Description: "D1", Address: "A1",
		ContactEmail: "m1@example.com", ContactPhone: "111", Status: "active",
	})
	m2, _ := s.merchantService.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: s.userID, Name: "M2", Description: "D2", Address: "A2",
		ContactEmail: "m2@example.com", ContactPhone: "222", Status: "active",
	})

	s.merchantService.MerchantCommand.TrashedMerchant(ctx, int(m1.MerchantID))
	s.merchantService.MerchantCommand.TrashedMerchant(ctx, int(m2.MerchantID))

	resRestoreAll, err := s.merchantService.MerchantCommand.RestoreAllMerchant(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.merchantService.MerchantCommand.TrashedMerchant(ctx, int(m1.MerchantID))
	s.merchantService.MerchantCommand.TrashedMerchant(ctx, int(m2.MerchantID))

	resDeleteAll, err := s.merchantService.MerchantCommand.DeleteAllMerchantPermanent(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestMerchantServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantServiceTestSuite))
}
