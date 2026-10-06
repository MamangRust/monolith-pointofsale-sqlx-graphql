package merchant_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrFindAllMerchants = errors.ErrInternal.WithMessage("Failed to find all merchants")
	ErrFindByActive     = errors.ErrInternal.WithMessage("Failed to find active merchants")
	ErrFindByTrashed    = errors.ErrInternal.WithMessage("Failed to find trashed merchants")
	ErrFindById         = errors.ErrInternal.WithMessage("Failed to find merchant by ID")

	ErrCreateMerchant             = errors.ErrInternal.WithMessage("Failed to create merchant")
	ErrUpdateMerchant             = errors.ErrInternal.WithMessage("Failed to update merchant")
	ErrUpdateMerchantStatusFailed = errors.ErrInternal.WithMessage("Failed to update merchant status")

	ErrTrashedMerchant            = errors.ErrInternal.WithMessage("Failed to move merchant to trash")
	ErrRestoreMerchant            = errors.ErrInternal.WithMessage("Failed to restore merchant from trash")
	ErrDeleteMerchantPermanent    = errors.ErrInternal.WithMessage("Failed to permanently delete merchant")
	ErrRestoreAllMerchant         = errors.ErrInternal.WithMessage("Failed to restore all trashed merchants")
	ErrDeleteAllMerchantPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all trashed merchants")
)
