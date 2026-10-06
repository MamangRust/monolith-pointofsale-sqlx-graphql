package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-pointofsale-cashier/service"
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
	"go.uber.org/zap"
)

type cashierQueryHandleGrpc struct {
	pbcashier.UnimplementedCashierQueryServiceServer
	cashierQuery service.CashierQueryService
	logger       logger.LoggerInterface
}

func NewCashierQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcashier.CashierQueryServiceServer {
	return &cashierQueryHandleGrpc{
		cashierQuery: service.CashierQuery,
		logger:       logger,
	}
}

func (s *cashierQueryHandleGrpc) FindAll(ctx context.Context, request *pbcashier.FindAllCashierRequest) (*pbcashier.ApiResponsePaginationCashier, error) {
	s.logger.Info("FindAll cashier called", zap.Int32("page", request.GetPage()), zap.Int32("pageSize", request.GetPageSize()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCashiers{
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	}

	cashier, totalRecords, err := s.cashierQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll cashier success", zap.Int("count", len(cashier)))

	return &pbcashier.ApiResponsePaginationCashier{
		Status:     "success",
		Message:    "Successfully fetched cashier",
		Data:       mapResponsesCashier(cashier),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *cashierQueryHandleGrpc) FindById(ctx context.Context, request *pbcashier.FindByIdCashierRequest) (*pbcashier.ApiResponseCashier, error) {
	s.logger.Info("FindById cashier called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	cashier, err := s.cashierQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById cashier success", zap.Int("id", id))

	return &pbcashier.ApiResponseCashier{
		Status:  "success",
		Message: "Successfully fetched categories",
		Data:    mapResponseCashier(cashier),
	}, nil
}

func (s *cashierQueryHandleGrpc) FindByActive(ctx context.Context, request *pbcashier.FindAllCashierRequest) (*pbcashier.ApiResponsePaginationCashierDeleteAt, error) {
	s.logger.Info("FindByActive cashier called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCashiers{
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	}

	cashier, totalRecords, err := s.cashierQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive cashier success", zap.Int("count", len(cashier)))

	return &pbcashier.ApiResponsePaginationCashierDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active cashier",
		Data:       mapResponsesCashierActive(cashier),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *cashierQueryHandleGrpc) FindByMerchant(ctx context.Context, request *pbcashier.FindByMerchantCashierRequest) (*pbcashier.ApiResponsePaginationCashier, error) {
	s.logger.Info("FindByMerchant cashier called", zap.Int32("merchantId", request.GetMerchantId()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchant_id := int(request.GetMerchantId())

	if merchant_id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCashierMerchant{
		Search:     search,
		Page:       page,
		PageSize:   pageSize,
		MerchantID: merchant_id,
	}

	cashier, totalRecords, err := s.cashierQuery.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByMerchant cashier success", zap.Int("count", len(cashier)))

	return &pbcashier.ApiResponsePaginationCashier{
		Status:     "success",
		Message:    "Successfully fetched cashier",
		Data:       mapResponsesCashierByMerchant(cashier),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
