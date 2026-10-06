package product_cache

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

type ProductQueryCache interface {
	GetCachedProducts(ctx context.Context, req *model.FindAllProductInput) (*model.APIResponsePaginationProduct, bool)
	SetCachedProducts(ctx context.Context, req *model.FindAllProductInput, res *model.APIResponsePaginationProduct)

	GetCachedProductsByMerchant(ctx context.Context, req *model.FindAllProductMerchantInput) (*model.APIResponsePaginationProduct, bool)
	SetCachedProductsByMerchant(ctx context.Context, req *model.FindAllProductMerchantInput, res *model.APIResponsePaginationProduct)

	GetCachedProductsByCategory(ctx context.Context, req *model.FindAllProductCategoryInput) (*model.APIResponsePaginationProduct, bool)
	SetCachedProductsByCategory(ctx context.Context, req *model.FindAllProductCategoryInput, res *model.APIResponsePaginationProduct)

	GetCachedProductActive(ctx context.Context, req *model.FindAllProductInput) (*model.APIResponsePaginationProductDeleteAt, bool)
	SetCachedProductActive(ctx context.Context, req *model.FindAllProductInput, res *model.APIResponsePaginationProductDeleteAt)

	GetCachedProductTrashed(ctx context.Context, req *model.FindAllProductInput) (*model.APIResponsePaginationProductDeleteAt, bool)
	SetCachedProductTrashed(ctx context.Context, req *model.FindAllProductInput, res *model.APIResponsePaginationProductDeleteAt)

	GetCachedProduct(ctx context.Context, productID int) (*model.APIResponseProduct, bool)
	SetCachedProduct(ctx context.Context, res *model.APIResponseProduct)
}

type ProductCommandCache interface {
	DeleteCachedProduct(ctx context.Context, productID int)
}
