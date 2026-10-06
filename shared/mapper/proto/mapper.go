package protomapper

type ProtoMapper struct {
	AuthProtoMapper             AuthProtoMapper
	RoleProtoMapper             RoleProtoMapper
	UserProtoMapper             UserProtoMapper
	CategoryProtoMapper         CategoryProtoMapper
	CashierProtoMapper          CashierProtoMapper
	MerchantProtoMapper         MerchantProtoMapper
	OrderItemProtoMapper        OrderItemProtoMapper
	OrderProtoMapper            OrderProtoMapper
	ProductProtoMapper          ProductProtoMapper
	TransactionProtoMapper      TransactionProtoMapper
	MerchantDocumentProtoMapper MerchantDocumentProtoMapper
}

func NewProtoMapper() *ProtoMapper {
	return &ProtoMapper{
		AuthProtoMapper:             NewAuthProtoMapper(),
		RoleProtoMapper:             NewRoleProtoMapper(),
		UserProtoMapper:             NewUserProtoMapper(),
		CategoryProtoMapper:         NewCategoryProtoMapper(),
		CashierProtoMapper:          NewCashierProtoMapper(),
		MerchantProtoMapper:         NewMerchantProtoMaper(),
		OrderItemProtoMapper:        NewOrderItemProtoMapper(),
		OrderProtoMapper:            NewOrderProtoMapper(),
		ProductProtoMapper:          NewProductProtoMapper(),
		TransactionProtoMapper:      NewTransactionProtoMapper(),
		MerchantDocumentProtoMapper: NewMerchantDocumentProtoMapper(),
	}
}
