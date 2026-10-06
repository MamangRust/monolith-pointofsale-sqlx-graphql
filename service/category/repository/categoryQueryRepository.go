package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/category_errors"
	"github.com/jmoiron/sqlx"
)

const getCategories = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM categories
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR slug_category ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getCategoriesActive = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM categories
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR slug_category ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getCategoriesTrashed = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM categories
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR slug_category ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getCategoryByID = `SELECT *
FROM categories
WHERE category_id = $1
  AND deleted_at IS NULL`

const getCategoryByIDTrashed = `SELECT *
FROM categories
WHERE category_id = $1
  AND deleted_at IS NOT NULL`

const getCategoryByName = `SELECT *
FROM categories
WHERE name = $1
  AND deleted_at IS NULL`

const getCategoryByNameAndId = `SELECT *
FROM categories
WHERE name = $1
  AND category_id = $2
  AND deleted_at IS NULL`

type categoryQueryRepository struct {
	db *sqlx.DB
}

func NewCategoryQueryRepository(db *sqlx.DB) CategoryQueryRepository {
	return &categoryQueryRepository{
		db: db,
	}
}

func (r *categoryQueryRepository) FindAllCategory(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var categories []*db.GetCategoriesRow
	err := r.db.SelectContext(ctx, &categories, getCategories,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, category_errors.ErrFindAllCategory
	}

	totalCount := 0
	if len(categories) > 0 {
		totalCount = int(categories[0].TotalCount)
	}

	return categories, &totalCount, nil
}

func (r *categoryQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var categories []*db.GetCategoriesActiveRow
	err := r.db.SelectContext(ctx, &categories, getCategoriesActive,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, category_errors.ErrFindByActive
	}

	totalCount := 0
	if len(categories) > 0 {
		totalCount = int(categories[0].TotalCount)
	}

	return categories, &totalCount, nil
}

func (r *categoryQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var categories []*db.GetCategoriesTrashedRow
	err := r.db.SelectContext(ctx, &categories, getCategoriesTrashed,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, category_errors.ErrFindByTrashed
	}

	totalCount := 0
	if len(categories) > 0 {
		totalCount = int(categories[0].TotalCount)
	}

	return categories, &totalCount, nil
}

func (r *categoryQueryRepository) FindById(ctx context.Context, category_id int) (*db.Category, error) {
	var category db.Category
	err := r.db.GetContext(ctx, &category, getCategoryByID, int32(category_id))
	if err != nil {
		return nil, category_errors.ErrFindById
	}

	return &category, nil
}

func (r *categoryQueryRepository) FindByIdTrashed(ctx context.Context, category_id int) (*db.Category, error) {
	var category db.Category
	err := r.db.GetContext(ctx, &category, getCategoryByIDTrashed, int32(category_id))
	if err != nil {
		return nil, category_errors.ErrFindById
	}

	return &category, nil
}

func (r *categoryQueryRepository) FindByName(ctx context.Context, name string) (*db.Category, error) {
	var category db.Category
	err := r.db.GetContext(ctx, &category, getCategoryByName, name)
	if err != nil {
		return nil, category_errors.ErrFindByName
	}

	return &category, nil
}

func (r *categoryQueryRepository) FindByNameAndId(ctx context.Context, req *requests.CategoryNameAndId) (*db.Category, error) {
	var category db.Category
	err := r.db.GetContext(ctx, &category, getCategoryByNameAndId,
		req.Name,
		int32(req.CategoryID),
	)
	if err != nil {
		return nil, category_errors.ErrFindByNameAndId
	}

	return &category, nil
}
