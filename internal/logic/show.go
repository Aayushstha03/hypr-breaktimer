package logic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	debugLog := strings.EqualFold(cfg.Debug.LogLevel, "debug")
	if mode == showDueOnly {
		if st.DoNotDisturb {
			if debugLog {
				fmt.Fprintln(os.Stderr, "tick: blocked_by=dnd")
			}
			return nil
		}
		if st.BlockedUntil != nil && now.Before(*st.BlockedUntil) {
			if debugLog {
				fmt.Fprintf(os.Stderr, "tick: blocked_by=blocked_until until=%s\n", st.BlockedUntil.Format(time.RFC3339))
			}
			return nil
		}

		quietWindows, err := config.ParseQuietWindows(cfg.QuietHours.Windows)
		if err != nil {
			return err
		}
		if cfg.QuietHours.Enabled && config.InQuietHours(now, quietWindows) {
			if debugLog {
				fmt.Fprintln(os.Stderr, "tick: blocked_by=quiet_hours")
			}
			return nil
		}
		if st.SnoozedUntil != nil && now.Before(*st.SnoozedUntil) {
			if debugLog {
				fmt.Fprintf(os.Stderr, "tick: blocked_by=snoozed until=%s\n", st.SnoozedUntil.Format(time.RFC3339))
			}
			return nil
		}

		// First run: establish a reference point so the interval can elapse.
		if st.LastBreakCompletedAt == nil && st.LastBreakStartedAt == nil && st.LastPopupShownAt == nil {
			st.LastPopupShownAt = &now
			if debugLog {
				fmt.Fprintf(os.Stderr, "tick: first_run set_last_popup_shown_at=%s\n", now.Format(time.RFC3339))
			}
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
		if now.Before(nextDue) {
			if debugLog {
				fmt.Fprintf(os.Stderr, "tick: not_due now=%s next_due=%s\n", now.Format(time.RFC3339), nextDue.Format(time.RFC3339))
			}
			return nil
		}

		if st.LastPopupShownAt != nil {
			minGap := cfg.Schedule.MinTimeBetweenPopups.Duration()
			if minGap > 0 && now.Sub(*st.LastPopupShownAt) < minGap {
				if debugLog {
					fmt.Fprintf(os.Stderr, "tick: blocked_by=min_time_between_popups last_popup=%s min_gap=%s\n", st.LastPopupShownAt.Format(time.RFC3339), minGap)
				}
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
	if debugLog {
		fmt.Fprintf(os.Stderr, "tick: launching app_id=%q title=%q\n", cfg.Launch.AppID, cfg.Launch.Title)
	}

	exe, err := os.Executable()
	if err != nil || exe == "" {
		exe, err = exec.LookPath("hypr-breaktimer")
		if err != nil {
			return errors.New("cannot find hypr-breaktimer executable")
		}
	}

	if cfg.Debug.DryRun {
		return nil
	}

	debugLaunch := strings.EqualFold(cfg.Debug.LogLevel, "debug")
	if err := launch.InDefaultTerminal(ctx, launch.Options{AppID: cfg.Launch.AppID, Title: cfg.Launch.Title, Debug: debugLaunch}, []string{exe, "__popup"}); err != nil {
		return fmt.Errorf("launch popup: %w", err)
	}
	return nil
}
