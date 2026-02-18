# go-iproto

Go implementation of the iproto binary protocol client and server.

## Packages

- **iproto** — iproto protocol client/server: Channel (single connection), Pool (connection pool), Server, binary packet reader/writer
- **netutil** — Dialer with retry/backoff and background dial support
- **syncutil** — concurrency primitives: TaskGroup, TaskRunner, Throttle

## Installation

```bash
go get github.com/Educentr/go-iproto
```

## License

MIT
