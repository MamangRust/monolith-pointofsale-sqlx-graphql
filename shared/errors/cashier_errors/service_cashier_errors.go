package cashier_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrFailedFindMonthlyTotalSales           = errors.ErrInternal.WithMessage("Failed to find monthly total sales")
	ErrFailedFindYearlyTotalSales            = errors.ErrInternal.WithMessage("Failed to find yearly total sales")
	ErrFailedFindMonthlyTotalSalesById       = errors.ErrInternal.WithMessage("Failed to find monthly total sales by ID")
	ErrFailedFindYearlyTotalSalesById        = errors.ErrInternal.WithMessage("Failed to find yearly total sales by ID")
	ErrFailedFindMonthlyTotalSalesByMerchant = errors.ErrInternal.WithMessage("Failed to find monthly total sales by merchant")
	ErrFailedFindYearlyTotalSalesByMerchant  = errors.ErrInternal.WithMessage("Failed to find yearly total sales by merchant")

	ErrFailedFindMonthlySales             = errors.ErrInternal.WithMessage("Failed to find monthly sales")
	ErrFailedFindYearlySales              = errors.ErrInternal.WithMessage("Failed to find yearly sales")
	ErrFailedFindMonthlyCashierByMerchant = errors.ErrInternal.WithMessage("Failed to find monthly cashier sales by merchant")
	ErrFailedFindYearlyCashierByMerchant  = errors.ErrInternal.WithMessage("Failed to find yearly cashier sales by merchant")
	ErrFailedFindMonthlyCashierById       = errors.ErrInternal.WithMessage("Failed to find monthly cashier sales by ID")
	ErrFailedFindYearlyCashierById        = errors.ErrInternal.WithMessage("Failed to find yearly cashier sales by ID")

	ErrFailedFindAllCashiers       = errors.ErrInternal.WithMessage("Failed to find all cashiers")
	ErrFailedFindCashierById       = errors.ErrInternal.WithMessage("Failed to find cashier by ID")
	ErrFailedFindCashierByActive   = errors.ErrInternal.WithMessage("Failed to find active cashiers")
	ErrFailedFindCashierByTrashed  = errors.ErrInternal.WithMessage("Failed to find trashed cashiers")
	ErrFailedFindCashierByMerchant = errors.ErrInternal.WithMessage("Failed to find cashiers by merchant")

	ErrFailedCreateCashier             = errors.ErrInternal.WithMessage("Failed to create cashier")
	ErrFailedUpdateCashier             = errors.ErrInternal.WithMessage("Failed to update cashier")
	ErrFailedTrashedCashier            = errors.ErrInternal.WithMessage("Failed to trash cashier")
	ErrFailedRestoreCashier            = errors.ErrInternal.WithMessage("Failed to restore cashier")
	ErrFailedDeleteCashierPermanent    = errors.ErrInternal.WithMessage("Failed to permanently delete cashier")
	ErrFailedRestoreAllCashiers        = errors.ErrInternal.WithMessage("Failed to restore all cashiers")
	ErrFailedDeleteAllCashierPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all cashiers")
)
