//go:build !solution

package printqueue

import (
	"context"
	"database/sql"
)

// Consume takes the requested number of pages from a printer's stock.
func Consume(ctx context.Context, db *sql.DB, printerID, pages int64) error {
	if pages <= 0 {
		return ErrInvalidPages
	}
	res, err := db.ExecContext(ctx, "UPDATE printers SET pages = pages - ? WHERE id = ? AND pages >= ?", pages, printerID, pages)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrUnavailable
	}
	return nil
}

// Refill adds pages to a printer's stock and records the delivery.
func Refill(
	ctx context.Context, db *sql.DB, printerID int64, deliveryID string, pages int64,
) error {
	if pages <= 0 {
		return ErrInvalidPages
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, "UPDATE printers SET pages = pages + ? WHERE id = ?", pages, printerID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO deliveries (delivery_id, printer_id, pages) VALUES (?, ?, ?)", deliveryID, printerID, pages); err != nil {
		return err
	}
	return tx.Commit()
}
