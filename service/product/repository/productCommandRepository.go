package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/product_errors"
	"github.com/jmoiron/sqlx"
)

const createProduct = `INSERT INTO products (merchant_id, category_id, name, description, price, count_in_stock, brand, weight, slug_product, image_product, barcode)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *`

const updateProduct = `UPDATE products
SET category_id = $2,
    name = $3,
    description = $4,
    price = $5,
    count_in_stock = $6,
    brand = $7,
    weight = $8,
    image_product = $9,
    barcode = $10,
    updated_at = CURRENT_TIMESTAMP
WHERE product_id = $1
  AND deleted_at IS NULL
  RETURNING *`

const updateProductCountStock = `UPDATE products
SET count_in_stock = $2
WHERE product_id = $1
    AND deleted_at IS NULL
RETURNING *`

const adjustProductStock = `UPDATE products
SET count_in_stock = count_in_stock + $1,
    updated_at = CURRENT_TIMESTAMP
WHERE product_id = $2
    AND deleted_at IS NULL
    AND count_in_stock + $1 >= 0
RETURNING *`

const trashProduct = `UPDATE products
SET
    deleted_at = current_timestamp
WHERE
    product_id = $1
    AND deleted_at IS NULL
    RETURNING *`

const restoreProduct = `UPDATE products
SET
    deleted_at = NULL
WHERE
    product_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *`

const deleteProductPermanently = `DELETE FROM products WHERE product_id = $1 AND deleted_at IS NOT NULL`

const restoreAllProducts = `UPDATE products
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL`

const deleteAllPermanentProducts = `DELETE FROM products
WHERE
    deleted_at IS NOT NULL`

type productCommandRepository struct {
	db *sqlx.DB
}

func NewProductCommandRepository(db *sqlx.DB) *productCommandRepository {
	return &productCommandRepository{
		db: db,
	}
}

func (r *productCommandRepository) CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*db.Product, error) {
	weightVal := int32(request.Weight)

	var product db.Product
	err := r.db.GetContext(ctx, &product, createProduct,
		int32(request.MerchantID),
		int32(request.CategoryID),
		request.Name,
		&request.Description,
		int32(request.Price),
		int32(request.CountInStock),
		&request.Brand,
		&weightVal,
		request.SlugProduct,
		&request.ImageProduct,
		nil,
	)
	if err != nil {
		return nil, product_errors.ErrCreateProduct
	}

	return &product, nil
}

func (r *productCommandRepository) UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*db.Product, error) {
	weightVal := int32(request.Weight)

	var product db.Product
	err := r.db.GetContext(ctx, &product, updateProduct,
		int32(*request.ProductID),
		int32(request.CategoryID),
		request.Name,
		&request.Description,
		int32(request.Price),
		int32(request.CountInStock),
		&request.Brand,
		&weightVal,
		&request.ImageProduct,
		nil,
	)
	if err != nil {
		return nil, product_errors.ErrUpdateProduct
	}

	return &product, nil
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.Product, error) {
	var product db.Product
	err := r.db.GetContext(ctx, &product, updateProductCountStock,
		int32(product_id),
		int32(stock),
	)
	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return &product, nil
}

func (r *productCommandRepository) AdjustProductStock(ctx context.Context, product_id int, delta int) (*db.Product, error) {
	var product db.Product
	err := r.db.GetContext(ctx, &product, adjustProductStock,
		int32(delta),
		int32(product_id),
	)
	if err != nil {
		// The atomic guard `count_in_stock + delta >= 0` returned no rows:
		// the adjustment would drive stock negative → 400 Bad Request.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, product_errors.ErrInsufficientProductStock
		}
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return &product, nil
}

func (r *productCommandRepository) TrashedProduct(ctx context.Context, product_id int) (*db.Product, error) {
	var product db.Product
	err := r.db.GetContext(ctx, &product, trashProduct, int32(product_id))
	if err != nil {
		return nil, product_errors.ErrTrashedProduct
	}

	return &product, nil
}

func (r *productCommandRepository) RestoreProduct(ctx context.Context, product_id int) (*db.Product, error) {
	var product db.Product
	err := r.db.GetContext(ctx, &product, restoreProduct, int32(product_id))
	if err != nil {
		return nil, product_errors.ErrRestoreProduct
	}

	return &product, nil
}

func (r *productCommandRepository) DeleteProductPermanent(ctx context.Context, product_id int) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteProductPermanently, int32(product_id)); err != nil {
		return false, product_errors.ErrDeleteProductPermanent
	}

	return true, nil
}

func (r *productCommandRepository) RestoreAllProducts(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, restoreAllProducts); err != nil {
		return false, product_errors.ErrRestoreAllProducts
	}

	return true, nil
}

func (r *productCommandRepository) DeleteAllProductPermanent(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteAllPermanentProducts); err != nil {
		return false, product_errors.ErrDeleteAllProductPermanent
	}

	return true, nil
}
