package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type ProductQueryRepository interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*db.GetProductsByMerchantRow, *int, error)
	FindByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*db.GetProductsByCategoryNameRow, *int, error)
	FindById(ctx context.Context, product_id int) (*db.Product, error)
	FindByIdTrashed(ctx context.Context, id int) (*db.Product, error)
}

type ProductCommandRepository interface {
	CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*db.Product, error)
	UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*db.Product, error)
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.Product, error)
	AdjustProductStock(ctx context.Context, product_id int, delta int) (*db.Product, error)
	TrashedProduct(ctx context.Context, product_id int) (*db.Product, error)
	RestoreProduct(ctx context.Context, product_id int) (*db.Product, error)
	DeleteProductPermanent(ctx context.Context, product_id int) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductPermanent(ctx context.Context) (bool, error)
}
