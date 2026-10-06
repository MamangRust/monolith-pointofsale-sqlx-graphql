package handler

import (
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pbmerchant_document "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
)

type MerchantQueryHandleGrpc interface {
	pbmerchant.MerchantQueryServiceServer
}

type MerchantCommandHandleGrpc interface {
	pbmerchant.MerchantCommandServiceServer
}

type MerchantDocumentQueryHandleGrpc interface {
	pbmerchant_document.MerchantDocumentQueryServiceServer
}

type MerchantDocumentCommandHandleGrpc interface {
	pbmerchant_document.MerchantDocumentCommandServiceServer
}
