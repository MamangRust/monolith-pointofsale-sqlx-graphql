package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type CashierStatsRepository interface {
	GetMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*db.GetMonthlyTotalSalesCashierRow, error)
	GetYearlyTotalSales(ctx context.Context, year int) ([]*db.GetYearlyTotalSalesCashierRow, error)

	GetMonthyCashier(ctx context.Context, year int) ([]*db.GetMonthlyCashierRow, error)
	GetYearlyCashier(ctx context.Context, year int) ([]*db.GetYearlyCashierRow, error)
}

type CashierStatByIdRepository interface {
	GetMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*db.GetMonthlyTotalSalesByIdRow, error)
	GetYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*db.GetYearlyTotalSalesByIdRow, error)

	GetMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*db.GetMonthlyCashierByCashierIdRow, error)
	GetYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*db.GetYearlyCashierByCashierIdRow, error)
}

type CashierStatByMerchantRepository interface {
	GetMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*db.GetMonthlyTotalSalesByMerchantRow, error)
	GetYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*db.GetYearlyTotalSalesByMerchantRow, error)

	GetMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*db.GetMonthlyCashierByMerchantRow, error)
	GetYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*db.GetYearlyCashierByMerchantRow, error)
}

type CashierQueryRepository interface {
	FindAllCashiers(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersRow, *int, error)
	FindById(ctx context.Context, cashier_id int) (*db.Cashier, error)
	FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*db.GetCashiersByMerchantRow, *int, error)
}

type CashierCommandRepository interface {
	CreateCashier(ctx context.Context, request *requests.CreateCashierRequest) (*db.Cashier, error)
	UpdateCashier(ctx context.Context, request *requests.UpdateCashierRequest) (*db.Cashier, error)
	TrashedCashier(ctx context.Context, cashier_id int) (*db.Cashier, error)
	RestoreCashier(ctx context.Context, cashier_id int) (*db.Cashier, error)
	DeleteCashierPermanent(ctx context.Context, cashier_id int) (bool, error)
	RestoreAllCashier(ctx context.Context) (bool, error)
	DeleteAllCashierPermanent(ctx context.Context) (bool, error)
}
