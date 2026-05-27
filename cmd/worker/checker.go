package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

// --- Ping logic ---
func ping(ctx context.Context, s Server) Result {
	start := time.Now()

	dialer := net.Dialer{Timeout: 2 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))

	latency := time.Since(start)

	if err != nil {
		return Result{
			ServerID: s.ID,
			Status:   0,
			Latency:  latency,
		}
	}

	_ = conn.Close()

	return Result{
		ServerID: s.ID,
		Status:   1,
		Latency:  latency,
	}
}
