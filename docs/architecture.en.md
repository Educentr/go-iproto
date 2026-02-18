# Architecture

## Protocol Format

The iproto protocol uses a simple binary format with a 12-byte header followed by a variable-length payload.

### Header (12 bytes)

```
+--------+--------+--------+--------+
| Msg (4 bytes, little-endian)      |  Message/method code
+--------+--------+--------+--------+
| Len (4 bytes, little-endian)      |  Payload length in bytes
+--------+--------+--------+--------+
| Sync (4 bytes, little-endian)     |  Request-response correlation ID
+--------+--------+--------+--------+
| Data (Len bytes)                  |  Payload
+--------+--------+--------+--------+
```

- **Msg** — identifies the request type (method code). Used by `ServeMux` for routing.
- **Len** — length of the data payload.
- **Sync** — auto-incrementing ID that correlates requests with responses.

### Control Messages

| Code | Name | Purpose |
|------|------|---------|
| `0xff00` | `MessagePing` | Keepalive ping |
| `0xff01` | `MessageShutdown` | Graceful shutdown handshake |

## Channel

`Channel` manages a single TCP connection. It runs two goroutines — one for reading and one for writing packets.

```
┌──────────────────────────────────────┐
│              Channel                 │
│                                      │
│  ┌──────────┐     ┌──────────────┐   │
│  │  Reader   │     │   Writer     │   │
│  │ goroutine │     │  goroutine   │   │
│  └─────┬─────┘     └──────┬───────┘   │
│        │                  │           │
│        ▼                  ▲           │
│   ┌─────────┐      ┌─────────┐       │
│   │ Handler │      │   out   │       │
│   │  (user) │      │  chan   │       │
│   └─────────┘      └─────────┘       │
│                                      │
│   ┌─────────────────────────┐        │
│   │   Pending Store         │        │
│   │   map[msg<<32|sync] fn  │        │
│   └─────────────────────────┘        │
└──────────────────────────────────────┘
```

**Request-response flow:**

1. `Call()` registers a callback in the pending store keyed by `(Msg << 32) | Sync`.
2. The packet is sent to the `out` channel for the writer goroutine.
3. When a response arrives, the reader matches it by `(Msg, Sync)` and invokes the callback.

**Lifecycle:**

- `NewChannel(conn, config)` — create a channel.
- `Init()` — start reader/writer goroutines.
- `Shutdown()` — graceful shutdown with a bidirectional handshake.
- `Close()` — immediate close.
- `Done()` — channel that signals when fully closed.

**Features:**

- Configurable ping/keepalive with `PingInterval`.
- Idle timeout with `IdleTimeout`.
- Size limit per packet with `SizeLimit` (default 100 MB).
- Connection hijacking with `Hijack()`.

## Pool

`Pool` manages multiple Channels to a single server address.

```
┌──────────────────────────────────┐
│              Pool                │
│                                  │
│  ┌──────────┐  ┌──────────┐     │
│  │ Channel  │  │ Channel  │ ... │
│  │    #0    │  │    #1    │     │
│  └──────────┘  └──────────┘     │
│                                  │
│  Round-robin load balancing      │
│  Automatic reconnect             │
│  Rate limiting (token bucket)    │
└──────────────────────────────────┘
```

**Key behaviors:**

- **Load balancing** — round-robin via `NextChannel()`.
- **Auto reconnect** — uses `BackgroundDialer` with configurable backoff.
- **Rate limiting** — optional token bucket limiter with shaping or policing mode.
- **FailOnCutoff** — when enabled, calls fail immediately if no channels are online.

## Server

`Server` accepts incoming TCP connections and creates a Channel for each.

```go
srv := &iproto.Server{
    ChannelConfig: &iproto.ChannelConfig{
        Handler: myHandler,
    },
}
srv.ListenAndServe(ctx, "tcp", ":3301")
```

The `Accept` field allows customizing Channel creation (e.g., per-connection handler pools).

## ServeMux

`ServeMux` routes incoming packets to handlers based on `Header.Msg`:

```go
mux := iproto.NewServeMux()
mux.Handle(1, echoHandler)
mux.Handle(2, timeHandler)
```

When a packet arrives, `ServeMux.ServeIProto` looks up the handler by message code and dispatches.

## Serialization

### Default Mode (Little-Endian Fixed-Width)

| Type | Wire Size |
|------|-----------|
| `uint8` / `int8` | 1 byte |
| `uint16` / `int16` | 2 bytes |
| `uint32` / `int32` / `int` / `uint` | 4 bytes |
| `uint64` / `int64` | 8 bytes |
| `string` / `[]byte` | 4-byte length prefix + data |
| `struct` | fields encoded in order |
| `slice` | 4-byte count prefix + elements |

### BER Mode (Variable-Length)

BER encoding uses 7 bits per byte, with the high bit indicating continuation. Small values take fewer bytes:

- `5` → `05` (1 byte)
- `128` → `81 00` (2 bytes)
- `100000` → `86 8d 20` (3 bytes)

Use struct tags to mix modes:

```go
type Message struct {
    Count uint32 `iproto:"ber"`  // BER-encoded
    Name  string                 // Default (LE fixed)
}
```

## syncutil

Concurrency primitives used internally and available for user code:

- **TaskGroup** — execute up to N tasks concurrently. Returns channels for results.
- **TaskRunner** — run exactly one task; multiple callers subscribe to the same result.
- **Throttle** — rate-limit function execution to at most once per time period.
- **Multitask / Every / Each** — run N tasks in parallel with cancellation.
