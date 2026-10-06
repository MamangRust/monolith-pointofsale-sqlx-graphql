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

func TestCategoryCommand_Create(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	req := &requests.CreateCategoryRequest{
		Name:        "Electronics",
		Description: "Electronic items and gadgets",
	}

	category, err := repos.CategoryCommand.CreateCategory(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, category)
	assert.Equal(t, "Electronics", category.Name)
	assert.Equal(t, "Electronic items and gadgets", *category.Description)
	assert.NotZero(t, category.CategoryID)
}

func TestCategoryQuery_FindByID(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:        "Books",
		Description: "All kinds of books",
	})
	require.NoError(t, err)

	found, err := repos.CategoryQuery.FindById(ctx, int(created.CategoryID))
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.CategoryID, found.CategoryID)
	assert.Equal(t, "Books", found.Name)
}

func TestCategoryCommand_Update(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:        "OldName",
		Description: "Old description",
	})
	require.NoError(t, err)

	updated, err := repos.CategoryCommand.UpdateCategory(ctx, &requests.UpdateCategoryRequest{
		CategoryID:  ptr(int(created.CategoryID)),
		Name:        "NewName",
		Description: "New description",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "NewName", updated.Name)
	assert.Equal(t, "New description", *updated.Description)
}

func TestCategoryQuery_FindAll(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := repos.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:        "Cat1",
		Description: "Desc1",
	})
	require.NoError(t, err)

	_, err = repos.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:        "Cat2",
		Description: "Desc2",
	})
	require.NoError(t, err)

	results, total, err := repos.CategoryQuery.FindAllCategory(ctx, &requests.FindAllCategory{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 2)
	assert.GreaterOrEqual(t, len(results), 2)
}

func TestCategoryCommand_Trash(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:        "TrashMe",
		Description: "Will be trashed",
	})
	require.NoError(t, err)

	trashed, err := repos.CategoryCommand.TrashedCategory(ctx, int(created.CategoryID))
	require.NoError(t, err)
	require.NotNil(t, trashed)
	assert.True(t, trashed.DeletedAt.Valid)
}

func TestCategoryQuery_FindByTrashed(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:        "TrashFind",
		Description: "Find in trash",
	})
	require.NoError(t, err)

	_, err = repos.CategoryCommand.TrashedCategory(ctx, int(created.CategoryID))
	require.NoError(t, err)

	results, total, err := repos.CategoryQuery.FindByTrashed(ctx, &requests.FindAllCategory{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestCategoryCommand_Restore(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:        "RestoreMe",
		Description: "Will be restored",
	})
	require.NoError(t, err)

	_, err = repos.CategoryCommand.TrashedCategory(ctx, int(created.CategoryID))
	require.NoError(t, err)

	restored, err := repos.CategoryCommand.RestoreCategory(ctx, int(created.CategoryID))
	require.NoError(t, err)
	require.NotNil(t, restored)
	assert.False(t, restored.DeletedAt.Valid)
}

func TestCategoryQuery_FindByActive(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:        "ActiveCat",
		Description: "Active category",
	})
	require.NoError(t, err)

	results, total, err := repos.CategoryQuery.FindByActive(ctx, &requests.FindAllCategory{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)

	found := false
	for _, r := range results {
		if int(r.CategoryID) == int(created.CategoryID) {
			found = true
			break
		}
	}
	assert.True(t, found, "created category should be in active results")
}

func TestCategoryCommand_DeletePermanent(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.CategoryCommand.CreateCategory(ctx, &requests.CreateCategoryRequest{
		Name:        "PermanentDelete",
		Description: "Will be permanently deleted",
	})
	require.NoError(t, err)

	_, err = repos.CategoryCommand.TrashedCategory(ctx, int(created.CategoryID))
	require.NoError(t, err)

	deleted, err := repos.CategoryCommand.DeleteCategoryPermanently(ctx, int(created.CategoryID))
	require.NoError(t, err)
	assert.True(t, deleted)

	_, err = repos.CategoryQuery.FindById(ctx, int(created.CategoryID))
	assert.Error(t, err)
}
