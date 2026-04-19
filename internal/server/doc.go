// Package server provides a gRPC server implementing the gRPC Health Checking Protocol.
// It wraps the health.Checker to expose service health status over gRPC,
// returning SERVING, NOT_SERVING, or SERVICE_UNKNOWN based on probe results.
package server
