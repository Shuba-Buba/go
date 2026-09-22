//go:build !solution

package divmod

import "errors"

// ErrDivisionByZero сообщает о попытке деления на ноль.
var ErrDivisionByZero = errors.New("division by zero")

// DivMod возвращает частное и остаток от деления dividend на divisor.
func DivMod(dividend, divisor int) (quotient, remainder int, err error) {
	if divisor == 0 {
		return 0, 0, ErrDivisionByZero
	}

	return dividend / divisor, dividend % divisor, nil
}
