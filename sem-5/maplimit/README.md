# maplimit

Реализуйте `MapLimit` — parallel loop с ограничением параллелизма, как
`makeThumbnails4/5` и counting semaphore в лекции.

```go
func MapLimit(items []int, limit int, fn func(int) (int, error)) ([]int, error)
```

- Для каждого элемента `items` вызовите `fn` в отдельной горутине.
- Одновременно выполняется не больше `limit` вызовов.
- Результаты верните **в том же порядке**, что и `items`.
- Если хотя бы один вызов вернул ошибку, дождитесь завершения всех
  горутин и верните `(nil, err)` — одну из ошибок (какую, не важно).
- Пустой `items` → `([]int{}, nil)`. `limit` всегда положительный.

```shell
go test ./maplimit/...
```
