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

type cashierStatsByIdDeps struct {
	Cache         mencache.CashierStatsByIdCache
	CashierStats  repository.CashierStatByIdRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type cashierStatsByIdService struct {
	mencache      mencache.CashierStatsByIdCache
	cashierStats  repository.CashierStatByIdRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewCashierStatsByIdService(params *cashierStatsByIdDeps) CashierStatsByIdService {
	return &cashierStatsByIdService{
		mencache:      params.Cache,
		cashierStats:  params.CashierStats,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *cashierStatsByIdService) FindMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*db.GetMonthlyTotalSalesByIdRow, error) {
	const method = "FindMonthlyTotalSalesById"
	month := req.Month
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalSalesByIdCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total sales by ID from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyTotalSalesById(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTotalSalesByIdRow](
			s.logger,
			cashier_errors.ErrFailedFindMonthlyTotalSalesById,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyTotalSalesByIdCache(ctx, req, res)
	logSuccess("Successfully fetched monthly total sales by ID", zap.Int("year", year), zap.Int("month", month))
	return res, nil
}

func (s *cashierStatsByIdService) FindYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*db.GetYearlyTotalSalesByIdRow, error) {
	const method = "FindYearlyTotalSalesById"
	year := req.Year
	cashier_id := req.CashierID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("cashier_id", cashier_id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTotalSalesByIdCache(ctx, req); found {
		logSuccess("Successfully fetched yearly total sales by ID from cache", zap.Int("year", year), zap.Int("cashier_id", cashier_id))
		return data, nil
	}

	res, err := s.cashierStats.GetYearlyTotalSalesById(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTotalSalesByIdRow](
			s.logger,
			cashier_errors.ErrFailedFindYearlyTotalSalesById,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyTotalSalesByIdCache(ctx, req, res)
	logSuccess("Successfully fetched yearly total sales by ID", zap.Int("year", year), zap.Int("cashier_id", cashier_id))
	return res, nil
}

func (s *cashierStatsByIdService) FindMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*db.GetMonthlyCashierByCashierIdRow, error) {
	const method = "FindMonthlyCashierById"
	year := req.Year
	cashier_id := req.CashierID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("cashier.id", cashier_id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyCashierByIdCache(ctx, req); found {
		logSuccess("Successfully fetched monthly cashier sales by ID from cache", zap.Int("year", year), zap.Int("cashier_id", cashier_id))
		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyCashierById(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyCashierByCashierIdRow](
			s.logger,
			cashier_errors.ErrFailedFindMonthlyCashierById,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyCashierByIdCache(ctx, req, res)
	logSuccess("Successfully fetched monthly cashier sales by ID", zap.Int("year", year), zap.Int("cashier_id", cashier_id))
	return res, nil
}

func (s *cashierStatsByIdService) FindYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*db.GetYearlyCashierByCashierIdRow, error) {
	const method = "FindYearlyCashierById"
	year := req.Year
	cashier_id := req.CashierID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("cashier.id", cashier_id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyCashierByIdCache(ctx, req); found {
		logSuccess("Successfully fetched yearly cashier sales by ID from cache", zap.Int("year", year), zap.Int("cashier_id", cashier_id))
		return data, nil
	}

	res, err := s.cashierStats.GetYearlyCashierById(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyCashierByCashierIdRow](
			s.logger,
			cashier_errors.ErrFailedFindYearlyCashierById,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyCashierByIdCache(ctx, req, res)
	logSuccess("Successfully fetched yearly cashier sales by ID", zap.Int("year", year), zap.Int("cashier_id", cashier_id))
	return res, nil
}
