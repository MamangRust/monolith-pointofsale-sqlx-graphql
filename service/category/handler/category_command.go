package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-pointofsale-category/service"
	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/category_errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type categoryCommandHandleGrpc struct {
	pbcategory.UnimplementedCategoryCommandServiceServer
	categoryQuery   service.CategoryQueryService
	categoryCommand service.CategoryCommandService
	logger          logger.LoggerInterface
}

func NewCategoryCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcategory.CategoryCommandServiceServer {
	return &categoryCommandHandleGrpc{
		categoryQuery:   service.CategoryQuery,
		categoryCommand: service.CategoryCommand,
		logger:          logger,
	}
}

func (s *categoryCommandHandleGrpc) FindByTrashed(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategoryDeleteAt, error) {
	s.logger.Info("FindByTrashed categories called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed categories success")

	return &pbcategory.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed categories",
		Data:       mapResponsesCategoryTrashed(categories),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *categoryCommandHandleGrpc) Create(ctx context.Context, request *pbcategory.CreateCategoryRequest) (*pbcategory.ApiResponseCategory, error) {
	s.logger.Info("Create category called", zap.String("name", request.GetName()))

	req := &requests.CreateCategoryRequest{
		Name:        request.GetName(),
		Description: request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateCreateCategory
	}

	category, err := s.categoryCommand.CreateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Create category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create category success")

	return &pbcategory.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully created category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryCommandHandleGrpc) Update(ctx context.Context, request *pbcategory.UpdateCategoryRequest) (*pbcategory.ApiResponseCategory, error) {
	s.logger.Info("Update category called", zap.Int32("id", request.GetCategoryId()))

	id := int(request.GetCategoryId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateCategoryRequest{
		CategoryID:  &id,
		Name:        request.GetName(),
		Description: request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateUpdateCategory
	}

	category, err := s.categoryCommand.UpdateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Update category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update category success")

	return &pbcategory.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully updated category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryCommandHandleGrpc) TrashedCategory(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategoryDeleteAt, error) {
	s.logger.Info("TrashedCategory called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryCommand.TrashedCategory(ctx, id)
	if err != nil {
		s.logger.Error("TrashedCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedCategory success")

	return &pbcategory.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully trashed category",
		Data:    mapResponseCategoryDeleteAt(category),
	}, nil
}

func (s *categoryCommandHandleGrpc) RestoreCategory(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategoryDeleteAt, error) {
	s.logger.Info("RestoreCategory called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryCommand.RestoreCategory(ctx, id)
	if err != nil {
		s.logger.Error("RestoreCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreCategory success")

	return &pbcategory.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully restored category",
		Data:    mapResponseCategoryDeleteAt(category),
	}, nil
}

func (s *categoryCommandHandleGrpc) DeleteCategoryPermanent(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategoryDelete, error) {
	s.logger.Info("DeleteCategoryPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.categoryCommand.DeleteCategoryPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteCategoryPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteCategoryPermanent success")

	return &pbcategory.ApiResponseCategoryDelete{
		Status:  "success",
		Message: "Successfully deleted category permanently",
	}, nil
}

func (s *categoryCommandHandleGrpc) RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pbcategory.ApiResponseCategoryAll, error) {
	s.logger.Info("RestoreAllCategory called")

	_, err := s.categoryCommand.RestoreAllCategories(ctx)
	if err != nil {
		s.logger.Error("RestoreAllCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllCategory success")

	return &pbcategory.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully restore all category",
	}, nil
}

func (s *categoryCommandHandleGrpc) DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pbcategory.ApiResponseCategoryAll, error) {
	s.logger.Info("DeleteAllCategoryPermanent called")

	_, err := s.categoryCommand.DeleteAllCategoriesPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllCategoryPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllCategoryPermanent success")

	return &pbcategory.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully delete category permanen",
	}, nil
}
