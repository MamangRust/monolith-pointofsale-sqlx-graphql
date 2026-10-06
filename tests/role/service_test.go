package role_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	role_cache "github.com/MamangRust/monolith-pointofsale-role/cache"
	"github.com/MamangRust/monolith-pointofsale-role/repository"
	"github.com/MamangRust/monolith-pointofsale-role/service"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
	tests "github.com/MamangRust/monolith-pointofsale-test"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type RoleServiceTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	roleService *service.Service
	roleID      int
}

func (s *RoleServiceTestSuite) SetupSuite() {
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
	mencache := role_cache.NewMencache(cacheStore)

	obs, _ := observability.NewObservability("test", log)
	s.roleService = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache,
		Observability: obs,
	})
}

func (s *RoleServiceTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *RoleServiceTestSuite) TestRoleLifecycle() {
	ctx := context.Background()

	// 1. Create
	req := &requests.CreateRoleRequest{
		Name: "Initial Service Role",
	}
	created, err := s.roleService.RoleCommand.CreateRole(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	roleID := int(created.RoleID)

	// 2. FindByID
	found, err := s.roleService.RoleQuery.FindById(ctx, roleID)
	s.Require().NoError(err)
	s.Equal(req.Name, found.RoleName)

	// 3. Update
	updateReq := &requests.UpdateRoleRequest{
		ID:   &roleID,
		Name: "Updated Service Role",
	}
	updated, err := s.roleService.RoleCommand.UpdateRole(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(updateReq.Name, updated.RoleName)

	// 4. FindAll
	_, total, err := s.roleService.RoleQuery.FindAll(ctx, &requests.FindAllRoles{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.roleService.RoleCommand.TrashedRole(ctx, roleID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.roleService.RoleQuery.FindByTrashedRole(ctx, &requests.FindAllRoles{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.roleService.RoleQuery.FindByActiveRole(ctx, &requests.FindAllRoles{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, r := range active {
		s.NotEqual(roleID, int(r.RoleID))
	}

	// 8. Restore
	_, err = s.roleService.RoleCommand.RestoreRole(ctx, roleID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.roleService.RoleCommand.TrashedRole(ctx, roleID)
	s.Require().NoError(err)
	success, err := s.roleService.RoleCommand.DeleteRolePermanent(ctx, roleID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	r1, _ := s.roleService.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{Name: "R1"})
	r2, _ := s.roleService.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{Name: "R2"})

	s.roleService.RoleCommand.TrashedRole(ctx, int(r1.RoleID))
	s.roleService.RoleCommand.TrashedRole(ctx, int(r2.RoleID))

	resRestoreAll, err := s.roleService.RoleCommand.RestoreAllRole(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.roleService.RoleCommand.TrashedRole(ctx, int(r1.RoleID))
	s.roleService.RoleCommand.TrashedRole(ctx, int(r2.RoleID))

	resDeleteAll, err := s.roleService.RoleCommand.DeleteAllRolePermanent(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestRoleServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleServiceTestSuite))
}
