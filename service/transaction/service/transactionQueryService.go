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

type transactionQueryService struct {
	mencache                   mencache.TransactionQueryCache
	transactionQueryRepository repository.TransactionQueryRepository
	logger                     logger.LoggerInterface
	observability              observability.TraceLoggerObservability
}

func NewTransactionQueryService(
	mencache mencache.TransactionQueryCache,
	transactionQueryRepository repository.TransactionQueryRepository,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *transactionQueryService {
	return &transactionQueryService{
		mencache:                   mencache,
		transactionQueryRepository: transactionQueryRepository,
		logger:                     logger,
		observability:              obs,
	}
}

func (s *transactionQueryService) FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsRow, *int, error) {
	const method = "FindAllTransactions"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionsCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactions(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetTransactionsRow](
			s.logger,
			transaction_errors.ErrFailedFindAllTransactions.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedTransactionsCache(ctx, req, transactions, totalRecords)

	logSuccess("Successfully fetched all transactions", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return transactions, totalRecords, nil
}

func (s *transactionQueryService) FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*db.GetTransactionByMerchantRow, *int, error) {
	const method = "FindByMerchant"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.Int("merchant_id", merchantID))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionByMerchant(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetTransactionByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindTransactionsByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedTransactionByMerchant(ctx, req, transactions, totalRecords)

	logSuccess("Successfully fetched all transactions", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return transactions, totalRecords, nil
}

func (s *transactionQueryService) FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsActiveRow, *int, error) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionActiveCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetTransactionsActiveRow](
			s.logger,
			transaction_errors.ErrFailedFindTransactionsByActive.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedTransactionActiveCache(ctx, req, transactions, totalRecords)

	logSuccess("Successfully fetched active transactions", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return transactions, totalRecords, nil
}

func (s *transactionQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsTrashedRow, *int, error) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionTrashedCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandlerErrorPagination[[]*db.GetTransactionsTrashedRow](
			s.logger,
			transaction_errors.ErrFailedFindTransactionsByTrashed.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedTransactionTrashedCache(ctx, req, transactions, totalRecords)

	logSuccess("Successfully fetched trashed transactions", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return transactions, totalRecords, nil
}

func (s *transactionQueryService) FindById(ctx context.Context, transactionID int) (*db.Transaction, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTransactionCache(ctx, transactionID); found {
		logSuccess("Successfully fetched transaction from cache", zap.Int("transaction.id", transactionID))
		return data, nil
	}

	transaction, err := s.transactionQueryRepository.FindById(ctx, transactionID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedFindTransactionById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedTransactionCache(ctx, transaction)

	logSuccess("Successfully fetched transaction", zap.Int("transaction.id", transactionID))

	return transaction, nil
}

func (s *transactionQueryService) FindByOrderId(ctx context.Context, orderID int) (*db.Transaction, error) {
	const method = "FindByOrderId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order.id", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTransactionByOrderId(ctx, orderID); found {
		logSuccess("Successfully fetched transaction from cache", zap.Int("order.id", orderID))
		return data, nil
	}

	transaction, err := s.transactionQueryRepository.FindByOrderId(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedFindTransactionByOrderId.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedTransactionByOrderId(ctx, orderID, transaction)

	logSuccess("Successfully fetched transaction", zap.Int("order.id", orderID))

	return transaction, nil
}

func (s *transactionQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
