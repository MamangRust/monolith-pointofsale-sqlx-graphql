package product_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-product/repository"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-pointofsale-test"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type ProductRepositoryTestSuite struct {
	suite.Suite
	ts        *tests.TestSuite
	repo      *repository.Repositories
	productID int
}

func (s *ProductRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)

	queries := pool
	// DB-only suite: category/merchant adapters are never invoked here, so nil
	// clients are safe (guard nil = passthrough, calls happen only if used).
	s.repo = repository.NewRepositories(queries, nil, nil)

	// Seed a merchant and category for product tests
	var userID, categoryID int
	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ($1, $2, $3, $4, 'test-verify', true) RETURNING user_id`,
		"Prod", "Repo", "prod.repo@example.com", "password123",
	).Scan(&userID)
	s.Require().NoError(err)

	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING category_id`,
		"Test Category", "Category for product tests",
	).Scan(&categoryID)
	s.Require().NoError(err)

	err = pool.QueryRowContext(s.ts.Ctx,
		`INSERT INTO merchants (user_id, name, description, address, contact_email, contact_phone, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING merchant_id`,
		userID, "Test Merchant", "Desc", "Addr", "pm@example.com", "123", "active",
	).Scan(&userID)
	s.Require().NoError(err)
}

func (s *ProductRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *ProductRepositoryTestSuite) Test1_CreateProduct() {
	ctx := context.Background()

	req := &requests.CreateProductRequest{
		MerchantID:   1,
		CategoryID:   1,
		Name:         "Test Product",
		Description:  "Product description",
		Price:        100,
		CountInStock: 50,
		Brand:        "Test Brand",
		Weight:       1,
		ImageProduct: "test.jpg",
	}

	res, err := s.repo.ProductCommand.CreateProduct(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Name, res.Name)
	s.productID = int(res.ProductID)
}

func (s *ProductRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	found, err := s.repo.ProductQuery.FindById(ctx, s.productID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.productID, int(found.ProductID))
}

func (s *ProductRepositoryTestSuite) Test3_UpdateProduct() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	req := &requests.UpdateProductRequest{
		ProductID:    &s.productID,
		MerchantID:   1,
		CategoryID:   1,
		Name:         "Updated Product",
		Description:  "Updated description",
		Price:        200,
		CountInStock: 100,
		Brand:        "Updated Brand",
		Weight:       2,
		ImageProduct: "updated.jpg",
	}

	res, err := s.repo.ProductCommand.UpdateProduct(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("Updated Product", res.Name)
}

func (s *ProductRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	// Trash
	trashed, err := s.repo.ProductCommand.TrashedProduct(ctx, s.productID)
	s.NoError(err)
	s.NotNil(trashed)

	// Restore
	restored, err := s.repo.ProductCommand.RestoreProduct(ctx, s.productID)
	s.NoError(err)
	s.NotNil(restored)

	// Verify restored
	found, err := s.repo.ProductQuery.FindById(ctx, s.productID)
	s.NoError(err)
	s.NotNil(found)
}

func (s *ProductRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	// Must be trashed first for permanent delete
	_, err := s.repo.ProductCommand.TrashedProduct(ctx, s.productID)
	s.NoError(err)

	success, err := s.repo.ProductCommand.DeleteProductPermanent(ctx, s.productID)
	s.NoError(err)
	s.True(success)
}

func TestProductRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductRepositoryTestSuite))
}
