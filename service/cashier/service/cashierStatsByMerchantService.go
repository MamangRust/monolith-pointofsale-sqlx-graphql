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

type cashierStatsByMerchantDeps struct {
	Cache         mencache.CashierStatsByMerchantCache
	CashierStats  repository.CashierStatByMerchantRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type cashierStatsByMerchantService struct {
	mencache      mencache.CashierStatsByMerchantCache
	cashierStats  repository.CashierStatByMerchantRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewCashierStatsByMerchantService(params *cashierStatsByMerchantDeps) CashierStatsByMerchant {
	return &cashierStatsByMerchantService{
		mencache:      params.Cache,
		cashierStats:  params.CashierStats,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *cashierStatsByMerchantService) FindMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*db.GetMonthlyTotalSalesByMerchantRow, error) {
	const method = "FindMonthlyTotalSalesByMerchant"
	month := req.Month
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalSalesByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total sales by merchant from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyTotalSalesByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTotalSalesByMerchantRow](
			s.logger,
			cashier_errors.ErrFailedFindMonthlyTotalSalesByMerchant,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyTotalSalesByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched monthly total sales by merchant", zap.Int("year", year), zap.Int("month", month))
	return res, nil
}

func (s *cashierStatsByMerchantService) FindYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*db.GetYearlyTotalSalesByMerchantRow, error) {
	const method = "FindYearlyTotalSalesByMerchant"
	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("merchant_id", merchant_id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTotalSalesByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly total sales by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
		return data, nil
	}

	res, err := s.cashierStats.GetYearlyTotalSalesByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTotalSalesByMerchantRow](
			s.logger,
			cashier_errors.ErrFailedFindYearlyTotalSalesByMerchant,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyTotalSalesByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched yearly total sales by merchant", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
	return res, nil
}

func (s *cashierStatsByMerchantService) FindMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*db.GetMonthlyCashierByMerchantRow, error) {
	const method = "FindMonthlyCashierByMerchant"
	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("merchant.id", merchant_id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyCashierByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly cashier sales by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyCashierByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyCashierByMerchantRow](
			s.logger,
			cashier_errors.ErrFailedFindMonthlyCashierByMerchant,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyCashierByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched monthly cashier sales by merchant", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
	return res, nil
}

func (s *cashierStatsByMerchantService) FindYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*db.GetYearlyCashierByMerchantRow, error) {
	const method = "FindYearlyCashierByMerchant"
	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("merchant.id", merchant_id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyCashierByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly cashier sales by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
		return data, nil
	}

	res, err := s.cashierStats.GetYearlyCashierByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyCashierByMerchantRow](
			s.logger,
			cashier_errors.ErrFailedFindYearlyCashierByMerchant,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyCashierByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched yearly cashier sales by merchant", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
	return res, nil
}
