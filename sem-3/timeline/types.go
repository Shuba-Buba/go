package timeline

import "errors"

// ErrNegativeValue сообщает о нарушении инварианта неотрицательного значения.
var ErrNegativeValue = errors.New("timeline value cannot be negative")

// Change задаёт изменение в конкретный день.
type Change struct {
	Day   int
	Delta int
}

// Snapshot хранит значение после применения всех изменений дня.
type Snapshot struct {
	Day   int
	Value int
}
