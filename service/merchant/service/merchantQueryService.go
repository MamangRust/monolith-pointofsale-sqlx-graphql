package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-merchant/cache"
	"github.com/MamangRust/monolith-pointofsale-merchant/repository"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantQueryDeps struct {
	Cache         mencache.MerchantQueryCache
	MerchantQuery repository.MerchantQueryRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type merchantQueryService struct {
	mencache                mencache.MerchantQueryCache
	merchantQueryRepository repository.MerchantQueryRepository
	logger                  logger.LoggerInterface
	observability           observability.TraceLoggerObservability
}

func NewMerchantQueryService(params *merchantQueryDeps) MerchantQueryService {
	return &merchantQueryService{
		mencache:                params.Cache,
		merchantQueryRepository: params.MerchantQuery,
		logger:                  params.Logger,
		observability:           params.Observability,
	}
}

func (s *merchantQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, *int, error) {
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

	if data, total, found := s.mencache.GetCachedMerchants(ctx, req); found {
		logSuccess("Successfully fetched merchants from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindAllMerchants(ctx, req)
	if err != nil {
		status = "error"
		_, handledErr := sharederrorhandler.HandleError[any](
			s.logger,
			merchant_errors.ErrFailedFindAllMerchants,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, handledErr
	}

	s.mencache.SetCachedMerchants(ctx, req, merchants, totalRecords)
	logSuccess("Successfully fetched merchants", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return merchants, totalRecords, nil
}

func (s *merchantQueryService) FindById(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", merchant_id),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMerchant(ctx, merchant_id); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", merchant_id))
		return data, nil
	}

	res, err := s.merchantQueryRepository.FindById(ctx, merchant_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMerchant(ctx, res)
	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchant_id))
	return res, nil
}

func (s *merchantQueryService) FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsActiveRow, *int, error) {
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

	if data, total, found := s.mencache.GetCachedMerchantActive(ctx, req); found {
		logSuccess("Successfully fetched active merchants from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		_, handledErr := sharederrorhandler.HandleError[any](
			s.logger,
			merchant_errors.ErrFailedFindMerchantsByActive,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, handledErr
	}

	s.mencache.SetCachedMerchantActive(ctx, req, merchants, totalRecords)
	logSuccess("Successfully fetched active merchants", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return merchants, totalRecords, nil
}

func (s *merchantQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsTrashedRow, *int, error) {
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

	if data, total, found := s.mencache.GetCachedMerchantTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed merchants from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		_, handledErr := sharederrorhandler.HandleError[any](
			s.logger,
			merchant_errors.ErrFailedFindMerchantsByTrashed,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, handledErr
	}

	s.mencache.SetCachedMerchantTrashed(ctx, req, merchants, totalRecords)
	logSuccess("Successfully fetched trashed merchants", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return merchants, totalRecords, nil
}

func (s *merchantQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
