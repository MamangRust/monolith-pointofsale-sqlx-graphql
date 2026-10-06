package service

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	mencache "github.com/MamangRust/monolith-pointofsale-product/cache"
	"github.com/MamangRust/monolith-pointofsale-product/repository"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type productQueryService struct {
	mencache               mencache.ProductQueryCache
	productQueryRepository repository.ProductQueryRepository
	logger                 logger.LoggerInterface
	observability          observability.TraceLoggerObservability
}

func NewProductQueryService(
	mencache mencache.ProductQueryCache,
	productQueryRepository repository.ProductQueryRepository,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *productQueryService {
	return &productQueryService{
		mencache:               mencache,
		productQueryRepository: productQueryRepository,
		logger:                 logger,
		observability:          obs,
	}
}

func (s *productQueryService) FindAll(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsRow, *int, error) {
	const method = "FindAll"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProducts(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	products, totalRecords, err := s.productQueryRepository.FindAllProducts(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	s.mencache.SetCachedProducts(ctx, req, products, totalRecords)

	logSuccess("Successfully fetched all products", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return products, totalRecords, nil
}

func (s *productQueryService) FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*db.GetProductsByMerchantRow, *int, error) {
	const method = "FindByMerchant"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProductsByMerchant(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Int("merchant.id", merchantID))
		return data, total, nil
	}

	products, totalRecords, err := s.productQueryRepository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	s.mencache.SetCachedProductsByMerchant(ctx, req, products, totalRecords)

	logSuccess("Successfully fetched all products by merchant", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Int("merchant.id", merchantID))

	return products, totalRecords, nil
}

func (s *productQueryService) FindByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*db.GetProductsByCategoryNameRow, *int, error) {
	const method = "FindByCategory"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	categoryName := req.CategoryName

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.String("category.name", categoryName))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProductsByCategory(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.String("category.name", categoryName))
		return data, total, nil
	}

	products, totalRecords, err := s.productQueryRepository.FindByCategory(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	s.mencache.SetCachedProductsByCategory(ctx, req, products, totalRecords)

	logSuccess("Successfully fetched all products by category", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.String("category.name", categoryName))

	return products, totalRecords, nil
}

func (s *productQueryService) FindById(ctx context.Context, productID int) (*db.Product, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("product.id", productID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedProduct(ctx, productID); found {
		logSuccess("Data found in cache", zap.Int("product.id", productID))
		return data, nil
	}

	product, err := s.productQueryRepository.FindById(ctx, productID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, err, method, span, zap.Error(err))
	}

	s.mencache.SetCachedProduct(ctx, product)

	logSuccess("Successfully fetched product by id", zap.Int("product.id", productID))

	return product, nil
}

func (s *productQueryService) FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsActiveRow, *int, error) {
	const method = "FindByActive"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProductActive(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	products, totalRecords, err := s.productQueryRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	s.mencache.SetCachedProductActive(ctx, req, products, totalRecords)

	logSuccess("Successfully fetched all products", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return products, totalRecords, nil
}

func (s *productQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*db.GetProductsTrashedRow, *int, error) {
	const method = "FindByTrashed"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProductTrashed(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	products, totalRecords, err := s.productQueryRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	s.mencache.SetCachedProductTrashed(ctx, req, products, totalRecords)

	logSuccess("Successfully fetched all products", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return products, totalRecords, nil
}

func (s *productQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
