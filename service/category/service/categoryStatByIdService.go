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

type categoryStatsByIdDeps struct {
	Cache                       mencache.CategoryStatsByIdCache
	CategoryStatsByIdRepository repository.CategoryStatsByIdRepository
	Logger                      logger.LoggerInterface
	Observability               observability.TraceLoggerObservability
}

type categoryStatsByIdService struct {
	mencache                    mencache.CategoryStatsByIdCache
	categoryStatsByIdRepository repository.CategoryStatsByIdRepository
	logger                      logger.LoggerInterface
	observability               observability.TraceLoggerObservability
}

func NewCategoryStatsByIdService(params *categoryStatsByIdDeps) CategoryStatsByIdService {
	return &categoryStatsByIdService{
		mencache:                    params.Cache,
		categoryStatsByIdRepository: params.CategoryStatsByIdRepository,
		logger:                      params.Logger,
		observability:               params.Observability,
	}
}

func (s *categoryStatsByIdService) FindMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*db.GetMonthlyTotalPriceByIdRow, error) {
	const method = "FindMonthlyTotalPriceById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
		attribute.Int("category.id", req.CategoryID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthTotalPriceByIdCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total price by ID from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetMonthlyTotalPriceById(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTotalPriceByIdRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthTotalPriceByIdCache(ctx, req, res)
	logSuccess("Successfully fetched monthly total price by ID", zap.Int("year", req.Year), zap.Int("month", req.Month))
	return res, nil
}

func (s *categoryStatsByIdService) FindYearlyTotalPriceById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*db.GetYearlyTotalPriceByIdRow, error) {
	const method = "FindYearlyTotalPriceById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("category.id", req.CategoryID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearTotalPriceByIdCache(ctx, req); found {
		logSuccess("Successfully fetched yearly total price by ID from cache", zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetYearlyTotalPricesById(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTotalPriceByIdRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearTotalPriceByIdCache(ctx, req, res)
	logSuccess("Successfully fetched yearly total price by ID", zap.Int("year", req.Year))
	return res, nil
}

func (s *categoryStatsByIdService) FindMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*db.GetMonthlyCategoryByIdRow, error) {
	const method = "FindMonthPriceById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("category.id", req.CategoryID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthPriceByIdCache(ctx, req); found {
		logSuccess("Successfully fetched monthly category prices by ID from cache", zap.Int("year", req.Year), zap.Int("category.id", req.CategoryID))
		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetMonthPriceById(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyCategoryByIdRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthPriceByIdCache(ctx, req, res)
	logSuccess("Successfully fetched monthly category prices by ID", zap.Int("year", req.Year), zap.Int("category.id", req.CategoryID))
	return res, nil
}

func (s *categoryStatsByIdService) FindYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*db.GetYearlyCategoryByIdRow, error) {
	const method = "FindYearPriceById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("category.id", req.CategoryID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearPriceByIdCache(ctx, req); found {
		logSuccess("Successfully fetched yearly category prices by ID from cache", zap.Int("year", req.Year), zap.Int("category.id", req.CategoryID))
		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetYearPriceById(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyCategoryByIdRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearPriceByIdCache(ctx, req, res)
	logSuccess("Successfully fetched yearly category prices by ID", zap.Int("year", req.Year), zap.Int("category.id", req.CategoryID))
	return res, nil
}
