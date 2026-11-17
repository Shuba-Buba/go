//go:build !solution

package printqueue

import (
	"context"
	"database/sql"
)

// Consume takes the requested number of pages from a printer's stock.
func Consume(ctx context.Context, db *sql.DB, printerID, pages int64) error {
	_, _, _, _ = ctx, db, printerID, pages

	panic("TODO: implement Consume")
}

// Refill adds pages to a printer's stock and records the delivery.
func Refill(
	ctx context.Context, db *sql.DB, printerID int64, deliveryID string, pages int64,
) error {
	_, _, _, _, _ = ctx, db, printerID, deliveryID, pages

	panic("TODO: implement Refill")
}
