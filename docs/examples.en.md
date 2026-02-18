# Examples

All examples are in the [`examples/`](https://github.com/Educentr/go-iproto/tree/main/examples) directory and can be run with `go run`.

## Echo Server

A basic echo server that returns any data it receives.

```bash
go run ./examples/echo
```

```go
// Server handler: echo data back.
handler := iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
    _ = c.Send(ctx, iproto.ResponseTo(p, p.Data))
})

srv := &iproto.Server{
    ChannelConfig: &iproto.ChannelConfig{Handler: handler},
}
go srv.Serve(ctx, ln)

// Client: send and receive.
pool, _ := iproto.Dial(ctx, "tcp", addr, &iproto.PoolConfig{Size: 1})
data := iproto.PackString(nil, "Hello, iproto!", iproto.ModeDefault)
resp, _ := pool.Call(ctx, 1, data)
```

## Connection Pool

Multiple goroutines sending requests through a pool of 4 connections. The server doubles each `uint32`.

```bash
go run ./examples/pool
```

```go
// Server: double the input.
handler := iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
    var in uint32
    iproto.UnpackUint32(bytes.NewReader(p.Data), &in, 0)
    _ = c.Send(ctx, iproto.ResponseTo(p, iproto.PackUint32(nil, in*2, 0)))
})

// Pool with 4 connections, 8 concurrent workers.
pool, _ := iproto.Dial(ctx, "tcp", addr, &iproto.PoolConfig{Size: 4})

var wg sync.WaitGroup
for w := range 8 {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        for i := range uint32(100) {
            resp, _ := pool.Call(ctx, 1, iproto.PackUint32(nil, i, 0))
            // verify result...
        }
    }(w)
}
wg.Wait()
fmt.Println("stats:", pool.Stats())
```

## Binary Packing

Pack and unpack primitive types, BER encoding, and structs with tags.

```bash
go run ./examples/packing
```

```go
// Primitive types.
data, _ := iproto.Pack(uint32(258))
// -> 02010000

// BER encoding (variable length).
data, _ = iproto.PackBER(uint32(128))
// -> 8100

// Struct with mixed modes.
type Message struct {
    Count uint32 `iproto:"ber"`
    Name  string
}
data, _ = iproto.Pack(Message{Count: 258, Name: "world"})

var msg Message
iproto.Unpack(data, &msg)
```

## ServeMux Routing

Register handlers for different message codes. The client calls each one.

```bash
go run ./examples/mux
```

```go
mux := iproto.NewServeMux()

// Echo handler on message code 1.
mux.Handle(1, iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
    _ = c.Send(ctx, iproto.ResponseTo(p, p.Data))
}))

// Time handler on message code 2.
mux.Handle(2, iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
    now := time.Now().Format(time.RFC3339)
    _ = c.Send(ctx, iproto.ResponseTo(p, iproto.PackString(nil, now, iproto.ModeDefault)))
}))

srv := &iproto.Server{
    ChannelConfig: &iproto.ChannelConfig{Handler: mux},
}

// Client calls different methods.
pool, _ := iproto.Dial(ctx, "tcp", addr, nil)
resp, _ := pool.Call(ctx, 1, echoData)   // -> echo handler
resp, _ = pool.Call(ctx, 2, nil)          // -> time handler
```
