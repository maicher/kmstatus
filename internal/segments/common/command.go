package common

import (
	"context"
	"io"
	"os/exec"
	"time"
)

// Maximum time an external program is allowed to run.
// A hanging program would otherwise block the segment and the whole status.
const CommandTimeout = 2 * time.Second

// RunCommand runs an external program, writing its standard output to stdout.
// The program is killed if it does not finish within CommandTimeout.
func RunCommand(stdout io.Writer, name string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), CommandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = stdout
	// Do not wait for children which inherited stdout after the program is killed.
	cmd.WaitDelay = 100 * time.Millisecond

	return cmd.Run()
}
