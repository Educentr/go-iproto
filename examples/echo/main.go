// Example echo demonstrates a basic iproto echo server and client.
//
// The server echoes back any data it receives. The client sends a message
// and prints the response.
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Educentr/go-iproto/iproto"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start TCP listener on a random port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("listening on", ln.Addr())

	// Server: echo handler returns the received data as-is.
	srv := &iproto.Server{
		ChannelConfig: &iproto.ChannelConfig{
			Handler: iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
				_ = c.Send(ctx, iproto.ResponseTo(p, p.Data))
			}),
		},
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Serve(ctx, ln)
	}()

	// Client: dial the server through a connection pool.
	pool, err := iproto.Dial(ctx, "tcp", ln.Addr().String(), &iproto.PoolConfig{
		Size: 1,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Send a string message. Pack it as a length-prefixed string.
	msg := "Hello, iproto!"
	data := iproto.PackString(nil, msg, iproto.ModeDefault)

	resp, err := pool.Call(ctx, 1, data)
	if err != nil {
		log.Fatal("call error:", err)
	}

	// Unpack the echoed response.
	var result string
	if err := iproto.UnpackString(bytes.NewReader(resp), &result, iproto.ModeDefault); err != nil {
		log.Fatal("unpack error:", err)
	}

	fmt.Println("sent:    ", msg)
	fmt.Println("received:", result)

	// Graceful shutdown.
	pool.Shutdown()
	ln.Close()
	<-pool.Done()

	fmt.Println("done")
}
