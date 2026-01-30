hypr-break

Tiny Go + Bubble Tea break reminder.

Commands
- `hypr-break popup`: open the break popup (manual test)
- `hypr-break show`: force open popup in a new terminal
- `hypr-break tick`: headless scheduler entrypoint (spawns popup when due)
- `hypr-break status`: print current config/state and next due time
- `hypr-break init`: create XDG dirs and write default config if missing

Config
- Path: `~/.config/hypr-break/config.toml`
- Defaults: work interval 30m, break 5m, snooze 10m

State
- Path: `~/.local/state/hypr-break/state.json`
- Lock: `~/.local/state/hypr-break/lock`

Terminal launching
- `tick` needs a terminal to show the Bubble Tea UI.
- Launcher order:
  1) `xdg-terminal-exec` (recommended; respects your configured default terminal and supports app-id/title)
  2) `$TERMINAL` (fallback; executed as a command)

Recommended setup:
- Install `xdg-terminal-exec` and let `hypr-break` use it.

Hyprland
- When launched via `xdg-terminal-exec`, `tick` passes `--app-id=hypr-break` so you can match that in your window rules.
