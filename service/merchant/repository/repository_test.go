package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
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

	repos := NewRepositories(pool, nil)

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

func seedMerchantUser(t *testing.T, ctx context.Context, pool *sqlx.DB) int {
	t.Helper()

	var userID int
	err := pool.QueryRowContext(ctx, `INSERT INTO users (firstname, lastname, email, password, verification_code) 
		VALUES ('Merchant', 'Owner', 'merchantowner@example.com', 'password', 'CODE') RETURNING user_id`).Scan(&userID)
	require.NoError(t, err)
	return userID
}

func TestMerchantCommand_Create(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)

	merchant, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:       userID,
		Name:         "Test Merchant",
		Description:  "A test merchant",
		Address:      "123 Test St",
		ContactEmail: "test@merchant.com",
		ContactPhone: "1234567890",
		Status:       "active",
	})
	require.NoError(t, err)
	require.NotNil(t, merchant)
	assert.Equal(t, "Test Merchant", merchant.Name)
	assert.Equal(t, int32(userID), merchant.UserID)
	assert.NotZero(t, merchant.MerchantID)
}

func TestMerchantQuery_FindByID(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)

	created, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "Find Merchant", Description: "Desc",
		Address: "Addr", ContactEmail: "find@m.com", ContactPhone: "111", Status: "active",
	})
	require.NoError(t, err)

	found, err := repos.MerchantQuery.FindById(ctx, int(created.MerchantID))
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.MerchantID, found.MerchantID)
	assert.Equal(t, "Find Merchant", found.Name)
}

func TestMerchantCommand_Update(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)

	created, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "Old Merchant", Description: "Old",
		Address: "Old Addr", ContactEmail: "old@m.com", ContactPhone: "000", Status: "inactive",
	})
	require.NoError(t, err)

	updated, err := repos.MerchantCommand.UpdateMerchant(ctx, &requests.UpdateMerchantRequest{
		MerchantID:   ptr(int(created.MerchantID)),
		UserID:       userID,
		Name:         "New Merchant",
		Description:  "New",
		Address:      "New Addr",
		ContactEmail: "new@m.com",
		ContactPhone: "999",
		Status:       "active",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "New Merchant", updated.Name)
	assert.Equal(t, "active", updated.Status)
}

func TestMerchantCommand_UpdateStatus(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	created, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "Status Merchant", Description: "Desc",
		Address: "Addr", ContactEmail: "status@m.com", ContactPhone: "555", Status: "inactive",
	})
	require.NoError(t, err)

	updated, err := repos.MerchantCommand.UpdateMerchantStatus(ctx, &requests.UpdateMerchantStatusRequest{
		MerchantID: ptr(int(created.MerchantID)),
		Status:     "suspended",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "suspended", updated.Status)
}

func TestMerchantQuery_FindAll(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)

	_, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "M1", Description: "D1",
		Address: "A1", ContactEmail: "m1@m.com", ContactPhone: "111", Status: "active",
	})
	require.NoError(t, err)
	_, err = repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "M2", Description: "D2",
		Address: "A2", ContactEmail: "m2@m.com", ContactPhone: "222", Status: "active",
	})
	require.NoError(t, err)

	results, total, err := repos.MerchantQuery.FindAllMerchants(ctx, &requests.FindAllMerchants{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 2)
	assert.GreaterOrEqual(t, len(results), 2)
}

func TestMerchantCommand_Trash(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	created, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "TrashMerchant", Description: "D",
		Address: "A", ContactEmail: "t@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	trashed, err := repos.MerchantCommand.TrashedMerchant(ctx, int(created.MerchantID))
	require.NoError(t, err)
	require.NotNil(t, trashed)
	assert.True(t, trashed.DeletedAt.Valid)
}

func TestMerchantQuery_FindByTrashed(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	created, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "TrashFindM", Description: "D",
		Address: "A", ContactEmail: "tf@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	_, err = repos.MerchantCommand.TrashedMerchant(ctx, int(created.MerchantID))
	require.NoError(t, err)

	results, total, err := repos.MerchantQuery.FindByTrashed(ctx, &requests.FindAllMerchants{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestMerchantCommand_Restore(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	created, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "RestoreM", Description: "D",
		Address: "A", ContactEmail: "r@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	_, err = repos.MerchantCommand.TrashedMerchant(ctx, int(created.MerchantID))
	require.NoError(t, err)

	restored, err := repos.MerchantCommand.RestoreMerchant(ctx, int(created.MerchantID))
	require.NoError(t, err)
	require.NotNil(t, restored)
	assert.False(t, restored.DeletedAt.Valid)
}

func TestMerchantQuery_FindByActive(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	created, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "ActiveM", Description: "D",
		Address: "A", ContactEmail: "a@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	results, total, err := repos.MerchantQuery.FindByActive(ctx, &requests.FindAllMerchants{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)

	found := false
	for _, r := range results {
		if int(r.MerchantID) == int(created.MerchantID) {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestMerchantCommand_DeletePermanent(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	created, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "PermDeleteM", Description: "D",
		Address: "A", ContactEmail: "pd@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	_, err = repos.MerchantCommand.TrashedMerchant(ctx, int(created.MerchantID))
	require.NoError(t, err)

	deleted, err := repos.MerchantCommand.DeleteMerchantPermanent(ctx, int(created.MerchantID))
	require.NoError(t, err)
	assert.True(t, deleted)

	_, err = repos.MerchantQuery.FindById(ctx, int(created.MerchantID))
	assert.Error(t, err)
}

// ----- Merchant Document Tests -----

func TestMerchantDocumentCommand_Create(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	createdM, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "DocMerchant", Description: "D",
		Address: "A", ContactEmail: "doc@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	doc, err := repos.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID:   int(createdM.MerchantID),
		DocumentType: "license",
		DocumentUrl:  "https://example.com/doc.pdf",
	})
	require.NoError(t, err)
	require.NotNil(t, doc)
	assert.Equal(t, "license", doc.DocumentType)
	assert.Equal(t, "https://example.com/doc.pdf", doc.DocumentUrl)
	assert.NotZero(t, doc.DocumentID)
}

func TestMerchantDocumentQuery_FindByID(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	createdM, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "DocFindM", Description: "D",
		Address: "A", ContactEmail: "df@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	created, err := repos.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID:   int(createdM.MerchantID),
		DocumentType: "id_card",
		DocumentUrl:  "https://example.com/id.pdf",
	})
	require.NoError(t, err)

	found, err := repos.MerchantDocumentQuery.FindById(ctx, int(created.DocumentID))
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.DocumentID, found.DocumentID)
	assert.Equal(t, "id_card", found.DocumentType)
}

func TestMerchantDocumentCommand_Update(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	createdM, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "DocUpdateM", Description: "D",
		Address: "A", ContactEmail: "du@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	created, err := repos.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID:   int(createdM.MerchantID),
		DocumentType: "old_type",
		DocumentUrl:  "https://old.com/doc.pdf",
	})
	require.NoError(t, err)

	updated, err := repos.MerchantDocumentCommand.UpdateMerchantDocument(ctx, &requests.UpdateMerchantDocumentRequest{
		DocumentID:   ptr(int(created.DocumentID)),
		MerchantID:   int(createdM.MerchantID),
		DocumentType: "new_type",
		DocumentUrl:  "https://new.com/doc.pdf",
		Status:       "verified",
		Note:         "Approved",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "new_type", updated.DocumentType)
	assert.Equal(t, "verified", updated.Status)
}

func TestMerchantDocumentCommand_UpdateStatus(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	createdM, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "DocStatusM", Description: "D",
		Address: "A", ContactEmail: "ds@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	created, err := repos.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID:   int(createdM.MerchantID),
		DocumentType: "cert",
		DocumentUrl:  "https://example.com/cert.pdf",
	})
	require.NoError(t, err)

	updated, err := repos.MerchantDocumentCommand.UpdateMerchantDocumentStatus(ctx, &requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: ptr(int(created.DocumentID)),
		MerchantID: int(createdM.MerchantID),
		Status:     "rejected",
		Note:       "Invalid document",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "rejected", updated.Status)
}

func TestMerchantDocumentQuery_FindAll(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	createdM, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "DocAllM", Description: "D",
		Address: "A", ContactEmail: "da@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	_, err = repos.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID: int(createdM.MerchantID), DocumentType: "t1", DocumentUrl: "https://t1.pdf",
	})
	require.NoError(t, err)
	_, err = repos.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID: int(createdM.MerchantID), DocumentType: "t2", DocumentUrl: "https://t2.pdf",
	})
	require.NoError(t, err)

	results, total, err := repos.MerchantDocumentQuery.FindAllDocuments(ctx, &requests.FindAllMerchantDocuments{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 2)
	assert.GreaterOrEqual(t, len(results), 2)
}

func TestMerchantDocumentCommand_Trash(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	createdM, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "DocTrashM", Description: "D",
		Address: "A", ContactEmail: "dt@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	created, err := repos.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID: int(createdM.MerchantID), DocumentType: "trash", DocumentUrl: "https://trash.pdf",
	})
	require.NoError(t, err)

	trashed, err := repos.MerchantDocumentCommand.TrashedMerchantDocument(ctx, int(created.DocumentID))
	require.NoError(t, err)
	require.NotNil(t, trashed)
	assert.True(t, trashed.DeletedAt.Valid)
}

func TestMerchantDocumentCommand_Restore(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	createdM, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "DocRestoreM", Description: "D",
		Address: "A", ContactEmail: "dr@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	created, err := repos.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID: int(createdM.MerchantID), DocumentType: "restore", DocumentUrl: "https://restore.pdf",
	})
	require.NoError(t, err)

	_, err = repos.MerchantDocumentCommand.TrashedMerchantDocument(ctx, int(created.DocumentID))
	require.NoError(t, err)

	restored, err := repos.MerchantDocumentCommand.RestoreMerchantDocument(ctx, int(created.DocumentID))
	require.NoError(t, err)
	require.NotNil(t, restored)
	assert.False(t, restored.DeletedAt.Valid)
}

func TestMerchantDocumentCommand_DeletePermanent(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedMerchantUser(t, ctx, pool)
	createdM, err := repos.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID: userID, Name: "DocDelM", Description: "D",
		Address: "A", ContactEmail: "dd@m.com", ContactPhone: "000", Status: "active",
	})
	require.NoError(t, err)

	created, err := repos.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID: int(createdM.MerchantID), DocumentType: "perm", DocumentUrl: "https://perm.pdf",
	})
	require.NoError(t, err)

	_, err = repos.MerchantDocumentCommand.TrashedMerchantDocument(ctx, int(created.DocumentID))
	require.NoError(t, err)

	deleted, err := repos.MerchantDocumentCommand.DeleteMerchantDocumentPermanent(ctx, int(created.DocumentID))
	require.NoError(t, err)
	assert.True(t, deleted)
}
