# Практика тестирования

Каталог содержит пять упражнений после лекции 6. Начните с `tablecases`,
`testhelper` и `testcleanup`. Затем вместе с преподавателем выполните короткие
`parallelcase` и `fuzzcase`.

Команды запускаются из этого каталога, в котором находится `go.mod`:

```shell
go test ./tablecases
go test ./testhelper
go test ./testcleanup
go test ./parallelcase
go test ./fuzzcase
```

Начальные заглушки намеренно не проходят тесты. Каждый пакет независим:
для работы с одним упражнением не нужно сначала решать остальные.
Используется только стандартная библиотека Go.

Это новые упражнения на основе прежнего семинара Testing. Они тренируют
таблицы и покрытие для `tabletest`, helpers для `testequal`, cleanup для
`fileleak`, fuzzing для `iprange`. `parallelcase` помогает понять порядок
событий перед бонусным `tparallel`, но использует готовый пакет `testing`.
