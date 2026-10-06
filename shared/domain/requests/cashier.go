package requests

import "github.com/go-playground/validator/v10"

type FindAllCashiers struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllCashierMerchant struct {
	MerchantID int    `json:"merchant_id" validate:"required"`
	Search     string `json:"search" validate:"required"`
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
}

type MonthTotalSales struct {
	Year  int `json:"year" validate:"required"`
	Month int `json:"month" validate:"required"`
}

type MonthTotalSalesCashier struct {
	CashierID int `json:"cashier_id" validate:"required"`
	Year      int `json:"year" validate:"required"`
	Month     int `json:"month" validate:"required"`
}

type YearTotalSalesCashier struct {
	CashierID int `json:"cashier_id" validate:"required"`
	Year      int `json:"year" validate:"required"`
}

type MonthTotalSalesMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
	Month      int `json:"month" validate:"required"`
}

type YearTotalSalesMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
}

type MonthCashierMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
}

type YearCashierMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
}

type MonthCashierId struct {
	CashierID int `json:"cashier_id" validate:"required"`
	Year      int `json:"year" validate:"required"`
}

type YearCashierId struct {
	CashierID int `json:"cashier_id" validate:"required"`
	Year      int `json:"year" validate:"required"`
}

type CreateCashierRequest struct {
	MerchantID int    `json:"merchant_id" validate:"required"`
	UserID     int    `json:"user_id" validate:"required"`
	Name       string `json:"name" validate:"required"`
}

type UpdateCashierRequest struct {
	CashierID *int   `json:"cashier_id"`
	Name      string `json:"name" validate:"required"`
}

func (r *CreateCashierRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateCashierRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
