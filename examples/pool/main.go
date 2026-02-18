// Example pool demonstrates a connection pool with concurrent workers.
//
// The server doubles every uint32 it receives. Multiple goroutines send
// requests through a pool of 4 connections and verify results.
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Educentr/go-iproto/iproto"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("listening on", ln.Addr())

	// Server: double the input uint32.
	srv := &iproto.Server{
		ChannelConfig: &iproto.ChannelConfig{
			Handler: iproto.HandlerFunc(func(ctx context.Context, c iproto.Conn, p iproto.Packet) {
				var in uint32
				if err := iproto.UnpackUint32(bytes.NewReader(p.Data), &in, 0); err != nil {
					return
				}
				_ = c.Send(ctx, iproto.ResponseTo(p, iproto.PackUint32(nil, in*2, 0)))
			}),
		},
	}

	go func() {
		_ = srv.Serve(context.Background(), ln)
	}()

	// Create a pool with 4 connections.
	pool, err := iproto.Dial(context.Background(), "tcp", ln.Addr().String(), &iproto.PoolConfig{
		Size:           4,
		ConnectTimeout: 10 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Launch 8 concurrent workers, each sending 100 requests.
	const (
		workers  = 8
		requests = 100
	)

	var wg sync.WaitGroup

	for w := range workers {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			for i := range uint32(requests) {
				resp, err := pool.Call(context.Background(), 1, iproto.PackUint32(nil, i, 0))
				if err != nil {
					log.Printf("worker %d: call error: %v", workerID, err)
					return
				}

				var result uint32
				if err := iproto.UnpackUint32(bytes.NewReader(resp), &result, 0); err != nil {
					log.Printf("worker %d: unpack error: %v", workerID, err)
					return
				}

				if result != i*2 {
					log.Printf("worker %d: got %d, want %d", workerID, result, i*2)
					return
				}
			}
		}(w)
	}

	wg.Wait()

	// Print pool statistics.
	stats := pool.Stats()
	fmt.Printf("completed: %d workers x %d requests = %d total\n", workers, requests, workers*requests)
	fmt.Printf("stats: sent=%d received=%d online=%d\n",
		stats.PacketsSent, stats.PacketsReceived, stats.Online)

	pool.Shutdown()
	<-pool.Done()

	fmt.Println("done")
}
