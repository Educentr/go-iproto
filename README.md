# go-iproto

[![Go Reference](https://pkg.go.dev/badge/github.com/Educentr/go-iproto.svg)](https://pkg.go.dev/github.com/Educentr/go-iproto)
[![CI](https://github.com/Educentr/go-iproto/actions/workflows/ci.yml/badge.svg)](https://github.com/Educentr/go-iproto/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go implementation of the iproto binary protocol — a lightweight, high-performance client/server framework for binary RPC communication.

**Documentation:** [English](https://educentr.github.io/go-iproto/) | [Русский](https://educentr.github.io/go-iproto/ru/)

## Features

- **Channel** — single connection with request-response correlation, ping/keepalive, graceful shutdown
- **Pool** — connection pool with load balancing, automatic reconnect, and rate limiting
- **Server / ServeMux** — accept incoming connections with handler routing by message code
- **Pack / Unpack** — binary serialization for integers, strings, structs with BER encoding support
- **netutil** — Dialer with retry/backoff and background dial
- **syncutil** — concurrency primitives: TaskGroup, TaskRunner, Throttle, Multitask

## Installation

```bash
go get github.com/Educentr/go-iproto
```

Requires Go 1.24+.

## Quick Start

**Server:**

```go
package main

import (
    "context"
    "log"

    "github.com/Educentr/go-iproto/iproto"
)

func main() {
    handler := iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
        // Echo: respond with the same data.
        _ = c.Send(ctx, iproto.ResponseTo(p, p.Data))
    })

    log.Fatal(iproto.ListenAndServe(context.Background(), "tcp", ":3301", handler))
}
```

**Client:**

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

    fmt.Printf("response: %s\n", resp)
}
```

## Packages

| Package | Description |
|---------|-------------|
| [`iproto`](https://pkg.go.dev/github.com/Educentr/go-iproto/iproto) | Core protocol: Channel, Pool, Server, ServeMux, Pack/Unpack, stream reader/writer |
| [`netutil`](https://pkg.go.dev/github.com/Educentr/go-iproto/netutil) | Dialer with retry/backoff and BackgroundDialer |
| [`syncutil`](https://pkg.go.dev/github.com/Educentr/go-iproto/syncutil) | TaskGroup, TaskRunner, Throttle, Multitask, Every, Each |

## Examples

See [`examples/`](examples/) for runnable programs:

- **[echo](examples/echo/)** — basic echo server and client
- **[pool](examples/pool/)** — connection pool with concurrent workers
- **[packing](examples/packing/)** — binary serialization with Pack/Unpack
- **[mux](examples/mux/)** — ServeMux request routing

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT — see [LICENSE](LICENSE) for details.

---

## На русском

Go-реализация бинарного протокола iproto — легковесный высокопроизводительный фреймворк для бинарного RPC.

**Документация:** [English](https://educentr.github.io/go-iproto/) | [Русский](https://educentr.github.io/go-iproto/ru/)

### Возможности

- **Channel** — одно соединение с корреляцией запрос-ответ, ping/keepalive, graceful shutdown
- **Pool** — пул соединений с балансировкой, автоматическим переподключением и rate limiting
- **Server / ServeMux** — приём входящих соединений с маршрутизацией по коду сообщения
- **Pack / Unpack** — бинарная сериализация целых чисел, строк, структур с поддержкой BER
- **netutil** — Dialer с retry/backoff и фоновым подключением
- **syncutil** — примитивы конкурентности: TaskGroup, TaskRunner, Throttle, Multitask

### Установка

```bash
go get github.com/Educentr/go-iproto
```

Требуется Go 1.24+.
