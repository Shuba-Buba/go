package catalog

// Identified связывает тип записи с типом её идентификатора.
type Identified[K comparable] interface {
	ID() K
}
