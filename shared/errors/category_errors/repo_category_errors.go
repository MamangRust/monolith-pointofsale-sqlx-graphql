package category_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrGetMonthlyTotalPrice           = errors.ErrInternal.WithMessage("Failed to get monthly total price")
	ErrGetYearlyTotalPrices           = errors.ErrInternal.WithMessage("Failed to get yearly total prices")
	ErrGetMonthlyTotalPriceById       = errors.ErrInternal.WithMessage("Failed to get monthly total price by category ID")
	ErrGetYearlyTotalPricesById       = errors.ErrInternal.WithMessage("Failed to get yearly total prices by category ID")
	ErrGetMonthlyTotalPriceByMerchant = errors.ErrInternal.WithMessage("Failed to get monthly total price by merchant")
	ErrGetYearlyTotalPricesByMerchant = errors.ErrInternal.WithMessage("Failed to get yearly total prices by merchant")

	ErrGetMonthPrice           = errors.ErrInternal.WithMessage("Failed to get month price")
	ErrGetYearPrice            = errors.ErrInternal.WithMessage("Failed to get year price")
	ErrGetMonthPriceByMerchant = errors.ErrInternal.WithMessage("Failed to get month price by merchant")
	ErrGetYearPriceByMerchant  = errors.ErrInternal.WithMessage("Failed to get year price by merchant")
	ErrGetMonthPriceById       = errors.ErrInternal.WithMessage("Failed to get month price by category ID")
	ErrGetYearPriceById        = errors.ErrInternal.WithMessage("Failed to get year price by category ID")

	ErrFindAllCategory = errors.ErrInternal.WithMessage("Failed to find all categories")
	ErrFindById        = errors.ErrInternal.WithMessage("Failed to find category by ID")
	ErrFindByNameAndId = errors.ErrInternal.WithMessage("Failed to find category by name and ID")
	ErrFindByName      = errors.ErrInternal.WithMessage("Failed to find category by name")
	ErrFindByActive    = errors.ErrInternal.WithMessage("Failed to find active categories")
	ErrFindByTrashed   = errors.ErrInternal.WithMessage("Failed to find trashed categories")

	ErrCreateCategory               = errors.ErrInternal.WithMessage("Failed to create category")
	ErrUpdateCategory               = errors.ErrInternal.WithMessage("Failed to update category")
	ErrTrashedCategory              = errors.ErrInternal.WithMessage("Failed to trash category")
	ErrRestoreCategory              = errors.ErrInternal.WithMessage("Failed to restore category")
	ErrDeleteCategoryPermanently    = errors.ErrInternal.WithMessage("Failed to permanently delete category")
	ErrRestoreAllCategories         = errors.ErrInternal.WithMessage("Failed to restore all categories")
	ErrDeleteAllPermanentCategories = errors.ErrInternal.WithMessage("Failed to permanently delete all trashed categories")
)
