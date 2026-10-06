package category_test

import (
	"context"
	"testing"

	cat_cache "github.com/MamangRust/monolith-pointofsale-category/cache"
	"github.com/MamangRust/monolith-pointofsale-category/repository"
	"github.com/MamangRust/monolith-pointofsale-category/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	tests "github.com/MamangRust/monolith-pointofsale-test"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type CategoryServiceTestSuite struct {
	suite.Suite
	ts              *tests.TestSuite
	categoryService *service.Service
	categoryID      int
}

func (s *CategoryServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := sqlx.ConnectContext(s.ts.Ctx, "pgx", s.ts.DBURL)
	s.Require().NoError(err)

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(opts)

	queries := pool
	repos := repository.NewRepositories(queries)

	log, _ := logger.NewLogger("test", nil)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)
	mencache := cat_cache.NewMencache(cacheStore)

	obs, _ := observability.NewObservability("test", log)
	s.categoryService = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache,
		Observability: obs,
	})
}

func (s *CategoryServiceTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *CategoryServiceTestSuite) TestCategoryLifecycle() {
	ctx := context.Background()

	// 1. Create
	req := &requests.CreateCategoryRequest{
		Name:        "Initial Category",
		Description: "Testing service layer",
	}
	created, err := s.categoryService.CategoryCommand.CreateCategory(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	catID := int(created.CategoryID)

	// 2. FindByID
	found, err := s.categoryService.CategoryQuery.FindById(ctx, catID)
	s.Require().NoError(err)
	s.Equal(req.Name, found.Name)

	// 3. Update
	updateReq := &requests.UpdateCategoryRequest{
		CategoryID:  &catID,
		Name:        "Updated Category Name",
		Description: "Updated description",
	}
	updated, err := s.categoryService.CategoryCommand.UpdateCategory(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(updateReq.Name, updated.Name)

	// 4. FindAll
	_, total, err := s.categoryService.CategoryQuery.FindAll(ctx, &requests.FindAllCategory{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.categoryService.CategoryCommand.TrashedCategory(ctx, catID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.categoryService.CategoryQuery.FindByTrashed(ctx, &requests.FindAllCategory{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.categoryService.CategoryQuery.FindByActive(ctx, &requests.FindAllCategory{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, c := range active {
		s.NotEqual(catID, int(c.CategoryID))
	}

	// 8. Restore
	_, err = s.categoryService.CategoryCommand.RestoreCategory(ctx, catID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.categoryService.CategoryCommand.TrashedCategory(ctx, catID)
	s.Require().NoError(err)
	success, err := s.categoryService.CategoryCommand.DeleteCategoryPermanent(ctx, catID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	c1, _ := s.categoryService.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{Name: "C1", Description: "D1"})
	c2, _ := s.categoryService.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{Name: "C2", Description: "D2"})

	s.categoryService.CategoryCommand.TrashedCategory(ctx, int(c1.CategoryID))
	s.categoryService.CategoryCommand.TrashedCategory(ctx, int(c2.CategoryID))

	resRestoreAll, err := s.categoryService.CategoryCommand.RestoreAllCategories(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.categoryService.CategoryCommand.TrashedCategory(ctx, int(c1.CategoryID))
	s.categoryService.CategoryCommand.TrashedCategory(ctx, int(c2.CategoryID))

	resDeleteAll, err := s.categoryService.CategoryCommand.DeleteAllCategoriesPermanent(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestCategoryServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryServiceTestSuite))
}
