package service

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type OrderStatsService interface {
	FindMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*db.GetMonthlyTotalRevenueRow, error)
	FindYearlyTotalRevenue(ctx context.Context, year int) ([]*db.GetYearlyTotalRevenueRow, error)

	FindMonthlyOrder(ctx context.Context, year int) ([]*db.GetMonthlyOrderRow, error)
	FindYearlyOrder(ctx context.Context, year int) ([]*db.GetYearlyOrderRow, error)
}

type OrderStatByMerchantService interface {
	FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error)
	FindYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, error)

	FindMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, error)
	FindYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, error)
}

type OrderQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersRow, *int, error)
	FindById(ctx context.Context, orderID int) (*db.Order, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*db.GetOrdersByMerchantRow, *int, error)
}

type OrderCommandService interface {
	CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*db.Order, error)
	UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*db.Order, error)
	TrashedOrder(ctx context.Context, orderID int) (*db.Order, error)
	RestoreOrder(ctx context.Context, orderID int) (*db.Order, error)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAllOrder(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}
