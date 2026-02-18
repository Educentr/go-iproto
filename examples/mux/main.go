// Example mux demonstrates ServeMux request routing.
//
// The server registers multiple handlers for different message codes.
// The client sends requests to each handler and prints the responses.
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Educentr/go-iproto/iproto"
)

const (
	methodEcho = 1
	methodTime = 2
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("listening on", ln.Addr())

	// Create a ServeMux and register handlers for different message codes.
	mux := iproto.NewServeMux()

	// Echo handler: returns the data as-is.
	mux.Handle(methodEcho, iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
		_ = c.Send(ctx, iproto.ResponseTo(p, p.Data))
	}))

	// Time handler: returns the current server time as a string.
	mux.Handle(methodTime, iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
		now := time.Now().Format(time.RFC3339)
		_ = c.Send(ctx, iproto.ResponseTo(p, iproto.PackString(nil, now, iproto.ModeDefault)))
	}))

	srv := &iproto.Server{
		ChannelConfig: &iproto.ChannelConfig{
			Handler: mux,
		},
	}

	go func() {
		_ = srv.Serve(context.Background(), ln)
	}()

	// Client: connect and call each handler.
	pool, err := iproto.Dial(context.Background(), "tcp", ln.Addr().String(), &iproto.PoolConfig{
		Size: 1,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Call the echo handler.
	echoData := iproto.PackString(nil, "Hello, ServeMux!", iproto.ModeDefault)
	resp, err := pool.Call(context.Background(), methodEcho, echoData)
	if err != nil {
		log.Fatal("echo call error:", err)
	}

	var echoResult string
	if err := iproto.UnpackString(bytes.NewReader(resp), &echoResult, iproto.ModeDefault); err != nil {
		log.Fatal("echo unpack error:", err)
	}
	fmt.Printf("[echo]  -> %s\n", echoResult)

	// Call the time handler.
	resp, err = pool.Call(context.Background(), methodTime, nil)
	if err != nil {
		log.Fatal("time call error:", err)
	}

	var timeResult string
	if err := iproto.UnpackString(bytes.NewReader(resp), &timeResult, iproto.ModeDefault); err != nil {
		log.Fatal("time unpack error:", err)
	}
	fmt.Printf("[time]  -> %s\n", timeResult)

	// Shutdown.
	pool.Shutdown()
	<-pool.Done()
	ln.Close()

	fmt.Println("done")
}
