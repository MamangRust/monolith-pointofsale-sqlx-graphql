package category_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiCategoryInvalidMerchantId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant id", http.StatusBadRequest)
	}

	ErrApiCategoryInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid year", http.StatusBadRequest)
	}

	ErrApiCategoryInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid month", http.StatusBadRequest)
	}

	ErrApiCategoryFailedFindAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find all categories", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedFindById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find category by id", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedFindByActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find active categories", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedFindByTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find trashed categories", http.StatusInternalServerError)
	}

	ErrApiCategoryFailedMonthTotalPrice = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly total pricing", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedYearTotalPrice = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly total pricing", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedMonthTotalPriceByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly total pricing by merchant", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedYearTotalPriceByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly total pricing by merchant", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedMonthTotalPriceById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly total pricing by category id", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedYearTotalPriceById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly total pricing by category id", http.StatusInternalServerError)
	}

	ErrApiCategoryFailedMonthPrice = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly pricing", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedYearPrice = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly pricing", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedMonthPriceByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly pricing by merchant", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedYearPriceByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly pricing by merchant", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedMonthPriceById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly pricing by category id", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedYearPriceById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly pricing by category id", http.StatusInternalServerError)
	}

	ErrApiCategoryFailedCreate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create category", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedUpdate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update category", http.StatusInternalServerError)
	}

	ErrApiValidateCreateCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create category request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update category request", http.StatusBadRequest)
	}

	ErrApiBindCreateCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create category request", http.StatusBadRequest)
	}

	ErrApiBindUpdateCategory = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update category request", http.StatusBadRequest)
	}

	ErrApiCategoryFailedTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash category", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedRestore = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore category", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete category", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedRestoreAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all categories", http.StatusInternalServerError)
	}
	ErrApiCategoryFailedDeleteAllPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all categories", http.StatusInternalServerError)
	}

	ErrApiCategoryNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "category not found", http.StatusNotFound)
	}
	ErrApiCategoryInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid category id", http.StatusBadRequest)
	}
)
