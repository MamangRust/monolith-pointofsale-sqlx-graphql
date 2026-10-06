package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/product_errors"
	"github.com/jmoiron/sqlx"
)

const getProducts = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM products as p
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL
       OR p.name ILIKE '%' || $1 || '%'
       OR p.description ILIKE '%' || $1 || '%'
       OR p.brand ILIKE '%' || $1 || '%'
       OR p.slug_product ILIKE '%' || $1 || '%'
       OR p.barcode ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getProductsActive = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM products as p
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL
       OR p.name ILIKE '%' || $1 || '%'
       OR p.description ILIKE '%' || $1 || '%'
       OR p.brand ILIKE '%' || $1 || '%'
       OR p.slug_product ILIKE '%' || $1 || '%'
       OR p.barcode ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getProductsTrashed = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM products as p
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL
       OR p.name ILIKE '%' || $1 || '%'
       OR p.description ILIKE '%' || $1 || '%'
       OR p.brand ILIKE '%' || $1 || '%'
       OR p.slug_product ILIKE '%' || $1 || '%'
       OR p.barcode ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getProductsByMerchant = `WITH filtered_products AS (
    SELECT
        p.product_id,
        p.name,
        p.description,
        p.price,
        p.count_in_stock,
        p.brand,
        p.image_product,
        p.created_at,
        c.name AS category_name
    FROM
        products p
    JOIN
        categories c ON p.category_id = c.category_id
    WHERE
        p.deleted_at IS NULL
        AND p.merchant_id = $1
        AND (
            p.name ILIKE '%' || COALESCE($2, '') || '%'
            OR p.description ILIKE '%' || COALESCE($2, '') || '%'
            OR $2 IS NULL
        )
        AND (
            c.category_id = NULLIF($3, 0)
            OR NULLIF($3, 0) IS NULL
        )
        AND (
            p.price >= COALESCE(NULLIF($4, 0), 0)
            AND p.price <= COALESCE(NULLIF($5, 0), 999999999)
        )
)
SELECT
    (SELECT COUNT(*) FROM filtered_products) AS total_count,
    fp.*
FROM
    filtered_products fp
ORDER BY
    fp.created_at DESC
LIMIT $6 OFFSET $7`

const getProductsByCategoryName = `WITH filtered_products AS (
    SELECT
        p.product_id,
        p.merchant_id,
        p.category_id,
        p.slug_product,
        p.weight,
        p.name,
        p.description,
        p.price,
        p.count_in_stock,
        p.brand,
        p.image_product,
        p.barcode,
        p.created_at,
        p.updated_at,
        p.deleted_at,
        c.name AS category_name
    FROM
        products p
    JOIN
        categories c ON p.category_id = c.category_id
    WHERE
        p.deleted_at IS NULL
        AND c.name = $1
        AND (
            $2::text IS NULL
            OR p.name ILIKE '%' || $2::text || '%'
            OR p.description ILIKE '%' || $2::text || '%'
        )
        AND (
            ($3::int IS NULL OR p.price >= $3::int)
            AND ($4::int IS NULL OR p.price <= $4::int)
        )
)
SELECT
    (SELECT COUNT(*) FROM filtered_products) AS total_count,
    fp.*
FROM
    filtered_products fp
ORDER BY
    fp.created_at DESC
LIMIT $5 OFFSET $6`

const getProductByID = `SELECT *
FROM products
WHERE product_id = $1
  AND deleted_at IS NULL`

const getProductByIdTrashed = `SELECT * FROM products WHERE product_id = $1`

type productQueryRepository struct {
	db *sqlx.DB
}

func NewProductQueryRepository(db *sqlx.DB) *productQueryRepository {
	return &productQueryRepository{
		db: db,
	}
}

func (r *productQueryRepository) FindAllProducts(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var products []*db.GetProductsRow
	err := r.db.SelectContext(ctx, &products, getProducts,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, product_errors.ErrFindAllProducts
	}

	var totalCount int
	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	return products, &totalCount, nil
}

func (r *productQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var products []*db.GetProductsActiveRow
	err := r.db.SelectContext(ctx, &products, getProductsActive,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, product_errors.ErrFindByActive
	}

	var totalCount int
	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	return products, &totalCount, nil
}

func (r *productQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var products []*db.GetProductsTrashedRow
	err := r.db.SelectContext(ctx, &products, getProductsTrashed,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, product_errors.ErrFindByTrashed
	}

	var totalCount int
	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	return products, &totalCount, nil
}

func (r *productQueryRepository) FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*db.GetProductsByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var products []*db.GetProductsByMerchantRow
	err := r.db.SelectContext(ctx, &products, getProductsByMerchant,
		int32(req.MerchantID),
		&req.Search,
		req.CategoryID,
		req.MinPrice,
		req.MaxPrice,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, product_errors.ErrFindByMerchant
	}

	var totalCount int
	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	return products, &totalCount, nil
}

func (r *productQueryRepository) FindByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*db.GetProductsByCategoryNameRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var products []*db.GetProductsByCategoryNameRow
	err := r.db.SelectContext(ctx, &products, getProductsByCategoryName,
		req.CategoryName,
		req.Search,
		req.MinPrice,
		req.MaxPrice,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, product_errors.ErrFindByCategory.WithInternal(err)
	}

	var totalCount int
	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	return products, &totalCount, nil
}

func (r *productQueryRepository) FindById(ctx context.Context, product_id int) (*db.Product, error) {
	var product db.Product
	err := r.db.GetContext(ctx, &product, getProductByID, int32(product_id))
	if err != nil {
		return nil, product_errors.ErrFindById
	}

	return &product, nil
}

func (r *productQueryRepository) FindByIdTrashed(ctx context.Context, product_id int) (*db.Product, error) {
	var product db.Product
	err := r.db.GetContext(ctx, &product, getProductByIdTrashed, int32(product_id))
	if err != nil {
		return nil, product_errors.ErrFindByIdTrashed
	}

	return &product, nil
}
