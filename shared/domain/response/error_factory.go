package response

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"

	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewErrorResponse(message string, code int) *ErrorResponse {
	return &ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	}
}

func NewApiErrorResponse(c echo.Context, statusText string, message string, code int) error {
	return c.JSON(code, ErrorResponse{
		Status:  statusText,
		Message: message,
		Code:    code,
	})
}

func ToGrpcErrorFromErrorResponse(err *ErrorResponse) error {
	if err == nil {
		return nil
	}
	return status.Errorf(codes.Code(err.Code),
		errors.GrpcErrorToJson(&pbcommon.ErrorResponse{
			Status:  err.Status,
			Message: err.Message,
			Code:    int32(err.Code),
		}),
	)
}

func NewGrpcError(statusText string, message string, code int) error {
	return status.Errorf(codes.Code(code),
		errors.GrpcErrorToJson(&pbcommon.ErrorResponse{
			Status:  statusText,
			Message: message,
			Code:    int32(code),
		}),
	)
}
