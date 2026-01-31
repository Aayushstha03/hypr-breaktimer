package logic

import (
	"fmt"
	"os"
	"time"

	"github.com/Aayushstha03/hypr-breaktimer/internal/config"
	"github.com/Aayushstha03/hypr-breaktimer/internal/state"
	"github.com/Aayushstha03/hypr-breaktimer/internal/xdg"
)

func Status() error {
	configPath, err := xdg.ConfigFile()
	if err != nil {
		return err
	}
	statePath, err := xdg.StateFile()
	if err != nil {
		return err
	}

	_, cfgErr := os.Stat(configPath)
	_, stErr := os.Stat(statePath)

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	st, err := state.Load(statePath)
	if err != nil {
		return err
	}

	now := time.Now()
	quietWindows, err := config.ParseQuietWindows(cfg.QuietHours.Windows)
	if err != nil {
		return err
	}
	inQuiet := cfg.QuietHours.Enabled && config.InQuietHours(now, quietWindows)

	ref := (*time.Time)(nil)
	if st.LastBreakCompletedAt != nil {
		ref = st.LastBreakCompletedAt
	} else if st.LastBreakStartedAt != nil {
		ref = st.LastBreakStartedAt
	} else if st.LastPopupShownAt != nil {
		ref = st.LastPopupShownAt
	}

	var nextDue *time.Time
	if ref != nil {
		t := ref.Add(cfg.Schedule.WorkInterval.Duration())
		nextDue = &t
	}

	blocked := ""
	if inQuiet {
		blocked = "quiet_hours"
	} else if st.SnoozedUntil != nil && now.Before(*st.SnoozedUntil) {
		blocked = "snoozed"
	} else if st.LastPopupShownAt != nil {
		minGap := cfg.Schedule.MinTimeBetweenPopups.Duration()
		if minGap > 0 && now.Sub(*st.LastPopupShownAt) < minGap {
			blocked = "min_time_between_popups"
		}
	}
	if blocked == "" && st.LastAction == state.ActionDismissed && st.LastActionAt != nil {
		cd := cfg.Schedule.DismissCooldown.Duration()
		if cd > 0 && now.Sub(*st.LastActionAt) < cd {
			blocked = "dismiss_cooldown"
		}
	}

	fmt.Printf("now: %s\n", now.Format(time.RFC3339))
	if cfgErr == nil {
		fmt.Printf("config: %s (exists)\n", configPath)
	} else {
		fmt.Printf("config: %s (missing; using defaults)\n", configPath)
	}
	if stErr == nil {
		fmt.Printf("state:  %s (exists)\n", statePath)
	} else {
		fmt.Printf("state:  %s (missing)\n", statePath)
	}

	fmt.Printf("work_interval:  %s\n", cfg.Schedule.WorkInterval.Duration())
	fmt.Printf("break_duration: %s\n", cfg.Schedule.BreakDuration.Duration())
	fmt.Printf("snooze:         %s\n", cfg.Schedule.SnoozeDuration.Duration())

	if st.LastPopupShownAt != nil {
		fmt.Printf("last_popup_shown_at: %s\n", st.LastPopupShownAt.Format(time.RFC3339))
	}
	if st.LastBreakStartedAt != nil {
		fmt.Printf("last_break_started_at: %s\n", st.LastBreakStartedAt.Format(time.RFC3339))
	}
	if st.LastBreakCompletedAt != nil {
		fmt.Printf("last_break_completed_at: %s\n", st.LastBreakCompletedAt.Format(time.RFC3339))
	}
	if st.SnoozedUntil != nil {
		fmt.Printf("snoozed_until: %s\n", st.SnoozedUntil.Format(time.RFC3339))
	}
	if st.LastAction != "" {
		at := ""
		if st.LastActionAt != nil {
			at = st.LastActionAt.Format(time.RFC3339)
		}
		fmt.Printf("last_action: %s %s\n", st.LastAction, at)
	}

	if nextDue != nil {
		fmt.Printf("next_due: %s (in %s)\n", nextDue.Format(time.RFC3339), humanDuration(nextDue.Sub(now)))
	} else {
		fmt.Printf("next_due: unknown (no reference time yet)\n")
	}

	if blocked != "" {
		fmt.Printf("blocked_by: %s\n", blocked)
	}
	return nil
}

func humanDuration(d time.Duration) string {
	if d < 0 {
		d = -d
		return "-" + d.Truncate(time.Second).String()
	}
	return d.Truncate(time.Second).String()
}
