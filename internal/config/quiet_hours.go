package config

import (
	"fmt"
	"time"
)

type QuietWindow struct {
	StartMin int
	EndMin   int
}

func ParseQuietWindows(w []QuietHourWindow) ([]QuietWindow, error) {
	out := make([]QuietWindow, 0, len(w))
	for _, win := range w {
		s, err := parseClockMinutes(win.Start)
		if err != nil {
			return nil, fmt.Errorf("quiet_hours.start %q: %w", win.Start, err)
		}
		e, err := parseClockMinutes(win.End)
		if err != nil {
			return nil, fmt.Errorf("quiet_hours.end %q: %w", win.End, err)
		}
		out = append(out, QuietWindow{StartMin: s, EndMin: e})
	}
	return out, nil
}

func InQuietHours(now time.Time, windows []QuietWindow) bool {
	min := now.Hour()*60 + now.Minute()
	for _, w := range windows {
		if w.StartMin == w.EndMin {
			continue
		}
		if w.StartMin < w.EndMin {
			if min >= w.StartMin && min < w.EndMin {
				return true
			}
			continue
		}
		// wraps past midnight
		if min >= w.StartMin || min < w.EndMin {
			return true
		}
	}
	return false
}

func parseClockMinutes(s string) (int, error) {
	// Accept HH:MM.
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, err
	}
	return t.Hour()*60 + t.Minute(), nil
}
