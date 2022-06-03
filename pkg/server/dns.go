package server

import (
	"context"
	"net"
	"time"
)

func InitDNSResolver(addr string) *net.Resolver {
	var r = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 3 * time.Second,
			}
			return d.DialContext(ctx, "udp", addr)
		},
	}
	return r
}
