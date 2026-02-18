# netutil Package

`import "github.com/Educentr/go-iproto/netutil"`

Full API reference: [pkg.go.dev](https://pkg.go.dev/github.com/Educentr/go-iproto/netutil)

## Dialer

`Dialer` establishes a TCP connection with retry and exponential backoff.

```go
type Dialer struct {
    Network         string
    Addr            string
    Timeout         time.Duration        // Single dial timeout
    LoopInterval    time.Duration        // Initial backoff (default 50ms)
    MaxLoopInterval time.Duration        // Max backoff
    NetDial         func(ctx, net, addr) (net.Conn, error)
    OnAttempt       func(err error)      // Called after each attempt
    Logf            func(string, ...any)
    Debugf          func(string, ...any)
}
```

**Methods:**

| Method | Description |
|--------|-------------|
| `Dial(ctx) (net.Conn, error)` | Dial with retries until success or context cancellation |

## BackgroundDialer

`BackgroundDialer` wraps `Dialer` for concurrent background dial management. Multiple callers can subscribe to the same dial operation.

```go
type BackgroundDialer struct {
    Dialer     *Dialer
    TaskGroup  *syncutil.TaskGroup
    TaskRunner *syncutil.TaskRunner
}
```

**Methods:**

| Method | Description |
|--------|-------------|
| `Dial(ctx, callback) <-chan error` | Start or join background dial; callback receives the connection |
| `Cancel()` | Cancel current dial routine |
| `SetDeadline(t time.Time)` | Set absolute dial deadline |
| `SetDeadlineAtLeast(t time.Time) time.Time` | Set deadline if current is earlier |

### Usage Example

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
    // use conn
})

<-done
```
