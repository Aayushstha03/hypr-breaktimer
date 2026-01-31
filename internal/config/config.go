package config

import (
	"os"
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
			AppID: "hypr-breaktimer",
			Title: "hypr-breaktimer",
		},
		Debug: Debug{DryRun: false, LogLevel: "info"},
	}
}

func Load(path string) (Config, error) {
	cfg := Defaults()
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
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
