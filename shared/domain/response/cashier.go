package response

type CashierResponse struct {
	ID         int    `json:"id"`
	MerchantID int    `json:"merchant_id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type CashierResponseDeleteAt struct {
	ID         int     `json:"id"`
	MerchantID int     `json:"merchant_id"`
	Name       string  `json:"name"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at"`
}

type CashierResponseMonthSales struct {
	Month       string `json:"month"`
	CashierID   int    `json:"cashier_id"`
	CashierName string `json:"cashier_name"`
	OrderCount  int    `json:"order_count"`
	TotalSales  int    `json:"total_sales"`
}

type CashierResponseYearSales struct {
	Year        string `json:"year"`
	CashierID   int    `json:"cashier_id"`
	CashierName string `json:"cashier_name"`
	OrderCount  int    `json:"order_count"`
	TotalSales  int    `json:"total_sales"`
}

type CashierResponseMonthTotalSales struct {
	Year       string `json:"year"`
	Month      string `json:"month"`
	TotalSales int    `json:"total_sales"`
}

type CashierResponseYearTotalSales struct {
	Year       string `json:"year"`
	TotalSales int    `json:"total_sales"`
}

type ApiResponseCashier struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    *CashierResponse `json:"data"`
}

type ApiResponsesCashier struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    []*CashierResponse `json:"data"`
}

type ApiResponseCashierDeleteAt struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    *CashierResponseDeleteAt `json:"data"`
}

type ApiResponseCashierDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseCashierAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationCashierDeleteAt struct {
	Status     string                     `json:"status"`
	Message    string                     `json:"message"`
	Data       []*CashierResponseDeleteAt `json:"data"`
	Pagination PaginationMeta             `json:"pagination"`
}

type ApiResponsePaginationCashier struct {
	Status     string             `json:"status"`
	Message    string             `json:"message"`
	Data       []*CashierResponse `json:"data"`
	Pagination PaginationMeta     `json:"pagination"`
}

type ApiResponseCashierMonthSales struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    []*CashierResponseMonthSales `json:"data"`
}

type ApiResponseCashierYearSales struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []*CashierResponseYearSales `json:"data"`
}

type ApiResponseCashierMonthlyTotalSales struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message"`
	Data    []*CashierResponseMonthTotalSales `json:"data"`
}

type ApiResponseCashierYearlyTotalSales struct {
	Status  string                           `json:"status"`
	Message string                           `json:"message"`
	Data    []*CashierResponseYearTotalSales `json:"data"`
}
