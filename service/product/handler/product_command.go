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
	"google.golang.org/protobuf/types/known/emptypb"
)

type productCommandHandleGrpc struct {
	pbproduct.UnimplementedProductCommandServiceServer
	productQueryService   service.ProductQueryService
	productCommandService service.ProductCommandService
	logger                logger.LoggerInterface
}

func NewProductCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbproduct.ProductCommandServiceServer {
	return &productCommandHandleGrpc{
		productQueryService:   service.ProductQuery,
		productCommandService: service.ProductCommand,
		logger:                logger,
	}
}

func (s *productCommandHandleGrpc) FindByTrashed(ctx context.Context, request *pbproduct.FindAllProductRequest) (*pbproduct.ApiResponsePaginationProductDeleteAt, error) {
	s.logger.Info("FindByTrashed products called")

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

	users, totalRecords, err := s.productQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed products failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed products success")

	return &pbproduct.ApiResponsePaginationProductDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed product",
		Data:       mapResponsesProductTrashed(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *productCommandHandleGrpc) Create(ctx context.Context, request *pbproduct.CreateProductRequest) (*pbproduct.ApiResponseProduct, error) {
	s.logger.Info("Create product called", zap.String("name", request.GetName()))

	req := &requests.CreateProductRequest{
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		return nil, product_errors.ErrGrpcValidateCreateProduct
	}

	product, err := s.productCommandService.CreateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Create product failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create product success")

	return &pbproduct.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully created product",
		Data:    mapResponseProduct(product),
	}, nil
}

func (s *productCommandHandleGrpc) Update(ctx context.Context, request *pbproduct.UpdateProductRequest) (*pbproduct.ApiResponseProduct, error) {
	s.logger.Info("Update product called", zap.Int32("id", request.GetProductId()))

	id := int(request.GetProductId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateProductRequest{
		ProductID:    &id,
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		return nil, product_errors.ErrGrpcValidateUpdateProduct
	}

	product, err := s.productCommandService.UpdateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Update product failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update product success")

	return &pbproduct.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully updated product",
		Data:    mapResponseProduct(product),
	}, nil
}

func (s *productCommandHandleGrpc) AdjustProductStock(ctx context.Context, request *pbproduct.AdjustProductStockRequest) (*pbproduct.ApiResponseProduct, error) {
	s.logger.Info("Adjust product stock called",
		zap.Int32("product_id", request.GetProductId()),
		zap.Int32("delta", request.GetDelta()))

	product, err := s.productCommandService.AdjustProductStock(ctx, int(request.GetProductId()), int(request.GetDelta()))
	if err != nil {
		s.logger.Error("Adjust product stock failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Adjust product stock success")

	return &pbproduct.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully adjusted product stock",
		Data:    mapResponseProduct(product),
	}, nil
}

func (s *productCommandHandleGrpc) TrashedProduct(ctx context.Context, request *pbproduct.FindByIdProductRequest) (*pbproduct.ApiResponseProductDeleteAt, error) {
	s.logger.Info("TrashedProduct called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productCommandService.TrashProduct(ctx, id)
	if err != nil {
		s.logger.Error("TrashedProduct failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedProduct success")

	return &pbproduct.ApiResponseProductDeleteAt{
		Status:  "success",
		Message: "Successfully trashed product",
		Data:    mapResponseProductDeleteAt(product),
	}, nil
}

func (s *productCommandHandleGrpc) RestoreProduct(ctx context.Context, request *pbproduct.FindByIdProductRequest) (*pbproduct.ApiResponseProductDeleteAt, error) {
	s.logger.Info("RestoreProduct called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productCommandService.RestoreProduct(ctx, id)
	if err != nil {
		s.logger.Error("RestoreProduct failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreProduct success")

	return &pbproduct.ApiResponseProductDeleteAt{
		Status:  "success",
		Message: "Successfully restored product",
		Data:    mapResponseProductDeleteAt(product),
	}, nil
}

func (s *productCommandHandleGrpc) DeleteProductPermanent(ctx context.Context, request *pbproduct.FindByIdProductRequest) (*pbproduct.ApiResponseProductDelete, error) {
	s.logger.Info("DeleteProductPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	_, err := s.productCommandService.DeleteProductPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteProductPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteProductPermanent success")

	return &pbproduct.ApiResponseProductDelete{
		Status:  "success",
		Message: "Successfully deleted Product permanently",
	}, nil
}

func (s *productCommandHandleGrpc) RestoreAllProduct(ctx context.Context, _ *emptypb.Empty) (*pbproduct.ApiResponseProductAll, error) {
	s.logger.Info("RestoreAllProduct called")

	_, err := s.productCommandService.RestoreAllProducts(ctx)
	if err != nil {
		s.logger.Error("RestoreAllProduct failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllProduct success")

	return &pbproduct.ApiResponseProductAll{
		Status:  "success",
		Message: "Successfully restore all Product",
	}, nil
}

func (s *productCommandHandleGrpc) DeleteAllProductPermanent(ctx context.Context, _ *emptypb.Empty) (*pbproduct.ApiResponseProductAll, error) {
	s.logger.Info("DeleteAllProductPermanent called")

	_, err := s.productCommandService.DeleteAllProductsPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllProductPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllProductPermanent success")

	return &pbproduct.ApiResponseProductAll{
		Status:  "success",
		Message: "Successfully delete Product permanen",
	}, nil
}
