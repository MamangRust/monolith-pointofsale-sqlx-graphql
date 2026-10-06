package service

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type CashierStatsService interface {
	FindMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*db.GetMonthlyTotalSalesCashierRow, error)
	FindYearlyTotalSales(ctx context.Context, year int) ([]*db.GetYearlyTotalSalesCashierRow, error)

	FindMonthlySales(ctx context.Context, year int) ([]*db.GetMonthlyCashierRow, error)
	FindYearlySales(ctx context.Context, year int) ([]*db.GetYearlyCashierRow, error)
}

type CashierStatsByIdService interface {
	FindMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*db.GetMonthlyTotalSalesByIdRow, error)
	FindYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*db.GetYearlyTotalSalesByIdRow, error)
	FindMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*db.GetMonthlyCashierByCashierIdRow, error)
	FindYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*db.GetYearlyCashierByCashierIdRow, error)
}

type CashierStatsByMerchant interface {
	FindMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*db.GetMonthlyTotalSalesByMerchantRow, error)
	FindYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*db.GetYearlyTotalSalesByMerchantRow, error)

	FindMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*db.GetMonthlyCashierByMerchantRow, error)
	FindYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*db.GetYearlyCashierByMerchantRow, error)
}

type CashierQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersRow, *int, error)
	FindById(ctx context.Context, cashierID int) (*db.Cashier, error)
	FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*db.GetCashiersByMerchantRow, *int, error)
}

type CashierCommandService interface {
	CreateCashier(ctx context.Context, req *requests.CreateCashierRequest) (*db.Cashier, error)
	UpdateCashier(ctx context.Context, req *requests.UpdateCashierRequest) (*db.Cashier, error)
	TrashedCashier(ctx context.Context, cashierID int) (*db.Cashier, error)
	RestoreCashier(ctx context.Context, cashierID int) (*db.Cashier, error)
	DeleteCashierPermanent(ctx context.Context, cashierID int) (bool, error)
	RestoreAllCashier(ctx context.Context) (bool, error)
	DeleteAllCashierPermanent(ctx context.Context) (bool, error)
}
