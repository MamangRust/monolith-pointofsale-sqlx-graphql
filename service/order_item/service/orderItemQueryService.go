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

type orderItemQueryDeps struct {
	Cache         mencache.OrderItemQueryCache
	Repo          repository.OrderItemQueryRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type orderItemQueryService struct {
	mencache      mencache.OrderItemQueryCache
	repo          repository.OrderItemQueryRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewOrderItemQueryService(params *orderItemQueryDeps) OrderItemQueryService {
	return &orderItemQueryService{
		mencache:      params.Cache,
		repo:          params.Repo,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *orderItemQueryService) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsRow, *int, error) {
	const method = "FindAllOrderItems"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedOrderItemsAll(ctx, req); found {
		logSuccess("Successfully fetched order items from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, totalRecords, err := s.repo.FindAllOrderItems(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetOrderItemsRow](
			s.logger,
			orderitem_errors.ErrFailedFindAllOrderItems.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedOrderItemsAll(ctx, req, orderItems, totalRecords)
	logSuccess("Successfully fetched order items", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return orderItems, totalRecords, nil
}

func (s *orderItemQueryService) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsActiveRow, *int, error) {
	const method = "FindByActive"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedOrderItemActive(ctx, req); found {
		logSuccess("Successfully fetched active order items from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, totalRecords, err := s.repo.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetOrderItemsActiveRow](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemsByActive.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedOrderItemActive(ctx, req, orderItems, totalRecords)
	logSuccess("Successfully fetched active order items", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return orderItems, totalRecords, nil
}

func (s *orderItemQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsTrashedRow, *int, error) {
	const method = "FindByTrashed"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedOrderItemTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed order items from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, totalRecords, err := s.repo.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetOrderItemsTrashedRow](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemsByTrashed.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedOrderItemTrashed(ctx, req, orderItems, totalRecords)
	logSuccess("Successfully fetched trashed order items", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return orderItems, totalRecords, nil
}

func (s *orderItemQueryService) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*db.OrderItem, error) {
	const method = "FindOrderItemByOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", orderID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedOrderItems(ctx, orderID); found {
		logSuccess("Successfully fetched order items by order ID from cache", zap.Int("orderID", orderID))
		return data, nil
	}

	orderItems, err := s.repo.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.OrderItem](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedOrderItems(ctx, orderItems)
	logSuccess("Successfully fetched order items by order ID", zap.Int("orderID", orderID))
	return orderItems, nil
}

func (s *orderItemQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
