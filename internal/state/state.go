package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Action string

const (
	ActionStarted   Action = "started"
	ActionCompleted Action = "completed"
	ActionSnoozed   Action = "snoozed"
	ActionDismissed Action = "dismissed"
	ActionQuit      Action = "quit"
)

type State struct {
	LastBreakStartedAt   *time.Time `json:"last_break_started_at,omitempty"`
	LastBreakCompletedAt *time.Time `json:"last_break_completed_at,omitempty"`
	LastPopupShownAt     *time.Time `json:"last_popup_shown_at,omitempty"`
	SnoozedUntil         *time.Time `json:"snoozed_until,omitempty"`
	LastAction           Action     `json:"last_action,omitempty"`
	LastActionAt         *time.Time `json:"last_action_at,omitempty"`
}

func Load(path string) (State, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return State{}, nil
		}
		return State{}, err
	}
	var st State
	if err := json.Unmarshal(b, &st); err != nil {
		return State{}, err
	}
	return st, nil
}

func SaveAtomic(path string, st State) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
