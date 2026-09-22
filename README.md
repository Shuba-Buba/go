# Go — семинары

Репозиторий с решениями задач по Go.

**Перед первой сдачей прочитайте этот файл целиком.**  
Работа, сданная не по формату, **не проверяется**.

---

## Формат сдачи (обязательно)

| Что | Формат | Пример |
|---|---|---|
| Ветка | `sem-<номер>/<фамилия>` | `sem-2/ivanov` |
| Папка с кодом | `sem-<номер>/<задача>/` | `sem-2/counter/` |
| Коммит | `sem-<номер>: <задачи>` | `sem-2: counter, divmod` |
| PR / MR (заголовок) | `sem-<номер>/<фамилия> — <задачи>` | `sem-2/ivanov — counter, divmod` |

### Правила именования

- **Только латиница**, строчные буквы, без пробелов и кириллицы.
- **Одна ветка = один семинар = один PR/MR.**
- Не смешивайте разные семинары в одной ветке.
- **Не пушьте в `main`** — только в свою ветку.

---

## Пошаговая инструкция

### 1. Клонирование (один раз)

```bash
git clone <url-репозитория>
cd go
```

### 2. Перед каждым семинаром

```bash
git checkout main
git pull origin main
git checkout -b sem-N/фамилия
```

Пример:

```bash
git checkout -b sem-2/ivanov
```

### 3. Решение задач

Код кладите **только** в папку текущего семинара:

```text
sem-2/
  counter/
    counter.go
    counter_test.go
  divmod/
    divmod.go
    divmod_test.go
```

**Не меняйте** файлы вне своего семинара и не трогайте чужие решения.

Перед сдачей запустите тесты:

```bash
cd sem-2
go test ./...
```

Все тесты должны проходить (или явно укажите в PR/MR, что не проходит).

### 4. Коммит

```bash
git add sem-2/
git commit -m "sem-2: counter, divmod"
```

Формат сообщения коммита:

```text
sem-<номер>: <задача1>, <задача2>
```

Примеры:

- `sem-2: counter`
- `sem-2: counter, divmod, linkedlist`
- `sem-2: fix counter after review` — после правок по ревью

### 5. Push в свою ветку

```bash
git push -u origin sem-2/ivanov
```

### 6. Открыть PR / MR

| Платформа | Куда |
|---|---|
| **GitHub** | Pull requests → New pull request → base: `main`, compare: ваша ветка |
| **GitLab** | Merge requests → New merge request → source: ваша ветка, target: `main` |

**Заголовок PR/MR** (обязательно):

```text
sem-2/ivanov — counter, divmod
```

**Описание PR/MR** — скопируйте шаблон ниже и заполните:

```markdown
## Автор
Иванов Иван

## Семинар
sem-2

## Задачи
- [x] counter
- [x] divmod
- [ ] linkedlist

## Тесты
- [x] `go test ./...` проходит локально

## Комментарии
linkedlist не успел
```

**Без PR/MR с заполненным описанием работа не считается сданной.**

### 7. Правки после ревью

Правьте код **в той же ветке**, push обновит PR/MR автоматически:

```bash
git add sem-2/
git commit -m "sem-2: fix counter after review"
git push
```

---

## Чеклист перед сдачей

- [ ] Ветка названа `sem-N/фамилия`
- [ ] Код лежит в `sem-N/<задача>/`
- [ ] Коммит в формате `sem-N: ...`
- [ ] `go test ./...` проходит
- [ ] Push только в свою ветку, не в `main`
- [ ] Открыт PR/MR в `main`
- [ ] Заголовок PR/MR: `sem-N/фамилия — ...`
- [ ] Описание PR/MR заполнено по шаблону

---

## Частые ошибки

| Ошибка | Как правильно |
|---|---|
| Push в `main` | Push в `sem-N/фамилия`, потом PR/MR |
| Ветка `sem2-ivanov` или `Sem-2/Ivanov` | Только `sem-2/ivanov` |
| Нет PR/MR | Открыть PR/MR после push |
| Пустое описание PR/MR | Заполнить шаблон |
| Два семинара в одной ветке | Отдельная ветка на каждый семинар |

---

## Структура репозитория

```text
sem-2/
  counter/
  divmod/
  linkedlist/
  permissions/
  treetransform/
sem-3/
  dedup/
  slicestack/
  timeline/
  utf8sanitize/
```
