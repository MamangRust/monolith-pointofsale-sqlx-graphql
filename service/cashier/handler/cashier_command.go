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
	"google.golang.org/protobuf/types/known/emptypb"
)

type cashierCommandHandleGrpc struct {
	pbcashier.UnimplementedCashierCommandServiceServer
	cashierQuery   service.CashierQueryService
	cashierCommand service.CashierCommandService
	logger         logger.LoggerInterface
}

func NewCashierCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcashier.CashierCommandServiceServer {
	return &cashierCommandHandleGrpc{
		cashierQuery:   service.CashierQuery,
		cashierCommand: service.CashierCommand,
		logger:         logger,
	}
}

func (s *cashierCommandHandleGrpc) FindByTrashed(ctx context.Context, request *pbcashier.FindAllCashierRequest) (*pbcashier.ApiResponsePaginationCashierDeleteAt, error) {
	s.logger.Info("FindByTrashed cashier called")

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

	users, totalRecords, err := s.cashierQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed cashier success", zap.Int("count", len(users)))

	return &pbcashier.ApiResponsePaginationCashierDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed cashier",
		Data:       mapResponsesCashierTrashed(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *cashierCommandHandleGrpc) CreateCashier(ctx context.Context, request *pbcashier.CreateCashierRequest) (*pbcashier.ApiResponseCashier, error) {
	s.logger.Info("CreateCashier called", zap.String("name", request.GetName()))

	req := &requests.CreateCashierRequest{
		Name:       request.GetName(),
		MerchantID: int(request.GetMerchantId()),
		UserID:     int(request.GetUserId()),
	}

	if err := req.Validate(); err != nil {
		return nil, cashier_errors.ErrGrpcValidateCreateCashier
	}

	cashier, err := s.cashierCommand.CreateCashier(ctx, req)
	if err != nil {
		s.logger.Error("CreateCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("CreateCashier success")

	return &pbcashier.ApiResponseCashier{
		Status:  "success",
		Message: "Successfully created cashier",
		Data:    mapResponseCashier(cashier),
	}, nil
}

func (s *cashierCommandHandleGrpc) UpdateCashier(ctx context.Context, request *pbcashier.UpdateCashierRequest) (*pbcashier.ApiResponseCashier, error) {
	s.logger.Info("UpdateCashier called", zap.Int32("id", request.GetCashierId()))

	id := int(request.GetCashierId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateCashierRequest{
		CashierID: &id,
		Name:      request.GetName(),
	}

	if err := req.Validate(); err != nil {
		return nil, cashier_errors.ErrGrpcValidateUpdateCashier
	}

	cashier, err := s.cashierCommand.UpdateCashier(ctx, req)
	if err != nil {
		s.logger.Error("UpdateCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateCashier success")

	return &pbcashier.ApiResponseCashier{
		Status:  "success",
		Message: "Successfully updated cashier",
		Data:    mapResponseCashier(cashier),
	}, nil
}

func (s *cashierCommandHandleGrpc) TrashedCashier(ctx context.Context, request *pbcashier.FindByIdCashierRequest) (*pbcashier.ApiResponseCashierDeleteAt, error) {
	s.logger.Info("TrashedCashier called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	cashier, err := s.cashierCommand.TrashedCashier(ctx, id)
	if err != nil {
		s.logger.Error("TrashedCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedCashier success")

	return &pbcashier.ApiResponseCashierDeleteAt{
		Status:  "success",
		Message: "Successfully trashed cashier",
		Data:    mapResponseCashierDeleteAt(cashier),
	}, nil
}

func (s *cashierCommandHandleGrpc) RestoreCashier(ctx context.Context, request *pbcashier.FindByIdCashierRequest) (*pbcashier.ApiResponseCashierDeleteAt, error) {
	s.logger.Info("RestoreCashier called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	cashier, err := s.cashierCommand.RestoreCashier(ctx, id)
	if err != nil {
		s.logger.Error("RestoreCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreCashier success")

	return &pbcashier.ApiResponseCashierDeleteAt{
		Status:  "success",
		Message: "Successfully restored cashier",
		Data:    mapResponseCashierDeleteAt(cashier),
	}, nil
}

func (s *cashierCommandHandleGrpc) DeleteCashierPermanent(ctx context.Context, request *pbcashier.FindByIdCashierRequest) (*pbcashier.ApiResponseCashierDelete, error) {
	s.logger.Info("DeleteCashierPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.cashierCommand.DeleteCashierPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteCashierPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteCashierPermanent success")

	return &pbcashier.ApiResponseCashierDelete{
		Status:  "success",
		Message: "Successfully deleted cashier permanently",
	}, nil
}

func (s *cashierCommandHandleGrpc) RestoreAllCashier(ctx context.Context, _ *emptypb.Empty) (*pbcashier.ApiResponseCashierAll, error) {
	s.logger.Info("RestoreAllCashier called")

	_, err := s.cashierCommand.RestoreAllCashier(ctx)
	if err != nil {
		s.logger.Error("RestoreAllCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllCashier success")

	return &pbcashier.ApiResponseCashierAll{
		Status:  "success",
		Message: "Successfully restore all cashier",
	}, nil
}

func (s *cashierCommandHandleGrpc) DeleteAllCashierPermanent(ctx context.Context, _ *emptypb.Empty) (*pbcashier.ApiResponseCashierAll, error) {
	s.logger.Info("DeleteAllCashierPermanent called")

	_, err := s.cashierCommand.DeleteAllCashierPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllCashierPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllCashierPermanent success")

	return &pbcashier.ApiResponseCashierAll{
		Status:  "success",
		Message: "Successfully delete cashier permanen",
	}, nil
}
