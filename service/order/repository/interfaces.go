package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type OrderStatsRepository interface {
	GetMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*db.GetMonthlyTotalRevenueRow, error)
	GetYearlyTotalRevenue(ctx context.Context, year int) ([]*db.GetYearlyTotalRevenueRow, error)
	GetMonthlyOrder(ctx context.Context, year int) ([]*db.GetMonthlyOrderRow, error)
	GetYearlyOrder(ctx context.Context, year int) ([]*db.GetYearlyOrderRow, error)
}

type OrderStatByMerchantRepository interface {
	GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error)
	GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, error)
	GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, error)
	GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, error)
}

type OrderQueryRepository interface {
	FindAllOrders(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*db.GetOrdersByMerchantRow, *int, error)
	FindById(ctx context.Context, orderID int) (*db.Order, error)
	FindByTrashedId(ctx context.Context, orderID int) (*db.Order, error)
}

type OrderCommandRepository interface {
	DeleteOrder(ctx context.Context, orderID int) error
	CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*db.Order, error)
	UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*db.Order, error)
	FindAllTrashed(ctx context.Context) ([]*db.Order, error)
	TrashedOrder(ctx context.Context, orderID int) (*db.Order, error)
	RestoreOrder(ctx context.Context, orderID int) (*db.Order, error)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}
