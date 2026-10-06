package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-category/cache"
	"github.com/MamangRust/monolith-pointofsale-category/repository"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-pkg/utils"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type categoryCommandDeps struct {
	Cache           mencache.CategoryCommandCache
	CategoryQuery   repository.CategoryQueryRepository
	CategoryCommand repository.CategoryCommandRepository
	Logger          logger.LoggerInterface
	Observability   observability.TraceLoggerObservability
}

type categoryCommandService struct {
	mencache        mencache.CategoryCommandCache
	categoryQuery   repository.CategoryQueryRepository
	categoryCommand repository.CategoryCommandRepository
	logger          logger.LoggerInterface
	observability   observability.TraceLoggerObservability
}

func NewCategoryCommandService(params *categoryCommandDeps) CategoryCommandService {
	return &categoryCommandService{
		mencache:        params.Cache,
		categoryQuery:   params.CategoryQuery,
		categoryCommand: params.CategoryCommand,
		logger:          params.Logger,
		observability:   params.Observability,
	}
}

func (s *categoryCommandService) CreateCategory(ctx context.Context, req *requests.CreateCategoryRequest) (*db.Category, error) {
	const method = "CreateCategory"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("name", req.Name))
	defer func() {
		end(status)
	}()

	slug := utils.GenerateSlug(req.Name)
	req.SlugCategory = &slug

	category, err := s.categoryCommand.CreateCategory(ctx, req)
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

	logSuccess("Successfully created category", zap.Int("category.id", int(category.CategoryID)))
	return category, nil
}

func (s *categoryCommandService) UpdateCategory(ctx context.Context, req *requests.UpdateCategoryRequest) (*db.Category, error) {
	const method = "UpdateCategory"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("name", req.Name))
	defer func() {
		end(status)
	}()

	slug := utils.GenerateSlug(req.Name)
	req.SlugCategory = &slug

	category, err := s.categoryCommand.UpdateCategory(ctx, req)
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

	s.mencache.DeleteCachedCategoryCache(ctx, *req.CategoryID)
	logSuccess("Successfully updated category", zap.Int("category.id", int(category.CategoryID)))
	return category, nil
}

func (s *categoryCommandService) TrashedCategory(ctx context.Context, category_id int) (*db.Category, error) {
	const method = "TrashedCategory"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("category.id", category_id))
	defer func() {
		end(status)
	}()

	category, err := s.categoryCommand.TrashedCategory(ctx, category_id)
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

	s.mencache.DeleteCachedCategoryCache(ctx, category_id)
	logSuccess("Successfully trashed category", zap.Int("category.id", category_id))
	return category, nil
}

func (s *categoryCommandService) RestoreCategory(ctx context.Context, categoryID int) (*db.Category, error) {
	const method = "RestoreCategory"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("category.id", categoryID))
	defer func() {
		end(status)
	}()

	category, err := s.categoryCommand.RestoreCategory(ctx, categoryID)
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

	s.mencache.DeleteCachedCategoryCache(ctx, categoryID)
	logSuccess("Successfully restored category", zap.Int("category.id", categoryID))
	return category, nil
}

func (s *categoryCommandService) DeleteCategoryPermanent(ctx context.Context, categoryID int) (bool, error) {
	const method = "DeleteCategoryPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("category.id", categoryID))
	defer func() {
		end(status)
	}()

	_, err := s.categoryQuery.FindByIdTrashed(ctx, categoryID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	success, err := s.categoryCommand.DeleteCategoryPermanently(ctx, categoryID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedCategoryCache(ctx, categoryID)
	logSuccess("Successfully deleted category permanently", zap.Bool("success", success))
	return success, nil
}

func (s *categoryCommandService) RestoreAllCategories(ctx context.Context) (bool, error) {
	const method = "RestoreAllCategories"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.categoryCommand.RestoreAllCategories(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedCategoryAllCache(ctx)
	logSuccess("Successfully restored all trashed categories", zap.Bool("success", success))
	return success, nil
}

func (s *categoryCommandService) DeleteAllCategoriesPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllCategoriesPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.categoryCommand.DeleteAllPermanentCategories(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedCategoryAllCache(ctx)
	logSuccess("Successfully deleted all categories permanently", zap.Bool("success", success))
	return success, nil
}
