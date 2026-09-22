# utf8sanitize

Реализуйте функцию:

```go
func Sanitize(input string) (normalized string, invalidBytes int)
```

Функция должна вернуть корректную UTF-8 строку. Каждый byte входа, который не
начинает корректное UTF-8 encoding, заменяется на Unicode replacement character
U+FFFD. Второй результат — количество заменённых bytes.

Законный символ U+FFFD, уже находившийся во входной строке, сохраняется и не
увеличивает счётчик. Отличить его от ошибки можно по размеру, который возвращает
`utf8.DecodeRuneInString`: ошибка декодирования имеет `r == utf8.RuneError` и
`size == 1`.

Собирайте результат с помощью `strings.Builder`.

## Проверка

```shell
go test ./utf8sanitize/...
```
