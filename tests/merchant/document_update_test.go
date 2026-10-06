package merchant_test

import (
	"context"
	"net"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-merchant/repository"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"

	tests "github.com/MamangRust/monolith-pointofsale-test"

	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MerchantDocumentUpdateTestSuite is a regression suite for the merchant
// document update fix: UpdateMerchantDocument and UpdateMerchantDocumentStatus
// must target the request's DocumentID — never the MerchantID. The fixture
// creates two documents under one merchant and updates the *second* document,
// whose ID is guaranteed to differ from the merchant ID, so the old bug
// (using MerchantID as DocumentID) cannot pass by coincidence.
type MerchantDocumentUpdateTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	repo       *repository.Repositories
	userID     int
	merchantID int
	docA       *db.MerchantDocument
	docB       *db.MerchantDocument
}

func (s *MerchantDocumentUpdateTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)

	// Placeholder gRPC listener to satisfy NewRepositories
	// (UserQuery gRPC methods are not called during DB-only tests).
	userLis, _ := net.Listen("tcp", "localhost:0")
	userConn, _ := grpc.NewClient(userLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	queries := pool
	s.repo = repository.NewRepositories(queries, pbuser.NewUserQueryServiceClient(userConn))

	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ($1, $2, $3, $4, 'doc-verify', true) RETURNING user_id`,
		"Document", "Owner", "merchant.document@example.com", "password123",
	).Scan(&s.userID)
	s.Require().NoError(err)

	ctx := context.Background()

	merchant, err := s.repo.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:       s.userID,
		Name:         "Document Merchant",
		Description:  "Doc",
		Address:      "Addr",
		ContactEmail: "doc@m.com",
		ContactPhone: "123",
		Status:       "active",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	s.docA, err = s.repo.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID:   s.merchantID,
		DocumentType: "doc_a_type",
		DocumentUrl:  "https://example.com/a.pdf",
	})
	s.Require().NoError(err)

	s.docB, err = s.repo.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID:   s.merchantID,
		DocumentType: "doc_b_type",
		DocumentUrl:  "https://example.com/b.pdf",
	})
	s.Require().NoError(err)

	s.Require().NotEqual(s.docA.DocumentID, s.docB.DocumentID)
	// The regression hinges on this: the update target (docB) must be
	// distinguishable from the merchant ID.
	s.Require().NotEqual(int(s.docB.DocumentID), s.merchantID,
		"fixture requires the target document ID to differ from the merchant ID")
}

func (s *MerchantDocumentUpdateTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *MerchantDocumentUpdateTestSuite) Test1_UpdateMerchantDocument_UsesDocumentID() {
	ctx := context.Background()
	docBID := int(s.docB.DocumentID)

	updated, err := s.repo.MerchantDocumentCommand.UpdateMerchantDocument(ctx, &requests.UpdateMerchantDocumentRequest{
		DocumentID:   &docBID,
		MerchantID:   s.merchantID,
		DocumentType: "doc_b_updated",
		DocumentUrl:  "https://example.com/b-updated.pdf",
		Status:       "verified",
		Note:         "Approved",
	})
	s.Require().NoError(err)
	s.Require().NotNil(updated)

	// The row returned must be docB — keyed by the request DocumentID,
	// never by the merchant ID.
	s.Equal(s.docB.DocumentID, updated.DocumentID)
	s.Equal("doc_b_updated", updated.DocumentType)
	s.Equal("verified", updated.Status)

	// docA must be left untouched by the update.
	docA, err := s.repo.MerchantDocumentQuery.FindById(ctx, int(s.docA.DocumentID))
	s.Require().NoError(err)
	s.Require().NotNil(docA)
	s.Equal("doc_a_type", docA.DocumentType)
	s.NotEqual("doc_b_updated", docA.DocumentType)
}

func (s *MerchantDocumentUpdateTestSuite) Test2_UpdateMerchantDocumentStatus_UsesDocumentID() {
	ctx := context.Background()
	docBID := int(s.docB.DocumentID)

	updated, err := s.repo.MerchantDocumentCommand.UpdateMerchantDocumentStatus(ctx, &requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: &docBID,
		MerchantID: s.merchantID,
		Status:     "rejected",
		Note:       "Invalid document",
	})
	s.Require().NoError(err)
	s.Require().NotNil(updated)

	// The status change must land on docB, not on the merchant's own row.
	s.Equal(s.docB.DocumentID, updated.DocumentID)
	s.Equal("rejected", updated.Status)

	// docA (created with status "pending") must not have been rejected.
	docA, err := s.repo.MerchantDocumentQuery.FindById(ctx, int(s.docA.DocumentID))
	s.Require().NoError(err)
	s.Require().NotNil(docA)
	s.NotEqual("rejected", docA.Status)
}

func TestMerchantDocumentUpdateSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDocumentUpdateTestSuite))
}
