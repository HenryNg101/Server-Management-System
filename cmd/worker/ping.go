package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/HenryNg101/server-management-system/internal/model"
)

// --- Ping logic ---
func ping(ctx context.Context, s *model.Server) *model.Server {
	dialer := net.Dialer{Timeout: 2 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", s.IPv4Address, s.Port))

	if err != nil {
		s.Status = false
		return s
	}

	_ = conn.Close()

	s.Status = true
	return s
}
