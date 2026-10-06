package errors

import stderrors "errors"

// NewStdlibError returns a plain (non-AppError) error, used to verify that
// ToGrpcError falls back to codes.Internal for non-AppError values.
func NewStdlibError(msg string) error {
	return stderrors.New(msg)
}
