package cashier_cache

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

type CashierQueryCache interface {
	GetCachedCashiersCache(ctx context.Context, req *model.FindAllCashierRequest) (*model.APIResponsePaginationCashier, bool)
	SetCachedCashiersCache(ctx context.Context, req *model.FindAllCashierRequest, res *model.APIResponsePaginationCashier)

	GetCachedCashier(ctx context.Context, cashierID int) (*model.APIResponseCashier, bool)
	SetCachedCashier(ctx context.Context, res *model.APIResponseCashier)

	GetCachedCashiersActive(ctx context.Context, req *model.FindAllCashierRequest) (*model.APIResponsePaginationCashierDeleteAt, bool)
	SetCachedCashiersActive(ctx context.Context, req *model.FindAllCashierRequest, res *model.APIResponsePaginationCashierDeleteAt)

	GetCachedCashiersTrashed(ctx context.Context, req *model.FindAllCashierRequest) (*model.APIResponsePaginationCashierDeleteAt, bool)
	SetCachedCashiersTrashed(ctx context.Context, req *model.FindAllCashierRequest, res *model.APIResponsePaginationCashierDeleteAt)

	GetCachedCashiersByMerchant(ctx context.Context, req *model.FindByMerchantCashierRequest) (*model.APIResponsePaginationCashier, bool)
	SetCachedCashiersByMerchant(ctx context.Context, req *model.FindByMerchantCashierRequest, res *model.APIResponsePaginationCashier)
}

type CashierCommandCache interface {
	DeleteCashierCache(ctx context.Context, id int)
}

type CashierStatsCache interface {
	GetMonthlyTotalSalesCache(ctx context.Context, req *model.FindYearMonthTotalSales) (*model.APIResponseCashierMonthlyTotalSales, bool)
	SetMonthlyTotalSalesCache(ctx context.Context, req *model.FindYearMonthTotalSales, res *model.APIResponseCashierMonthlyTotalSales)

	GetYearlyTotalSalesCache(ctx context.Context, year int) (*model.APIResponseCashierYearlyTotalSales, bool)
	SetYearlyTotalSalesCache(ctx context.Context, year int, res *model.APIResponseCashierYearlyTotalSales)

	GetMonthlySalesCache(ctx context.Context, year int) (*model.APIResponseCashierMonthSales, bool)
	SetMonthlySalesCache(ctx context.Context, year int, res *model.APIResponseCashierMonthSales)

	GetYearlySalesCache(ctx context.Context, year int) (*model.APIResponseCashierYearSales, bool)
	SetYearlySalesCache(ctx context.Context, year int, res *model.APIResponseCashierYearSales)
}

type CashierStatsByIdCache interface {
	GetMonthlyTotalSalesByIdCache(ctx context.Context, req *model.FindYearMonthTotalSalesByID) (*model.APIResponseCashierMonthlyTotalSales, bool)
	SetMonthlyTotalSalesByIdCache(ctx context.Context, req *model.FindYearMonthTotalSalesByID, res *model.APIResponseCashierMonthlyTotalSales)

	GetYearlyTotalSalesByIdCache(ctx context.Context, req *model.FindYearTotalSalesByID) (*model.APIResponseCashierYearlyTotalSales, bool)
	SetYearlyTotalSalesByIdCache(ctx context.Context, req *model.FindYearTotalSalesByID, res *model.APIResponseCashierYearlyTotalSales)

	GetMonthlyCashierByIdCache(ctx context.Context, req *model.FindYearCashierByID) (*model.APIResponseCashierMonthSales, bool)
	SetMonthlyCashierByIdCache(ctx context.Context, req *model.FindYearCashierByID, res *model.APIResponseCashierMonthSales)

	GetYearlyCashierByIdCache(ctx context.Context, req *model.FindYearCashierByID) (*model.APIResponseCashierYearSales, bool)
	SetYearlyCashierByIdCache(ctx context.Context, req *model.FindYearCashierByID, res *model.APIResponseCashierYearSales)
}

type CashierStatsByMerchantCache interface {
	GetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *model.FindYearMonthTotalSalesByMerchant) (*model.APIResponseCashierMonthlyTotalSales, bool)
	SetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *model.FindYearMonthTotalSalesByMerchant, res *model.APIResponseCashierMonthlyTotalSales)

	GetYearlyTotalSalesByMerchantCache(ctx context.Context, req *model.FindYearTotalSalesByMerchant) (*model.APIResponseCashierYearlyTotalSales, bool)
	SetYearlyTotalSalesByMerchantCache(ctx context.Context, req *model.FindYearTotalSalesByMerchant, res *model.APIResponseCashierYearlyTotalSales)

	GetMonthlyCashierByMerchantCache(ctx context.Context, req *model.FindYearCashierByMerchant) (*model.APIResponseCashierMonthSales, bool)
	SetMonthlyCashierByMerchantCache(ctx context.Context, req *model.FindYearCashierByMerchant, res *model.APIResponseCashierMonthSales)

	GetYearlyCashierByMerchantCache(ctx context.Context, req *model.FindYearCashierByMerchant) (*model.APIResponseCashierYearSales, bool)
	SetYearlyCashierByMerchantCache(ctx context.Context, req *model.FindYearCashierByMerchant, res *model.APIResponseCashierYearSales)
}
