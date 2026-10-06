package record

type MerchantRecord struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Address      string  `json:"address"`
	ContactEmail string  `json:"contact_email"`
	ContactPhone string  `json:"contact_phone"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DeletedAt    *string `json:"deleted_at"`
}
