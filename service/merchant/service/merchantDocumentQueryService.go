package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-merchant/cache"
	"github.com/MamangRust/monolith-pointofsale-merchant/repository"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	merchantdocument_errors "github.com/MamangRust/monolith-pointofsale-shared/errors/merchant_document_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantDocumentQueryDeps struct {
	Cache                 mencache.MerchantDocumentQueryCache
	MerchantDocumentQuery repository.MerchantDocumentQueryRepository
	Logger                logger.LoggerInterface
	Observability         observability.TraceLoggerObservability
}

type merchantDocumentQueryService struct {
	mencache                        mencache.MerchantDocumentQueryCache
	merchantDocumentQueryRepository repository.MerchantDocumentQueryRepository
	logger                          logger.LoggerInterface
	observability                   observability.TraceLoggerObservability
}

func NewMerchantDocumentQueryService(params *merchantDocumentQueryDeps) MerchantDocumentQueryService {
	return &merchantDocumentQueryService{
		mencache:                        params.Cache,
		merchantDocumentQueryRepository: params.MerchantDocumentQuery,
		logger:                          params.Logger,
		observability:                   params.Observability,
	}
}

func (s *merchantDocumentQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, error) {
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

	if data, total, found := s.mencache.GetCachedMerchantDocuments(ctx, req); found {
		logSuccess("Successfully fetched merchant documents from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search))
		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindAllDocuments(ctx, req)
	if err != nil {
		status = "error"
		_, handledErr := sharederrorhandler.HandleError[any](
			s.logger,
			merchantdocument_errors.ErrFindAllMerchantDocumentsFailed,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, handledErr
	}

	s.mencache.SetCachedMerchantDocuments(ctx, req, merchantDocuments, total)
	logSuccess("Successfully fetched merchant documents", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return merchantDocuments, total, nil
}

func (s *merchantDocumentQueryService) FindById(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantDocument.id", documentID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMerchantDocument(ctx, documentID); found {
		logSuccess("Successfully fetched merchant document from cache", zap.Int("merchantDocument.id", documentID))
		return data, nil
	}

	merchantDocument, err := s.merchantDocumentQueryRepository.FindById(ctx, documentID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			merchantdocument_errors.ErrFindMerchantDocumentByIdFailed,
			method,
			span,
			zap.Int("merchantDocument.id", documentID),
		)
	}

	s.mencache.SetCachedMerchantDocument(ctx, merchantDocument)
	logSuccess("Successfully fetched merchant document", zap.Int("merchantDocument.id", documentID))
	return merchantDocument, nil
}

func (s *merchantDocumentQueryService) FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, error) {
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

	if data, total, found := s.mencache.GetCachedMerchantDocumentsActive(ctx, req); found {
		logSuccess("Successfully fetched active merchant documents from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search))
		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		_, handledErr := sharederrorhandler.HandleError[any](
			s.logger,
			merchantdocument_errors.ErrFindActiveMerchantDocumentsFailed,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, handledErr
	}

	s.mencache.SetCachedMerchantDocumentsActive(ctx, req, merchantDocuments, total)
	logSuccess("Successfully fetched active merchant documents", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return merchantDocuments, total, nil
}

func (s *merchantDocumentQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, error) {
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

	if data, total, found := s.mencache.GetCachedMerchantDocumentsTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed merchant documents from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", req.Search))
		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		_, handledErr := sharederrorhandler.HandleError[any](
			s.logger,
			merchantdocument_errors.ErrFindTrashedMerchantDocumentsFailed,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, handledErr
	}

	s.mencache.SetCachedMerchantDocumentsTrashed(ctx, req, merchantDocuments, total)
	logSuccess("Successfully fetched trashed merchant documents", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return merchantDocuments, total, nil
}

func (s *merchantDocumentQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
