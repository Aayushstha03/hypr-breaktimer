package logic

import (
	"context"
	"os/exec"
)

func notify(ctx context.Context, title, body string) error {
	cmd := exec.CommandContext(ctx, "notify-send", "-t", "5000", title, body)
	return cmd.Run()
}
