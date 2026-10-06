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

type categoryStatsDeps struct {
	Cache                   mencache.CategoryStatsCache
	CategoryStatsRepository repository.CategoryStatsRepository
	Logger                  logger.LoggerInterface
	Observability           observability.TraceLoggerObservability
}

type categoryStatsService struct {
	mencache                mencache.CategoryStatsCache
	categoryStatsRepository repository.CategoryStatsRepository
	logger                  logger.LoggerInterface
	observability           observability.TraceLoggerObservability
}

func NewCategoryStatsService(params *categoryStatsDeps) CategoryStatsService {
	return &categoryStatsService{
		mencache:                params.Cache,
		categoryStatsRepository: params.CategoryStatsRepository,
		logger:                  params.Logger,
		observability:           params.Observability,
	}
}

func (s *categoryStatsService) FindMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*db.GetMonthlyTotalPriceRow, error) {
	const method = "FindMonthlyTotalPrice"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthTotalPriceCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total price from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.categoryStatsRepository.GetMonthlyTotalPrice(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTotalPriceRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthTotalPriceCache(ctx, req, res)
	logSuccess("Successfully fetched monthly total price", zap.Int("year", req.Year), zap.Int("month", req.Month))
	return res, nil
}

func (s *categoryStatsService) FindYearlyTotalPrice(ctx context.Context, year int) ([]*db.GetYearlyTotalPriceRow, error) {
	const method = "FindYearlyTotalPrice"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearTotalPriceCache(ctx, year); found {
		logSuccess("Successfully fetched yearly total price from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.categoryStatsRepository.GetYearlyTotalPrices(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTotalPriceRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearTotalPriceCache(ctx, year, res)
	logSuccess("Successfully fetched yearly total price", zap.Int("year", year))
	return res, nil
}

func (s *categoryStatsService) FindMonthPrice(ctx context.Context, year int) ([]*db.GetMonthlyCategoryRow, error) {
	const method = "FindMonthPrice"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthPriceCache(ctx, year); found {
		logSuccess("Successfully fetched monthly category prices from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.categoryStatsRepository.GetMonthPrice(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyCategoryRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthPriceCache(ctx, year, res)
	logSuccess("Successfully fetched monthly category prices", zap.Int("year", year))
	return res, nil
}

func (s *categoryStatsService) FindYearPrice(ctx context.Context, year int) ([]*db.GetYearlyCategoryRow, error) {
	const method = "FindYearPrice"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearPriceCache(ctx, year); found {
		logSuccess("Successfully fetched yearly category prices from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.categoryStatsRepository.GetYearPrice(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyCategoryRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearPriceCache(ctx, year, res)
	logSuccess("Successfully fetched yearly category prices", zap.Int("year", year))
	return res, nil
}
