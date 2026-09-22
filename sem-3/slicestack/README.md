# slicestack

Реализуйте стек целых чисел поверх slice:

```go
func Push(stack []int, value int) []int
func Peek(stack []int) (value int, err error)
func Pop(stack []int) (value int, rest []int, err error)
```

- `Push` добавляет значение на вершину и возвращает новый slice header.
- `Peek` возвращает верхнее значение, не меняя длину стека.
- `Pop` возвращает верхнее значение и стек без него.
- `Peek` на пустом стеке возвращает `(0, ErrEmpty)`.
- `Pop` на пустом стеке возвращает `(0, stack, ErrEmpty)`, сохраняя исходный
  пустой slice: для `nil` остаётся `nil`, для `[]int{}` — не-`nil` empty slice.

Поскольку slice header передаётся по значению, результат `Push` и оставшийся
стек из `Pop` нужно использовать у вызывающего кода.

## Проверка

```shell
go test ./slicestack/...
```
