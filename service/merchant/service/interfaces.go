package service

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type MerchantQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsTrashedRow, *int, error)
	FindById(ctx context.Context, merchantID int) (*db.Merchant, error)
}

type MerchantCommandService interface {
	CreateMerchant(ctx context.Context, req *requests.CreateMerchantRequest) (*db.Merchant, error)
	UpdateMerchant(ctx context.Context, req *requests.UpdateMerchantRequest) (*db.Merchant, error)
	UpdateMerchantStatus(ctx context.Context, req *requests.UpdateMerchantStatusRequest) (*db.Merchant, error)
	TrashedMerchant(ctx context.Context, merchantID int) (*db.Merchant, error)
	RestoreMerchant(ctx context.Context, merchantID int) (*db.Merchant, error)
	DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, error)
	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type MerchantDocumentQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, error)
	FindById(ctx context.Context, documentID int) (*db.MerchantDocument, error)
}

type MerchantDocumentCommandService interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.MerchantDocument, error)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.MerchantDocument, error)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.MerchantDocument, error)
	TrashedMerchantDocument(ctx context.Context, documentID int) (*db.MerchantDocument, error)
	RestoreMerchantDocument(ctx context.Context, documentID int) (*db.MerchantDocument, error)
	DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error)
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}
