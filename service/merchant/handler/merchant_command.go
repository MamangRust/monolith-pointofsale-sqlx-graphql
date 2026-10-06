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
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantCommandHandleGrpc struct {
	pbmerchant.UnimplementedMerchantCommandServiceServer
	merchantQuery   service.MerchantQueryService
	merchantCommand service.MerchantCommandService
	logger          logger.LoggerInterface
}

func NewMerchantCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbmerchant.MerchantCommandServiceServer {
	return &merchantCommandHandleGrpc{
		merchantQuery:   service.MerchantQuery,
		merchantCommand: service.MerchantCommand,
		logger:          logger,
	}
}

func (s *merchantCommandHandleGrpc) FindByTrashed(ctx context.Context, req *pbmerchant.FindAllMerchantRequest) (*pbmerchant.ApiResponsePaginationMerchantDeleteAt, error) {
	s.logger.Info("FindByTrashed merchants called")

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

	res, totalRecords, err := s.merchantQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed merchants failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed merchants success")

	return &pbmerchant.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       mapResponsesGetMerchantsTrashedRow(res),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantCommandHandleGrpc) Create(ctx context.Context, request *pbmerchant.CreateMerchantRequest) (*pbmerchant.ApiResponseMerchant, error) {
	s.logger.Info("Create merchant called", zap.String("name", request.GetName()))

	req := &requests.CreateMerchantRequest{
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateCreateMerchant
	}

	merchant, err := s.merchantCommand.CreateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Create merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create merchant success")

	return &pbmerchant.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully created merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantCommandHandleGrpc) Update(ctx context.Context, request *pbmerchant.UpdateMerchantRequest) (*pbmerchant.ApiResponseMerchant, error) {
	s.logger.Info("Update merchant called", zap.Int32("id", request.GetMerchantId()))

	id := int(request.GetMerchantId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateMerchantRequest{
		MerchantID:   &id,
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchant
	}

	merchant, err := s.merchantCommand.UpdateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Update merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update merchant success")

	return &pbmerchant.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantCommandHandleGrpc) UpdateMerchantStatus(ctx context.Context, req *pbmerchant.UpdateMerchantStatusRequest) (*pbmerchant.ApiResponseMerchant, error) {
	s.logger.Info("UpdateMerchantStatus merchant called", zap.Int32("id", req.GetMerchantId()))

	id := int(req.GetMerchantId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	request := requests.UpdateMerchantStatusRequest{
		MerchantID: &id,
		Status:     req.GetStatus(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchantStatus
	}

	merchant, err := s.merchantCommand.UpdateMerchantStatus(ctx, &request)
	if err != nil {
		s.logger.Error("UpdateMerchantStatus merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateMerchantStatus merchant success")

	return &pbmerchant.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant status",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantCommandHandleGrpc) TrashedMerchant(ctx context.Context, request *pbmerchant.FindByIdMerchantRequest) (*pbmerchant.ApiResponseMerchantDeleteAt, error) {
	s.logger.Info("TrashedMerchant called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	merchant, err := s.merchantCommand.TrashedMerchant(ctx, id)
	if err != nil {
		s.logger.Error("TrashedMerchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedMerchant success")

	return &pbmerchant.ApiResponseMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    mapResponseMerchantDeleteAt(merchant),
	}, nil
}

func (s *merchantCommandHandleGrpc) RestoreMerchant(ctx context.Context, request *pbmerchant.FindByIdMerchantRequest) (*pbmerchant.ApiResponseMerchant, error) {
	s.logger.Info("RestoreMerchant called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	merchant, err := s.merchantCommand.RestoreMerchant(ctx, id)
	if err != nil {
		s.logger.Error("RestoreMerchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreMerchant success")

	return &pbmerchant.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantCommandHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pbmerchant.FindByIdMerchantRequest) (*pbmerchant.ApiResponseMerchantDelete, error) {
	s.logger.Info("DeleteMerchantPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	_, err := s.merchantCommand.DeleteMerchantPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteMerchantPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteMerchantPermanent success")

	return &pbmerchant.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant permanently",
	}, nil
}

func (s *merchantCommandHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	s.logger.Info("RestoreAllMerchant called")

	_, err := s.merchantCommand.RestoreAllMerchant(ctx)
	if err != nil {
		s.logger.Error("RestoreAllMerchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllMerchant success")

	return &pbmerchant.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restore all merchant",
	}, nil
}

func (s *merchantCommandHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	s.logger.Info("DeleteAllMerchantPermanent called")

	_, err := s.merchantCommand.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllMerchantPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllMerchantPermanent success")

	return &pbmerchant.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully delete merchant permanen",
	}, nil
}
