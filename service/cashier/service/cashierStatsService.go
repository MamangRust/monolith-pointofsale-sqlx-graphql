package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-cashier/cache"
	"github.com/MamangRust/monolith-pointofsale-cashier/repository"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type cashierStatsDeps struct {
	Cache         mencache.CashierStatsCache
	CashierStats  repository.CashierStatsRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type cashierStatsService struct {
	mencache      mencache.CashierStatsCache
	cashierStats  repository.CashierStatsRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewCashierStatsService(params *cashierStatsDeps) CashierStatsService {
	return &cashierStatsService{
		mencache:      params.Cache,
		cashierStats:  params.CashierStats,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *cashierStatsService) FindMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*db.GetMonthlyTotalSalesCashierRow, error) {
	const method = "FindMonthlyTotalSales"
	month := req.Month
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("month", month), attribute.Int("year", year))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalSalesCache(ctx, req); found {
		logSuccess("Fetched monthly total sales from cache", zap.Int("month", month), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyTotalSales(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTotalSalesCashierRow](
			s.logger,
			cashier_errors.ErrFailedFindMonthlyTotalSales,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyTotalSalesCache(ctx, req, res)
	logSuccess("Fetched monthly total sales from DB", zap.Int("month", month), zap.Int("year", year))
	return res, nil
}

func (s *cashierStatsService) FindYearlyTotalSales(ctx context.Context, year int) ([]*db.GetYearlyTotalSalesCashierRow, error) {
	const method = "FindYearlyTotalSales"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTotalSalesCache(ctx, year); found {
		logSuccess("Fetched yearly total sales from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cashierStats.GetYearlyTotalSales(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTotalSalesCashierRow](
			s.logger,
			cashier_errors.ErrFailedFindYearlyTotalSales,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyTotalSalesCache(ctx, year, res)
	logSuccess("Fetched yearly total sales from DB", zap.Int("year", year))
	return res, nil
}

func (s *cashierStatsService) FindMonthlySales(ctx context.Context, year int) ([]*db.GetMonthlyCashierRow, error) {
	const method = "FindMonthlySales"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlySalesCache(ctx, year); found {
		logSuccess("Fetched monthly sales from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cashierStats.GetMonthyCashier(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyCashierRow](
			s.logger,
			cashier_errors.ErrFailedFindMonthlySales,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlySalesCache(ctx, year, res)
	logSuccess("Fetched monthly sales from DB", zap.Int("year", year))
	return res, nil
}

func (s *cashierStatsService) FindYearlySales(ctx context.Context, year int) ([]*db.GetYearlyCashierRow, error) {
	const method = "FindYearlySales"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlySalesCache(ctx, year); found {
		logSuccess("Fetched yearly sales from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cashierStats.GetYearlyCashier(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyCashierRow](
			s.logger,
			cashier_errors.ErrFailedFindYearlySales,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlySalesCache(ctx, year, res)
	logSuccess("Fetched yearly sales from DB", zap.Int("year", year))
	return res, nil
}
