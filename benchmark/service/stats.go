package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/stats"
)

type statsHandler struct {
	last       int
	rps        [60]int
	totalCount uint64
}

// HandleConn handle the connection
func (c *statsHandler) HandleConn(ctx context.Context, cs stats.ConnStats) {}

// TagConn exists to satisfy gRPC stats.Handler.
func (c *statsHandler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	return ctx
}

// TagRPC implements per-RPC context management.
func (c *statsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (c *statsHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	switch rs.(type) {
	case *stats.InPayload:
		s := rs.(*stats.InPayload)

		c.totalCount++
		idx := s.RecvTime.Second()
		if c.last != idx {
			c.rps[idx] = 0
			c.last = idx
		}

		c.rps[idx]++

	}
}

// print prints stats every 1 second
func (c *statsHandler) print() {
	go func() {
		for range time.Tick(1 * time.Second) {
			idx := time.Now().Second() - 1
			if idx == -1 {
				idx = 59
			}

			fmt.Println(c.rps[idx], "req/s")
		}
	}()
}
