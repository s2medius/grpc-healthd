package probe

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// ExecProbe runs a local command and reports healthy if it exits with code 0.
type ExecProbe struct {
	command string
	args    []string
	timeout time.Duration
}

// NewExecProbe creates an ExecProbe. If timeout is zero, DefaultTimeout is used.
func NewExecProbe(command string, args []string, timeout time.Duration) *ExecProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &ExecProbe{
		command: command,
		args:    args,
		timeout: timeout,
	}
}

// Probe executes the command and returns a Result.
func (e *ExecProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.command, e.args...)
	err := cmd.Run()
	duration := time.Since(start)

	if err != nil {
		msg := err.Error()
		if ctx.Err() == context.DeadlineExceeded {
			msg = fmt.Sprintf("command timed out after %s", e.timeout)
		}
		return Result{
			Status:   StatusUnhealthy,
			Duration: duration,
			Message:  msg,
		}
	}

	return Result{
		Status:   StatusHealthy,
		Duration: duration,
		Message:  "exit 0",
	}
}
