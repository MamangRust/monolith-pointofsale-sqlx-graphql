package mencache

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type MerchantQueryCache interface {
	GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, *int, bool)
	SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants, data []*db.GetMerchantsRow, total *int)

	GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsActiveRow, *int, bool)
	SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants, data []*db.GetMerchantsActiveRow, total *int)

	GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsTrashedRow, *int, bool)
	SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants, data []*db.GetMerchantsTrashedRow, total *int)

	GetCachedMerchant(ctx context.Context, id int) (*db.Merchant, bool)
	SetCachedMerchant(ctx context.Context, data *db.Merchant)

	GetCachedMerchantsByUserId(ctx context.Context, id int) ([]*db.Merchant, bool)
	SetCachedMerchantsByUserId(ctx context.Context, userId int, data []*db.Merchant)
}

type MerchantCommandCache interface {
	DeleteCachedMerchant(ctx context.Context, id int)
	DeleteCachedMerchantAllCache(ctx context.Context)
}

type MerchantDocumentQueryCache interface {
	GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, bool)
	SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*db.GetMerchantDocumentsRow, total *int)

	GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, bool)
	SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*db.GetActiveMerchantDocumentsRow, total *int)

	GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, bool)
	SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*db.GetTrashedMerchantDocumentsRow, total *int)

	GetCachedMerchantDocument(ctx context.Context, id int) (*db.MerchantDocument, bool)
	SetCachedMerchantDocument(ctx context.Context, data *db.MerchantDocument)
}

type MerchantDocumentCommandCache interface {
	DeleteCachedMerchantDocuments(ctx context.Context, id int)
	DeleteCachedMerchantDocumentsAllCache(ctx context.Context)
}
