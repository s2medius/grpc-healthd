package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/grpc-healthd/internal/config"
	"github.com/yourorg/grpc-healthd/internal/health"
	"github.com/yourorg/grpc-healthd/internal/probe"
	"github.com/yourorg/grpc-healthd/internal/server"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	checker := health.NewChecker()
	scheduler := health.NewScheduler(checker)

	for _, pc := range cfg.Probes {
		if err := pc.Validate(); err != nil {
			log.Printf("skipping invalid probe %q: %v", pc.Name, err)
			continue
		}
		p, err := probe.FromConfig(pc)
		if err != nil {
			log.Printf("skipping probe %q: %v", pc.Name, err)
			continue
		}
		scheduler.Register(pc.Name, p, pc.Interval)
	}

	scheduler.Start()
	defer scheduler.Stop()

	grpcSrv := server.NewHealthServer(checker)
	metricsSrv := server.NewMetricsServer(cfg.Metrics.Addr)

	go func() {
		if err := grpcSrv.ListenAndServe(cfg.GRPC.Addr); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	go func() {
		if err := metricsSrv.ListenAndServe(); err != nil {
			log.Fatalf("metrics server error: %v", err)
		}
	}()

	log.Printf("grpc-healthd started (grpc=%s metrics=%s)", cfg.GRPC.Addr, cfg.Metrics.Addr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("shutting down")
}
