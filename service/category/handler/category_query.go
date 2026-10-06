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
)

type categoryQueryHandleGrpc struct {
	pbcategory.UnimplementedCategoryQueryServiceServer
	categoryQuery service.CategoryQueryService
	logger        logger.LoggerInterface
}

func NewCategoryQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcategory.CategoryQueryServiceServer {
	return &categoryQueryHandleGrpc{
		categoryQuery: service.CategoryQuery,
		logger:        logger,
	}
}

func (s *categoryQueryHandleGrpc) FindAll(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategory, error) {
	s.logger.Info("FindAll categories called", zap.Int32("page", request.GetPage()))

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

	category, totalRecords, err := s.categoryQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll categories success", zap.Int("count", len(category)))

	return &pbcategory.ApiResponsePaginationCategory{
		Status:     "success",
		Message:    "Successfully fetched categories",
		Data:       mapResponsesCategory(category),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *categoryQueryHandleGrpc) FindById(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategory, error) {
	s.logger.Info("FindById category called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById category success", zap.Int("id", id))

	return &pbcategory.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully fetched category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryQueryHandleGrpc) FindByActive(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategoryDeleteAt, error) {
	s.logger.Info("FindByActive categories called", zap.Int32("page", request.GetPage()))

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

	categories, totalRecords, err := s.categoryQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive categories success")

	return &pbcategory.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active categories",
		Data:       mapResponsesCategoryActive(categories),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
