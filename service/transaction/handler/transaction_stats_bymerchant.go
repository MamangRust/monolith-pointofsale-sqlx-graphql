package handler

import (
	"context"

	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-pointofsale-transacton/service"
	"go.uber.org/zap"
)

type transactionStatsByMerchantHandleGrpc struct {
	pbtransaction.UnimplementedTransactionStatsByMerchantServiceServer
	transactionStatsByMerchant service.TransactionStatsByMerchantService
	logger                     logger.LoggerInterface
}

func NewTransactionStatsByMerchantHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) *transactionStatsByMerchantHandleGrpc {
	return &transactionStatsByMerchantHandleGrpc{
		transactionStatsByMerchant: service.TransactionStatsByMerchant,
		logger:                     logger,
	}
}

func (s *transactionStatsByMerchantHandleGrpc) FindMonthStatusSuccessByMerchant(ctx context.Context, request *pbtransaction.FindMonthlyTransactionStatusByMerchant) (*pbtransaction.ApiResponseTransactionMonthAmountSuccess, error) {
	s.logger.Info("FindMonthStatusSuccessByMerchant transactions called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	month := int(request.GetMonth())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.MonthAmountTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindMonthlyAmountSuccessByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthStatusSuccessByMerchant transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthStatusSuccessByMerchant transactions success")

	return &pbtransaction.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Merchant monthly success data retrieved successfully",
		Data:    mapResponsesTransactionMonthlyAmountSuccessByMerchant(res),
	}, nil
}

func (s *transactionStatsByMerchantHandleGrpc) FindYearStatusSuccessByMerchant(ctx context.Context, request *pbtransaction.FindYearlyTransactionStatusByMerchant) (*pbtransaction.ApiResponseTransactionYearAmountSuccess, error) {
	s.logger.Info("FindYearStatusSuccessByMerchant transactions called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearAmountTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindYearlyAmountSuccessByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearStatusSuccessByMerchant transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearStatusSuccessByMerchant transactions success")

	return &pbtransaction.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Merchant yearly success data retrieved successfully",
		Data:    mapResponsesTransactionYearlyAmountSuccessByMerchant(res),
	}, nil
}

func (s *transactionStatsByMerchantHandleGrpc) FindMonthStatusFailedByMerchant(ctx context.Context, request *pbtransaction.FindMonthlyTransactionStatusByMerchant) (*pbtransaction.ApiResponseTransactionMonthAmountFailed, error) {
	s.logger.Info("FindMonthStatusFailedByMerchant transactions called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	month := int(request.GetMonth())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.MonthAmountTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindMonthlyAmountFailedByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthStatusFailedByMerchant transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthStatusFailedByMerchant transactions success")

	return &pbtransaction.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Merchant monthly failed data retrieved successfully",
		Data:    mapResponsesTransactionMonthlyAmountFailedByMerchant(res),
	}, nil
}

func (s *transactionStatsByMerchantHandleGrpc) FindYearStatusFailedByMerchant(ctx context.Context, request *pbtransaction.FindYearlyTransactionStatusByMerchant) (*pbtransaction.ApiResponseTransactionYearAmountFailed, error) {
	s.logger.Info("FindYearStatusFailedByMerchant transactions called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearAmountTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindYearlyAmountFailedByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearStatusFailedByMerchant transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearStatusFailedByMerchant transactions success")

	return &pbtransaction.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Merchant yearly failed data retrieved successfully",
		Data:    mapResponsesTransactionYearlyAmountFailedByMerchant(res),
	}, nil
}

func (s *transactionStatsByMerchantHandleGrpc) FindMonthMethodByMerchantSuccess(ctx context.Context, req *pbtransaction.MonthTransactionMethodByMerchant) (*pbtransaction.ApiResponseTransactionMonthPaymentMethod, error) {
	s.logger.Info("FindMonthMethodByMerchantSuccess transactions called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
		Month:      month,
	}

	methods, err := s.transactionStatsByMerchant.FindMonthlyMethodByMerchantSuccess(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthMethodByMerchantSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthMethodByMerchantSuccess transactions success")

	return &pbtransaction.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    mapResponsesTransactionMonthlyMethodByMerchantSuccess(methods),
	}, nil
}

func (s *transactionStatsByMerchantHandleGrpc) FindYearMethodByMerchantSuccess(ctx context.Context, req *pbtransaction.YearTransactionMethodByMerchant) (*pbtransaction.ApiResponseTransactionYearPaymentmethod, error) {
	s.logger.Info("FindYearMethodByMerchantSuccess transactions called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.transactionStatsByMerchant.FindYearlyMethodByMerchantSuccess(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearMethodByMerchantSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearMethodByMerchantSuccess transactions success")

	return &pbtransaction.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesTransactionYearlyMethodByMerchantSuccess(methods),
	}, nil
}

func (s *transactionStatsByMerchantHandleGrpc) FindMonthMethodByMerchantFailed(ctx context.Context, req *pbtransaction.MonthTransactionMethodByMerchant) (*pbtransaction.ApiResponseTransactionMonthPaymentMethod, error) {
	s.logger.Info("FindMonthMethodByMerchantFailed transactions called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
		Month:      month,
	}

	methods, err := s.transactionStatsByMerchant.FindMonthlyMethodByMerchantFailed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthMethodByMerchantFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthMethodByMerchantFailed transactions success")

	return &pbtransaction.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    mapResponsesTransactionMonthlyMethodByMerchantFailed(methods),
	}, nil
}

func (s *transactionStatsByMerchantHandleGrpc) FindYearMethodByMerchantFailed(ctx context.Context, req *pbtransaction.YearTransactionMethodByMerchant) (*pbtransaction.ApiResponseTransactionYearPaymentmethod, error) {
	s.logger.Info("FindYearMethodByMerchantFailed transactions called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.transactionStatsByMerchant.FindYearlyMethodByMerchantFailed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyMethodByMerchantFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyMethodByMerchantFailed transactions success")

	return &pbtransaction.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesTransactionYearlyMethodByMerchantFailed(methods),
	}, nil
}
