# iproto Package

`import "github.com/Educentr/go-iproto/iproto"`

Full API reference: [pkg.go.dev](https://pkg.go.dev/github.com/Educentr/go-iproto/iproto)

## Types

### Header

```go
type Header struct {
    Msg  uint32 // Message/method code
    Len  uint32 // Payload length
    Sync uint32 // Request-response correlation ID
}
```

### Packet

```go
type Packet struct {
    Header Header
    Data   []byte
}
```

### Handler

```go
type Handler interface {
    ServeIProto(ctx context.Context, c Conn, p Packet)
}

type HandlerFunc func(context.Context, Conn, Packet)
```

### Conn

```go
type Conn interface {
    Sender
    Closer
    GetBytes(n int) []byte
    PutBytes(p []byte)
    RemoteAddr() net.Addr
    LocalAddr() net.Addr
}
```

### Sender

```go
type Sender interface {
    Call(ctx context.Context, message uint32, data []byte) ([]byte, error)
    Notify(ctx context.Context, message uint32, data []byte) error
    Send(ctx context.Context, packet Packet) error
}
```

## Channel

Single iproto connection.

```go
func NewChannel(conn net.Conn, config *ChannelConfig) *Channel
func RunChannel(conn net.Conn, config *ChannelConfig) (*Channel, error)
```

**Key methods:**

| Method | Description |
|--------|-------------|
| `Init()` | Start reader/writer goroutines |
| `Call(ctx, msg, data)` | Send request, wait for response |
| `Notify(ctx, msg, data)` | Send fire-and-forget notification |
| `Send(ctx, packet)` | Send raw packet |
| `Shutdown()` | Graceful shutdown |
| `Close()` | Immediate close |
| `Done()` | Closure signal channel |
| `Stats()` | Get `ChannelStats` |
| `Hijack()` | Take over the connection |

### ChannelConfig

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Handler` | `Handler` | — | Incoming packet handler |
| `RequestTimeout` | `time.Duration` | 5s | Call timeout |
| `NoticeTimeout` | `time.Duration` | 1s | Notify timeout |
| `PingInterval` | `time.Duration` | 1m | Keepalive interval |
| `IdleTimeout` | `time.Duration` | 0 | Idle connection timeout |
| `ShutdownTimeout` | `time.Duration` | 5s | Shutdown handshake timeout |
| `SizeLimit` | `uint32` | 1e8 | Max packet size |
| `WriteQueueSize` | `int` | 50 | Output buffer size |
| `Init` | `func(context.Context, *Channel) error` | nil | Post-init callback |

## Pool

Connection pool to a single server.

```go
func Dial(ctx context.Context, network, addr string, config *PoolConfig) (*Pool, error)
func NewPool(network, addr string, config *PoolConfig) *Pool
```

**Key methods:**

| Method | Description |
|--------|-------------|
| `Init(ctx)` | Initialize one connection |
| `InitAll(ctx)` | Initialize all connections |
| `Call(ctx, msg, data)` | Call via any available channel |
| `Notify(ctx, msg, data)` | Notify via any channel |
| `NextChannel(ctx)` | Get next channel (round-robin) |
| `Online()` | List of active channels |
| `Stats()` | Get `PoolStats` |
| `Shutdown()` | Graceful shutdown |
| `Close()` | Immediate close |

### PoolConfig

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Size` | `int` | 1 | Number of connections |
| `ChannelConfig` | `*ChannelConfig` | — | Config for each channel |
| `ConnectTimeout` | `time.Duration` | 1s | Max wait for available channel |
| `RedialInterval` | `time.Duration` | 10ms | Initial retry delay |
| `RedialTimeout` | `time.Duration` | 1h | Max total dial time |
| `RedialForever` | `bool` | false | Disable dial timeout |
| `FailOnCutoff` | `bool` | false | Fail fast when offline |
| `RateLimit` | `rate.Limit` | 0 | Token bucket rate |
| `RateBurst` | `int` | 0 | Token bucket burst (0 = disabled) |
| `RateWait` | `bool` | false | Shape (true) or police (false) |

## Server

```go
type Server struct {
    Accept        AcceptFn
    ChannelConfig *ChannelConfig
    Log           Logger
    OnClose       []func()
    OnShutdown    []func()
}

func ListenAndServe(ctx context.Context, network, addr string, h Handler) error
```

## ServeMux

```go
func NewServeMux() *ServeMux
func (s *ServeMux) Handle(message uint32, handler Handler)
func (s *ServeMux) Handler(message uint32) Handler
```

## Pack / Unpack

```go
func Pack(value ...any) ([]byte, error)
func PackBER(value ...any) ([]byte, error)
func Append(data []byte, value ...any) ([]byte, error)
func Unpack(data []byte, value ...any) ([]byte, error)
func UnpackBER(data []byte, value ...any) ([]byte, error)
```

**Per-type functions:**

```go
func PackUint8(w []byte, v uint8, mode PackMode) []byte
func PackUint16(w []byte, v uint16, mode PackMode) []byte
func PackUint32(w []byte, v uint32, mode PackMode) []byte
func PackUint64(w []byte, v uint64, mode PackMode) []byte
func PackString(w []byte, v string, mode PackMode) []byte
func PackBytes(w []byte, v []byte, mode PackMode) []byte
```

**Custom serialization interfaces:**

```go
type Packer interface {
    IprotoPack(w []byte, mode PackMode) ([]byte, error)
}

type Unpacker interface {
    IprotoUnpack(r *bytes.Reader, mode PackMode) error
}
```

## Errors

| Variable | Description |
|----------|-------------|
| `ErrTimeout` | Request exceeded timeout |
| `ErrDroppedConn` | Connection was dropped |
| `ErrStopped` | Channel/pool is stopped |
| `ErrHijacked` | Channel was hijacked |
| `ErrCutoff` | Pool offline (FailOnCutoff mode) |
| `ErrPoolFull` | No space in pool |
| `ErrNoChannel` | Timeout getting channel |
| `ErrPolicied` | Rate limit rejected |

## Utility

```go
func ResponseTo(p Packet, b []byte) Packet
func CopyChannelConfig(c *ChannelConfig) *ChannelConfig
func CopyPoolConfig(c *PoolConfig) *PoolConfig
```
