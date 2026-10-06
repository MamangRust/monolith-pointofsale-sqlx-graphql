package response_api

type ResponseApiMapper struct {
	AuthResponseMapper        AuthResponseMapper
	RoleResponseMapper        RoleResponseMapper
	UserResponseMapper        UserResponseMapper
	CategoryResponseMapper    CategoryResponseMapper
	CashierResponseMapper     CashierResponseMapper
	MerchantResponseMapper    MerchantResponseMapper
	OrderItemResponseMapper   OrderItemResponseMapper
	OrderResponseMapper       OrderResponseMapper
	ProductResponseMapper     ProductResponseMapper
	TransactionResponseMapper TransactionResponseMapper
	MerchantDocumentProMapper MerchantDocumentResponseMapper
}

func NewResponseApiMapper() *ResponseApiMapper {
	return &ResponseApiMapper{
		AuthResponseMapper:        NewAuthResponseMapper(),
		UserResponseMapper:        NewUserResponseMapper(),
		RoleResponseMapper:        NewRoleResponseMapper(),
		CategoryResponseMapper:    NewCategoryResponseMapper(),
		CashierResponseMapper:     NewCashierResponseMapper(),
		MerchantResponseMapper:    NewMerchantResponseMapper(),
		OrderItemResponseMapper:   NewOrderItemResponseMapper(),
		OrderResponseMapper:       NewOrderResponseMapper(),
		ProductResponseMapper:     NewProductResponseMapper(),
		TransactionResponseMapper: NewTransactionResponseMapper(),
		MerchantDocumentProMapper: NewMerchantDocumentResponseMapper(),
	}
}
