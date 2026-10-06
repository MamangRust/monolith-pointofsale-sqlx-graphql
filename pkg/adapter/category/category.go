// Package category adapts the Category service gRPC API into the sqlc row types
// used by the consumer repositories.
package category

import (
	"context"

	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter/convert"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/category_errors"
)

// Repository reads categories from the Category service over gRPC.
type Repository interface {
	FindById(ctx context.Context, id int) (*db.Category, error)
}

type repository struct {
	client pbcategory.CategoryQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a category adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough.
func New(client pbcategory.CategoryQueryServiceClient, opts ...adapter.GuardOption) Repository {
	r := &repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// SetGuard lets adapter.WithDependencyGuard attach a guard after construction.
func (r *repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

var _ Repository = (*repository)(nil)

func (r *repository) FindById(ctx context.Context, id int) (*db.Category, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbcategory.ApiResponseCategory, error) {
		return r.client.FindById(callCtx, &pbcategory.FindByIdCategoryRequest{Id: int32(id)})
	})
	if err != nil {
		return nil, category_errors.ErrFindById.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, category_errors.ErrFindById
	}

	return mapCategory(resp.Data), nil
}

func mapCategory(c *pbcategory.CategoryResponse) *db.Category {
	if c == nil {
		return nil
	}
	return &db.Category{
		CategoryID:   c.Id,
		Name:         c.Name,
		Description:  convert.NullableString(c.Description),
		SlugCategory: convert.NullableString(c.SlugCategory),
		CreatedAt:    convert.Timestamp(c.CreatedAt),
		UpdatedAt:    convert.Timestamp(c.UpdatedAt),
	}
}
