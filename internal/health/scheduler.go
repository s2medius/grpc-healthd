package health

import (
	"context"
	"log"
	"time"

	"github.com/yourorg/grpc-healthd/internal/config"
	"github.com/yourorg/grpc-healthd/internal/probe"
)

// Scheduler periodically runs probes and updates the Checker.
type Scheduler struct {
	checker  *Checker
	probes   []config.ProbeConfig
	stopCh   chan struct{}
}

// NewScheduler creates a Scheduler for the given probes and checker.
func NewScheduler(checker *Checker, probes []config.ProbeConfig) *Scheduler {
	return &Scheduler{
		checker: checker,
		probes:  probes,
		stopCh:  make(chan struct{}),
	}
}

// Start launches a goroutine per probe that ticks on the configured interval.
func (s *Scheduler) Start(ctx context.Context) {
	for _, pc := range s.probes {
		pc := pc
		p, err := probe.FromConfig(pc)
		if err != nil {
			log.Printf("scheduler: skipping probe %q: %v", pc.Name, err)
			continue
		}
		s.checker.Register(pc.Name)
		go s.runLoop(ctx, pc.Name, pc.Interval, p)
	}
}

func (s *Scheduler) runLoop(ctx context.Context, name string, interval time.Duration, p probe.Probe) {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case <-ticker.C:
			result := p.Probe(ctx)
			s.checker.Update(name, result)
			log.Printf("scheduler: probe %q status=%s duration=%s", name, result.Status, result.Duration)
		}
	}
}

// Stop signals all probe loops to exit.
func (s *Scheduler) Stop() {
	close(s.stopCh)
}
