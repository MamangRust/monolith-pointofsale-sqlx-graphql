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

type orderStatsByMerchantDeps struct {
	Cache                          mencache.OrderStatsByMerchantCache
	OrderStatsByMerchantRepository repository.OrderStatByMerchantRepository
	Logger                         logger.LoggerInterface
	Observability                  observability.TraceLoggerObservability
}

type orderStatsByMerchantService struct {
	mencache                       mencache.OrderStatsByMerchantCache
	orderStatsByMerchantRepository repository.OrderStatByMerchantRepository
	logger                         logger.LoggerInterface
	observability                  observability.TraceLoggerObservability
}

func NewOrderStatsByMerchantService(params *orderStatsByMerchantDeps) OrderStatByMerchantService {
	return &orderStatsByMerchantService{
		mencache:                       params.Cache,
		orderStatsByMerchantRepository: params.OrderStatsByMerchantRepository,
		logger:                         params.Logger,
		observability:                  params.Observability,
	}
}

func (s *orderStatsByMerchantService) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error) {
	const method = "FindMonthlyTotalRevenueByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalRevenueByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total revenue from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetMonthlyTotalRevenueByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTotalRevenueByMerchantRow](
			s.logger,
			order_errors.ErrFailedFindMonthlyTotalRevenueByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyTotalRevenueByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched monthly total revenue", zap.Int("year", req.Year), zap.Int("month", req.Month))
	return res, nil
}

func (s *orderStatsByMerchantService) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, error) {
	const method = "FindYearlyTotalRevenueByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("year", req.Year),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTotalRevenueByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly total revenue from cache", zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetYearlyTotalRevenueByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTotalRevenueByMerchantRow](
			s.logger,
			order_errors.ErrFailedFindYearlyTotalRevenueByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyTotalRevenueByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched yearly total revenue", zap.Int("year", req.Year))
	return res, nil
}

func (s *orderStatsByMerchantService) FindMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, error) {
	const method = "FindMonthlyOrderByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("year", req.Year),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyOrderByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly orders from cache", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetMonthlyOrderByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyOrderByMerchantRow](
			s.logger,
			order_errors.ErrFailedFindMonthlyOrderByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyOrderByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched monthly orders", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
	return res, nil
}

func (s *orderStatsByMerchantService) FindYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, error) {
	const method = "FindYearlyOrderByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("year", req.Year),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyOrderByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly orders from cache", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetYearlyOrderByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyOrderByMerchantRow](
			s.logger,
			order_errors.ErrFailedFindYearlyOrderByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyOrderByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched yearly orders", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
	return res, nil
}
