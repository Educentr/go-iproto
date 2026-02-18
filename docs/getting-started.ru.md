# Быстрый старт

## Установка

```bash
go get github.com/Educentr/go-iproto
```

Требуется Go 1.24 или новее.

## Первый сервер

Простой echo-сервер, который возвращает полученные данные:

```go
package main

import (
    "context"
    "log"

    "github.com/Educentr/go-iproto/iproto"
)

func main() {
    handler := iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
        _ = c.Send(ctx, iproto.ResponseTo(p, p.Data))
    })

    log.Fatal(iproto.ListenAndServe(context.Background(), "tcp", ":3301", handler))
}
```

## Первый клиент

Подключение к серверу и отправка запроса:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Educentr/go-iproto/iproto"
)

func main() {
    pool, err := iproto.Dial(context.Background(), "tcp", "127.0.0.1:3301", nil)
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    resp, err := pool.Call(context.Background(), 1, []byte("hello"))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ответ: %s\n", resp)
}
```

## Структура пакетов

| Пакет | Путь импорта | Описание |
|-------|-------------|----------|
| iproto | `github.com/Educentr/go-iproto/iproto` | Ядро протокола: Channel, Pool, Server, ServeMux, Pack/Unpack |
| netutil | `github.com/Educentr/go-iproto/netutil` | Dialer с retry/backoff, BackgroundDialer |
| syncutil | `github.com/Educentr/go-iproto/syncutil` | TaskGroup, TaskRunner, Throttle, Multitask |

## Ключевые концепции

- **Channel** управляет одним TCP-соединением. Обрабатывает чтение/запись пакетов, ping/keepalive и корреляцию запрос-ответ через `Sync` ID.
- **Pool** управляет несколькими Channel к одному серверу. Балансировка нагрузки (round-robin), автоматическое переподключение, опциональный rate limiting.
- **Server** принимает входящие TCP-соединения и создаёт Channel для каждого.
- **ServeMux** маршрутизирует входящие пакеты к обработчикам по коду сообщения (`Header.Msg`).
- **Pack/Unpack** обеспечивает бинарную сериализацию. Целые числа — little-endian по умолчанию, с опциональным BER-кодированием (переменная длина).

## Далее

- [Архитектура](architecture.md) — устройство протокола
- [API-справочник](api/iproto.md) — подробная документация API
- [Примеры](examples.md) — запускаемые примеры
