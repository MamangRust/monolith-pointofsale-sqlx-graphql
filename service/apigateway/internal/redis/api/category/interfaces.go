package category_cache

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

type CategoryQueryCache interface {
	GetCachedCategoriesCache(ctx context.Context, req *model.FindAllCategoryRequest) (*model.APIResponsePaginationCategory, bool)
	SetCachedCategoriesCache(ctx context.Context, req *model.FindAllCategoryRequest, res *model.APIResponsePaginationCategory)

	GetCachedCategoryActiveCache(ctx context.Context, req *model.FindAllCategoryRequest) (*model.APIResponsePaginationCategoryDeleteAt, bool)
	SetCachedCategoryActiveCache(ctx context.Context, req *model.FindAllCategoryRequest, res *model.APIResponsePaginationCategoryDeleteAt)

	GetCachedCategoryTrashedCache(ctx context.Context, req *model.FindAllCategoryRequest) (*model.APIResponsePaginationCategoryDeleteAt, bool)
	SetCachedCategoryTrashedCache(ctx context.Context, req *model.FindAllCategoryRequest, res *model.APIResponsePaginationCategoryDeleteAt)

	GetCachedCategoryCache(ctx context.Context, id int) (*model.APIResponseCategory, bool)
	SetCachedCategoryCache(ctx context.Context, res *model.APIResponseCategory)
}

type CategoryCommandCache interface {
	DeleteCachedCategoryCache(ctx context.Context, id int)
}

type CategoryStatsCache interface {
	GetCachedMonthTotalPriceCache(ctx context.Context, req *model.FindYearMonthTotalPrices) (*model.APIResponseCategoryMonthlyTotalPrice, bool)
	SetCachedMonthTotalPriceCache(ctx context.Context, req *model.FindYearMonthTotalPrices, res *model.APIResponseCategoryMonthlyTotalPrice)

	GetCachedYearTotalPriceCache(ctx context.Context, year int) (*model.APIResponseCategoryYearlyTotalPrice, bool)
	SetCachedYearTotalPriceCache(ctx context.Context, year int, res *model.APIResponseCategoryYearlyTotalPrice)

	GetCachedMonthPriceCache(ctx context.Context, year int) (*model.APIResponseCategoryMonthPrice, bool)
	SetCachedMonthPriceCache(ctx context.Context, year int, res *model.APIResponseCategoryMonthPrice)

	GetCachedYearPriceCache(ctx context.Context, year int) (*model.APIResponseCategoryYearPrice, bool)
	SetCachedYearPriceCache(ctx context.Context, year int, res *model.APIResponseCategoryYearPrice)
}

type CategoryStatsByIdCache interface {
	GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *model.FindYearMonthTotalPriceByID) (*model.APIResponseCategoryMonthlyTotalPrice, bool)
	SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *model.FindYearMonthTotalPriceByID, res *model.APIResponseCategoryMonthlyTotalPrice)

	GetCachedYearTotalPriceByIdCache(ctx context.Context, req *model.FindYearTotalPriceByID) (*model.APIResponseCategoryYearlyTotalPrice, bool)
	SetCachedYearTotalPriceByIdCache(ctx context.Context, req *model.FindYearTotalPriceByID, res *model.APIResponseCategoryYearlyTotalPrice)

	GetCachedMonthPriceByIdCache(ctx context.Context, req *model.FindYearCategoryByID) (*model.APIResponseCategoryMonthPrice, bool)
	SetCachedMonthPriceByIdCache(ctx context.Context, req *model.FindYearCategoryByID, res *model.APIResponseCategoryMonthPrice)

	GetCachedYearPriceByIdCache(ctx context.Context, req *model.FindYearCategoryByID) (*model.APIResponseCategoryYearPrice, bool)
	SetCachedYearPriceByIdCache(ctx context.Context, req *model.FindYearCategoryByID, res *model.APIResponseCategoryYearPrice)
}

type CategoryStatsByMerchantCache interface {
	GetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *model.FindYearMonthTotalPriceByMerchant) (*model.APIResponseCategoryMonthlyTotalPrice, bool)
	SetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *model.FindYearMonthTotalPriceByMerchant, res *model.APIResponseCategoryMonthlyTotalPrice)

	GetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *model.FindYearTotalPriceByMerchant) (*model.APIResponseCategoryYearlyTotalPrice, bool)
	SetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *model.FindYearTotalPriceByMerchant, res *model.APIResponseCategoryYearlyTotalPrice)

	GetCachedMonthPriceByMerchantCache(ctx context.Context, req *model.FindYearCategoryByMerchant) (*model.APIResponseCategoryMonthPrice, bool)
	SetCachedMonthPriceByMerchantCache(ctx context.Context, req *model.FindYearCategoryByMerchant, res *model.APIResponseCategoryMonthPrice)

	GetCachedYearPriceByMerchantCache(ctx context.Context, req *model.FindYearCategoryByMerchant) (*model.APIResponseCategoryYearPrice, bool)
	SetCachedYearPriceByMerchantCache(ctx context.Context, req *model.FindYearCategoryByMerchant, res *model.APIResponseCategoryYearPrice)
}
