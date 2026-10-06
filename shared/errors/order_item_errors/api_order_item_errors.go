package orderitem_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiOrderItemFailedFindAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find all order items", http.StatusInternalServerError)
	}
	ErrApiOrderItemFailedFindByOrderId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find order items by order ID", http.StatusInternalServerError)
	}
	ErrApiOrderItemFailedFindByActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find active order items", http.StatusInternalServerError)
	}
	ErrApiOrderItemFailedFindByTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find trashed order items", http.StatusInternalServerError)
	}

	ErrApiOrderItemNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "order item not found", http.StatusNotFound)
	}
	ErrApiOrderItemInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid order item ID", http.StatusBadRequest)
	}
)
