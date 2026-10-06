package handler

import (
	"github.com/MamangRust/monolith-pointofsale-merchant/service"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pbmerchant_document "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Merchant                pbmerchant.MerchantQueryServiceServer
	MerchantCommand         pbmerchant.MerchantCommandServiceServer
	MerchantDocument        pbmerchant_document.MerchantDocumentQueryServiceServer
	MerchantDocumentCommand pbmerchant_document.MerchantDocumentCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Merchant:                NewMerchantQueryHandleGrpc(deps.Service, deps.Logger),
		MerchantCommand:         NewMerchantCommandHandleGrpc(deps.Service, deps.Logger),
		MerchantDocument:        NewMerchantDocumentQueryHandleGrpc(deps.Service, deps.Logger),
		MerchantDocumentCommand: NewMerchantDocumentCommandHandleGrpc(deps.Service, deps.Logger),
	}
}
