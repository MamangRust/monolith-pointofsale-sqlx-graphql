package handler

import (
	"context"
	"math"

	pbcommon "github.com/MamangRust/monolith-pointofsale-pb/common"
	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-pointofsale-transacton/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionCommandHandleGrpc struct {
	pbtransaction.UnimplementedTransactionCommandServiceServer
	transactionQuery   service.TransactionQueryService
	transactionCommand service.TransactionCommandService
	logger             logger.LoggerInterface
}

func NewTransactionCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbtransaction.TransactionCommandServiceServer {
	return &transactionCommandHandleGrpc{
		transactionQuery:   service.TransactionQuery,
		transactionCommand: service.TransactionCommand,
		logger:             logger,
	}
}

func (s *transactionCommandHandleGrpc) FindByTrashed(ctx context.Context, request *pbtransaction.FindAllTransactionRequest) (*pbtransaction.ApiResponsePaginationTransactionDeleteAt, error) {
	s.logger.Info("FindByTrashed transactions called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed transactions success")

	return &pbtransaction.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed transaction",
		Data:       mapResponsesTransactionTrashed(transaction),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *transactionCommandHandleGrpc) Create(ctx context.Context, request *pbtransaction.CreateTransactionRequest) (*pbtransaction.ApiResponseTransaction, error) {
	s.logger.Info("Create transaction called", zap.Int32("orderId", request.GetOrderId()))

	req := &requests.CreateTransactionRequest{
		CashierID:     int(request.GetCashierId()),
		OrderID:       int(request.GetOrderId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateCreateTransaction
	}

	transaction, err := s.transactionCommand.CreateTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Create transaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create transaction success")

	return &pbtransaction.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully created transaction",
		Data:    mapResponseTransaction(transaction),
	}, nil
}

func (s *transactionCommandHandleGrpc) Update(ctx context.Context, request *pbtransaction.UpdateTransactionRequest) (*pbtransaction.ApiResponseTransaction, error) {
	s.logger.Info("Update transaction called", zap.Int32("id", request.GetTransactionId()))

	id := int(request.GetTransactionId())
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateTransactionRequest{
		TransactionID: &id,
		OrderID:       int(request.GetOrderId()),
		CashierID:     int(request.GetCashierId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateUpdateTransaction
	}

	transaction, err := s.transactionCommand.UpdateTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Update transaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update transaction success")

	return &pbtransaction.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully updated transaction",
		Data:    mapResponseTransaction(transaction),
	}, nil
}

func (s *transactionCommandHandleGrpc) TrashedTransaction(ctx context.Context, request *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDeleteAt, error) {
	s.logger.Info("TrashedTransaction called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionCommand.TrashedTransaction(ctx, id)
	if err != nil {
		s.logger.Error("TrashedTransaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedTransaction success")

	return &pbtransaction.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transaction",
		Data:    mapResponseTransactionDeleteAt(transaction),
	}, nil
}

func (s *transactionCommandHandleGrpc) RestoreTransaction(ctx context.Context, request *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDeleteAt, error) {
	s.logger.Info("RestoreTransaction called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionCommand.RestoreTransaction(ctx, id)
	if err != nil {
		s.logger.Error("RestoreTransaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreTransaction success")

	return &pbtransaction.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully restored transaction",
		Data:    mapResponseTransactionDeleteAt(transaction),
	}, nil
}

func (s *transactionCommandHandleGrpc) DeleteTransactionPermanent(ctx context.Context, request *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDelete, error) {
	s.logger.Info("DeleteTransactionPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	_, err := s.transactionCommand.DeleteTransactionPermanently(ctx, id)
	if err != nil {
		s.logger.Error("DeleteTransactionPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteTransactionPermanent success")

	return &pbtransaction.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted Transaction permanently",
	}, nil
}

func (s *transactionCommandHandleGrpc) RestoreAllTransaction(ctx context.Context, _ *emptypb.Empty) (*pbtransaction.ApiResponseTransactionAll, error) {
	s.logger.Info("RestoreAllTransaction called")

	_, err := s.transactionCommand.RestoreAllTransactions(ctx)
	if err != nil {
		s.logger.Error("RestoreAllTransaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllTransaction success")

	return &pbtransaction.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully restore all Transaction",
	}, nil
}

func (s *transactionCommandHandleGrpc) DeleteAllTransactionPermanent(ctx context.Context, _ *emptypb.Empty) (*pbtransaction.ApiResponseTransactionAll, error) {
	s.logger.Info("DeleteAllTransactionPermanent called")

	_, err := s.transactionCommand.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllTransactionPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllTransactionPermanent success")

	return &pbtransaction.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully delete Transaction permanen",
	}, nil
}
