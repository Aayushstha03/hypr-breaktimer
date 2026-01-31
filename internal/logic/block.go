package logic

import (
	"errors"
	"strconv"
	"time"

	"github.com/Aayushstha03/hypr-breaktimer/internal/state"
	"github.com/Aayushstha03/hypr-breaktimer/internal/xdg"
)

// Block suppresses scheduled popups (tick) for a duration.
//
// Passing 0 enables do-not-disturb mode until Unblock is called.
func Block(arg string) error {
	statePath, err := xdg.StateFile()
	if err != nil {
		return err
	}
	lockPath, err := xdg.LockFile()
	if err != nil {
		return err
	}

	lock, ok, err := state.Acquire(lockPath)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	defer func() { _ = lock.Release() }()

	d, err := parseDurationOrMinutes(arg)
	if err != nil {
		return err
	}
	if d < 0 {
		return errors.New("duration must be >= 0")
	}

	st, err := state.Load(statePath)
	if err != nil {
		return err
	}

	now := time.Now()
	if d == 0 {
		st.DoNotDisturb = true
		st.BlockedUntil = nil
		st.LastAction = state.ActionDndOn
		st.LastActionAt = &now
		return state.SaveAtomic(statePath, st)
	}

	until := now.Add(d)
	st.DoNotDisturb = false
	st.BlockedUntil = &until
	st.LastAction = state.ActionBlocked
	st.LastActionAt = &now
	return state.SaveAtomic(statePath, st)
}

// Unblock disables do-not-disturb and clears any active block window.
func Unblock() error {
	statePath, err := xdg.StateFile()
	if err != nil {
		return err
	}
	lockPath, err := xdg.LockFile()
	if err != nil {
		return err
	}

	lock, ok, err := state.Acquire(lockPath)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	defer func() { _ = lock.Release() }()

	st, err := state.Load(statePath)
	if err != nil {
		return err
	}

	now := time.Now()
	st.DoNotDisturb = false
	st.BlockedUntil = nil
	st.LastAction = state.ActionUnblocked
	st.LastActionAt = &now
	return state.SaveAtomic(statePath, st)
}

func parseDurationOrMinutes(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("missing duration")
	}
	if mins, err := strconv.Atoi(s); err == nil {
		return time.Duration(mins) * time.Minute, nil
	}
	return time.ParseDuration(s)
}
