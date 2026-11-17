package readinglist

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func openBooks(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	// Each SQLite connection to :memory: has its own database.
	db.SetMaxOpenConns(1)
	_, err = db.ExecContext(t.Context(), `
		CREATE TABLE books (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			UNIQUE (author, title)
		)
	`)
	require.NoError(t, err)

	return db
}

func TestReadingList(t *testing.T) {
	db := openBooks(t)
	ctx := t.Context()

	lastID, err := Add(ctx, db, "The Last Gift", "Abdulrazak Gurnah")
	require.NoError(t, err)
	paradiseID, err := Add(ctx, db, "Paradise", "Abdulrazak Gurnah")
	require.NoError(t, err)
	otherID, err := Add(ctx, db, "The River Between", "Ngugi wa Thiong'o")
	require.NoError(t, err)
	require.NotZero(t, lastID)
	require.NotZero(t, paradiseID)
	require.NotZero(t, otherID)
	require.NotEqual(t, lastID, paradiseID)

	_, err = Add(ctx, db, "Paradise", "Abdulrazak Gurnah")
	require.Error(t, err)

	books, err := ByAuthor(ctx, db, "Abdulrazak Gurnah")
	require.NoError(t, err)
	require.Equal(t, []Book{
		{ID: paradiseID, Title: "Paradise"},
		{ID: lastID, Title: "The Last Gift"},
	}, books)

	books, err = ByAuthor(ctx, db, "Ngugi wa Thiong'o")
	require.NoError(t, err)
	require.Equal(t, []Book{{ID: otherID, Title: "The River Between"}}, books)

	books, err = ByAuthor(ctx, db, "' OR 1=1 --")
	require.NoError(t, err)
	require.Empty(t, books)

	require.NoError(t, Remove(ctx, db, paradiseID))
	require.ErrorIs(t, Remove(ctx, db, paradiseID), sql.ErrNoRows)

	books, err = ByAuthor(ctx, db, "Abdulrazak Gurnah")
	require.NoError(t, err)
	require.Equal(t, []Book{{ID: lastID, Title: "The Last Gift"}}, books)
}

func TestAddCanceled(t *testing.T) {
	db := openBooks(t)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := Add(ctx, db, "Paradise", "Abdulrazak Gurnah")
	require.ErrorIs(t, err, context.Canceled)

	books, err := ByAuthor(t.Context(), db, "Abdulrazak Gurnah")
	require.NoError(t, err)
	require.Empty(t, books)
}
