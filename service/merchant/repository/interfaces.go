package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/jmoiron/sqlx"
)

type MerchantDocumentQueryRepository interface {
	FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, error)
	FindById(ctx context.Context, id int) (*db.MerchantDocument, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, error)
}

type MerchantDocumentCommandRepository interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.MerchantDocument, error)
	CreateMerchantDocumentInTx(ctx context.Context, tx *sqlx.Tx, request *requests.CreateMerchantDocumentRequest) (*db.MerchantDocument, error)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.MerchantDocument, error)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.MerchantDocument, error)
	UpdateMerchantDocumentStatusInTx(ctx context.Context, tx *sqlx.Tx, request *requests.UpdateMerchantDocumentStatusRequest) (*db.MerchantDocument, error)
	TrashedMerchantDocument(ctx context.Context, merchantDocumentID int) (*db.MerchantDocument, error)
	RestoreMerchantDocument(ctx context.Context, merchantDocumentID int) (*db.MerchantDocument, error)
	DeleteMerchantDocumentPermanent(ctx context.Context, merchantDocumentID int) (bool, error)
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsTrashedRow, *int, error)
	FindById(ctx context.Context, userID int) (*db.Merchant, error)
}

type MerchantCommandRepository interface {
	CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*db.Merchant, error)
	CreateMerchantInTx(ctx context.Context, tx *sqlx.Tx, request *requests.CreateMerchantRequest) (*db.Merchant, error)
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.Merchant, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.Merchant, error)
	UpdateMerchantStatusInTx(ctx context.Context, tx *sqlx.Tx, request *requests.UpdateMerchantStatusRequest) (*db.Merchant, error)
	TrashedMerchant(ctx context.Context, merchantID int) (*db.Merchant, error)
	RestoreMerchant(ctx context.Context, merchantID int) (*db.Merchant, error)
	DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, error)
	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}
