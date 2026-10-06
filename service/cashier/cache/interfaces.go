package mencache

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type CashierQueryCache interface {
	GetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersRow, *int, bool)
	SetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers, res []*db.GetCashiersRow, total *int)

	GetCachedCashier(ctx context.Context, cashierID int) (*db.Cashier, bool)
	SetCachedCashier(ctx context.Context, res *db.Cashier)

	GetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersActiveRow, *int, bool)
	SetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers, res []*db.GetCashiersActiveRow, total *int)

	GetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersTrashedRow, *int, bool)
	SetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers, res []*db.GetCashiersTrashedRow, total *int)

	GetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*db.GetCashiersByMerchantRow, *int, bool)
	SetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant, res []*db.GetCashiersByMerchantRow, total *int)
}

type CashierCommandCache interface {
	DeleteCashierCache(ctx context.Context, id int)
	DeleteCashierListCache(ctx context.Context)
}

type CashierStatsCache interface {
	GetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales) ([]*db.GetMonthlyTotalSalesCashierRow, bool)
	SetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales, res []*db.GetMonthlyTotalSalesCashierRow)

	GetYearlyTotalSalesCache(ctx context.Context, year int) ([]*db.GetYearlyTotalSalesCashierRow, bool)
	SetYearlyTotalSalesCache(ctx context.Context, year int, res []*db.GetYearlyTotalSalesCashierRow)

	GetMonthlySalesCache(ctx context.Context, year int) ([]*db.GetMonthlyCashierRow, bool)
	SetMonthlySalesCache(ctx context.Context, year int, res []*db.GetMonthlyCashierRow)

	GetYearlySalesCache(ctx context.Context, year int) ([]*db.GetYearlyCashierRow, bool)
	SetYearlySalesCache(ctx context.Context, year int, res []*db.GetYearlyCashierRow)
}

type CashierStatsByIdCache interface {
	GetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*db.GetMonthlyTotalSalesByIdRow, bool)
	SetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier, res []*db.GetMonthlyTotalSalesByIdRow)

	GetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*db.GetYearlyTotalSalesByIdRow, bool)
	SetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier, res []*db.GetYearlyTotalSalesByIdRow)

	GetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId) ([]*db.GetMonthlyCashierByCashierIdRow, bool)
	SetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId, res []*db.GetMonthlyCashierByCashierIdRow)

	GetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId) ([]*db.GetYearlyCashierByCashierIdRow, bool)
	SetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId, res []*db.GetYearlyCashierByCashierIdRow)
}

type CashierStatsByMerchantCache interface {
	GetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*db.GetMonthlyTotalSalesByMerchantRow, bool)
	SetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant, res []*db.GetMonthlyTotalSalesByMerchantRow)

	GetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*db.GetYearlyTotalSalesByMerchantRow, bool)
	SetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant, res []*db.GetYearlyTotalSalesByMerchantRow)

	GetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant) ([]*db.GetMonthlyCashierByMerchantRow, bool)
	SetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant, res []*db.GetMonthlyCashierByMerchantRow)

	GetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant) ([]*db.GetYearlyCashierByMerchantRow, bool)
	SetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant, res []*db.GetYearlyCashierByMerchantRow)
}
