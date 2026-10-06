package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-cashier/cache"
	"github.com/MamangRust/monolith-pointofsale-cashier/repository"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type cashierQueryDeps struct {
	Cache         mencache.CashierQueryCache
	CashierQuery  repository.CashierQueryRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type cashierQueryService struct {
	mencache      mencache.CashierQueryCache
	cashierQuery  repository.CashierQueryRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewCashierQueryService(params *cashierQueryDeps) CashierQueryService {
	return &cashierQueryService{
		mencache:      params.Cache,
		cashierQuery:  params.CashierQuery,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *cashierQueryService) FindAll(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersRow, *int, error) {
	const method = "FindAll"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCashiersCache(ctx, req); found {
		logSuccess("Successfully fetched cashier from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	cashier, totalRecords, err := s.cashierQuery.FindAllCashiers(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[[]*db.GetCashiersRow](
			s.logger,
			cashier_errors.ErrFailedFindAllCashiers,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, err
	}

	s.mencache.SetCachedCashiersCache(ctx, req, cashier, totalRecords)
	logSuccess("Successfully fetched cashier", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return cashier, totalRecords, nil
}

func (s *cashierQueryService) FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*db.GetCashiersByMerchantRow, *int, error) {
	const method = "FindByMerchant"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
		attribute.Int("merchant.id", req.MerchantID),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCashiersByMerchant(ctx, req); found {
		logSuccess("Successfully fetched cashier from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	cashier, totalRecords, err := s.cashierQuery.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[[]*db.GetCashiersByMerchantRow](
			s.logger,
			cashier_errors.ErrFailedFindCashierByMerchant,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, err
	}

	s.mencache.SetCachedCashiersByMerchant(ctx, req, cashier, totalRecords)
	logSuccess("Successfully fetched cashier", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return cashier, totalRecords, nil
}

func (s *cashierQueryService) FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersActiveRow, *int, error) {
	const method = "FindByActive"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCashiersActive(ctx, req); found {
		logSuccess("Successfully fetched cashier from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	cashier, totalRecords, err := s.cashierQuery.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[[]*db.GetCashiersActiveRow](
			s.logger,
			cashier_errors.ErrFailedFindCashierByActive,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, err
	}

	s.mencache.SetCachedCashiersActive(ctx, req, cashier, totalRecords)
	logSuccess("Successfully fetched cashier", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return cashier, totalRecords, nil
}

func (s *cashierQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersTrashedRow, *int, error) {
	const method = "FindByTrashed"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCashiersTrashed(ctx, req); found {
		logSuccess("Successfully fetched cashier from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	cashier, totalRecords, err := s.cashierQuery.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[[]*db.GetCashiersTrashedRow](
			s.logger,
			cashier_errors.ErrFailedFindCashierByTrashed,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, err
	}

	s.mencache.SetCachedCashiersTrashed(ctx, req, cashier, totalRecords)
	logSuccess("Successfully fetched cashier", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return cashier, totalRecords, nil
}

func (s *cashierQueryService) FindById(ctx context.Context, cashier_id int) (*db.Cashier, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("cashier.id", cashier_id),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedCashier(ctx, cashier_id); found {
		logSuccess("Successfully fetched cashier from cache", zap.Int("cashier.id", cashier_id))
		return data, nil
	}

	cashier, err := s.cashierQuery.FindById(ctx, cashier_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Cashier](
			s.logger,
			cashier_errors.ErrFailedFindCashierById,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedCashier(ctx, cashier)
	logSuccess("Successfully fetched cashier", zap.Int("cashier.id", cashier_id))
	return cashier, nil
}

func (s *cashierQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
