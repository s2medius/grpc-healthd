# grpc-healthd

Lightweight gRPC health check daemon with Prometheus metrics and configurable probes.

## Installation

```bash
go install github.com/grpc-healthd/grpc-healthd@latest
```

Or build from source:

```bash
git clone https://github.com/grpc-healthd/grpc-healthd.git && cd grpc-healthd && make build
```

## Usage

Define your probes in a YAML config file:

```yaml
probes:
  - name: my-service
    address: localhost:50051
    interval: 10s
    timeout: 2s

metrics:
  port: 9090
```

Run the daemon:

```bash
grpc-healthd --config config.yaml
```

Prometheus metrics will be available at `http://localhost:9090/metrics`.

```
# EXAMPLE METRICS
grpc_healthd_probe_status{service="my-service"} 1
grpc_healthd_probe_duration_seconds{service="my-service"} 0.004
```

## Configuration

| Field | Description | Default |
|-------|-------------|---------|
| `address` | gRPC server address | required |
| `interval` | Check interval | `30s` |
| `timeout` | Probe timeout | `5s` |

## License

MIT © grpc-healthd contributors