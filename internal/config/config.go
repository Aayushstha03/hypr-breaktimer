package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Schedule   Schedule   `toml:"schedule"`
	QuietHours QuietHours `toml:"quiet_hours"`
	Popup      Popup      `toml:"popup"`
	Launch     Launch     `toml:"launch"`
	Debug      Debug      `toml:"debug"`
}

type Schedule struct {
	WorkInterval         Duration `toml:"work_interval"`
	BreakDuration        Duration `toml:"break_duration"`
	SnoozeDuration       Duration `toml:"snooze_duration"`
	MinTimeBetweenPopups Duration `toml:"min_time_between_popups"`
	DismissCooldown      Duration `toml:"dismiss_cooldown"`
}

type QuietHours struct {
	Enabled bool              `toml:"enabled"`
	Windows []QuietHourWindow `toml:"windows"`
}

type QuietHourWindow struct {
	Start string `toml:"start"`
	End   string `toml:"end"`
}

type Popup struct {
	Title          string `toml:"title"`
	Message        string `toml:"message"`
	AutoStartBreak bool   `toml:"auto_start_break"`
}

type Launch struct {
	AppID string `toml:"app_id"`
	Title string `toml:"title"`
}

type Debug struct {
	DryRun   bool   `toml:"dry_run"`
	LogLevel string `toml:"log_level"`
}

func Defaults() Config {
	return Config{
		Schedule: Schedule{
			WorkInterval:         MustDuration(30 * time.Minute),
			BreakDuration:        MustDuration(5 * time.Minute),
			SnoozeDuration:       MustDuration(10 * time.Minute),
			MinTimeBetweenPopups: MustDuration(10 * time.Minute),
			DismissCooldown:      MustDuration(15 * time.Minute),
		},
		QuietHours: QuietHours{Enabled: false},
		Popup: Popup{
			Title:          "Take a break",
			Message:        "Stand up, look away from the screen, and stretch.",
			AutoStartBreak: false,
		},
		Launch: Launch{
			AppID: "hypr-break",
			Title: "hypr-break",
		},
		Debug: Debug{DryRun: false, LogLevel: "info"},
	}
}

func Load(path string) (Config, error) {
	cfg := Defaults()
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}
	if _, err := toml.Decode(string(b), &cfg); err != nil {
		return Config{}, err
	}
	applyDefaults(&cfg)
	return cfg, nil
}

func WriteDefaultIfMissing(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return false, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}
	def := Defaults()
	buf := defaultTOML(def)
	if err := os.WriteFile(path, []byte(buf), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

func defaultTOML(cfg Config) string {
	// Keep this stable and human-editable.
	return "" +
		"[schedule]\n" +
		"work_interval = \"" + cfg.Schedule.WorkInterval.Duration().String() + "\"\n" +
		"break_duration = \"" + cfg.Schedule.BreakDuration.Duration().String() + "\"\n" +
		"snooze_duration = \"" + cfg.Schedule.SnoozeDuration.Duration().String() + "\"\n" +
		"min_time_between_popups = \"" + cfg.Schedule.MinTimeBetweenPopups.Duration().String() + "\"\n" +
		"dismiss_cooldown = \"" + cfg.Schedule.DismissCooldown.Duration().String() + "\"\n" +
		"\n" +
		"[quiet_hours]\n" +
		"enabled = false\n" +
		"windows = []\n" +
		"\n" +
		"[popup]\n" +
		"title = \"" + escapeTOMLString(cfg.Popup.Title) + "\"\n" +
		"message = \"" + escapeTOMLString(cfg.Popup.Message) + "\"\n" +
		"auto_start_break = false\n" +
		"\n" +
		"[launch]\n" +
		"app_id = \"" + escapeTOMLString(cfg.Launch.AppID) + "\"\n" +
		"title = \"" + escapeTOMLString(cfg.Launch.Title) + "\"\n" +
		"\n" +
		"[debug]\n" +
		"dry_run = false\n" +
		"log_level = \"info\"\n"
}

func escapeTOMLString(s string) string {
	// Minimal escaping for default config generation.
	// TOML basic strings use backslash for escapes.
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch b {
		case '\\':
			out = append(out, '\\', '\\')
		case '"':
			out = append(out, '\\', '"')
		case '\n':
			out = append(out, '\\', 'n')
		default:
			out = append(out, b)
		}
	}
	return string(out)
}

func applyDefaults(cfg *Config) {
	def := Defaults()
	if cfg.Schedule.WorkInterval.Duration() <= 0 {
		cfg.Schedule.WorkInterval = def.Schedule.WorkInterval
	}
	if cfg.Schedule.BreakDuration.Duration() <= 0 {
		cfg.Schedule.BreakDuration = def.Schedule.BreakDuration
	}
	if cfg.Schedule.SnoozeDuration.Duration() <= 0 {
		cfg.Schedule.SnoozeDuration = def.Schedule.SnoozeDuration
	}
	if cfg.Schedule.MinTimeBetweenPopups.Duration() < 0 {
		cfg.Schedule.MinTimeBetweenPopups = def.Schedule.MinTimeBetweenPopups
	}
	if cfg.Schedule.DismissCooldown.Duration() < 0 {
		cfg.Schedule.DismissCooldown = def.Schedule.DismissCooldown
	}
	if cfg.Popup.Title == "" {
		cfg.Popup.Title = def.Popup.Title
	}
	if cfg.Popup.Message == "" {
		cfg.Popup.Message = def.Popup.Message
	}
	if cfg.Launch.AppID == "" {
		cfg.Launch.AppID = def.Launch.AppID
	}
	if cfg.Launch.Title == "" {
		cfg.Launch.Title = def.Launch.Title
	}
	if cfg.Debug.LogLevel == "" {
		cfg.Debug.LogLevel = def.Debug.LogLevel
	}
}
