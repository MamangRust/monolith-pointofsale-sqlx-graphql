package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-order-item/cache"
	"github.com/MamangRust/monolith-pointofsale-order-item/repository"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	orderitem_errors "github.com/MamangRust/monolith-pointofsale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderItemCommandDeps struct {
	Cache         mencache.OrderItemCommandCache
	Repo          repository.OrderItemCommandRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type orderItemCommandService struct {
	mencache      mencache.OrderItemCommandCache
	repo          repository.OrderItemCommandRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewOrderItemCommandService(params *orderItemCommandDeps) OrderItemCommandService {
	return &orderItemCommandService{
		mencache:      params.Cache,
		repo:          params.Repo,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *orderItemCommandService) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.OrderItem, error) {
	const method = "CreateOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order.id", req.OrderID),
		attribute.Int("product.id", req.ProductID),
		attribute.Int("quantity", req.Quantity),
	)
	defer func() {
		end(status)
	}()

	item, err := s.repo.CreateOrderItem(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.OrderItem](
			s.logger,
			orderitem_errors.ErrFailedCreateOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedOrderItems(ctx, req.OrderID)
	logSuccess("Successfully created order item", zap.Int32("order_item.id", item.OrderItemID))

	return item, nil
}

func (s *orderItemCommandService) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.OrderItem, error) {
	const method = "UpdateOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_item.id", req.OrderItemID),
		attribute.Int("order.id", req.OrderID),
	)
	defer func() {
		end(status)
	}()

	if req.OrderItemID <= 0 {
		status = "error"
		return sharederrorhandler.HandleError[*db.OrderItem](
			s.logger,
			orderitem_errors.ErrFailedOrderItemEmptyId,
			method,
			span,
		)
	}

	item, err := s.repo.UpdateOrderItem(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.OrderItem](
			s.logger,
			orderitem_errors.ErrFailedUpdateOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedOrderItems(ctx, req.OrderID)
	logSuccess("Successfully updated order item", zap.Int32("order_item.id", item.OrderItemID))

	return item, nil
}

func (s *orderItemCommandService) DeleteOrderItem(ctx context.Context, orderItemID int) error {
	const method = "DeleteOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_item.id", orderItemID),
	)
	defer func() {
		end(status)
	}()

	if orderItemID <= 0 {
		status = "error"
		_, err := sharederrorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedOrderItemEmptyId,
			method,
			span,
		)
		return err
	}

	if err := s.repo.DeleteOrderItem(ctx, orderItemID); err != nil {
		status = "error"
		_, handledErr := sharederrorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedDeleteOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
		return handledErr
	}

	s.mencache.DeleteCachedOrderItemsAll(ctx)
	logSuccess("Successfully deleted order item", zap.Int("order_item.id", orderItemID))

	return nil
}
