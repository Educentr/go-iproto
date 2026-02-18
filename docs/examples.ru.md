# Примеры

Все примеры находятся в каталоге [`examples/`](https://github.com/Educentr/go-iproto/tree/main/examples) и запускаются через `go run`.

## Echo-сервер

Базовый echo-сервер, который возвращает полученные данные.

```bash
go run ./examples/echo
```

```go
// Обработчик сервера: эхо-ответ.
handler := iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
    _ = c.Send(ctx, iproto.ResponseTo(p, p.Data))
})

srv := &iproto.Server{
    ChannelConfig: &iproto.ChannelConfig{Handler: handler},
}
go srv.Serve(ctx, ln)

// Клиент: отправка и получение.
pool, _ := iproto.Dial(ctx, "tcp", addr, &iproto.PoolConfig{Size: 1})
data := iproto.PackString(nil, "Hello, iproto!", iproto.ModeDefault)
resp, _ := pool.Call(ctx, 1, data)
```

## Пул соединений

Несколько горутин отправляют запросы через пул из 4 соединений. Сервер удваивает каждый `uint32`.

```bash
go run ./examples/pool
```

```go
// Сервер: удвоение числа.
handler := iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
    var in uint32
    iproto.UnpackUint32(bytes.NewReader(p.Data), &in, 0)
    _ = c.Send(ctx, iproto.ResponseTo(p, iproto.PackUint32(nil, in*2, 0)))
})

// Пул из 4 соединений, 8 конкурентных воркеров.
pool, _ := iproto.Dial(ctx, "tcp", addr, &iproto.PoolConfig{Size: 4})

var wg sync.WaitGroup
for w := range 8 {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        for i := range uint32(100) {
            resp, _ := pool.Call(ctx, 1, iproto.PackUint32(nil, i, 0))
            // проверка результата...
        }
    }(w)
}
wg.Wait()
fmt.Println("статистика:", pool.Stats())
```

## Бинарная сериализация

Упаковка и распаковка примитивных типов, BER-кодирование и структуры с тегами.

```bash
go run ./examples/packing
```

```go
// Примитивные типы.
data, _ := iproto.Pack(uint32(258))
// -> 02010000

// BER-кодирование (переменная длина).
data, _ = iproto.PackBER(uint32(128))
// -> 8100

// Структура со смешанными режимами.
type Message struct {
    Count uint32 `iproto:"ber"`
    Name  string
}
data, _ = iproto.Pack(Message{Count: 258, Name: "world"})

var msg Message
iproto.Unpack(data, &msg)
```

## Маршрутизация ServeMux

Регистрация обработчиков для разных кодов сообщений. Клиент вызывает каждый из них.

```bash
go run ./examples/mux
```

```go
mux := iproto.NewServeMux()

// Echo-обработчик на коде 1.
mux.Handle(1, iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
    _ = c.Send(ctx, iproto.ResponseTo(p, p.Data))
}))

// Time-обработчик на коде 2.
mux.Handle(2, iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
    now := time.Now().Format(time.RFC3339)
    _ = c.Send(ctx, iproto.ResponseTo(p, iproto.PackString(nil, now, iproto.ModeDefault)))
}))

srv := &iproto.Server{
    ChannelConfig: &iproto.ChannelConfig{Handler: mux},
}

// Клиент вызывает разные методы.
pool, _ := iproto.Dial(ctx, "tcp", addr, nil)
resp, _ := pool.Call(ctx, 1, echoData)   // -> echo-обработчик
resp, _ = pool.Call(ctx, 2, nil)          // -> time-обработчик
```
