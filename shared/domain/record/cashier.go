package record

type CashierRecord struct {
	ID         int     `json:"id"`
	MerchantID int     `json:"merchant_id"`
	UserID     int     `json:"user_id"`
	Name       string  `json:"name"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at"`
}

type CashierRecordMonthSales struct {
	Month       string `json:"month"`
	CashierID   int    `json:"cashier_id"`
	CashierName string `json:"cashier_name"`
	OrderCount  int    `json:"order_count"`
	TotalSales  int    `json:"total_sales"`
}

type CashierRecordYearSales struct {
	Year        string `json:"year"`
	CashierID   int    `json:"cashier_id"`
	CashierName string `json:"cashier_name"`
	OrderCount  int    `json:"order_count"`
	TotalSales  int    `json:"total_sales"`
}

type CashierRecordMonthTotalSales struct {
	Year       string `json:"year"`
	Month      string `json:"month"`
	TotalSales int    `json:"total_sales"`
}

type CashierRecordYearTotalSales struct {
	Year       string `json:"year"`
	TotalSales int    `json:"total_sales"`
}
