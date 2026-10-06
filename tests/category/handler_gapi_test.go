package category_test

import (
	"context"
	"testing"

	cat_cache "github.com/MamangRust/monolith-pointofsale-category/cache"
	cat_handler "github.com/MamangRust/monolith-pointofsale-category/handler"
	cat_repo "github.com/MamangRust/monolith-pointofsale-category/repository"
	cat_service "github.com/MamangRust/monolith-pointofsale-category/service"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryGapiTestSuite struct {
	tests.BaseTestSuite
	clientQuery           pbcategory.CategoryQueryServiceClient
	clientCommand         pbcategory.CategoryCommandServiceClient
	clientStats           pbcategory.CategoryStatsServiceClient
	clientStatsById       pbcategory.CategoryStatsByIdServiceClient
	clientStatsByMerchant pbcategory.CategoryStatsByMerchantServiceClient
}

func (s *CategoryGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := s.SQLxDB()

	// Category dependencies
	mencache := cat_cache.NewMencache(cacheStore)
	repos := cat_repo.NewRepositories(queries)
	svc := cat_service.NewService(&cat_service.Deps{
		Ctx:           context.Background(),
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := cat_handler.NewHandler(&cat_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// GRPC Server
	server := grpc.NewServer()
	pbcategory.RegisterCategoryQueryServiceServer(server, handler.Category)
	pbcategory.RegisterCategoryCommandServiceServer(server, handler.CategoryCommand)
	pbcategory.RegisterCategoryStatsServiceServer(server, handler.CategoryStats)
	pbcategory.RegisterCategoryStatsByIdServiceServer(server, handler.CategoryStatsById)
	pbcategory.RegisterCategoryStatsByMerchantServiceServer(server, handler.CategoryStatsByMerchant)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.clientQuery = pbcategory.NewCategoryQueryServiceClient(conn)
	s.clientCommand = pbcategory.NewCategoryCommandServiceClient(conn)
	s.clientStats = pbcategory.NewCategoryStatsServiceClient(conn)
	s.clientStatsById = pbcategory.NewCategoryStatsByIdServiceClient(conn)
	s.clientStatsByMerchant = pbcategory.NewCategoryStatsByMerchantServiceClient(conn)
}

func (s *CategoryGapiTestSuite) TestCategoryGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createRes, err := s.clientCommand.Create(ctx, &pbcategory.CreateCategoryRequest{
		Name:        "GAPI Category",
		Description: "Testing via GRPC",
	})
	s.NoError(err)
	s.NotNil(createRes)
	catID := createRes.Data.Id

	// 2. FindById
	getRes, err := s.clientQuery.FindById(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)
	s.Equal("GAPI Category", getRes.Data.Name)

	// 3. FindAll
	allRes, err := s.clientQuery.FindAll(ctx, &pbcategory.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.clientQuery.FindByActive(ctx, &pbcategory.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, err := s.clientCommand.Update(ctx, &pbcategory.UpdateCategoryRequest{
		CategoryId:  catID,
		Name:        "GAPI Category Updated",
		Description: "Updated via GRPC",
	})
	s.NoError(err)
	s.Equal("GAPI Category Updated", updateRes.Data.Name)

	// 6. Trash
	_, err = s.clientCommand.TrashedCategory(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.clientCommand.FindByTrashed(ctx, &pbcategory.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.clientCommand.RestoreCategory(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 9. DeletePermanent
	_, _ = s.clientCommand.TrashedCategory(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	_, err = s.clientCommand.DeleteCategoryPermanent(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 10. RestoreAll
	_, err = s.clientCommand.RestoreAllCategory(ctx, &emptypb.Empty{})
	s.NoError(err)

	// 11. DeleteAll
	_, err = s.clientCommand.DeleteAllCategoryPermanent(ctx, &emptypb.Empty{})
	s.NoError(err)
}

func TestCategoryGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryGapiTestSuite))
}
