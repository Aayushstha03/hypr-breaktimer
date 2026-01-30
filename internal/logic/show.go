package logic

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"time"

	"github.com/Aayushstha03/hypr-timer/internal/config"
	"github.com/Aayushstha03/hypr-timer/internal/launch"
	"github.com/Aayushstha03/hypr-timer/internal/state"
	"github.com/Aayushstha03/hypr-timer/internal/xdg"
)

// Show forces the popup to open in a new terminal.
func Show(ctx context.Context) error {
	configPath, err := xdg.ConfigFile()
	if err != nil {
		return err
	}
	statePath, err := xdg.StateFile()
	if err != nil {
		return err
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	// Record that we attempted to show a popup.
	st, _ := state.Load(statePath)
	now := time.Now()
	st.LastPopupShownAt = &now
	_ = state.SaveAtomic(statePath, st)

	exe, err := os.Executable()
	if err != nil || exe == "" {
		exe, err = exec.LookPath("hypr-break")
		if err != nil {
			return errors.New("cannot find hypr-break executable")
		}
	}

	return launch.InDefaultTerminal(ctx, launch.Options{AppID: cfg.Launch.AppID, Title: cfg.Launch.Title}, []string{exe, "popup"})
}
