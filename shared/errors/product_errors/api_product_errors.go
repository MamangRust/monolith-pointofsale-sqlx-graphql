package product_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidBodyCreate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid request body", http.StatusBadRequest)
	}
	ErrApiInvalidBodyUpdate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid request body", http.StatusBadRequest)
	}

	ErrApiInvalidUploadCreate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid upload file", http.StatusBadRequest)
	}
	ErrApiInvalidUploadUpdate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid upload file", http.StatusBadRequest)
	}

	ErrApiProductInvalidMerchantId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant id", http.StatusBadRequest)
	}
	ErrApiProductInvalidCategoryName = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid category name", http.StatusBadRequest)
	}

	ErrApiProductFailedFindAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find all products", http.StatusInternalServerError)
	}
	ErrApiProductFailedFindById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find product by ID", http.StatusInternalServerError)
	}
	ErrApiProductFailedFindByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find product by merchant", http.StatusInternalServerError)
	}
	ErrApiProductFailedFindByCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find product by category", http.StatusInternalServerError)
	}
	ErrApiProductFailedFindByActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find active products", http.StatusInternalServerError)
	}
	ErrApiProductFailedFindByTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find trashed products", http.StatusInternalServerError)
	}

	ErrApiProductFailedCreate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create product", http.StatusInternalServerError)
	}
	ErrApiProductFailedUpdate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update product", http.StatusInternalServerError)
	}

	ErrApiValidateCreateProduct = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create bank request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateProduct = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update bank request", http.StatusBadRequest)
	}

	ErrApiBindCreateProduct = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create bank request", http.StatusBadRequest)
	}

	ErrApiBindUpdateProduct = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update bank request", http.StatusBadRequest)
	}

	ErrApiProductFailedTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trashed product", http.StatusInternalServerError)
	}
	ErrApiProductFailedRestore = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore product", http.StatusInternalServerError)
	}
	ErrApiProductFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete product", http.StatusInternalServerError)
	}
	ErrApiProductFailedRestoreAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all products", http.StatusInternalServerError)
	}
	ErrApiProductFailedDeleteAllPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all products", http.StatusInternalServerError)
	}

	ErrApiProductNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "product not found", http.StatusNotFound)
	}
	ErrApiProductInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid product ID", http.StatusBadRequest)
	}
)
