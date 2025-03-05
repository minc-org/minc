package exec

import "context"

// DefaultCmder is a LocalCmder instance used for convenience, packages
// originally using os/exec.Command can instead use pkg/exec/exec.Command
// which forwards to this instance
var DefaultCmder = &LocalCmder{}

// Command is a convenience wrapper over DefaultCmder.Command
func Command(command string, args ...string) Cmd {
	return DefaultCmder.Command(command, args...)
}

// CommandContext is a convenience wrapper over DefaultCmder.CommandContext
func CommandContext(ctx context.Context, command string, args ...string) Cmd {
	return DefaultCmder.CommandContext(ctx, command, args...)
}
