package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-order/cache"
	"github.com/MamangRust/monolith-pointofsale-order/repository"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/order_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderQueryDeps struct {
	Cache                mencache.OrderQueryCache
	OrderQueryRepository repository.OrderQueryRepository
	Logger               logger.LoggerInterface
	Observability        observability.TraceLoggerObservability
}

type orderQueryService struct {
	mencache             mencache.OrderQueryCache
	orderQueryRepository repository.OrderQueryRepository
	logger               logger.LoggerInterface
	observability        observability.TraceLoggerObservability
}

func NewOrderQueryService(params *orderQueryDeps) OrderQueryService {
	return &orderQueryService{
		mencache:             params.Cache,
		orderQueryRepository: params.OrderQueryRepository,
		logger:               params.Logger,
		observability:        params.Observability,
	}
}

func (s *orderQueryService) FindAll(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersRow, *int, error) {
	const method = "FindAll"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetOrderAllCache(ctx, req); found {
		logSuccess("Successfully fetched order from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search))
		return data, total, nil
	}

	orders, totalRecords, err := s.orderQueryRepository.FindAllOrders(ctx, req)
	if err != nil {
		status = "error"
		resErr, err := sharederrorhandler.HandleError[[]*db.GetOrdersRow](
			s.logger,
			order_errors.ErrFailedFindAllOrders.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
		return resErr, nil, err
	}

	s.mencache.SetOrderAllCache(ctx, req, orders, totalRecords)
	logSuccess("Successfully fetched order", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search))
	return orders, totalRecords, nil
}

func (s *orderQueryService) FindById(ctx context.Context, orderID int) (*db.Order, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order.id", orderID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedOrderCache(ctx, orderID); found {
		logSuccess("Successfully fetched order from cache", zap.Int("order.id", orderID))
		return data, nil
	}

	order, err := s.orderQueryRepository.FindById(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedFindOrderById.WithInternal(err),
			method,
			span,
			zap.Int("order_id", orderID),
			zap.Error(err),
		)
	}

	s.mencache.SetCachedOrderCache(ctx, order)
	logSuccess("Successfully fetched order", zap.Int("order.id", orderID))
	return order, nil
}

func (s *orderQueryService) FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersActiveRow, *int, error) {
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

	if data, total, found := s.mencache.GetOrderActiveCache(ctx, req); found {
		logSuccess("Successfully fetched active order from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search))
		return data, total, nil
	}

	orders, totalRecords, err := s.orderQueryRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		resErr, err := sharederrorhandler.HandleError[[]*db.GetOrdersActiveRow](
			s.logger,
			order_errors.ErrFailedFindOrdersByActive.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
		return resErr, nil, err
	}

	s.mencache.SetOrderActiveCache(ctx, req, orders, totalRecords)
	logSuccess("Successfully fetched active order", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search))
	return orders, totalRecords, nil
}

func (s *orderQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersTrashedRow, *int, error) {
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

	if data, total, found := s.mencache.GetOrderTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched trashed order from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search))
		return data, total, nil
	}

	orders, totalRecords, err := s.orderQueryRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		resErr, err := sharederrorhandler.HandleError[[]*db.GetOrdersTrashedRow](
			s.logger,
			order_errors.ErrFailedFindOrdersByTrashed.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
		return resErr, nil, err
	}

	s.mencache.SetOrderTrashedCache(ctx, req, orders, totalRecords)
	logSuccess("Successfully fetched trashed order", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search))
	return orders, totalRecords, nil
}

func (s *orderQueryService) FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*db.GetOrdersByMerchantRow, *int, error) {
	const method = "FindByMerchant"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
		attribute.Int("merchant.id", req.MerchantID),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedOrderMerchant(ctx, req); found {
		logSuccess("Successfully fetched order from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search), zap.Int("merchant.id", req.MerchantID))
		return data, total, nil
	}

	orders, totalRecords, err := s.orderQueryRepository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		resErr, err := sharederrorhandler.HandleError[[]*db.GetOrdersByMerchantRow](
			s.logger,
			order_errors.ErrFailedFindOrdersByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
		return resErr, nil, err
	}

	s.mencache.SetCachedOrderMerchant(ctx, req, orders, totalRecords)
	logSuccess("Successfully fetched order", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search), zap.Int("merchant.id", req.MerchantID))
	return orders, totalRecords, nil
}

func (s *orderQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
