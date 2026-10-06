package merchant_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrFailedFindAllMerchants            = errors.ErrInternal.WithMessage("Failed to find all merchants")
	ErrFailedFindMerchantsByActive       = errors.ErrInternal.WithMessage("Failed to find active merchants")
	ErrFailedFindMerchantsByTrashed      = errors.ErrInternal.WithMessage("Failed to find trashed merchants")
	ErrFailedFindMerchantById            = errors.ErrInternal.WithMessage("Failed to find merchant by ID")
	ErrFailedCreateMerchant              = errors.ErrInternal.WithMessage("Failed to create merchant")
	ErrFailedUpdateMerchant              = errors.ErrInternal.WithMessage("Failed to update merchant")
	ErrFailedTrashMerchant               = errors.ErrInternal.WithMessage("Failed to trash merchant")
	ErrFailedRestoreMerchant             = errors.ErrInternal.WithMessage("Failed to restore merchant")
	ErrFailedDeleteMerchantPermanent     = errors.ErrInternal.WithMessage("Failed to permanently delete merchant")
	ErrFailedRestoreAllMerchants         = errors.ErrInternal.WithMessage("Failed to restore all merchants")
	ErrFailedDeleteAllMerchantsPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all merchants")

	ErrFailedSendEmail = errors.ErrInternal.WithMessage("Failed to send email")
)
