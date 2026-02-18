# Пакет syncutil

`import "github.com/Educentr/go-iproto/syncutil"`

Полный API-справочник: [pkg.go.dev](https://pkg.go.dev/github.com/Educentr/go-iproto/syncutil)

## TaskGroup

Выполнение до N конкурентных задач. Слоты переиспользуются по мере завершения.

```go
type TaskGroup struct {
    N    int    // Макс. конкурентных задач (по умолчанию 1)
    Goer GoerFn // Пользовательский запуск горутин (опционально)
}
```

**Методы:**

| Метод | Описание |
|-------|----------|
| `Do(ctx, n, task) []<-chan error` | Запуск до n задач; возвращает каналы результатов |
| `Cancel()` | Отмена всех выполняющихся задач |

### Пример

```go
tg := &syncutil.TaskGroup{N: 4}

results := tg.Do(ctx, 4, func(ctx context.Context, i int) error {
    // задача i выполняется конкурентно
    return nil
})

for _, ch := range results {
    if err := <-ch; err != nil {
        log.Fatal(err)
    }
}
```

## TaskRunner

Запуск ровно одной задачи. Несколько вызывающих подписываются на один результат.

```go
type TaskRunner struct{}
```

**Методы:**

| Метод | Описание |
|-------|----------|
| `Do(ctx, task) <-chan error` | Запуск задачи (или подписка, если уже выполняется) |
| `Cancel()` | Отмена текущей задачи |

### Пример

```go
var runner syncutil.TaskRunner

// Обе горутины получат одинаковый результат.
ch1 := runner.Do(ctx, func(ctx context.Context) error {
    return expensiveOperation(ctx)
})
ch2 := runner.Do(ctx, func(ctx context.Context) error {
    return expensiveOperation(ctx) // не запускается — подписывается на существующую
})

err1 := <-ch1
err2 := <-ch2 // то же значение, что и err1
```

## Throttle

Ограничение выполнения функции — не чаще одного раза за период.

```go
type Throttle struct {
    Period time.Duration
}
```

**Методы:**

| Метод | Описание |
|-------|----------|
| `Next() bool` | Возвращает true, если период прошёл |
| `Reset()` | Сброс throttle |
| `Set(time.Time)` | Установка конкретного времени |

## Multitask

Запуск N задач параллельно с поддержкой отмены.

```go
type Multitask struct {
    ContinueOnError bool   // Не отменять остальные при ошибке Goer
    Goer            GoerFn // Пользовательский запуск горутин
}
```

**Методы:**

| Метод | Описание |
|-------|----------|
| `Do(ctx, n, actor) error` | Запуск actor n раз параллельно |

## Вспомогательные функции

```go
// Every запускает n задач; отменяет все при первой ошибке.
func Every(ctx context.Context, n int, actor func(context.Context, int) error) error

// Each запускает n задач; ждёт завершения всех.
func Each(ctx context.Context, n int, actor func(context.Context, int))
```
