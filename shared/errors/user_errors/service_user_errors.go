package user_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrUserNotFoundRes       = errors.ErrNotFound.WithMessage("User not found")
	ErrAccountLocked         = errors.ErrBadRequest.WithMessage("Account is locked")
	ErrUserEmailAlready      = errors.ErrBadRequest.WithMessage("User email already exists")
	ErrUserPassword          = errors.ErrBadRequest.WithMessage("Failed invalid password")
	ErrFailedPasswordNoMatch = errors.ErrBadRequest.WithMessage("Failed password not match")
	ErrFailedFindAll         = errors.ErrInternal.WithMessage("Failed to fetch users")
	ErrFailedFindActive      = errors.ErrInternal.WithMessage("Failed to fetch active users")
	ErrFailedFindTrashed     = errors.ErrInternal.WithMessage("Failed to fetch trashed users")
	ErrInternalServerError   = errors.ErrInternal.WithMessage("Internal server error")

	ErrFailedSendEmail = errors.ErrInternal.WithMessage("Failed to send email")

	ErrFailedCreateUser = errors.ErrInternal.WithMessage("Failed to create user")
	ErrFailedUpdateUser = errors.ErrInternal.WithMessage("Failed to update user")

	ErrFailedTrashedUser     = errors.ErrInternal.WithMessage("Failed to move user to trash")
	ErrFailedRestoreUser     = errors.ErrInternal.WithMessage("Failed to restore user")
	ErrFailedDeletePermanent = errors.ErrInternal.WithMessage("Failed to delete user permanently")

	ErrFailedRestoreAll = errors.ErrInternal.WithMessage("Failed to restore all users")
	ErrFailedDeleteAll  = errors.ErrInternal.WithMessage("Failed to delete all users permanently")
)
