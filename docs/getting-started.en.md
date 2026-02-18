# Getting Started

## Installation

```bash
go get github.com/Educentr/go-iproto
```

Requires Go 1.24 or later.

## Your First Server

Create a simple echo server that returns any data it receives:

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

## Your First Client

Connect to the server and send a request:

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

## Package Structure

| Package | Import Path | Description |
|---------|------------|-------------|
| iproto | `github.com/Educentr/go-iproto/iproto` | Core protocol: Channel, Pool, Server, ServeMux, Pack/Unpack |
| netutil | `github.com/Educentr/go-iproto/netutil` | Dialer with retry/backoff, BackgroundDialer |
| syncutil | `github.com/Educentr/go-iproto/syncutil` | TaskGroup, TaskRunner, Throttle, Multitask |

## Key Concepts

- **Channel** manages a single TCP connection. It handles reading/writing packets, ping/keepalive, and request-response correlation via a `Sync` ID.
- **Pool** manages multiple Channels to one server. It provides load balancing (round-robin), automatic reconnection, and optional rate limiting.
- **Server** accepts incoming TCP connections and creates a Channel for each one.
- **ServeMux** routes incoming packets to handlers based on the message code (`Header.Msg`).
- **Pack/Unpack** provides binary serialization. Integers are little-endian by default, with optional BER (variable-length) encoding.

## Next Steps

- [Architecture](architecture.md) — protocol internals
- [API Reference](api/iproto.md) — detailed API documentation
- [Examples](examples.md) — more runnable examples
