package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) (*sqlx.DB, *Repositories, *sqlx.DB, func()) {
	t.Helper()

	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := sqlx.Connect("pgx", connStr)
	require.NoError(t, err)

	runMigrations(t, ctx, pool, "../../../pkg/database/migrations")

	repos := NewRepositories(pool)

	cleanup := func() {
		pool.Close()
		pgContainer.Terminate(ctx)
	}

	return pool, repos, pool, cleanup
}

func runMigrations(t *testing.T, ctx context.Context, pool *sqlx.DB, migrationDir string) {
	t.Helper()

	migrations := []string{
		"20250204052452_add_pgcrypto_and_uuid_columns.sql",
		"20250204052518_create_users.sql",
		"20250204052555_create_role.sql",
		"20250204052605_create_user_role.sql",
		"20250204052613_create_refreshtoken.sql",
		"20250204052614_create_reset_token.sql",
		"20250206050402_create_merchants.sql",
		"20250206050403_create_merchant_documents.sql",
		"20250206050448_create_cashiers.sql",
		"20250206050455_create_categories.sql",
		"20250206050505_create_products.sql",
		"20250206050519_create_orders.sql",
		"20250206050525_create_order_items.sql",
		"20250206050554_create_transactions.sql",
	}

	for _, file := range migrations {
		path := migrationDir + "/" + file
		sql, err := os.ReadFile(path)
		require.NoError(t, err)

		_, err = pool.ExecContext(ctx, string(sql))
		if err != nil {
			t.Logf("Warning executing %s: %v", file, err)
		}
	}
}

func ptr[T any](v T) *T {
	return &v
}

func TestRoleCommand_Create(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	role, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{
		Name: "admin",
	})
	require.NoError(t, err)
	require.NotNil(t, role)
	assert.Equal(t, "admin", role.RoleName)
	assert.NotZero(t, role.RoleID)
}

func TestRoleQuery_FindByID(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{
		Name: "editor",
	})
	require.NoError(t, err)

	found, err := repos.RoleQuery.FindById(ctx, int(created.RoleID))
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.RoleID, found.RoleID)
	assert.Equal(t, "editor", found.RoleName)
}

func TestRoleQuery_FindByName(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{
		Name: "viewer",
	})
	require.NoError(t, err)

	found, err := repos.RoleQuery.FindByName(ctx, "viewer")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "viewer", found.RoleName)
}

func TestRoleCommand_Update(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{
		Name: "oldrole",
	})
	require.NoError(t, err)

	updated, err := repos.RoleCommand.UpdateRole(ctx, &requests.UpdateRoleRequest{
		ID:   ptr(int(created.RoleID)),
		Name: "newrole",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "newrole", updated.RoleName)
}

func TestRoleQuery_FindAll(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{Name: "role1"})
	require.NoError(t, err)
	_, err = repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{Name: "role2"})
	require.NoError(t, err)

	results, total, err := repos.RoleQuery.FindAllRoles(ctx, &requests.FindAllRoles{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 2)
	assert.GreaterOrEqual(t, len(results), 2)
}

func TestRoleCommand_Trash(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{Name: "temporary"})
	require.NoError(t, err)

	trashed, err := repos.RoleCommand.TrashedRole(ctx, int(created.RoleID))
	require.NoError(t, err)
	require.NotNil(t, trashed)
	assert.True(t, trashed.DeletedAt.Valid)
}

func TestRoleQuery_FindByTrashed(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{Name: "trashrole"})
	require.NoError(t, err)

	_, err = repos.RoleCommand.TrashedRole(ctx, int(created.RoleID))
	require.NoError(t, err)

	results, total, err := repos.RoleQuery.FindByTrashedRole(ctx, &requests.FindAllRoles{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestRoleCommand_Restore(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{Name: "restorerole"})
	require.NoError(t, err)

	_, err = repos.RoleCommand.TrashedRole(ctx, int(created.RoleID))
	require.NoError(t, err)

	restored, err := repos.RoleCommand.RestoreRole(ctx, int(created.RoleID))
	require.NoError(t, err)
	require.NotNil(t, restored)
	assert.False(t, restored.DeletedAt.Valid)
}

func TestRoleQuery_FindByActive(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{Name: "activerole"})
	require.NoError(t, err)

	results, total, err := repos.RoleQuery.FindByActiveRole(ctx, &requests.FindAllRoles{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)

	found := false
	for _, r := range results {
		if r.RoleID == created.RoleID {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestRoleCommand_DeletePermanent(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{Name: "deletepermanent"})
	require.NoError(t, err)

	_, err = repos.RoleCommand.TrashedRole(ctx, int(created.RoleID))
	require.NoError(t, err)

	deleted, err := repos.RoleCommand.DeleteRolePermanent(ctx, int(created.RoleID))
	require.NoError(t, err)
	assert.True(t, deleted)

	_, err = repos.RoleQuery.FindById(ctx, int(created.RoleID))
	assert.Error(t, err)
}

func TestRoleQuery_FindByUserId(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Seed a user and a role, then assign the role
	var userID int
	err := pool.QueryRowContext(ctx, `INSERT INTO users (firstname, lastname, email, password, verification_code)
		VALUES ('Test', 'User', 'testuser@example.com', 'password', 'CODE') RETURNING user_id`).Scan(&userID)
	require.NoError(t, err)

	created, err := repos.RoleCommand.CreateRole(ctx, &requests.CreateRoleRequest{Name: "userole"})
	require.NoError(t, err)

	_, err = pool.ExecContext(ctx, `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, userID, created.RoleID)
	require.NoError(t, err)

	roles, err := repos.RoleQuery.FindByUserId(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, roles)
	assert.GreaterOrEqual(t, len(roles), 1)
	assert.Equal(t, "userole", roles[0].RoleName)
}
