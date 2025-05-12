package app

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func (a *Abdd) ExecuteCommand(t *Test) error {
	if t.Command == nil {
		return nil
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", t.Command.Command)
	} else {
		cmd = exec.Command("sh", "-c", t.Command.Command)
	}

	if t.Command.Directory != "" {
		cmd.Dir = t.Command.Directory
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("command execution failed: %s: %w", stderr.String(), err)
		}
		return fmt.Errorf("command execution failed: %w", err)
	}

	if t.Command.As != "" {
		a.Store[t.Command.As] = strings.TrimSpace(stdout.String())
	}
	return nil
}
