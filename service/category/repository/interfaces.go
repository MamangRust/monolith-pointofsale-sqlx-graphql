package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type CategoryStatsRepository interface {
	GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*db.GetMonthlyTotalPriceRow, error)
	GetYearlyTotalPrices(ctx context.Context, year int) ([]*db.GetYearlyTotalPriceRow, error)

	GetMonthPrice(ctx context.Context, year int) ([]*db.GetMonthlyCategoryRow, error)
	GetYearPrice(ctx context.Context, year int) ([]*db.GetYearlyCategoryRow, error)
}

type CategoryStatsByIdRepository interface {
	GetMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*db.GetMonthlyTotalPriceByIdRow, error)
	GetYearlyTotalPricesById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*db.GetYearlyTotalPriceByIdRow, error)

	GetMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*db.GetMonthlyCategoryByIdRow, error)
	GetYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*db.GetYearlyCategoryByIdRow, error)
}

type CategoryStatsByMerchantRepository interface {
	GetMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*db.GetMonthlyTotalPriceByMerchantRow, error)
	GetYearlyTotalPricesByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*db.GetYearlyTotalPriceByMerchantRow, error)

	GetMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*db.GetMonthlyCategoryByMerchantRow, error)
	GetYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*db.GetYearlyCategoryByMerchantRow, error)
}

type CategoryQueryRepository interface {
	FindAllCategory(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, *int, error)
	FindById(ctx context.Context, category_id int) (*db.Category, error)
	FindByNameAndId(ctx context.Context, req *requests.CategoryNameAndId) (*db.Category, error)
	FindByName(ctx context.Context, name string) (*db.Category, error)

	FindByIdTrashed(ctx context.Context, category_id int) (*db.Category, error)

	FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, *int, error)
}

type CategoryCommandRepository interface {
	CreateCategory(ctx context.Context, request *requests.CreateCategoryRequest) (*db.Category, error)
	UpdateCategory(ctx context.Context, request *requests.UpdateCategoryRequest) (*db.Category, error)
	TrashedCategory(ctx context.Context, category_id int) (*db.Category, error)
	RestoreCategory(ctx context.Context, category_id int) (*db.Category, error)
	DeleteCategoryPermanently(ctx context.Context, category_id int) (bool, error)
	RestoreAllCategories(ctx context.Context) (bool, error)
	DeleteAllPermanentCategories(ctx context.Context) (bool, error)
}
