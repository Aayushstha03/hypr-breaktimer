package xdg

import (
	"errors"
	"os"
	"path/filepath"
)

func ConfigFile() (string, error) {
	base, err := configHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "hypr-break", "config.toml"), nil
}

func StateFile() (string, error) {
	base, err := stateHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "hypr-break", "state.json"), nil
}

func LockFile() (string, error) {
	base, err := stateHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "hypr-break", "lock"), nil
}

func configHome() (string, error) {
	if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
		return v, nil
	}
	v, err := os.UserConfigDir()
	if err == nil && v != "" {
		return v, nil
	}
	h, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if h == "" {
		return "", errors.New("cannot determine home directory")
	}
	return filepath.Join(h, ".config"), nil
}

func stateHome() (string, error) {
	if v := os.Getenv("XDG_STATE_HOME"); v != "" {
		return v, nil
	}
	h, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if h == "" {
		return "", errors.New("cannot determine home directory")
	}
	return filepath.Join(h, ".local", "state"), nil
}
