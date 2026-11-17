package printqueue

import (
	"database/sql"
	"errors"
	"net/url"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func openPrinters(t *testing.T) *sql.DB {
	t.Helper()

	// Set the timeout on each connection and use the same file for the whole pool.
	dsn := url.URL{Scheme: "file", Path: filepath.Join(t.TempDir(), "printers.sqlite")}
	dsn.RawQuery = url.Values{"_pragma": {"busy_timeout(5000)"}}.Encode()
	db, err := sql.Open("sqlite", dsn.String())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	db.SetMaxOpenConns(4)
	_, err = db.ExecContext(t.Context(), `
		CREATE TABLE printers (
			id INTEGER PRIMARY KEY,
			pages INTEGER NOT NULL CHECK (pages >= 0)
		);
		CREATE TABLE deliveries (
			delivery_id TEXT NOT NULL PRIMARY KEY,
			printer_id INTEGER NOT NULL,
			pages INTEGER NOT NULL CHECK (pages > 0)
		);
		INSERT INTO printers (id, pages) VALUES (1, 8);
	`)
	require.NoError(t, err)

	return db
}

func checkPages(t *testing.T, db *sql.DB, want int64) {
	t.Helper()

	var pages int64

	err := db.QueryRowContext(t.Context(), `SELECT pages FROM printers WHERE id = 1`).Scan(&pages)
	require.NoError(t, err)
	require.Equal(t, want, pages)
}

func TestConsume(t *testing.T) {
	db := openPrinters(t)
	ctx := t.Context()

	require.ErrorIs(t, Consume(ctx, db, 1, 0), ErrInvalidPages)
	require.ErrorIs(t, Consume(ctx, db, 1, -1), ErrInvalidPages)
	require.ErrorIs(t, Consume(ctx, db, 2, 1), ErrUnavailable)
	require.ErrorIs(t, Consume(ctx, db, 1, 9), ErrUnavailable)
	checkPages(t, db, 8)

	require.NoError(t, Consume(ctx, db, 1, 8))
	checkPages(t, db, 0)
}

func TestConsumeConcurrent(t *testing.T) {
	db := openPrinters(t)
	ctx := t.Context()

	const calls = 10

	start := make(chan struct{})
	results := make(chan error, calls)

	var wg sync.WaitGroup

	for range calls {
		wg.Go(func() {
			<-start

			results <- Consume(ctx, db, 1, 2)
		})
	}

	close(start)
	wg.Wait()
	close(results)

	success := 0
	unavailable := 0

	for err := range results {
		switch {
		case err == nil:
			success++
		case errors.Is(err, ErrUnavailable):
			unavailable++
		default:
			t.Fatalf("unexpected Consume error: %v", err)
		}
	}

	require.Equal(t, 4, success)
	require.Equal(t, 6, unavailable)
	checkPages(t, db, 0)
}

func TestRefill(t *testing.T) {
	db := openPrinters(t)
	ctx := t.Context()

	require.ErrorIs(t, Refill(ctx, db, 1, "empty", 0), ErrInvalidPages)
	require.ErrorIs(t, Refill(ctx, db, 1, "negative", -1), ErrInvalidPages)
	require.ErrorIs(t, Refill(ctx, db, 2, "missing", 3), sql.ErrNoRows)
	checkPages(t, db, 8)

	require.NoError(t, Refill(ctx, db, 1, "box'1", 3))
	checkPages(t, db, 11)
	require.Error(t, Refill(ctx, db, 1, "box'1", 4))
	checkPages(t, db, 11)

	var (
		deliveryID string
		printerID  int64
		pages      int64
	)

	err := db.QueryRowContext(ctx, `SELECT delivery_id, printer_id, pages FROM deliveries`).
		Scan(&deliveryID, &printerID, &pages)
	require.NoError(t, err)
	require.Equal(t, "box'1", deliveryID)
	require.Equal(t, int64(1), printerID)
	require.Equal(t, int64(3), pages)

	var count int

	err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM deliveries`).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}
