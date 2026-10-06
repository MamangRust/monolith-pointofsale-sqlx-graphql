package order_cache

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

type OrderStatsCache interface {
	GetMonthlyTotalRevenueCache(ctx context.Context, req *model.FindYearMonthTotalRevenueInput) (*model.APIResponseOrderMonthlyTotalRevenue, bool)
	SetMonthlyTotalRevenueCache(ctx context.Context, req *model.FindYearMonthTotalRevenueInput, res *model.APIResponseOrderMonthlyTotalRevenue)

	GetYearlyTotalRevenueCache(ctx context.Context, year int) (*model.APIResponseOrderYearlyTotalRevenue, bool)
	SetYearlyTotalRevenueCache(ctx context.Context, year int, res *model.APIResponseOrderYearlyTotalRevenue)

	GetMonthlyOrderCache(ctx context.Context, year int) (*model.APIResponseOrderMonthly, bool)
	SetMonthlyOrderCache(ctx context.Context, year int, res *model.APIResponseOrderMonthly)

	GetYearlyOrderCache(ctx context.Context, year int) (*model.APIResponseOrderYearly, bool)
	SetYearlyOrderCache(ctx context.Context, year int, res *model.APIResponseOrderYearly)
}

type OrderStatsByMerchantCache interface {
	GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *model.FindYearMonthTotalRevenueByMerchantInput) (*model.APIResponseOrderMonthlyTotalRevenue, bool)
	SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *model.FindYearMonthTotalRevenueByMerchantInput, res *model.APIResponseOrderMonthlyTotalRevenue)

	GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *model.FindYearTotalRevenueByMerchantInput) (*model.APIResponseOrderYearlyTotalRevenue, bool)
	SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *model.FindYearTotalRevenueByMerchantInput, res *model.APIResponseOrderYearlyTotalRevenue)

	GetMonthlyOrderByMerchantCache(ctx context.Context, req *model.FindYearOrderByMerchantInput) (*model.APIResponseOrderMonthly, bool)
	SetMonthlyOrderByMerchantCache(ctx context.Context, req *model.FindYearOrderByMerchantInput, res *model.APIResponseOrderMonthly)

	GetYearlyOrderByMerchantCache(ctx context.Context, req *model.FindYearOrderByMerchantInput) (*model.APIResponseOrderYearly, bool)
	SetYearlyOrderByMerchantCache(ctx context.Context, req *model.FindYearOrderByMerchantInput, res *model.APIResponseOrderYearly)
}

type OrderQueryCache interface {
	GetOrderAllCache(ctx context.Context, req *model.FindAllOrderInput) (*model.APIResponsePaginationOrder, bool)
	SetOrderAllCache(ctx context.Context, req *model.FindAllOrderInput, res *model.APIResponsePaginationOrder)

	GetCachedOrderCache(ctx context.Context, orderID int) (*model.APIResponseOrder, bool)
	SetCachedOrderCache(ctx context.Context, res *model.APIResponseOrder)

	GetCachedOrderMerchant(ctx context.Context, req *model.FindAllOrderMerchantInput) (*model.APIResponsePaginationOrder, bool)
	SetCachedOrderMerchant(ctx context.Context, req *model.FindAllOrderMerchantInput, res *model.APIResponsePaginationOrder)

	GetOrderActiveCache(ctx context.Context, req *model.FindAllOrderInput) (*model.APIResponsePaginationOrderDeleteAt, bool)
	SetOrderActiveCache(ctx context.Context, req *model.FindAllOrderInput, res *model.APIResponsePaginationOrderDeleteAt)

	GetOrderTrashedCache(ctx context.Context, req *model.FindAllOrderInput) (*model.APIResponsePaginationOrderDeleteAt, bool)
	SetOrderTrashedCache(ctx context.Context, req *model.FindAllOrderInput, res *model.APIResponsePaginationOrderDeleteAt)
}

type OrderCommandCache interface {
	DeleteOrderCache(ctx context.Context, id int)
}
