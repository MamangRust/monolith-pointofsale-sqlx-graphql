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

type categoryQueryDeps struct {
	Cache         mencache.CategoryQueryCache
	CategoryQuery repository.CategoryQueryRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type categoryQueryService struct {
	mencache      mencache.CategoryQueryCache
	categoryQuery repository.CategoryQueryRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewCategoryQueryService(params *categoryQueryDeps) CategoryQueryService {
	return &categoryQueryService{
		mencache:      params.Cache,
		categoryQuery: params.CategoryQuery,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *categoryQueryService) FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, *int, error) {
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

	if data, total, found := s.mencache.GetCachedCategoriesCache(ctx, req); found {
		logSuccess("Successfully fetched categories from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	category, totalRecords, err := s.categoryQueryRepository().FindAllCategory(ctx, req)
	if err != nil {
		status = "error"
		_, mappedErr := sharederrorhandler.HandleError[[]*db.GetCategoriesRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, mappedErr
	}

	s.mencache.SetCachedCategoriesCache(ctx, req, category, totalRecords)
	logSuccess("Successfully fetched categories", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return category, totalRecords, nil
}

func (s *categoryQueryService) FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, *int, error) {
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

	if data, total, found := s.mencache.GetCachedCategoryActiveCache(ctx, req); found {
		logSuccess("Successfully fetched active categories from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	category, totalRecords, err := s.categoryQueryRepository().FindByActive(ctx, req)
	if err != nil {
		status = "error"
		_, mappedErr := sharederrorhandler.HandleError[[]*db.GetCategoriesActiveRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, mappedErr
	}

	s.mencache.SetCachedCategoryActiveCache(ctx, req, category, totalRecords)
	logSuccess("Successfully fetched active categories", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return category, totalRecords, nil
}

func (s *categoryQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, *int, error) {
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

	if data, total, found := s.mencache.GetCachedCategoryTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched trashed categories from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	categories, totalRecords, err := s.categoryQueryRepository().FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		_, mappedErr := sharederrorhandler.HandleError[[]*db.GetCategoriesTrashedRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, mappedErr
	}

	s.mencache.SetCachedCategoryTrashedCache(ctx, req, categories, totalRecords)
	logSuccess("Successfully fetched trashed categories", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return categories, totalRecords, nil
}

func (s *categoryQueryService) FindById(ctx context.Context, category_id int) (*db.Category, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("category.id", category_id),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedCategoryCache(ctx, category_id); found {
		logSuccess("Successfully fetched category from cache", zap.Int("category.id", category_id))
		return data, nil
	}

	category, err := s.categoryQueryRepository().FindById(ctx, category_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Category](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedCategoryCache(ctx, category)
	logSuccess("Successfully fetched category", zap.Int("category.id", category_id))
	return category, nil
}

func (s *categoryQueryService) categoryQueryRepository() repository.CategoryQueryRepository {
	return s.categoryQuery
}

func (s *categoryQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
