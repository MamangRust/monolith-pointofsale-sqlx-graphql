package service

import (
	"context"
	"os"

	categoryadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/category"
	merchantadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/merchant"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-pkg/utils"
	mencache "github.com/MamangRust/monolith-pointofsale-product/cache"
	"github.com/MamangRust/monolith-pointofsale-product/repository"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-pointofsale-shared/errorhandler"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/product_errors"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type productCommandService struct {
	mencache                 mencache.ProductCommandCache
	categoryRepository       categoryadapter.Repository
	merchantRepository       merchantadapter.Repository
	productQueryRepository   repository.ProductQueryRepository
	productCommandRepository repository.ProductCommandRepository
	logger                   logger.LoggerInterface
	observability            observability.TraceLoggerObservability
}

func NewProductCommandService(
	mencache mencache.ProductCommandCache,
	categoryRepository categoryadapter.Repository,
	merchantRepository merchantadapter.Repository,
	productQueryRepository repository.ProductQueryRepository,
	productCommandRepository repository.ProductCommandRepository,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *productCommandService {
	return &productCommandService{
		mencache:                 mencache,
		categoryRepository:       categoryRepository,
		merchantRepository:       merchantRepository,
		productQueryRepository:   productQueryRepository,
		productCommandRepository: productCommandRepository,
		logger:                   logger,
		observability:            obs,
	}
}

func (s *productCommandService) CreateProduct(ctx context.Context, req *requests.CreateProductRequest) (*db.Product, error) {
	const method = "CreateProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("product.name", req.Name), attribute.Int("product.category_id", req.CategoryID), attribute.Int("product.merchant_id", req.MerchantID))

	defer func() {
		end(status)
	}()

	_, err := s.categoryRepository.FindById(ctx, req.CategoryID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, err, method, span, zap.Error(err))
	}

	_, err = s.merchantRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, err, method, span, zap.Error(err))
	}

	slug := utils.GenerateSlug(req.Name)
	req.SlugProduct = &slug

	product, err := s.productCommandRepository.CreateProduct(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, err, method, span, zap.Error(err))
	}

	logSuccess("Successfully created product", zap.Int32("productID", product.ProductID), zap.Bool("success", true))

	return product, nil
}

func (s *productCommandService) UpdateProduct(ctx context.Context, req *requests.UpdateProductRequest) (*db.Product, error) {
	const method = "UpdateProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.categoryRepository.FindById(ctx, req.CategoryID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, err, method, span, zap.Error(err))
	}

	_, err = s.merchantRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, err, method, span, zap.Error(err))
	}

	slug := utils.GenerateSlug(req.Name)
	req.SlugProduct = &slug

	product, err := s.productCommandRepository.UpdateProduct(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, err, method, span, zap.Error(err))
	}

	s.mencache.DeleteCachedProduct(ctx, *req.ProductID)

	logSuccess("Successfully updated product", zap.Int("product.id", *req.ProductID), zap.Bool("success", true))

	return product, nil
}

func (s *productCommandService) AdjustProductStock(ctx context.Context, productID int, delta int) (*db.Product, error) {
	const method = "AdjustProductStock"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("product.id", productID), attribute.Int("delta", delta))

	defer func() {
		end(status)
	}()

	if productID <= 0 || delta == 0 {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, product_errors.ErrFailedCountStock, method, span)
	}

	product, err := s.productCommandRepository.AdjustProductStock(ctx, productID, delta)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, err, method, span, zap.Error(err))
	}

	s.mencache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully adjusted product stock", zap.Int("product.id", productID), zap.Int("delta", delta))

	return product, nil
}

func (s *productCommandService) TrashProduct(ctx context.Context, productID int) (*db.Product, error) {
	const method = "TrashProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	product, err := s.productCommandRepository.TrashedProduct(ctx, productID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, err, method, span, zap.Error(err))
	}

	s.mencache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully trashed product", zap.Int("product.id", productID), zap.Bool("success", true))

	return product, nil
}

func (s *productCommandService) RestoreProduct(ctx context.Context, productID int) (*db.Product, error) {
	const method = "RestoreProduct"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	product, err := s.productCommandRepository.RestoreProduct(ctx, productID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Product](s.logger, err, method, span, zap.Error(err))
	}

	s.mencache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully restored product", zap.Int("product.id", productID), zap.Bool("success", true))

	return product, nil
}

func (s *productCommandService) DeleteProductPermanent(ctx context.Context, productID int) (bool, error) {
	const method = "DeleteProductPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.productQueryRepository.FindByIdTrashed(ctx, productID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.Error(err))
	}

	if res.ImageProduct != nil && *res.ImageProduct != "" {
		err := os.Remove(*res.ImageProduct)
		if err != nil {
			if !os.IsNotExist(err) {
				status = "error"
				return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.String("image_path", *res.ImageProduct), zap.Error(err))
			}
			s.logger.Error("Failed to delete product image (not exist)",
				zap.String("image_path", *res.ImageProduct),
				zap.Error(err))
		} else {
			s.logger.Debug("Successfully deleted product image",
				zap.String("image_path", *res.ImageProduct))
		}
	}

	_, err = s.productCommandRepository.DeleteProductPermanent(ctx, productID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.Error(err))
	}

	s.mencache.DeleteCachedProduct(ctx, productID)

	logSuccess("Product deleted permanently successfully", zap.Int("product.id", productID), zap.Bool("success", true))

	return true, nil
}

func (s *productCommandService) RestoreAllProducts(ctx context.Context) (bool, error) {
	const method = "RestoreAllProducts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.productCommandRepository.RestoreAllProducts(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.Error(err))
	}

	s.mencache.DeleteCachedProductAllCache(ctx)
	logSuccess("All trashed products restored successfully", zap.Bool("success", success))

	return success, nil
}

func (s *productCommandService) DeleteAllProductsPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllProductsPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.productCommandRepository.DeleteAllProductPermanent(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.Error(err))
	}

	s.mencache.DeleteCachedProductAllCache(ctx)
	logSuccess("All trashed products deleted permanently successfully", zap.Bool("success", success))

	return success, nil
}
