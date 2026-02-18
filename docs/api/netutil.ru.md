# Пакет netutil

`import "github.com/Educentr/go-iproto/netutil"`

Полный API-справочник: [pkg.go.dev](https://pkg.go.dev/github.com/Educentr/go-iproto/netutil)

## Dialer

`Dialer` устанавливает TCP-соединение с повторными попытками и экспоненциальным backoff.

```go
type Dialer struct {
    Network         string
    Addr            string
    Timeout         time.Duration        // Таймаут одной попытки
    LoopInterval    time.Duration        // Начальный backoff (по умолчанию 50мс)
    MaxLoopInterval time.Duration        // Максимальный backoff
    NetDial         func(ctx, net, addr) (net.Conn, error)
    OnAttempt       func(err error)      // Вызывается после каждой попытки
    Logf            func(string, ...any)
    Debugf          func(string, ...any)
}
```

**Методы:**

| Метод | Описание |
|-------|----------|
| `Dial(ctx) (net.Conn, error)` | Подключение с повторами до успеха или отмены контекста |

## BackgroundDialer

`BackgroundDialer` оборачивает `Dialer` для конкурентного фонового подключения. Несколько вызывающих могут подписаться на одну операцию подключения.

```go
type BackgroundDialer struct {
    Dialer     *Dialer
    TaskGroup  *syncutil.TaskGroup
    TaskRunner *syncutil.TaskRunner
}
```

**Методы:**

| Метод | Описание |
|-------|----------|
| `Dial(ctx, callback) <-chan error` | Запуск или присоединение к фоновому подключению |
| `Cancel()` | Отмена текущей операции подключения |
| `SetDeadline(t time.Time)` | Установка абсолютного дедлайна |
| `SetDeadlineAtLeast(t time.Time) time.Time` | Установка дедлайна, если текущий раньше |

### Пример использования

```go
d := &netutil.BackgroundDialer{
    Dialer: &netutil.Dialer{
        Network: "tcp",
        Addr:    "127.0.0.1:3301",
        Timeout: 5 * time.Second,
    },
    TaskRunner: &syncutil.TaskRunner{},
}

done := d.Dial(ctx, func(conn net.Conn, err error) {
    if err != nil {
        log.Fatal(err)
    }
    // использование conn
})

<-done
```
