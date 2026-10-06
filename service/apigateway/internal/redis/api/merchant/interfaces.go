package merchant_cache

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

type MerchantQueryCache interface {
	GetCachedMerchants(ctx context.Context, req *model.FindAllMerchantInput) (*model.APIResponsePaginationMerchant, bool)
	SetCachedMerchants(ctx context.Context, req *model.FindAllMerchantInput, res *model.APIResponsePaginationMerchant)

	GetCachedMerchantActive(ctx context.Context, req *model.FindAllMerchantInput) (*model.APIResponsePaginationMerchantDeleteAt, bool)
	SetCachedMerchantActive(ctx context.Context, req *model.FindAllMerchantInput, res *model.APIResponsePaginationMerchantDeleteAt)

	GetCachedMerchantTrashed(ctx context.Context, req *model.FindAllMerchantInput) (*model.APIResponsePaginationMerchantDeleteAt, bool)
	SetCachedMerchantTrashed(ctx context.Context, req *model.FindAllMerchantInput, res *model.APIResponsePaginationMerchantDeleteAt)

	GetCachedMerchant(ctx context.Context, id int) (*model.APIResponseMerchant, bool)
	SetCachedMerchant(ctx context.Context, res *model.APIResponseMerchant)

	GetCachedMerchantsByUserId(ctx context.Context, id int) (*model.APIResponsesMerchant, bool)
	SetCachedMerchantsByUserId(ctx context.Context, userId int, res *model.APIResponsesMerchant)
}

type MerchantCommandCache interface {
	DeleteCachedMerchant(ctx context.Context, id int)
}
