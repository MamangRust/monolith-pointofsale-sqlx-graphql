package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-pointofsale-order/service"
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/order_errors"
	"go.uber.org/zap"
)

type orderQueryHandleGrpc struct {
	pborder.UnimplementedOrderQueryServiceServer
	orderQuery service.OrderQueryService
	logger     logger.LoggerInterface
}

func NewOrderQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborder.OrderQueryServiceServer {
	return &orderQueryHandleGrpc{
		orderQuery: service.OrderQuery,
		logger:     logger,
	}
}

func (s *orderQueryHandleGrpc) FindAll(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrder, error) {
	s.logger.Info("FindAll orders called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll orders success")

	return &pborder.ApiResponsePaginationOrder{
		Status:     "success",
		Message:    "Successfully fetched order",
		Data:       mapResponsesOrder(merchant),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderQueryHandleGrpc) FindByMerchant(ctx context.Context, request *pborder.FindAllOrderMerchantRequest) (*pborder.ApiResponsePaginationOrder, error) {
	s.logger.Info("FindByMerchant orders called", zap.Int32("merchant_id", request.GetMerchantId()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderMerchant{
		MerchantID: int(request.GetMerchantId()),
		Search:     search,
		Page:       page,
		PageSize:   pageSize,
	}

	orders, totalRecords, err := s.orderQuery.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByMerchant orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByMerchant orders success")

	return &pborder.ApiResponsePaginationOrder{
		Status:     "success",
		Message:    "Successfully fetched order",
		Data:       mapResponsesOrderByMerchant(orders),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderQueryHandleGrpc) FindById(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrder, error) {
	s.logger.Info("FindById order called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById order success")

	return &pborder.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully fetched order",
		Data:    mapResponseOrder(merchant),
	}, nil
}

func (s *orderQueryHandleGrpc) FindByActive(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrderDeleteAt, error) {
	s.logger.Info("FindByActive orders called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive orders success")

	return &pborder.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order",
		Data:       mapResponsesOrderActive(merchant),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
