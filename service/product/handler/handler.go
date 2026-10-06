package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb/product"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-product/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Product        pb.ProductQueryServiceServer
	ProductCommand pb.ProductCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Product:        NewProductQueryHandleGrpc(deps.Service, deps.Logger),
		ProductCommand: NewProductCommandHandleGrpc(deps.Service, deps.Logger),
	}
}
