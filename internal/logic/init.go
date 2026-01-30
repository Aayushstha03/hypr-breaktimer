package logic

import (
	"os"
	"path/filepath"

	"github.com/Aayushstha03/hypr-timer/internal/config"
	"github.com/Aayushstha03/hypr-timer/internal/xdg"
)

func Init() error {
	configPath, err := xdg.ConfigFile()
	if err != nil {
		return err
	}
	statePath, err := xdg.StateFile()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return err
	}
	_, err = config.WriteDefaultIfMissing(configPath)
	if err != nil {
		return err
	}
	return nil
}
