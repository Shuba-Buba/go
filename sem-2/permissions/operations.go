//go:build !solution

package permissions

// Grant добавляет права added к current.
func Grant(current, added Permission) Permission {
	return current | added
}

// Revoke удаляет права removed из current.
func Revoke(current, removed Permission) Permission {
	return current &^ removed
}

// Has проверяет наличие всех прав required в current.
func Has(current, required Permission) bool {
	return current&required == required
}
