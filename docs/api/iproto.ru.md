# Пакет iproto

`import "github.com/Educentr/go-iproto/iproto"`

Полный API-справочник: [pkg.go.dev](https://pkg.go.dev/github.com/Educentr/go-iproto/iproto)

## Типы

### Header

```go
type Header struct {
    Msg  uint32 // Код сообщения/метода
    Len  uint32 // Длина данных
    Sync uint32 // ID корреляции запрос-ответ
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

Одно iproto-соединение.

```go
func NewChannel(conn net.Conn, config *ChannelConfig) *Channel
func RunChannel(conn net.Conn, config *ChannelConfig) (*Channel, error)
```

**Основные методы:**

| Метод | Описание |
|-------|----------|
| `Init()` | Запуск горутин чтения/записи |
| `Call(ctx, msg, data)` | Отправка запроса, ожидание ответа |
| `Notify(ctx, msg, data)` | Отправка уведомления (fire-and-forget) |
| `Send(ctx, packet)` | Отправка сырого пакета |
| `Shutdown()` | Graceful shutdown |
| `Close()` | Немедленное закрытие |
| `Done()` | Канал сигнала закрытия |
| `Stats()` | Получить `ChannelStats` |
| `Hijack()` | Перехват соединения |

### ChannelConfig

| Поле | Тип | По умолчанию | Описание |
|------|-----|-------------|----------|
| `Handler` | `Handler` | — | Обработчик входящих пакетов |
| `RequestTimeout` | `time.Duration` | 5с | Таймаут Call |
| `NoticeTimeout` | `time.Duration` | 1с | Таймаут Notify |
| `PingInterval` | `time.Duration` | 1м | Интервал keepalive |
| `IdleTimeout` | `time.Duration` | 0 | Таймаут неактивного соединения |
| `ShutdownTimeout` | `time.Duration` | 5с | Таймаут handshake shutdown |
| `SizeLimit` | `uint32` | 1e8 | Максимальный размер пакета |
| `WriteQueueSize` | `int` | 50 | Размер буфера вывода |
| `Init` | `func(context.Context, *Channel) error` | nil | Колбэк после инициализации |

## Pool

Пул соединений к одному серверу.

```go
func Dial(ctx context.Context, network, addr string, config *PoolConfig) (*Pool, error)
func NewPool(network, addr string, config *PoolConfig) *Pool
```

**Основные методы:**

| Метод | Описание |
|-------|----------|
| `Init(ctx)` | Инициализация одного соединения |
| `InitAll(ctx)` | Инициализация всех соединений |
| `Call(ctx, msg, data)` | Вызов через любой доступный канал |
| `Notify(ctx, msg, data)` | Уведомление через любой канал |
| `NextChannel(ctx)` | Следующий канал (round-robin) |
| `Online()` | Список активных каналов |
| `Stats()` | Получить `PoolStats` |
| `Shutdown()` | Graceful shutdown |
| `Close()` | Немедленное закрытие |

### PoolConfig

| Поле | Тип | По умолчанию | Описание |
|------|-----|-------------|----------|
| `Size` | `int` | 1 | Количество соединений |
| `ChannelConfig` | `*ChannelConfig` | — | Конфигурация каждого канала |
| `ConnectTimeout` | `time.Duration` | 1с | Макс. ожидание доступного канала |
| `RedialInterval` | `time.Duration` | 10мс | Начальная задержка повтора |
| `RedialTimeout` | `time.Duration` | 1ч | Макс. общее время дозвона |
| `RedialForever` | `bool` | false | Отключить таймаут дозвона |
| `FailOnCutoff` | `bool` | false | Быстрый отказ при оффлайне |
| `RateLimit` | `rate.Limit` | 0 | Скорость token bucket |
| `RateBurst` | `int` | 0 | Всплеск token bucket (0 = выкл.) |
| `RateWait` | `bool` | false | Shaping (true) или policing (false) |

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

**Функции по типам:**

```go
func PackUint8(w []byte, v uint8, mode PackMode) []byte
func PackUint16(w []byte, v uint16, mode PackMode) []byte
func PackUint32(w []byte, v uint32, mode PackMode) []byte
func PackUint64(w []byte, v uint64, mode PackMode) []byte
func PackString(w []byte, v string, mode PackMode) []byte
func PackBytes(w []byte, v []byte, mode PackMode) []byte
```

**Интерфейсы пользовательской сериализации:**

```go
type Packer interface {
    IprotoPack(w []byte, mode PackMode) ([]byte, error)
}

type Unpacker interface {
    IprotoUnpack(r *bytes.Reader, mode PackMode) error
}
```

## Ошибки

| Переменная | Описание |
|-----------|----------|
| `ErrTimeout` | Запрос превысил таймаут |
| `ErrDroppedConn` | Соединение было сброшено |
| `ErrStopped` | Канал/пул остановлен |
| `ErrHijacked` | Канал перехвачен |
| `ErrCutoff` | Пул оффлайн (режим FailOnCutoff) |
| `ErrPoolFull` | Нет места в пуле |
| `ErrNoChannel` | Таймаут получения канала |
| `ErrPolicied` | Отклонено лимитом скорости |

## Утилиты

```go
func ResponseTo(p Packet, b []byte) Packet
func CopyChannelConfig(c *ChannelConfig) *ChannelConfig
func CopyPoolConfig(c *PoolConfig) *PoolConfig
```
