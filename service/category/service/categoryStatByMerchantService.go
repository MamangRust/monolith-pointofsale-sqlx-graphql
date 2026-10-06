package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-category/cache"
	"github.com/MamangRust/monolith-pointofsale-category/repository"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type categoryStatsByMerchantDeps struct {
	Cache                             mencache.CategoryStatsByMerchantCache
	CategoryStatsByMerchantRepository repository.CategoryStatsByMerchantRepository
	Logger                            logger.LoggerInterface
	Observability                     observability.TraceLoggerObservability
}

type categoryStatsByMerchantService struct {
	mencache                          mencache.CategoryStatsByMerchantCache
	categoryStatsByMerchantRepository repository.CategoryStatsByMerchantRepository
	logger                            logger.LoggerInterface
	observability                     observability.TraceLoggerObservability
}

func NewCategoryStatsByMerchantService(params *categoryStatsByMerchantDeps) CategoryStatsByMerchantService {
	return &categoryStatsByMerchantService{
		mencache:                          params.Cache,
		categoryStatsByMerchantRepository: params.CategoryStatsByMerchantRepository,
		logger:                            params.Logger,
		observability:                     params.Observability,
	}
}

func (s *categoryStatsByMerchantService) FindMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*db.GetMonthlyTotalPriceByMerchantRow, error) {
	const method = "FindMonthlyTotalPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
		attribute.Int("merchant.id", req.MerchantID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthTotalPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total price by merchant from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetMonthlyTotalPriceByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTotalPriceByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthTotalPriceByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched monthly total price by merchant", zap.Int("year", req.Year), zap.Int("month", req.Month))
	return res, nil
}

func (s *categoryStatsByMerchantService) FindYearlyTotalPriceByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*db.GetYearlyTotalPriceByMerchantRow, error) {
	const method = "FindYearlyTotalPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant.id", req.MerchantID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearTotalPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly total price by merchant from cache", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetYearlyTotalPricesByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTotalPriceByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearTotalPriceByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched yearly total price by merchant", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
	return res, nil
}

func (s *categoryStatsByMerchantService) FindMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*db.GetMonthlyCategoryByMerchantRow, error) {
	const method = "FindMonthPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant.id", req.MerchantID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly category prices by merchant from cache", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetMonthPriceByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyCategoryByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthPriceByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched monthly category prices by merchant", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
	return res, nil
}

func (s *categoryStatsByMerchantService) FindYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*db.GetYearlyCategoryByMerchantRow, error) {
	const method = "FindYearPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant.id", req.MerchantID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly category prices by merchant from cache", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetYearPriceByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyCategoryByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearPriceByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched yearly category prices by merchant", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
	return res, nil
}
