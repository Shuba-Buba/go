//go:build !solution

package readinglist

import (
	"context"
	"database/sql"
)

// Add saves a book and returns its database ID.
func Add(ctx context.Context, db *sql.DB, title, author string) (int64, error) {
	_, _, _, _ = ctx, db, title, author

	panic("TODO: implement Add")
}

// ByAuthor returns the author's books ordered by title.
func ByAuthor(ctx context.Context, db *sql.DB, author string) ([]Book, error) {
	_, _, _ = ctx, db, author

	panic("TODO: implement ByAuthor")
}

// Remove deletes the book with the given ID.
func Remove(ctx context.Context, db *sql.DB, id int64) error {
	_, _, _ = ctx, db, id

	panic("TODO: implement Remove")
}
