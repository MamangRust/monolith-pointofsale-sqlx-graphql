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

type transactionStatsByMerchantService struct {
	mencache                             mencache.TransactionStatsByMerchantCache
	transactionStatsByMerchantRepository repository.TransactionStatsByMerchantRepository
	logger                               logger.LoggerInterface
	observability                        observability.TraceLoggerObservability
}

func NewTransactionStatsByMerchantService(
	mencache mencache.TransactionStatsByMerchantCache,
	transactionStatsByMerchantRepository repository.TransactionStatsByMerchantRepository,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *transactionStatsByMerchantService {
	return &transactionStatsByMerchantService{
		mencache:                             mencache,
		transactionStatsByMerchantRepository: transactionStatsByMerchantRepository,
		logger:                               logger,
		observability:                        obs,
	}
}

func (s *transactionStatsByMerchantService) FindMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, error) {
	const method = "FindMonthlyAmountSuccessByMerchant"

	year := req.Year
	month := req.Month
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthAmountSuccessByMerchantCached(ctx, req); found {
		logSuccess("Successfully fetched monthly successful transactions by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyAmountSuccessByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyAmountTransactionSuccessByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountSuccessByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthAmountSuccessByMerchantCached(ctx, req, res)

	logSuccess("Successfully fetched monthly successful transactions by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, error) {
	const method = "FindYearlyAmountSuccessByMerchant"

	year := req.Year
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearAmountSuccessByMerchantCached(ctx, req); found {
		logSuccess("Successfully fetched yearly successful transactions by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyAmountSuccessByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyAmountTransactionSuccessByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountSuccessByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearAmountSuccessByMerchantCached(ctx, req, res)

	logSuccess("Successfully fetched yearly successful transactions by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, error) {
	const method = "FindMonthlyAmountFailedByMerchant"

	year := req.Year
	month := req.Month
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthAmountFailedByMerchantCached(ctx, req); found {
		logSuccess("Successfully fetched monthly failed transactions by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyAmountFailedByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyAmountTransactionFailedByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountFailedByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthAmountFailedByMerchantCached(ctx, req, res)

	logSuccess("Successfully fetched monthly failed transactions by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, error) {
	const method = "FindYearlyAmountFailedByMerchant"

	year := req.Year
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearAmountFailedByMerchantCached(ctx, req); found {
		logSuccess("Successfully fetched yearly failed transactions by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyAmountFailedByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyAmountTransactionFailedByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountFailedByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearAmountFailedByMerchantCached(ctx, req, res)

	logSuccess("Successfully fetched yearly failed transactions by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, error) {
	const method = "FindMonthlyMethodByMerchantSuccess"

	year := req.Year
	month := req.Month
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthMethodSuccessByMerchantCached(ctx, req); found {
		logSuccess("Successfully found monthly successful transaction methods by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyTransactionMethodByMerchantSuccess(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyMethodByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthMethodSuccessByMerchantCached(ctx, req, res)

	logSuccess("Successfully found monthly successful transaction methods by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, error) {
	const method = "FindYearlyMethodByMerchantSuccess"

	year := req.Year
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearMethodSuccessByMerchantCached(ctx, req); found {
		logSuccess("Successfully found yearly successful transaction methods by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyTransactionMethodByMerchantSuccess(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTransactionMethodsByMerchantSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyMethodByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearMethodSuccessByMerchantCached(ctx, req, res)

	logSuccess("Successfully found yearly successful transaction methods by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, error) {
	const method = "FindMonthlyMethodByMerchantFailed"

	year := req.Year
	month := req.Month
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthMethodFailedByMerchantCached(ctx, req); found {
		logSuccess("Successfully found monthly failed transaction methods by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyTransactionMethodByMerchantFailed(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsByMerchantFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyMethodByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthMethodFailedByMerchantCached(ctx, req, res)

	logSuccess("Successfully found monthly failed transaction methods by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, error) {
	const method = "FindYearlyMethodByMerchantFailed"

	year := req.Year
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearMethodFailedByMerchantCached(ctx, req); found {
		logSuccess("Successfully found yearly failed transaction methods by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyTransactionMethodByMerchantFailed(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*db.GetYearlyTransactionMethodsByMerchantFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyMethodByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearMethodFailedByMerchantCached(ctx, req, res)

	logSuccess("Successfully found yearly failed transaction methods by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantID))

	return res, nil
}
