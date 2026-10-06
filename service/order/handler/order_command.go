package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-pointofsale-order/service"
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/order_errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderCommandHandleGrpc struct {
	pborder.UnimplementedOrderCommandServiceServer
	orderQuery   service.OrderQueryService
	orderCommand service.OrderCommandService
	logger       logger.LoggerInterface
}

func NewOrderCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborder.OrderCommandServiceServer {
	return &orderCommandHandleGrpc{
		orderQuery:   service.OrderQuery,
		orderCommand: service.OrderCommand,
		logger:       logger,
	}
}

func (s *orderCommandHandleGrpc) FindByTrashed(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrderDeleteAt, error) {
	s.logger.Info("FindByTrashed orders called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.orderQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed orders success")

	return &pborder.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order",
		Data:       mapResponsesOrderTrashed(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderCommandHandleGrpc) Create(ctx context.Context, request *pborder.CreateOrderRequest) (*pborder.ApiResponseOrder, error) {
	s.logger.Info("Create order called", zap.Int32("merchantId", request.GetMerchantId()))

	req := &requests.CreateOrderRequest{
		MerchantID: int(request.GetMerchantId()),
		CashierID:  int(request.GetCashierId()),
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.CreateOrderItemRequest{
			ProductID: int(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
		})
	}

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateCreateOrder
	}

	order, err := s.orderCommand.CreateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Create order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create order success")

	return &pborder.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully created order",
		Data:    mapResponseOrder(order),
	}, nil
}

func (s *orderCommandHandleGrpc) Update(ctx context.Context, request *pborder.UpdateOrderRequest) (*pborder.ApiResponseOrder, error) {
	s.logger.Info("Update order called", zap.Int32("id", request.GetOrderId()))

	id := int(request.GetOrderId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateOrderRequest{
		OrderID: &id,
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.UpdateOrderItemRequest{
			OrderItemID: int(item.GetOrderItemId()),
			ProductID:   int(item.GetProductId()),
			Quantity:    int(item.GetQuantity()),
		})
	}

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateUpdateOrder
	}

	order, err := s.orderCommand.UpdateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Update order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update order success")

	return &pborder.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully updated order",
		Data:    mapResponseOrder(order),
	}, nil
}

func (s *orderCommandHandleGrpc) TrashedOrder(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDeleteAt, error) {
	s.logger.Info("TrashedOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderCommand.TrashedOrder(ctx, id)
	if err != nil {
		s.logger.Error("TrashedOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedOrder success")

	return &pborder.ApiResponseOrderDeleteAt{
		Status:  "success",
		Message: "Successfully trashed order",
		Data:    mapResponseOrderDeleteAt(merchant),
	}, nil
}

func (s *orderCommandHandleGrpc) RestoreOrder(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDeleteAt, error) {
	s.logger.Info("RestoreOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderCommand.RestoreOrder(ctx, id)
	if err != nil {
		s.logger.Error("RestoreOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreOrder success")

	return &pborder.ApiResponseOrderDeleteAt{
		Status:  "success",
		Message: "Successfully restored order",
		Data:    mapResponseOrderDeleteAt(merchant),
	}, nil
}

func (s *orderCommandHandleGrpc) DeleteOrderPermanent(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDelete, error) {
	s.logger.Info("DeleteOrderPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.orderCommand.DeleteOrderPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteOrderPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteOrderPermanent success")

	return &pborder.ApiResponseOrderDelete{
		Status:  "success",
		Message: "Successfully deleted order permanently",
	}, nil
}

func (s *orderCommandHandleGrpc) RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pborder.ApiResponseOrderAll, error) {
	s.logger.Info("RestoreAllOrder called")

	_, err := s.orderCommand.RestoreAllOrder(ctx)
	if err != nil {
		s.logger.Error("RestoreAllOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllOrder success")

	return &pborder.ApiResponseOrderAll{
		Status:  "success",
		Message: "Successfully restore all order",
	}, nil
}

func (s *orderCommandHandleGrpc) DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pborder.ApiResponseOrderAll, error) {
	s.logger.Info("DeleteAllOrderPermanent called")

	_, err := s.orderCommand.DeleteAllOrderPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllOrderPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllOrderPermanent success")

	return &pborder.ApiResponseOrderAll{
		Status:  "success",
		Message: "Successfully delete order permanen",
	}, nil
}
