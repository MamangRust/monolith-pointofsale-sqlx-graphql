package mencache

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type ProductQueryCache interface {
	GetCachedProducts(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsRow, *int, bool)
	SetCachedProducts(ctx context.Context, req *requests.FindAllProducts, data []*db.GetProductsRow, total *int)

	GetCachedProductsByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*db.GetProductsByMerchantRow, *int, bool)
	SetCachedProductsByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest, data []*db.GetProductsByMerchantRow, total *int)

	GetCachedProductsByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*db.GetProductsByCategoryNameRow, *int, bool)
	SetCachedProductsByCategory(ctx context.Context, req *requests.ProductByCategoryRequest, data []*db.GetProductsByCategoryNameRow, total *int)

	GetCachedProductActive(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsActiveRow, *int, bool)
	SetCachedProductActive(ctx context.Context, req *requests.FindAllProducts, data []*db.GetProductsActiveRow, total *int)

	GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsTrashedRow, *int, bool)
	SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProducts, data []*db.GetProductsTrashedRow, total *int)

	GetCachedProduct(ctx context.Context, productID int) (*db.Product, bool)
	SetCachedProduct(ctx context.Context, data *db.Product)
}

type ProductCommandCache interface {
	DeleteCachedProduct(ctx context.Context, productID int)
	DeleteCachedProductAllCache(ctx context.Context)
}
