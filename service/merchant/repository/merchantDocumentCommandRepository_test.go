package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-pointofsale-shared/errors"
)

// fakeDB is a minimal database/sql driver that records the executed SQL and
// its positional args and returns a canned row for queries. It allows testing
// the repository param building without a real database.
type fakeDB struct {
	lastSQL  string
	lastArgs []driver.Value

	columns []string
	row     []driver.Value
	rowErr  error
}

func (f *fakeDB) Connect(ctx context.Context) (driver.Conn, error) { return f, nil }

func (f *fakeDB) Driver() driver.Driver { return f }

func (f *fakeDB) Open(name string) (driver.Conn, error) { return f, nil }

func (f *fakeDB) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("fakeDB: Prepare not supported")
}

func (f *fakeDB) Close() error { return nil }

func (f *fakeDB) Begin() (driver.Tx, error) {
	return nil, errors.New("fakeDB: Begin not supported")
}

func (f *fakeDB) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	f.lastSQL = query
	f.lastArgs = namedValues(args)
	return driver.RowsAffected(0), nil
}

func (f *fakeDB) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	f.lastSQL = query
	f.lastArgs = namedValues(args)
	if f.rowErr != nil {
		return nil, f.rowErr
	}
	return &fakeDriverRows{columns: f.columns, vals: f.row}, nil
}

func namedValues(args []driver.NamedValue) []driver.Value {
	vals := make([]driver.Value, len(args))
	for i, a := range args {
		vals[i] = a.Value
	}
	return vals
}

// fakeDriverRows returns a single canned row.
type fakeDriverRows struct {
	columns []string
	vals    []driver.Value
	done    bool
}

func (r *fakeDriverRows) Columns() []string { return r.columns }

func (r *fakeDriverRows) Close() error { return nil }

func (r *fakeDriverRows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}
	r.done = true
	copy(dest, r.vals)
	return nil
}

// newFakeSQLxDB wraps fakeDB in a real *sqlx.DB so the repository exercises the
// exact sqlx code path used in production.
func newFakeSQLxDB(t *testing.T, row []driver.Value, rowErr error) (*sqlx.DB, *fakeDB) {
	t.Helper()
	f := &fakeDB{
		columns: []string{
			"document_id", "merchant_id", "document_type", "document_url",
			"status", "note", "uploaded_at", "created_at", "updated_at", "deleted_at",
		},
		row:    row,
		rowErr: rowErr,
	}
	return sqlx.NewDb(sql.OpenDB(f), "postgres"), f
}

// merchantDocumentRowFixture mirrors the column order of the
// merchant_documents RETURNING * scan.
func merchantDocumentRowFixture(docID, merchantID int) []driver.Value {
	return []driver.Value{
		int64(docID),
		int64(merchantID),
		"updated_type",
		"https://example.com/updated.pdf",
		"verified",
		"Approved",
		nil,
		nil,
		nil,
		nil,
	}
}

// Regression tests for the DocumentID-vs-MerchantID fix: the update must be
// keyed on the request's DocumentID — never on the MerchantID. The fixture
// deliberately uses different values (DocumentID=42, MerchantID=999) so a
// regression (using MerchantID as the WHERE key) cannot pass by coincidence.

func TestUpdateMerchantDocument_UsesDocumentID_NotMerchantID(t *testing.T) {
	const (
		docID      = 42
		merchantID = 999
	)

	fakeConn, fake := newFakeSQLxDB(t, merchantDocumentRowFixture(docID, merchantID), nil)
	repo := NewMerchantDocumentCommandRepository(fakeConn)

	updated, err := repo.UpdateMerchantDocument(context.Background(), &requests.UpdateMerchantDocumentRequest{
		DocumentID:   ptr(docID),
		MerchantID:   merchantID,
		DocumentType: "updated_type",
		DocumentUrl:  "https://example.com/updated.pdf",
		Status:       "verified",
		Note:         "Approved",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, int32(docID), updated.DocumentID)
	assert.Equal(t, int32(merchantID), updated.MerchantID)

	// The first positional arg ($1 — the WHERE key) must be the DocumentID.
	require.Len(t, fake.lastArgs, 5)
	assert.Equal(t, int64(docID), fake.lastArgs[0],
		"WHERE key must be the request DocumentID, never the MerchantID")
	assert.NotEqual(t, int64(merchantID), fake.lastArgs[0])

	// All params map to the query in order: DocumentID, DocumentType,
	// DocumentUrl, Status, Note.
	assert.Equal(t, "updated_type", fake.lastArgs[1])
	assert.Equal(t, "https://example.com/updated.pdf", fake.lastArgs[2])
	assert.Equal(t, "verified", fake.lastArgs[3])
	assert.Equal(t, "Approved", fake.lastArgs[4])

	// The SQL must filter by document_id — a "merchant_id = $n"
	// predicate is exactly the bug this fix removes.
	assert.Contains(t, fake.lastSQL, "document_id = $1 AND deleted_at IS NULL")
	assert.NotContains(t, fake.lastSQL, "merchant_id = $")
}

func TestUpdateMerchantDocumentStatus_UsesDocumentID_NotMerchantID(t *testing.T) {
	const (
		docID      = 42
		merchantID = 999
	)

	fakeConn, fake := newFakeSQLxDB(t, merchantDocumentRowFixture(docID, merchantID), nil)
	repo := NewMerchantDocumentCommandRepository(fakeConn)

	updated, err := repo.UpdateMerchantDocumentStatus(context.Background(), &requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: ptr(docID),
		MerchantID: merchantID,
		Status:     "rejected",
		Note:       "Invalid document",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, int32(docID), updated.DocumentID)

	require.Len(t, fake.lastArgs, 3)
	assert.Equal(t, int64(docID), fake.lastArgs[0],
		"WHERE key must be the request DocumentID, never the MerchantID")
	assert.NotEqual(t, int64(merchantID), fake.lastArgs[0])

	// Remaining params map in order: Status, Note.
	assert.Equal(t, "rejected", fake.lastArgs[1])
	assert.Equal(t, "Invalid document", fake.lastArgs[2])

	assert.Contains(t, fake.lastSQL, "document_id = $1 AND deleted_at IS NULL")
	assert.NotContains(t, fake.lastSQL, "merchant_id = $")
}

func TestUpdateMerchantDocument_NilDocumentID_ReturnsBadRequest(t *testing.T) {
	fakeConn, fake := newFakeSQLxDB(t, nil, nil)
	repo := NewMerchantDocumentCommandRepository(fakeConn)

	_, err := repo.UpdateMerchantDocument(context.Background(), &requests.UpdateMerchantDocumentRequest{
		MerchantID: 999,
	})
	require.Error(t, err)

	var appErr *sharedErrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, sharedErrors.ErrorTypeBadRequest, appErr.Type)
	assert.Contains(t, appErr.Message, "document id is required")

	// The DB must never be touched when DocumentID is missing.
	assert.Empty(t, fake.lastSQL)
}

func TestUpdateMerchantDocumentStatus_NilDocumentID_ReturnsBadRequest(t *testing.T) {
	fakeConn, fake := newFakeSQLxDB(t, nil, nil)
	repo := NewMerchantDocumentCommandRepository(fakeConn)

	_, err := repo.UpdateMerchantDocumentStatus(context.Background(), &requests.UpdateMerchantDocumentStatusRequest{
		MerchantID: 999,
	})
	require.Error(t, err)

	var appErr *sharedErrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, sharedErrors.ErrorTypeBadRequest, appErr.Type)
	assert.Contains(t, appErr.Message, "document id is required")

	assert.Empty(t, fake.lastSQL)
}

func TestUpdateMerchantDocument_DBError_WrapsInternal(t *testing.T) {
	inner := errors.New("connection refused")
	fakeConn, _ := newFakeSQLxDB(t, nil, inner)
	repo := NewMerchantDocumentCommandRepository(fakeConn)

	_, err := repo.UpdateMerchantDocument(context.Background(), &requests.UpdateMerchantDocumentRequest{
		DocumentID: ptr(42),
		MerchantID: 999,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, sharedErrors.ErrInternal)
	assert.ErrorIs(t, err, inner)
}

func TestUpdateMerchantDocumentStatus_DBError_WrapsInternal(t *testing.T) {
	inner := errors.New("connection refused")
	fakeConn, _ := newFakeSQLxDB(t, nil, inner)
	repo := NewMerchantDocumentCommandRepository(fakeConn)

	_, err := repo.UpdateMerchantDocumentStatus(context.Background(), &requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: ptr(42),
		MerchantID: 999,
		Status:     "approved",
		Note:       "OK",
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, sharedErrors.ErrInternal)
	assert.ErrorIs(t, err, inner)
}
