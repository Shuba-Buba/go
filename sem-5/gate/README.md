# gate

Реализуйте counting semaphore на буферизованном канале, как в лекции про
web crawler.

```go
func NewGate(capacity int) *Gate
func (g *Gate) Acquire()
func (g *Gate) Release()
```

`NewGate(n)` ограничивает число одновременных `Acquire` без `Release` значением
`n`. `Release` освобождает слот. `capacity` всегда положительный.

```shell
go test ./gate/...
```
