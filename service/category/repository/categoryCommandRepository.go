package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/category_errors"
	"github.com/jmoiron/sqlx"
)

const createCategory = `INSERT INTO categories (name, description, slug_category)
VALUES ($1, $2, $3)
  RETURNING *`

const updateCategory = `UPDATE categories
SET name = $2,
    description = $3,
    slug_category = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE category_id = $1
  AND deleted_at IS NULL
  RETURNING *`

const trashCategory = `UPDATE categories
SET
    deleted_at = current_timestamp
WHERE
    category_id = $1
    AND deleted_at IS NULL
    RETURNING *`

const restoreCategory = `UPDATE categories
SET
    deleted_at = NULL
WHERE
    category_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *`

const deleteCategoryPermanently = `DELETE FROM categories WHERE category_id = $1 AND deleted_at IS NOT NULL`

const restoreAllCategories = `UPDATE categories
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL`

const deleteAllPermanentCategories = `DELETE FROM categories
WHERE
    deleted_at IS NOT NULL`

type categoryCommandRepository struct {
	db *sqlx.DB
}

func NewCategoryCommandRepository(db *sqlx.DB) CategoryCommandRepository {
	return &categoryCommandRepository{
		db: db,
	}
}

func (r *categoryCommandRepository) CreateCategory(ctx context.Context, request *requests.CreateCategoryRequest) (*db.Category, error) {
	var category db.Category
	err := r.db.GetContext(ctx, &category, createCategory,
		request.Name,
		&request.Description,
		request.SlugCategory,
	)
	if err != nil {
		return nil, category_errors.ErrCreateCategory
	}

	return &category, nil
}

func (r *categoryCommandRepository) UpdateCategory(ctx context.Context, request *requests.UpdateCategoryRequest) (*db.Category, error) {
	var category db.Category
	err := r.db.GetContext(ctx, &category, updateCategory,
		int32(*request.CategoryID),
		request.Name,
		&request.Description,
		request.SlugCategory,
	)
	if err != nil {
		return nil, category_errors.ErrUpdateCategory
	}

	return &category, nil
}

func (r *categoryCommandRepository) TrashedCategory(ctx context.Context, category_id int) (*db.Category, error) {
	var category db.Category
	err := r.db.GetContext(ctx, &category, trashCategory, int32(category_id))
	if err != nil {
		return nil, category_errors.ErrTrashedCategory
	}

	return &category, nil
}

func (r *categoryCommandRepository) RestoreCategory(ctx context.Context, category_id int) (*db.Category, error) {
	var category db.Category
	err := r.db.GetContext(ctx, &category, restoreCategory, int32(category_id))
	if err != nil {
		return nil, category_errors.ErrRestoreCategory
	}

	return &category, nil
}

func (r *categoryCommandRepository) DeleteCategoryPermanently(ctx context.Context, category_id int) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteCategoryPermanently, int32(category_id)); err != nil {
		return false, category_errors.ErrDeleteCategoryPermanently
	}

	return true, nil
}

func (r *categoryCommandRepository) RestoreAllCategories(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, restoreAllCategories); err != nil {
		return false, category_errors.ErrRestoreAllCategories
	}

	return true, nil
}

func (r *categoryCommandRepository) DeleteAllPermanentCategories(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteAllPermanentCategories); err != nil {
		return false, category_errors.ErrDeleteAllPermanentCategories
	}

	return true, nil
}
