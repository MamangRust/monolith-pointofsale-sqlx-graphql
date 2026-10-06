package transaction_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrGetMonthlyAmountSuccess = errors.ErrInternal.WithMessage("Failed to get monthly amount success")
	ErrGetYearlyAmountSuccess  = errors.ErrInternal.WithMessage("Failed to get yearly amount success")
	ErrGetMonthlyAmountFailed  = errors.ErrInternal.WithMessage("Failed to get monthly amount failed")
	ErrGetYearlyAmountFailed   = errors.ErrInternal.WithMessage("Failed to get yearly amount failed")

	ErrGetMonthlyAmountSuccessByMerchant = errors.ErrInternal.WithMessage("Failed to get monthly amount success by merchant")
	ErrGetYearlyAmountSuccessByMerchant  = errors.ErrInternal.WithMessage("Failed to get yearly amount success by merchant")
	ErrGetMonthlyAmountFailedByMerchant  = errors.ErrInternal.WithMessage("Failed to get monthly amount failed by merchant")
	ErrGetYearlyAmountFailedByMerchant   = errors.ErrInternal.WithMessage("Failed to get yearly amount failed by merchant")

	ErrGetMonthlyTransactionMethod           = errors.ErrInternal.WithMessage("Failed to get monthly transaction method")
	ErrGetYearlyTransactionMethod            = errors.ErrInternal.WithMessage("Failed to get yearly transaction method")
	ErrGetMonthlyTransactionMethodByMerchant = errors.ErrInternal.WithMessage("Failed to get monthly transaction method by merchant")
	ErrGetYearlyTransactionMethodByMerchant  = errors.ErrInternal.WithMessage("Failed to get yearly transaction method by merchant")

	ErrFindAllTransactions = errors.ErrInternal.WithMessage("Failed to find all transactions")
	ErrFindByActive        = errors.ErrInternal.WithMessage("Failed to find active transactions")
	ErrFindByTrashed       = errors.ErrInternal.WithMessage("Failed to find trashed transactions")
	ErrFindByMerchant      = errors.ErrInternal.WithMessage("Failed to find transactions by merchant")
	ErrFindById            = errors.ErrInternal.WithMessage("Failed to find transaction by ID")
	ErrFindByOrderId       = errors.ErrInternal.WithMessage("Failed to find transaction by order ID")

	ErrCreateTransaction             = errors.ErrInternal.WithMessage("Failed to create transaction")
	ErrUpdateTransaction             = errors.ErrInternal.WithMessage("Failed to update transaction")
	ErrTrashTransaction              = errors.ErrInternal.WithMessage("Failed to move transaction to trash")
	ErrRestoreTransaction            = errors.ErrInternal.WithMessage("Failed to restore transaction")
	ErrDeleteTransactionPermanently  = errors.ErrInternal.WithMessage("Failed to permanently delete transaction")
	ErrRestoreAllTransactions        = errors.ErrInternal.WithMessage("Failed to restore all transactions")
	ErrDeleteAllTransactionPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all transactions")
)
