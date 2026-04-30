package probe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/grpc-healthd/internal/metrics"
)

// GraphQLProbe checks health by sending a simple introspection query to a GraphQL endpoint.
type GraphQLProbe struct {
	address string
	timeout time.Duration
}

// NewGraphQLProbe creates a new GraphQLProbe.
func NewGraphQLProbe(address string, timeout time.Duration) *GraphQLProbe {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &GraphQLProbe{address: address, timeout: timeout}
}

func (p *GraphQLProbe) Check(ctx context.Context) Status {
	start := time.Now()

	payload := map[string]string{"query": "{__typename}"}
	body, err := json.Marshal(payload)
	if err != nil {
		metrics.RecordProbe(p.address, "graphql", false, time.Since(start).Seconds())
		return StatusUnhealthy
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.address, bytes.NewReader(body))
	if err != nil {
		metrics.RecordProbe(p.address, "graphql", false, time.Since(start).Seconds())
		return StatusUnhealthy
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: p.timeout}
	resp, err := client.Do(req)
	duration := time.Since(start).Seconds()
	if err != nil {
		metrics.RecordProbe(p.address, "graphql", false, duration)
		return StatusUnhealthy
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		metrics.RecordProbe(p.address, "graphql", false, duration)
		return StatusUnhealthy
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		metrics.RecordProbe(p.address, "graphql", false, duration)
		return StatusUnhealthy
	}

	if _, hasErrors := result["errors"]; hasErrors {
		metrics.RecordProbe(p.address, "graphql", false, duration)
		return StatusUnhealthy
	}

	if _, hasData := result["data"]; !hasData {
		metrics.RecordProbe(p.address, "graphql", false, duration)
		return StatusUnhealthy
	}

	metrics.RecordProbe(p.address, "graphql", true, duration)
	return StatusHealthy
}

func (p *GraphQLProbe) String() string {
	return fmt.Sprintf("GraphQLProbe(%s)", p.address)
}
