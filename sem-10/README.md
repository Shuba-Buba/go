# Семинар 10. SQL из Go

1. [readinglist](readinglist/README.md) — вставка с получением ID, выборка и удаление.
2. [printqueue](printqueue/README.md) — условное обновление остатка и транзакция.

Запускайте `go test -timeout=1m ./readinglist` и `go test -timeout=1m ./printqueue`
из корня модуля. Для конкурентных списаний дополнительно запустите
`go test -race -timeout=1m ./printqueue`. До реализации тесты падают.
Используется SQLite, отдельный сервер запускать не нужно.
