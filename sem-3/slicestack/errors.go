package slicestack

import "errors"

// ErrEmpty сообщает об операции чтения из пустого стека.
var ErrEmpty = errors.New("empty stack")
