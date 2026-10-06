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

type orderStatsDeps struct {
	Cache                mencache.OrderStatsCache
	OrderStatsRepository repository.OrderStatsRepository
	Logger               logger.LoggerInterface
	Observability        observability.TraceLoggerObservability
}

type orderStatsService struct {
	mencache             mencache.OrderStatsCache
	orderStatsRepository repository.OrderStatsRepository
	logger               logger.LoggerInterface
	observability        observability.TraceLoggerObservability
}

func NewOrderStatsService(params *orderStatsDeps) OrderStatsService {
	return &orderStatsService{
		mencache:             params.Cache,
		orderStatsRepository: params.OrderStatsRepository,
		logger:               params.Logger,
		observability:        params.Observability,
	}
}

func (s *orderStatsService) FindMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*db.GetMonthlyTotalRevenueRow, error) {
	const method = "FindMonthlyTotalRevenue"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalRevenueCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total revenue from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.orderStatsRepository.GetMonthlyTotalRevenue(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTotalRevenueRow](
			s.logger,
			order_errors.ErrFailedFindMonthlyTotalRevenue.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyTotalRevenueCache(ctx, req, res)
	logSuccess("Successfully fetched monthly total revenue", zap.Int("year", req.Year), zap.Int("month", req.Month))
	return res, nil
}

func (s *orderStatsService) FindYearlyTotalRevenue(ctx context.Context, year int) ([]*db.GetYearlyTotalRevenueRow, error) {
	const method = "FindYearlyTotalRevenue"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTotalRevenueCache(ctx, year); found {
		logSuccess("Successfully fetched yearly total revenue from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.orderStatsRepository.GetYearlyTotalRevenue(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTotalRevenueRow](
			s.logger,
			order_errors.ErrFailedFindYearlyTotalRevenue.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyTotalRevenueCache(ctx, year, res)
	logSuccess("Successfully fetched yearly total revenue", zap.Int("year", year))
	return res, nil
}

func (s *orderStatsService) FindMonthlyOrder(ctx context.Context, year int) ([]*db.GetMonthlyOrderRow, error) {
	const method = "FindMonthlyOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyOrderCache(ctx, year); found {
		logSuccess("Successfully fetched monthly orders from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.orderStatsRepository.GetMonthlyOrder(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyOrderRow](
			s.logger,
			order_errors.ErrFailedFindMonthlyOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyOrderCache(ctx, year, res)
	logSuccess("Successfully fetched monthly orders", zap.Int("year", year))
	return res, nil
}

func (s *orderStatsService) FindYearlyOrder(ctx context.Context, year int) ([]*db.GetYearlyOrderRow, error) {
	const method = "FindYearlyOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyOrderCache(ctx, year); found {
		logSuccess("Successfully fetched yearly orders from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.orderStatsRepository.GetYearlyOrder(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyOrderRow](
			s.logger,
			order_errors.ErrFailedFindYearlyOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyOrderCache(ctx, year, res)
	logSuccess("Successfully fetched yearly orders", zap.Int("year", year))
	return res, nil
}
