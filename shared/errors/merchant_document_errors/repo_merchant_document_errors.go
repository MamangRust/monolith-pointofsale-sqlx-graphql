package merchantdocument_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrFindAllMerchantDocumentsFailed     = errors.ErrInternal.WithMessage("Failed to find all merchant documents")
	ErrFindActiveMerchantDocumentsFailed  = errors.ErrInternal.WithMessage("Failed to find active merchant documents")
	ErrFindTrashedMerchantDocumentsFailed = errors.ErrInternal.WithMessage("Failed to find trashed merchant documents")
	ErrFindMerchantDocumentByIdFailed     = errors.ErrInternal.WithMessage("Failed to find merchant document by ID")

	ErrCreateMerchantDocumentFailed       = errors.ErrInternal.WithMessage("Failed to create merchant document")
	ErrUpdateMerchantDocumentFailed       = errors.ErrInternal.WithMessage("Failed to update merchant document")
	ErrUpdateMerchantDocumentStatusFailed = errors.ErrInternal.WithMessage("Failed to update merchant document status")

	ErrTrashedMerchantDocumentFailed             = errors.ErrInternal.WithMessage("Failed to soft-delete (trash) merchant document")
	ErrRestoreMerchantDocumentFailed             = errors.ErrInternal.WithMessage("Failed to restore merchant document")
	ErrDeleteMerchantDocumentPermanentFailed     = errors.ErrInternal.WithMessage("Failed to permanently delete merchant document")
	ErrRestoreAllMerchantDocumentsFailed         = errors.ErrInternal.WithMessage("Failed to restore all merchant documents")
	ErrDeleteAllMerchantDocumentsPermanentFailed = errors.ErrInternal.WithMessage("Failed to permanently delete all merchant documents")
)
