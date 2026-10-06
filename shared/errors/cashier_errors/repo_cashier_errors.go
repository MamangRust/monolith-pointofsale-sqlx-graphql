package cashier_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrGetMonthlyTotalSales           = errors.ErrInternal.WithMessage("Failed to get monthly total sales")
	ErrGetYearlyTotalSales            = errors.ErrInternal.WithMessage("Failed to get yearly total sales")
	ErrGetMonthlyTotalSalesById       = errors.ErrInternal.WithMessage("Failed to get monthly total sales by cashier ID")
	ErrGetYearlyTotalSalesById        = errors.ErrInternal.WithMessage("Failed to get yearly total sales by cashier ID")
	ErrGetMonthlyTotalSalesByMerchant = errors.ErrInternal.WithMessage("Failed to get monthly total sales by merchant")
	ErrGetYearlyTotalSalesByMerchant  = errors.ErrInternal.WithMessage("Failed to get yearly total sales by merchant")

	ErrGetMonthlyCashier           = errors.ErrInternal.WithMessage("Failed to get monthly cashier sales")
	ErrGetYearlyCashier            = errors.ErrInternal.WithMessage("Failed to get yearly cashier sales")
	ErrGetMonthlyCashierByMerchant = errors.ErrInternal.WithMessage("Failed to get monthly cashier sales by merchant")
	ErrGetYearlyCashierByMerchant  = errors.ErrInternal.WithMessage("Failed to get yearly cashier sales by merchant")
	ErrGetMonthlyCashierById       = errors.ErrInternal.WithMessage("Failed to get monthly cashier sales by cashier ID")
	ErrGetYearlyCashierById        = errors.ErrInternal.WithMessage("Failed to get yearly cashier sales by cashier ID")

	ErrFindAllCashiers        = errors.ErrInternal.WithMessage("Failed to find all cashiers")
	ErrFindCashierById        = errors.ErrInternal.WithMessage("Failed to find cashier by ID")
	ErrFindActiveCashiers     = errors.ErrInternal.WithMessage("Failed to find active cashiers")
	ErrFindTrashedCashiers    = errors.ErrInternal.WithMessage("Failed to find trashed cashiers")
	ErrFindCashiersByMerchant = errors.ErrInternal.WithMessage("Failed to find cashiers by merchant")

	ErrCreateCashier              = errors.ErrInternal.WithMessage("Failed to create cashier")
	ErrUpdateCashier              = errors.ErrInternal.WithMessage("Failed to update cashier")
	ErrTrashedCashier             = errors.ErrInternal.WithMessage("Failed to move cashier to trash")
	ErrRestoreCashier             = errors.ErrInternal.WithMessage("Failed to restore cashier from trash")
	ErrDeleteCashierPermanent     = errors.ErrInternal.WithMessage("Failed to permanently delete cashier")
	ErrRestoreAllCashiers         = errors.ErrInternal.WithMessage("Failed to restore all cashiers")
	ErrDeleteAllCashiersPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all trashed cashiers")
)
