# go-iproto

Go implementation of the iproto binary protocol — a lightweight, high-performance framework for binary RPC communication.

## Features

- **Channel** — single connection with request-response correlation, ping/keepalive, graceful shutdown
- **Pool** — connection pool with load balancing, automatic reconnect, and rate limiting
- **Server / ServeMux** — accept connections with handler routing by message code
- **Pack / Unpack** — binary serialization for integers, strings, structs with BER encoding
- **netutil** — Dialer with retry/backoff and background dial
- **syncutil** — concurrency primitives: TaskGroup, TaskRunner, Throttle, Multitask

## Installation

```bash
go get github.com/Educentr/go-iproto
```

Requires Go 1.24+.

## Quick Example

```go
// Server
handler := iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
    _ = c.Send(ctx, iproto.ResponseTo(p, p.Data))
})
iproto.ListenAndServe(ctx, "tcp", ":3301", handler)

// Client
pool, _ := iproto.Dial(ctx, "tcp", "127.0.0.1:3301", nil)
resp, _ := pool.Call(ctx, 1, []byte("hello"))
```

## Links

- [Getting Started](getting-started.md) — step-by-step guide
- [Architecture](architecture.md) — protocol design and internals
- [API Reference](api/iproto.md) — types, functions, constants
- [Examples](examples.md) — runnable code examples
- [GitHub](https://github.com/Educentr/go-iproto)
- [pkg.go.dev](https://pkg.go.dev/github.com/Educentr/go-iproto)
