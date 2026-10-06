package order_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrGetMonthlyTotalRevenue           = errors.ErrInternal.WithMessage("Failed to get monthly total revenue")
	ErrGetYearlyTotalRevenue            = errors.ErrInternal.WithMessage("Failed to get yearly total revenue")
	ErrGetMonthlyTotalRevenueById       = errors.ErrInternal.WithMessage("Failed to get monthly total revenue by order ID")
	ErrGetYearlyTotalRevenueById        = errors.ErrInternal.WithMessage("Failed to get yearly total revenue by order ID")
	ErrGetMonthlyTotalRevenueByMerchant = errors.ErrInternal.WithMessage("Failed to get monthly total revenue by merchant")
	ErrGetYearlyTotalRevenueByMerchant  = errors.ErrInternal.WithMessage("Failed to get yearly total revenue by merchant")

	ErrGetMonthlyOrder           = errors.ErrInternal.WithMessage("Failed to get monthly orders")
	ErrGetYearlyOrder            = errors.ErrInternal.WithMessage("Failed to get yearly orders")
	ErrGetMonthlyOrderByMerchant = errors.ErrInternal.WithMessage("Failed to get monthly orders by merchant")
	ErrGetYearlyOrderByMerchant  = errors.ErrInternal.WithMessage("Failed to get yearly orders by merchant")

	ErrFindAllOrders                = errors.ErrInternal.WithMessage("Failed to find all orders")
	ErrFindByActive                 = errors.ErrInternal.WithMessage("Failed to find active orders")
	ErrFindByTrashed                = errors.ErrInternal.WithMessage("Failed to find trashed orders")
	ErrFindByMerchant               = errors.ErrInternal.WithMessage("Failed to find orders by merchant")
	ErrFindById                     = errors.ErrInternal.WithMessage("Failed to find order by ID")
	ErrFindByTrashedId              = errors.ErrInternal.WithMessage("Failed to find trashed order by ID")
	ErrFindAllTrashed               = errors.ErrInternal.WithMessage("Failed to find all trashed orders")
	ErrCreateOrder                  = errors.ErrInternal.WithMessage("Failed to create order")
	ErrUpdateOrder                  = errors.ErrInternal.WithMessage("Failed to update order")
	ErrTrashedOrder                 = errors.ErrInternal.WithMessage("Failed to move order to trash")
	ErrRestoreOrder                 = errors.ErrInternal.WithMessage("Failed to restore order from trash")
	ErrRestoreOrderNotFound         = errors.ErrNotFound.WithMessage("Order is no longer trashed")
	ErrDeleteOrderPermanentNotFound = errors.ErrNotFound.WithMessage("Order is no longer trashed")
	ErrDeleteOrderPermanent         = errors.ErrInternal.WithMessage("Failed to permanently delete order")
	ErrRestoreAllOrder              = errors.ErrInternal.WithMessage("Failed to restore all trashed orders")
	ErrDeleteAllOrderPermanent      = errors.ErrInternal.WithMessage("Failed to permanently delete all trashed orders")
)
