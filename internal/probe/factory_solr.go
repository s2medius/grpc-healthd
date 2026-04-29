package probe

import (
	"fmt"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func init() {
	registerFactory("solr", newSolrFromConfig)
}

// newSolrFromConfig constructs a SolrProbe from a ProbeConfig.
// It expects the address field to contain the Solr base URL (e.g. "http://localhost:8983").
// An optional timeout may be specified via the config; otherwise the default is used.
func newSolrFromConfig(cfg config.ProbeConfig) (Probe, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("solr probe %q: address is required", cfg.Name)
	}

	var opts []func(*SolrProbe)

	if cfg.Timeout != "" {
		d, err := time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("solr probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
		opts = append(opts, func(p *SolrProbe) {
			p.timeout = d
		})
	}

	return NewSolrProbe(cfg.Address, opts...), nil
}
