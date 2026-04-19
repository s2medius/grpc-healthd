package probe

import (
	"context"
	"fmt"
	"net"
	"time"
)

// DNSProbe checks that a hostname resolves via DNS.
type DNSProbe struct {
	host    string
	timeout time.Duration
	resolver *net.Resolver
}

// NewDNSProbe creates a DNSProbe for the given host.
// timeout of 0 defaults to 5 seconds.
func NewDNSProbe(host string, timeout time.Duration) *DNSProbe {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &DNSProbe{
		host:     host,
		timeout:  timeout,
		resolver: net.DefaultResolver,
	}
}

// Probe performs a DNS lookup and returns a Result.
func (d *DNSProbe) Probe() Result {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	addrs, err := d.resolver.LookupHost(ctx, d.host)
	duration := time.Since(start)

	if err != nil || len(addrs) == 0 {
		msg := fmt.Sprintf("dns lookup failed for %s", d.host)
		if err != nil {
			msg = err.Error()
		}
		return Result{Status: StatusUnhealthy, Duration: duration, Message: msg}
	}

	return Result{
		Status:   StatusHealthy,
		Duration: duration,
		Message:  fmt.Sprintf("resolved %s to %d address(es)", d.host, len(addrs)),
	}
}
