package requests

import "github.com/go-playground/validator/v10"

type FindAllOrders struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllOrderMerchant struct {
	MerchantID int    `json:"merchant_id" validate:"required"`
	Search     string `json:"search" validate:"required"`
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
}

type MonthTotalRevenue struct {
	Year  int `json:"year" validate:"required"`
	Month int `json:"month" validate:"required"`
}

type MonthTotalRevenueOrder struct {
	OrderID int `json:"order_id" validate:"required"`
	Year    int `json:"year" validate:"required"`
	Month   int `json:"month" validate:"required"`
}

type YearTotalRevenueOrder struct {
	OrderID int `json:"order_id" validate:"required"`
	Year    int `json:"year" validate:"required"`
}

type MonthTotalRevenueMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
	Month      int `json:"month" validate:"required"`
}

type YearTotalRevenueMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
}

type MonthOrderMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
}

type YearOrderMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
}

type CreateOrderRecordRequest struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	CashierID  int `json:"cashier_id"`
	TotalPrice int `json:"total_price"`
}

type UpdateOrderRecordRequest struct {
	OrderID    int `json:"order_id" validate:"required"`
	TotalPrice int `json:"total_price" validate:"required"`
}

type CreateOrderRequest struct {
	MerchantID int                      `json:"merchant_id" validate:"required"`
	CashierID  int                      `json:"cashier_id" validate:"required"`
	Items      []CreateOrderItemRequest `json:"items" validate:"required"`
}

type UpdateOrderRequest struct {
	OrderID *int                     `json:"order_id"`
	Items   []UpdateOrderItemRequest `json:"items" validate:"required"`
}

type CreateOrderItemRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
}

type UpdateOrderItemRequest struct {
	OrderItemID int `json:"order_item_id" validate:"required"`
	ProductID   int `json:"product_id" validate:"required"`
	Quantity    int `json:"quantity" validate:"required"`
}

func (r *CreateOrderRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateOrderRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
