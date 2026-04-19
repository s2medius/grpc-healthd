package health

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/yourorg/grpc-healthd/internal/metrics"
	"github.com/yourorg/grpc-healthd/internal/probe"
)

// ServiceStatus holds the latest probe result for a named service.
type ServiceStatus struct {
	Name   string
	Status probe.Status
	LastChecked time.Time
}

// Checker runs probes on a schedule and tracks service health.
type Checker struct {
	mu       sync.RWMutex
	statuses map[string]*ServiceStatus
	probes   map[string]probe.Probe
	interval time.Duration
}

// NewChecker creates a Checker with the given check interval.
func NewChecker(interval time.Duration) *Checker {
	return &Checker{
		statuses: make(map[string]*ServiceStatus),
		probes:   make(map[string]probe.Probe),
		interval: interval,
	}
}

// Register adds a named probe to the checker.
func (c *Checker) Register(name string, p probe.Probe) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.probes[name] = p
	c.statuses[name] = &ServiceStatus{Name: name, Status: probe.StatusUnknown}
}

// GetStatus returns the latest status for a service.
func (c *Checker) GetStatus(name string) (ServiceStatus, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.statuses[name]
	if !ok {
		return ServiceStatus{}, false
	}
	return *s, true
}

// Run starts the periodic health check loop until ctx is cancelled.
func (c *Checker) Run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.checkAll(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (c *Checker) checkAll(ctx context.Context) {
	c.mu.RLock()
	names := make([]string, 0, len(c.probes))
	for n := range c.probes {
		names = append(names, n)
	}
	c.mu.RUnlock()

	for _, name := range names {
		c.mu.RLock()
		p := c.probes[name]
		c.mu.RUnlock()

		result := p.Check(ctx)
		metrics.RecordProbe(name, result)
		log.Printf("health check: service=%s status=%s", name, result.Status)

		c.mu.Lock()
		c.statuses[name] = &ServiceStatus{
			Name:        name,
			Status:      result.Status,
			LastChecked: time.Now(),
		}
		c.mu.Unlock()
	}
}
