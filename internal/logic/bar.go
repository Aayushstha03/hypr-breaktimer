package logic

import (
	"fmt"
	"time"

	"github.com/Aayushstha03/hypr-breaktimer/internal/config"
	"github.com/Aayushstha03/hypr-breaktimer/internal/state"
	"github.com/Aayushstha03/hypr-breaktimer/internal/xdg"
)

// Bar prints a minimal one-line status intended for waybar.
//
// Output format:
//
//	<mode> <minutes>m
//
// Modes:
//
//	break|snoozed|blocked|working|dnd|quiet|unknown
//
// Notes:
// - "dnd" and "quiet" omit minutes.
// - "due" is not emitted; when a break is due, mode is "working" with 0m.
func Bar() error {
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

	if isBreaking(now, st, cfg.Schedule.BreakDuration.Duration()) {
		mins := minutesRemaining(now, st.LastBreakStartedAt.Add(cfg.Schedule.BreakDuration.Duration()))
		fmt.Printf("Break %dm\n", mins)
		return nil
	}
	if st.SnoozedUntil != nil && now.Before(*st.SnoozedUntil) {
		mins := minutesRemaining(now, *st.SnoozedUntil)
		fmt.Printf("Snoozed %dm\n", mins)
		return nil
	}
	if st.DoNotDisturb {
		fmt.Println("dnd")
		return nil
	}
	if st.BlockedUntil != nil && now.Before(*st.BlockedUntil) {
		mins := minutesRemaining(now, *st.BlockedUntil)
		fmt.Printf("Blocked %dm\n", mins)
		return nil
	}
	if inQuiet {
		fmt.Println("Quiet")
		return nil
	}

	ref := (*time.Time)(nil)
	if st.LastBreakCompletedAt != nil {
		ref = st.LastBreakCompletedAt
	} else if st.LastBreakStartedAt != nil {
		ref = st.LastBreakStartedAt
	} else if st.LastPopupShownAt != nil {
		ref = st.LastPopupShownAt
	}
	if ref == nil {
		// fmt.Println(" unknown")
		return nil
	}

	nextDue := ref.Add(cfg.Schedule.WorkInterval.Duration())
	mins := minutesRemaining(now, nextDue)
	fmt.Printf("Working %dm\n", mins)
	return nil
}

func isBreaking(now time.Time, st state.State, breakDuration time.Duration) bool {
	if st.LastBreakStartedAt == nil {
		return false
	}
	if breakDuration <= 0 {
		return false
	}
	if st.LastBreakCompletedAt != nil && !st.LastBreakStartedAt.After(*st.LastBreakCompletedAt) {
		return false
	}
	endsAt := st.LastBreakStartedAt.Add(breakDuration)
	return now.Before(endsAt)
}

func minutesRemaining(now, until time.Time) int {
	if !until.After(now) {
		return 0
	}
	d := until.Sub(now)
	// Ceil to minutes so waybar updates once per minute without going negative.
	mins := int((d + time.Minute - 1) / time.Minute)
	if mins < 0 {
		return 0
	}
	return mins
}
