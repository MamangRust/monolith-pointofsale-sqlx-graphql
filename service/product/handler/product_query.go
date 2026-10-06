package handler

import (
	"context"
	"math"

	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-product/service"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/product_errors"
	"go.uber.org/zap"
)

type productQueryHandleGrpc struct {
	pbproduct.UnimplementedProductQueryServiceServer
	productQueryService service.ProductQueryService
	logger              logger.LoggerInterface
}

func NewProductQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbproduct.ProductQueryServiceServer {
	return &productQueryHandleGrpc{
		productQueryService: service.ProductQuery,
		logger:              logger,
	}
}

func (s *productQueryHandleGrpc) FindAll(ctx context.Context, request *pbproduct.FindAllProductRequest) (*pbproduct.ApiResponsePaginationProduct, error) {
	s.logger.Info("FindAll products called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllProducts{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	product, totalRecords, err := s.productQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll products failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll products success")

	return &pbproduct.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched product",
		Data:       mapResponsesProduct(product),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *productQueryHandleGrpc) FindByMerchant(ctx context.Context, request *pbproduct.FindAllProductMerchantRequest) (*pbproduct.ApiResponsePaginationProduct, error) {
	s.logger.Info("FindByMerchant products called", zap.Int32("merchantId", request.GetMerchantId()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchant_id := int(request.GetMerchantId())
	min_price := int(request.GetMinPrice())
	max_price := int(request.GetMaxPrice())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if min_price <= 0 {
		min_price = 0
	}
	if max_price <= 0 {
		max_price = 0
	}

	reqService := requests.ProductByMerchantRequest{
		MerchantID: merchant_id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
		MinPrice:   &min_price,
		MaxPrice:   &max_price,
	}

	product, totalRecords, err := s.productQueryService.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByMerchant products failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByMerchant products success")

	return &pbproduct.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched product",
		Data:       mapResponsesProductByMerchant(product, int32(merchant_id)),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *productQueryHandleGrpc) FindByCategory(ctx context.Context, request *pbproduct.FindAllProductCategoryRequest) (*pbproduct.ApiResponsePaginationProduct, error) {
	s.logger.Info("FindByCategory products called", zap.String("categoryName", request.GetCategoryName()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	category_name := request.GetCategoryName()
	min_price := int(request.GetMinprice())
	max_price := int(request.GetMaxprice())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if min_price <= 0 {
		min_price = 0
	}
	if max_price <= 0 {
		max_price = 0
	}

	reqService := requests.ProductByCategoryRequest{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		CategoryName: category_name,
		MinPrice:     &min_price,
		MaxPrice:     &max_price,
	}

	product, totalRecords, err := s.productQueryService.FindByCategory(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByCategory products failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByCategory products success")

	return &pbproduct.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched product",
		Data:       mapResponsesProductByCategory(product),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *productQueryHandleGrpc) FindById(ctx context.Context, request *pbproduct.FindByIdProductRequest) (*pbproduct.ApiResponseProduct, error) {
	s.logger.Info("FindById product called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productQueryService.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById product failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById product success")

	return &pbproduct.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully fetched product",
		Data:    mapResponseProduct(product),
	}, nil
}

func (s *productQueryHandleGrpc) FindByActive(ctx context.Context, request *pbproduct.FindAllProductRequest) (*pbproduct.ApiResponsePaginationProductDeleteAt, error) {
	s.logger.Info("FindByActive products called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllProducts{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	product, totalRecords, err := s.productQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive products failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive products success")

	return &pbproduct.ApiResponsePaginationProductDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active product",
		Data:       mapResponsesProductActive(product),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
