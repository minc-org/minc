package exec

import (
	"context"
	"io"
)

// Cmd abstracts over running a command somewhere, this is useful for testing
type Cmd interface {
	// Run executes the command (like os/exec.Cmd.Run), it should return
	// a *RunError if there is any error
	Run() error
	// Each entry should be of the form "key=value"
	SetEnv(...string) Cmd
	SetStdin(io.Reader) Cmd
	SetStdout(io.Writer) Cmd
	SetStderr(io.Writer) Cmd
}

// Cmder abstracts over creating commands
type Cmder interface {
	// command, args..., just like os/exec.Cmd
	Command(string, ...string) Cmd
	CommandContext(context.Context, string, ...string) Cmd
}

// RunError represents an error running a Cmd
type RunError struct {
	Command []string // [Name Args...]
	Output  []byte   // Captured Stdout / Stderr of the command
	Inner   error    // Underlying error if any
}
