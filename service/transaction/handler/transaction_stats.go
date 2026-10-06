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

type transactionStatsHandleGrpc struct {
	*transactionStatsCoreHandleGrpc
	*transactionStatsByMerchantHandleGrpc
}

func NewTransactionStatsHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbtransaction.TransactionStatsServiceServer {
	return &transactionStatsHandleGrpc{
		transactionStatsCoreHandleGrpc:       NewTransactionStatsCoreHandleGrpc(service, logger),
		transactionStatsByMerchantHandleGrpc: NewTransactionStatsByMerchantHandleGrpc(service, logger),
	}
}

type transactionStatsCoreHandleGrpc struct {
	pbtransaction.UnimplementedTransactionStatsServiceServer
	transactionStats service.TransactionStatsService
	logger           logger.LoggerInterface
}

func NewTransactionStatsCoreHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) *transactionStatsCoreHandleGrpc {
	return &transactionStatsCoreHandleGrpc{
		transactionStats: service.TransactionStats,
		logger:           logger,
	}
}

func (s *transactionStatsCoreHandleGrpc) FindMonthStatusSuccess(ctx context.Context, request *pbtransaction.FindMonthlyTransactionStatus) (*pbtransaction.ApiResponseTransactionMonthAmountSuccess, error) {
	s.logger.Info("FindMonthStatusSuccess transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionStats.FindMonthlyAmountSuccess(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthStatusSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthStatusSuccess transactions success")

	return &pbtransaction.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Monthly success data retrieved successfully",
		Data:    mapResponsesTransactionMonthlyAmountSuccess(res),
	}, nil
}

func (s *transactionStatsCoreHandleGrpc) FindYearStatusSuccess(ctx context.Context, request *pbtransaction.FindYearlyTransactionStatus) (*pbtransaction.ApiResponseTransactionYearAmountSuccess, error) {
	s.logger.Info("FindYearStatusSuccess transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	res, err := s.transactionStats.FindYearlyAmountSuccess(ctx, year)
	if err != nil {
		s.logger.Error("FindYearStatusSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearStatusSuccess transactions success")

	return &pbtransaction.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Yearly success data retrieved successfully",
		Data:    mapResponsesTransactionYearlyAmountSuccess(res),
	}, nil
}

func (s *transactionStatsCoreHandleGrpc) FindMonthStatusFailed(ctx context.Context, request *pbtransaction.FindMonthlyTransactionStatus) (*pbtransaction.ApiResponseTransactionMonthAmountFailed, error) {
	s.logger.Info("FindMonthStatusFailed transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionStats.FindMonthlyAmountFailed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthStatusFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthStatusFailed transactions success")

	return &pbtransaction.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Monthly failed data retrieved successfully",
		Data:    mapResponsesTransactionMonthlyAmountFailed(res),
	}, nil
}

func (s *transactionStatsCoreHandleGrpc) FindYearStatusFailed(ctx context.Context, request *pbtransaction.FindYearlyTransactionStatus) (*pbtransaction.ApiResponseTransactionYearAmountFailed, error) {
	s.logger.Info("FindYearStatusFailed transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	res, err := s.transactionStats.FindYearlyAmountFailed(ctx, year)
	if err != nil {
		s.logger.Error("FindYearStatusFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearStatusFailed transactions success")

	return &pbtransaction.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Yearly failed data retrieved successfully",
		Data:    mapResponsesTransactionYearlyAmountFailed(res),
	}, nil
}

func (s *transactionStatsCoreHandleGrpc) FindMonthMethodSuccess(ctx context.Context, req *pbtransaction.MonthTransactionMethod) (*pbtransaction.ApiResponseTransactionMonthPaymentMethod, error) {
	s.logger.Info("FindMonthMethodSuccess transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	methods, err := s.transactionStats.FindMonthlyMethodSuccess(ctx, &requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})
	if err != nil {
		s.logger.Error("FindMonthMethodSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthMethodSuccess transactions success")

	return &pbtransaction.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    mapResponsesTransactionMonthlyMethodSuccess(methods),
	}, nil
}

func (s *transactionStatsCoreHandleGrpc) FindYearMethodSuccess(ctx context.Context, req *pbtransaction.YearTransactionMethod) (*pbtransaction.ApiResponseTransactionYearPaymentmethod, error) {
	s.logger.Info("FindYearMethodSuccess transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := s.transactionStats.FindYearlyMethodSuccess(ctx, year)
	if err != nil {
		s.logger.Error("FindYearMethodSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearMethodSuccess transactions success")

	return &pbtransaction.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponsesTransactionYearlyMethodSuccess(methods),
	}, nil
}

func (s *transactionStatsCoreHandleGrpc) FindMonthMethodFailed(ctx context.Context, req *pbtransaction.MonthTransactionMethod) (*pbtransaction.ApiResponseTransactionMonthPaymentMethod, error) {
	s.logger.Info("FindMonthMethodFailed transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	methods, err := s.transactionStats.FindMonthlyMethodFailed(ctx, &requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})
	if err != nil {
		s.logger.Error("FindMonthMethodFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthMethodFailed transactions success")

	return &pbtransaction.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    mapResponsesTransactionMonthlyMethodFailed(methods),
	}, nil
}

func (s *transactionStatsCoreHandleGrpc) FindYearMethodFailed(ctx context.Context, req *pbtransaction.YearTransactionMethod) (*pbtransaction.ApiResponseTransactionYearPaymentmethod, error) {
	s.logger.Info("FindYearMethodFailed transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := s.transactionStats.FindYearlyMethodFailed(ctx, year)
	if err != nil {
		s.logger.Error("FindYearMethodFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearMethodFailed transactions success")

	return &pbtransaction.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponsesTransactionYearlyMethodFailed(methods),
	}, nil
}
