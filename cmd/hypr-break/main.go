package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Aayushstha03/hypr-timer/internal/logic"
	"github.com/Aayushstha03/hypr-timer/internal/ui/popup"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"popup"}
	}

	switch strings.ToLower(args[0]) {
	case "popup":
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
	case "init":
		if err := logic.Init(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case "status":
		if err := logic.Status(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "usage: hypr-break [popup|show|tick|init|status]")
		os.Exit(2)
	}
}
