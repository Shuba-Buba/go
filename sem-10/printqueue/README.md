# printqueue

Учёт бумаги хранит остаток листов для каждого принтера и номера поставок.
Тесты передают готовый `*sql.DB` с таблицами:

```sql
CREATE TABLE printers (
    id INTEGER PRIMARY KEY,
    pages INTEGER NOT NULL CHECK (pages >= 0)
);

CREATE TABLE deliveries (
    delivery_id TEXT NOT NULL PRIMARY KEY,
    printer_id INTEGER NOT NULL,
    pages INTEGER NOT NULL CHECK (pages > 0)
);
```

- `Consume` списывает заданное число листов одним `UPDATE`. Если принтера нет
  или листов недостаточно, верните `ErrUnavailable`. При одновременном
  списании остаток не должен уходить в минус.
- `Refill` добавляет листы к остатку и записывает строку поставки в одной
  транзакции. Повторный номер поставки возвращает ошибку. При отсутствии
  принтера верните `sql.ErrNoRows`. После ошибки остаток и журнал должны
  остаться прежними.
- Неположительное число листов для обеих операций даёт `ErrInvalidPages`.
  Сценарии используют остатки и суммы, которые помещаются в `int64`.

Передавайте значения отдельными аргументами и используйте полученный `ctx`.
За закрытие `db` отвечает вызывающий код.

```sh
go test -timeout=1m ./printqueue
go test -race -timeout=1m ./printqueue
```
