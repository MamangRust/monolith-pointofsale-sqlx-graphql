package response_service

type ResponseServiceMapper struct {
	RoleResponseMapper             RoleResponseMapper
	RefreshTokenResponseMapper     RefreshTokenResponseMapper
	UserResponseMapper             UserResponseMapper
	CategoryResponseMapper         CategoryResponseMapper
	CashierResponseMapper          CashierResponseMapper
	MerchantResponseMapper         MerchantResponseMapper
	OrderResponseMapper            OrderResponseMapper
	OrderItemResponseMapper        OrderItemResponseMapper
	ProductResponseMapper          ProductResponseMapper
	TransactionResponseMapper      TransactionResponseMapper
	MerchantDocumentResponseMapper MerchantDocumentResponseMapper
}

func NewResponseServiceMapper() *ResponseServiceMapper {
	return &ResponseServiceMapper{
		UserResponseMapper:             NewUserResponseMapper(),
		RefreshTokenResponseMapper:     NewRefreshTokenResponseMapper(),
		RoleResponseMapper:             NewRoleResponseMapper(),
		CategoryResponseMapper:         NewCategoryResponseMapper(),
		CashierResponseMapper:          NewCashierResponseMapper(),
		MerchantResponseMapper:         NewMerchantResponseMapper(),
		OrderResponseMapper:            NewOrderResponseMapper(),
		OrderItemResponseMapper:        NewOrderItemResponseMapper(),
		ProductResponseMapper:          NewProductResponseMapper(),
		TransactionResponseMapper:      NewTransactionResponseMapper(),
		MerchantDocumentResponseMapper: NewMerchantDocumentResponseMapper(),
	}
}
