//go:build !solution

package readinglist

import (
	"context"
	"database/sql"
)

// Add saves a book and returns its database ID.
func Add(ctx context.Context, db *sql.DB, title, author string) (int64, error) {
	res, err := db.ExecContext(ctx, `INSERT INTO books (title, author) VALUES (?, ?)`, title, author)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// ByAuthor returns the author's books ordered by title.
func ByAuthor(ctx context.Context, db *sql.DB, author string) ([]Book, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, title FROM books WHERE author = ? ORDER BY title ASC`, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

// Remove deletes the book with the given ID.
func Remove(ctx context.Context, db *sql.DB, id int64) error {
	res, err := db.ExecContext(ctx, `DELETE FROM books WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
