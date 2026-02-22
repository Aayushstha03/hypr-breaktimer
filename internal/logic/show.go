package logic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Aayushstha03/hypr-breaktimer/internal/config"
	"github.com/Aayushstha03/hypr-breaktimer/internal/launch"
	"github.com/Aayushstha03/hypr-breaktimer/internal/state"
	"github.com/Aayushstha03/hypr-breaktimer/internal/xdg"
)

type showMode int

const (
	showForce showMode = iota
	showDueOnly
)

// Show forces the popup to open in a new terminal.
func Show(ctx context.Context) error {
	return show(ctx, showForce)
}

func show(ctx context.Context, mode showMode) error {
	configPath, err := xdg.ConfigFile()
	if err != nil {
		return err
	}
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

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	st, err := state.Load(statePath)
	if err != nil {
		return err
	}

	now := time.Now()
	if mode == showDueOnly {
		if st.DoNotDisturb {
			return nil
		}
		if st.BlockedUntil != nil && now.Before(*st.BlockedUntil) {
			return nil
		}

		quietWindows, err := config.ParseQuietWindows(cfg.QuietHours.Windows)
		if err != nil {
			return err
		}
		if cfg.QuietHours.Enabled && config.InQuietHours(now, quietWindows) {
			return nil
		}

		// First run: establish a reference point so the interval can elapse.
		if st.LastBreakCompletedAt == nil && st.LastBreakStartedAt == nil && st.LastPopupShownAt == nil {
			st.LastPopupShownAt = &now
			return state.SaveAtomic(statePath, st)
		}

		ref := now
		if st.LastBreakCompletedAt != nil {
			ref = *st.LastBreakCompletedAt
		} else if st.LastBreakStartedAt != nil {
			ref = *st.LastBreakStartedAt
		} else if st.LastPopupShownAt != nil {
			ref = *st.LastPopupShownAt
		}

		nextDue := ref.Add(cfg.Schedule.WorkInterval.Duration())

		// Occasional nudge notifications (not affected by snooze).
		// Suppressed by DND, quiet hours and blocking.
		if st.LastNudgeNotifiedAt == nil || now.Sub(*st.LastNudgeNotifiedAt) >= 10*time.Minute {
			nearPopup := false
			if st.LastPopupShownAt != nil && now.Sub(*st.LastPopupShownAt) < 3*time.Minute {
				nearPopup = true
			}
			// Avoid nudges right before a popup is due (unless snoozed).
			if (st.SnoozedUntil == nil || now.After(*st.SnoozedUntil)) && now.Before(nextDue) && nextDue.Sub(now) <= 3*time.Minute {
				nearPopup = true
			}
			if !nearPopup {
				name, msg := RandomNudgeMessage()
				if err := notify(ctx, name, msg); err == nil {
					st.LastNudgeNotifiedAt = &now
					_ = state.SaveAtomic(statePath, st)
				}
			}
		}

		// Snooze suppresses popup and break-soon reminder.
		if st.SnoozedUntil != nil && now.Before(*st.SnoozedUntil) {
			return nil
		}

		// Break-soon reminder: once, 2 minutes before due.
		soonAt := nextDue.Add(-2 * time.Minute)
		if (now.Equal(soonAt) || now.After(soonAt)) && now.Before(nextDue) {
			already := st.LastBreakSoonNotifiedFor != nil && st.LastBreakSoonNotifiedFor.Equal(nextDue)
			if !already {
				if err := notify(ctx, "hypr-breaktimer", "Break in 2 minutes"); err == nil {
					due := nextDue
					st.LastBreakSoonNotifiedFor = &due
					_ = state.SaveAtomic(statePath, st)
				}
			}
		}

		if now.Before(nextDue) {
			return nil
		}

		if st.LastPopupShownAt != nil {
			minGap := cfg.Schedule.MinTimeBetweenPopups.Duration()
			if minGap > 0 && now.Sub(*st.LastPopupShownAt) < minGap {
				return nil
			}
		}
		// no dismiss cooldown; user actions reset the schedule via break completion.
	}

	// Mark shown before spawning to prevent multiple popups.
	st.LastPopupShownAt = &now
	if err := state.SaveAtomic(statePath, st); err != nil {
		return err
	}

	exe, err := os.Executable()
	if err != nil || exe == "" {
		exe, err = exec.LookPath("hypr-breaktimer")
		if err != nil {
			return errors.New("cannot find hypr-breaktimer executable")
		}
	}

	if err := launch.InDefaultTerminal(ctx, launch.Options{AppID: cfg.Launch.AppID, Title: cfg.Launch.Title}, []string{exe, "__popup"}); err != nil {
		return fmt.Errorf("launch popup: %w", err)
	}
	return nil
}
