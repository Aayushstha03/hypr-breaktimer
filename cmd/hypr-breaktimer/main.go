package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Aayushstha03/hypr-breaktimer/internal/logic"
	"github.com/Aayushstha03/hypr-breaktimer/internal/ui/popup"
)

func usage(w *os.File) {
	fmt.Fprintln(w, "usage: hypr-breaktimer [show|tick|status|bar|block|unblock]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "commands:")
	fmt.Fprintln(w, "  show    force opens the popup")
	fmt.Fprintln(w, "  tick    headless scheduler entrypoint (spawns popup when due)")
	fmt.Fprintln(w, "  status  print current config/state and next due time")
	fmt.Fprintln(w, "  bar     print one-line status for waybar")
	fmt.Fprintln(w, "  block   suppress scheduled popups for a duration (0 = dnd)")
	fmt.Fprintln(w, "  unblock disable dnd and clear any active block")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "flags:")
	fmt.Fprintln(w, "  -h, --help  show this help")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"status"}
	}

	if len(args) > 0 {
		switch strings.ToLower(args[0]) {
		case "-h", "--help", "help":
			usage(os.Stdout)
			os.Exit(0)
		}
	}

	switch strings.ToLower(args[0]) {
	case "__popup":
		// Internal entrypoint used by `show`/`tick` when launching a new terminal.
		code, err := popup.Run(context.Background(), popup.Options{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		os.Exit(code)
	case "show":
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := logic.Show(ctx); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case "tick":
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := logic.Tick(ctx); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case "status":
		if err := logic.Status(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case "bar":
		if err := logic.Bar(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case "block":
		if len(args) != 2 {
			usage(os.Stderr)
			os.Exit(2)
		}
		if err := logic.Block(args[1]); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case "unblock":
		if len(args) != 1 {
			usage(os.Stderr)
			os.Exit(2)
		}
		if err := logic.Unblock(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	default:
		usage(os.Stderr)
		os.Exit(2)
	}
}
