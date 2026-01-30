package config

import (
	"fmt"
	"time"
)

type Duration struct{ d time.Duration }

func (d Duration) Duration() time.Duration { return d.d }

func MustDuration(v time.Duration) Duration { return Duration{d: v} }

func (d *Duration) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		d.d = 0
		return nil
	}
	v, err := time.ParseDuration(string(text))
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", string(text), err)
	}
	d.d = v
	return nil
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.d.String()), nil
}
