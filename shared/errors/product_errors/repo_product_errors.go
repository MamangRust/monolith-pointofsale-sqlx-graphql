package product_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrFindAllProducts         = errors.ErrInternal.WithMessage("Failed to find all products")
	ErrFindByActive            = errors.ErrInternal.WithMessage("Failed to find active products")
	ErrFindByTrashed           = errors.ErrInternal.WithMessage("Failed to find trashed products")
	ErrFindByMerchant          = errors.ErrInternal.WithMessage("Failed to find products by merchant")
	ErrFindByCategory          = errors.ErrInternal.WithMessage("Failed to find products by category")
	ErrFindById                = errors.ErrInternal.WithMessage("Failed to find product by ID")
	ErrFindByIdTrashed         = errors.ErrInternal.WithMessage("Failed to find trashed product by ID")
	ErrCreateProduct           = errors.ErrInternal.WithMessage("Failed to create product")
	ErrUpdateProduct           = errors.ErrInternal.WithMessage("Failed to update product")
	ErrUpdateProductCountStock = errors.ErrInternal.WithMessage("Failed to update product stock count")
	// ErrInsufficientProductStock: an atomic stock adjustment would drive
	// count_in_stock negative → caller sent a bad request, not a server fault.
	ErrInsufficientProductStock  = errors.ErrBadRequest.WithMessage("Insufficient product stock")
	ErrTrashedProduct            = errors.ErrInternal.WithMessage("Failed to move product to trash")
	ErrRestoreProduct            = errors.ErrInternal.WithMessage("Failed to restore product")
	ErrDeleteProductPermanent    = errors.ErrInternal.WithMessage("Failed to permanently delete product")
	ErrRestoreAllProducts        = errors.ErrInternal.WithMessage("Failed to restore all products")
	ErrDeleteAllProductPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all products")
)
