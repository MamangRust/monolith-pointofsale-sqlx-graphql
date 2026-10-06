package service

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	mencache "github.com/MamangRust/monolith-pointofsale-transacton/cache"
	"github.com/MamangRust/monolith-pointofsale-transacton/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionStatsService struct {
	mencache                   mencache.TransactionStatsCache
	transactionStatsRepository repository.TransactionStatsRepository
	logger                     logger.LoggerInterface
	observability              observability.TraceLoggerObservability
}

func NewTransactionStatsService(
	mencache mencache.TransactionStatsCache,
	transactionStatsRepository repository.TransactionStatsRepository,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *transactionStatsService {
	return &transactionStatsService{
		mencache:                   mencache,
		transactionStatsRepository: transactionStatsRepository,
		logger:                     logger,
		observability:              obs,
	}
}

func (s *transactionStatsService) FindMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionSuccessRow, error) {
	const method = "FindMonthlyAmountSuccess"

	year := req.Year
	month := req.Month

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthAmountSuccessCached(ctx, req); found {
		logSuccess("Successfully fetched monthly successful transaction amounts from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyAmountSuccess(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyAmountTransactionSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountSuccess.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthAmountSuccessCached(ctx, req, res)

	logSuccess("Successfully fetched monthly successful transaction amounts", zap.Int("year", year), zap.Int("month", month))

	return res, nil
}

func (s *transactionStatsService) FindYearlyAmountSuccess(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionSuccessRow, error) {
	const method = "FindYearlyAmountSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearAmountSuccessCached(ctx, year); found {
		logSuccess("Successfully fetched yearly successful transaction amounts from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyAmountSuccess(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyAmountTransactionSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountSuccess.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearAmountSuccessCached(ctx, year, res)

	logSuccess("Successfully fetched yearly successful transaction amounts", zap.Int("year", year))

	return res, nil
}

func (s *transactionStatsService) FindMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionFailedRow, error) {
	const method = "FindMonthlyAmountFailed"

	year := req.Year
	month := req.Month

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthAmountFailedCached(ctx, req); found {
		logSuccess("Successfully fetched monthly failed transaction amounts from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyAmountFailed(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyAmountTransactionFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountFailed.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthAmountFailedCached(ctx, req, res)

	logSuccess("Successfully fetched monthly failed transaction amounts", zap.Int("year", year), zap.Int("month", month))

	return res, nil
}

func (s *transactionStatsService) FindYearlyAmountFailed(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionFailedRow, error) {
	const method = "FindYearlyAmountFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearAmountFailedCached(ctx, year); found {
		logSuccess("Successfully fetched yearly failed transaction amounts from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyAmountFailed(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyAmountTransactionFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountFailed.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearAmountFailedCached(ctx, year, res)

	logSuccess("Successfully fetched yearly failed transaction amounts", zap.Int("year", year))

	return res, nil
}

func (s *transactionStatsService) FindMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsSuccessRow, error) {
	const method = "FindMonthlyMethodSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year), attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthMethodSuccessCached(ctx, req); found {
		logSuccess("Successfully fetched monthly successful transaction methods from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyTransactionMethodSuccess(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyMethod.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthMethodSuccessCached(ctx, req, res)

	logSuccess("Successfully fetched monthly successful transaction methods", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return res, nil
}

func (s *transactionStatsService) FindYearlyMethodSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsSuccessRow, error) {
	const method = "FindYearlyMethodSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearMethodSuccessCached(ctx, year); found {
		logSuccess("Successfully fetched yearly successful transaction methods from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyTransactionMethodSuccess(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTransactionMethodsSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyMethod.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearMethodSuccessCached(ctx, year, res)

	logSuccess("Successfully fetched yearly successful transaction methods", zap.Int("year", year))

	return res, nil
}

func (s *transactionStatsService) FindMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsFailedRow, error) {
	const method = "FindMonthlyMethodFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year), attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthMethodFailedCached(ctx, req); found {
		logSuccess("Successfully fetched monthly failed transaction methods from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyTransactionMethodFailed(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyMethod.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthMethodFailedCached(ctx, req, res)

	logSuccess("Successfully fetched monthly failed transaction methods", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return res, nil
}

func (s *transactionStatsService) FindYearlyMethodFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsFailedRow, error) {
	const method = "FindYearlyMethodFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearMethodFailedCached(ctx, year); found {
		logSuccess("Successfully fetched yearly failed transaction methods from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyTransactionMethodFailed(ctx, year)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTransactionMethodsFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyMethod.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearMethodFailedCached(ctx, year, res)

	logSuccess("Successfully fetched yearly failed transaction methods", zap.Int("year", year))

	return res, nil
}
