package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-pointofsale-merchant/service"
	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/merchant_errors"
	"go.uber.org/zap"
)

type merchantQueryHandleGrpc struct {
	pbmerchant.UnimplementedMerchantQueryServiceServer
	merchantQuery service.MerchantQueryService
	logger        logger.LoggerInterface
}

func NewMerchantQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbmerchant.MerchantQueryServiceServer {
	return &merchantQueryHandleGrpc{
		merchantQuery: service.MerchantQuery,
		logger:        logger,
	}
}

func (s *merchantQueryHandleGrpc) FindAll(ctx context.Context, req *pbmerchant.FindAllMerchantRequest) (*pbmerchant.ApiResponsePaginationMerchant, error) {
	s.logger.Info("FindAll merchant called", zap.Int32("page", req.GetPage()))

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll merchant success")

	return &pbmerchant.ApiResponsePaginationMerchant{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       mapResponsesGetMerchantsRow(merchants),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantQueryHandleGrpc) FindByActive(ctx context.Context, req *pbmerchant.FindAllMerchantRequest) (*pbmerchant.ApiResponsePaginationMerchantDeleteAt, error) {
	s.logger.Info("FindByActive merchants called", zap.Int32("page", req.GetPage()))

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive merchants failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive merchants success")

	return &pbmerchant.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       mapResponsesGetMerchantsActiveRow(res),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantQueryHandleGrpc) FindById(ctx context.Context, request *pbmerchant.FindByIdMerchantRequest) (*pbmerchant.ApiResponseMerchant, error) {
	s.logger.Info("FindById merchant called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	merchant, err := s.merchantQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById merchant success")

	return &pbmerchant.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}
