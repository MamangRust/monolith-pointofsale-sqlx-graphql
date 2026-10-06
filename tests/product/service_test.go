package product_test

import (
	"context"
	"testing"

	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	prod_cache "github.com/MamangRust/monolith-pointofsale-product/cache"
	"github.com/MamangRust/monolith-pointofsale-product/repository"
	"github.com/MamangRust/monolith-pointofsale-product/service"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	tests "github.com/MamangRust/monolith-pointofsale-test"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type ProductServiceTestSuite struct {
	tests.BaseTestSuite
	svc       *service.Service
	productID int
}

func (s *ProductServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()

	pool, err := sqlx.ConnectContext(s.Ctx, "pgx", s.DBURL())
	s.Require().NoError(err)

	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)

	queries := pool

	log, _ := logger.NewLogger("test", nil)
	mencache := prod_cache.NewMencache(cacheStore)

	obs, _ := observability.NewObservability("test", log)

	repos := repository.NewRepositories(queries,
		pbcategory.NewCategoryQueryServiceClient(s.Conns["category"]),
		pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)

	s.svc = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache,
		Observability: obs,
	})

	// Seed a merchant and category for product tests
	var userID, merchantID, categoryID int
	err = pool.QueryRowContext(s.Ctx,
		`INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ($1, $2, $3, $4, 'test-verify', true) RETURNING user_id`,
		"Prod", "Svc", "prod.svc@example.com", "password123",
	).Scan(&userID)
	s.Require().NoError(err)

	err = pool.QueryRowContext(s.Ctx,
		`INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING category_id`,
		"Svc Category", "Category for service tests",
	).Scan(&categoryID)
	s.Require().NoError(err)

	err = pool.QueryRowContext(s.Ctx,
		`INSERT INTO merchants (user_id, name, description, address, contact_email, contact_phone, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING merchant_id`,
		userID, "Svc Merchant", "Desc", "Addr", "ps@example.com", "123", "active",
	).Scan(&merchantID)
	s.Require().NoError(err)
}

func (s *ProductServiceTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *ProductServiceTestSuite) TestProductLifecycle() {
	ctx := context.Background()

	// 1. Create
	req := &requests.CreateProductRequest{
		MerchantID:   1,
		CategoryID:   1,
		Name:         "Test Product",
		Description:  "Product description",
		Price:        10000,
		CountInStock: 100,
		Brand:        "Test Brand",
		Weight:       1000,
		ImageProduct: "test.jpg",
	}
	created, err := s.svc.ProductCommand.CreateProduct(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	productID := int(created.ProductID)

	// 2. FindByID
	found, err := s.svc.ProductQuery.FindById(ctx, productID)
	s.Require().NoError(err)
	s.Equal(req.Name, found.Name)

	// 3. Update
	updateReq := &requests.UpdateProductRequest{
		ProductID:    &productID,
		MerchantID:   1,
		CategoryID:   1,
		Name:         "Updated Product Name",
		Description:  "Updated description",
		Price:        20000,
		CountInStock: 90,
		Brand:        "Updated Brand",
		Weight:       1000,
		ImageProduct: "updated.jpg",
	}
	updated, err := s.svc.ProductCommand.UpdateProduct(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(updateReq.Name, updated.Name)

	// 4. FindAll
	_, total, err := s.svc.ProductQuery.FindAll(ctx, &requests.FindAllProducts{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.svc.ProductCommand.TrashProduct(ctx, productID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.svc.ProductQuery.FindByTrashed(ctx, &requests.FindAllProducts{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.svc.ProductQuery.FindByActive(ctx, &requests.FindAllProducts{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, p := range active {
		s.NotEqual(productID, int(p.ProductID))
	}

	// 8. Restore
	_, err = s.svc.ProductCommand.RestoreProduct(ctx, productID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.svc.ProductCommand.TrashProduct(ctx, productID)
	s.Require().NoError(err)
	success, err := s.svc.ProductCommand.DeleteProductPermanent(ctx, productID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	p1, _ := s.svc.ProductCommand.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID: 1, CategoryID: 1, Name: "P1", Description: "D1", Price: 100,
		CountInStock: 10, Brand: "B", Weight: 1, ImageProduct: "i1.jpg",
	})
	p2, _ := s.svc.ProductCommand.CreateProduct(ctx, &requests.CreateProductRequest{
		MerchantID: 1, CategoryID: 1, Name: "P2", Description: "D2", Price: 200,
		CountInStock: 20, Brand: "B", Weight: 2, ImageProduct: "i2.jpg",
	})

	s.svc.ProductCommand.TrashProduct(ctx, int(p1.ProductID))
	s.svc.ProductCommand.TrashProduct(ctx, int(p2.ProductID))

	resRestoreAll, err := s.svc.ProductCommand.RestoreAllProducts(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.ProductCommand.TrashProduct(ctx, int(p1.ProductID))
	s.svc.ProductCommand.TrashProduct(ctx, int(p2.ProductID))

	resDeleteAll, err := s.svc.ProductCommand.DeleteAllProductsPermanent(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestProductServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductServiceTestSuite))
}
