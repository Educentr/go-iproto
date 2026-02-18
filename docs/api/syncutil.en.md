# syncutil Package

`import "github.com/Educentr/go-iproto/syncutil"`

Full API reference: [pkg.go.dev](https://pkg.go.dev/github.com/Educentr/go-iproto/syncutil)

## TaskGroup

Execute up to N concurrent tasks. Reuses slots as tasks complete.

```go
type TaskGroup struct {
    N    int    // Max concurrent tasks (default 1)
    Goer GoerFn // Custom goroutine launcher (optional)
}
```

**Methods:**

| Method | Description |
|--------|-------------|
| `Do(ctx, n, task) []<-chan error` | Start up to n tasks; returns result channels |
| `Cancel()` | Cancel all running tasks |

### Example

```go
tg := &syncutil.TaskGroup{N: 4}

results := tg.Do(ctx, 4, func(ctx context.Context, i int) error {
    // task i runs concurrently
    return nil
})

for _, ch := range results {
    if err := <-ch; err != nil {
        log.Fatal(err)
    }
}
```

## TaskRunner

Run exactly one task at a time. Multiple callers subscribe to the same result.

```go
type TaskRunner struct{}
```

**Methods:**

| Method | Description |
|--------|-------------|
| `Do(ctx, task) <-chan error` | Start task (or subscribe if already running) |
| `Cancel()` | Cancel current task |

### Example

```go
var runner syncutil.TaskRunner

// Both goroutines receive the same result.
ch1 := runner.Do(ctx, func(ctx context.Context) error {
    return expensiveOperation(ctx)
})
ch2 := runner.Do(ctx, func(ctx context.Context) error {
    return expensiveOperation(ctx) // not started — subscribes to existing
})

err1 := <-ch1
err2 := <-ch2 // same value as err1
```

## Throttle

Limit function execution to at most once per time period.

```go
type Throttle struct {
    Period time.Duration
}
```

**Methods:**

| Method | Description |
|--------|-------------|
| `Next() bool` | Returns true if the period has elapsed |
| `Reset()` | Clear the throttle |
| `Set(time.Time)` | Set specific throttle time |

## Multitask

Run N tasks in parallel with cancellation support.

```go
type Multitask struct {
    ContinueOnError bool   // Don't cancel others on Goer failure
    Goer            GoerFn // Custom goroutine launcher
}
```

**Methods:**

| Method | Description |
|--------|-------------|
| `Do(ctx, n, actor) error` | Run actor n times in parallel |

## Helper Functions

```go
// Every runs n tasks; cancels all on first error.
func Every(ctx context.Context, n int, actor func(context.Context, int) error) error

// Each runs n tasks; waits for all to complete.
func Each(ctx context.Context, n int, actor func(context.Context, int))
```
