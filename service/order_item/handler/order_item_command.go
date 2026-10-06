package handler

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-order-item/service"
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	orderitem_errors "github.com/MamangRust/monolith-pointofsale-shared/errors/order_item_errors"
	"go.uber.org/zap"
)

type orderItemCommandHandleGrpc struct {
	pb.UnimplementedOrderItemCommandServiceServer
	orderItemCommandService service.OrderItemCommandService
	logger                  logger.LoggerInterface
}

func NewOrderItemCommandHandlerGrpc(
	orderItemCommandService service.OrderItemCommandService,
	logger logger.LoggerInterface,
) OrderItemCommandHandlerGrpc {
	return &orderItemCommandHandleGrpc{
		orderItemCommandService: orderItemCommandService,
		logger:                  logger,
	}
}

func (s *orderItemCommandHandleGrpc) CreateOrderItem(ctx context.Context, request *pb.CreateOrderItemRecordRequest) (*pb.ApiResponseOrderItem, error) {
	s.logger.Info("CreateOrderItem called", zap.Int32("order_id", request.GetOrderId()))

	if request.GetOrderId() <= 0 || request.GetProductId() <= 0 ||
		request.GetQuantity() <= 0 || request.GetPrice() <= 0 {
		return nil, orderitem_errors.ErrGrpcValidateCreateOrderItem
	}

	req := &requests.CreateOrderItemRecordRequest{
		OrderID:   int(request.GetOrderId()),
		ProductID: int(request.GetProductId()),
		Quantity:  int(request.GetQuantity()),
		Price:     int(request.GetPrice()),
	}

	item, err := s.orderItemCommandService.CreateOrderItem(ctx, req)
	if err != nil {
		s.logger.Error("CreateOrderItem failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("CreateOrderItem success")

	return &pb.ApiResponseOrderItem{
		Status:  "success",
		Message: "Successfully created order item",
		Data:    mapOrderItemToProto(item),
	}, nil
}

func (s *orderItemCommandHandleGrpc) UpdateOrderItem(ctx context.Context, request *pb.UpdateOrderItemRecordRequest) (*pb.ApiResponseOrderItem, error) {
	s.logger.Info("UpdateOrderItem called", zap.Int32("order_item_id", request.GetOrderItemId()))

	if request.GetOrderItemId() <= 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateOrderItemRecordRequest{
		OrderItemID: int(request.GetOrderItemId()),
		OrderID:     int(request.GetOrderId()),
		ProductID:   int(request.GetProductId()),
		Quantity:    int(request.GetQuantity()),
		Price:       int(request.GetPrice()),
	}

	item, err := s.orderItemCommandService.UpdateOrderItem(ctx, req)
	if err != nil {
		s.logger.Error("UpdateOrderItem failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateOrderItem success")

	return &pb.ApiResponseOrderItem{
		Status:  "success",
		Message: "Successfully updated order item",
		Data:    mapOrderItemToProto(item),
	}, nil
}

func (s *orderItemCommandHandleGrpc) DeleteOrderItem(ctx context.Context, request *pb.DeleteOrderItemRecordRequest) (*pb.ApiResponseOrderItemDelete, error) {
	s.logger.Info("DeleteOrderItem called", zap.Int32("order_item_id", request.GetOrderItemId()))

	if request.GetOrderItemId() <= 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	if err := s.orderItemCommandService.DeleteOrderItem(ctx, int(request.GetOrderItemId())); err != nil {
		s.logger.Error("DeleteOrderItem failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteOrderItem success")

	return &pb.ApiResponseOrderItemDelete{
		Status:  "success",
		Message: "Successfully deleted order item",
	}, nil
}
