package mencache

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type OrderStatsCache interface {
	GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) ([]*db.GetMonthlyTotalRevenueRow, bool)
	SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, res []*db.GetMonthlyTotalRevenueRow)

	GetYearlyTotalRevenueCache(ctx context.Context, year int) ([]*db.GetYearlyTotalRevenueRow, bool)
	SetYearlyTotalRevenueCache(ctx context.Context, year int, res []*db.GetYearlyTotalRevenueRow)

	GetMonthlyOrderCache(ctx context.Context, year int) ([]*db.GetMonthlyOrderRow, bool)
	SetMonthlyOrderCache(ctx context.Context, year int, res []*db.GetMonthlyOrderRow)

	GetYearlyOrderCache(ctx context.Context, year int) ([]*db.GetYearlyOrderRow, bool)
	SetYearlyOrderCache(ctx context.Context, year int, res []*db.GetYearlyOrderRow)
}

type OrderStatsByMerchantCache interface {
	GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, bool)
	SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, res []*db.GetMonthlyTotalRevenueByMerchantRow)

	GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, bool)
	SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, res []*db.GetYearlyTotalRevenueByMerchantRow)

	GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, bool)
	SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, res []*db.GetMonthlyOrderByMerchantRow)

	GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, bool)
	SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, res []*db.GetYearlyOrderByMerchantRow)
}

type OrderQueryCache interface {
	GetOrderAllCache(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersRow, *int, bool)
	SetOrderAllCache(ctx context.Context, req *requests.FindAllOrders, data []*db.GetOrdersRow, total *int)

	GetCachedOrderCache(ctx context.Context, orderID int) (*db.Order, bool)
	SetCachedOrderCache(ctx context.Context, data *db.Order)

	GetCachedOrderMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*db.GetOrdersByMerchantRow, *int, bool)
	SetCachedOrderMerchant(ctx context.Context, req *requests.FindAllOrderMerchant, res []*db.GetOrdersByMerchantRow, total *int)

	GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersActiveRow, *int, bool)
	SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrders, data []*db.GetOrdersActiveRow, total *int)

	GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersTrashedRow, *int, bool)
	SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrders, data []*db.GetOrdersTrashedRow, total *int)
}

type OrderCommandCache interface {
	DeleteOrderCache(ctx context.Context, id int)
	DeleteOrderAllCache(ctx context.Context)
}
