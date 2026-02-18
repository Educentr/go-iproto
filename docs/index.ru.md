# go-iproto

Go-реализация бинарного протокола iproto — легковесный высокопроизводительный фреймворк для бинарного RPC.

## Возможности

- **Channel** — одно соединение с корреляцией запрос-ответ, ping/keepalive, graceful shutdown
- **Pool** — пул соединений с балансировкой, автоматическим переподключением и rate limiting
- **Server / ServeMux** — приём соединений с маршрутизацией по коду сообщения
- **Pack / Unpack** — бинарная сериализация целых чисел, строк, структур с поддержкой BER
- **netutil** — Dialer с retry/backoff и фоновым подключением
- **syncutil** — примитивы конкурентности: TaskGroup, TaskRunner, Throttle, Multitask

## Установка

```bash
go get github.com/Educentr/go-iproto
```

Требуется Go 1.24+.

## Быстрый пример

```go
// Сервер
handler := iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
    _ = c.Send(ctx, iproto.ResponseTo(p, p.Data))
})
iproto.ListenAndServe(ctx, "tcp", ":3301", handler)

// Клиент
pool, _ := iproto.Dial(ctx, "tcp", "127.0.0.1:3301", nil)
resp, _ := pool.Call(ctx, 1, []byte("hello"))
```

## Ссылки

- [Быстрый старт](getting-started.md) — пошаговое руководство
- [Архитектура](architecture.md) — устройство протокола
- [API-справочник](api/iproto.md) — типы, функции, константы
- [Примеры](examples.md) — запускаемые примеры кода
- [GitHub](https://github.com/Educentr/go-iproto)
- [pkg.go.dev](https://pkg.go.dev/github.com/Educentr/go-iproto)
