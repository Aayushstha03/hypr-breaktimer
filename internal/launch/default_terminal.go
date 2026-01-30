package launch

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
)

type Options struct {
	AppID string
	Title string
}

func InDefaultTerminal(ctx context.Context, opts Options, argv []string) error {
	if len(argv) == 0 {
		return errors.New("no command provided")
	}

	// Prefer xdg-terminal-exec if available: it respects the user's default terminal
	// and supports setting app-id/title for reliable compositor matching.
	if _, err := exec.LookPath("xdg-terminal-exec"); err == nil {
		args := []string{}
		if opts.AppID != "" {
			args = append(args, "--app-id="+opts.AppID)
		}
		if opts.Title != "" {
			args = append(args, "--title="+opts.Title)
		}
		args = append(args, "--")
		args = append(args, argv...)
		cmd := exec.CommandContext(ctx, "xdg-terminal-exec", args...)
		cmd.Env = os.Environ()
		return cmd.Start()
	}

	// Fallback: respect $TERMINAL if it points to a command.
	if term := strings.TrimSpace(os.Getenv("TERMINAL")); term != "" {
		fields := strings.Fields(term)
		if len(fields) > 0 {
			cmdName := fields[0]
			cmdArgs := append(fields[1:], argv...)
			cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
			cmd.Env = os.Environ()
			return cmd.Start()
		}
	}

	return errors.New("no terminal launcher found (install xdg-terminal-exec or set $TERMINAL)")
}
