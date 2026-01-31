## hypr-breaktimer

Tiny Go + Bubble Tea break reminder.

Commands

- `hypr-break show`: open the break timer popup
- `hypr-break tick`: headless scheduler entrypoint (spawns popup when due)
- `hypr-break status`: print current config/state and next due time

Config

- Path: `~/.config/hypr-break/config.toml`
- Defaults: work interval 30m, break 5m, snooze 10m

State

- Path: `~/.local/state/hypr-break/state.json`
- Lock: `~/.local/state/hypr-break/lock`

Terminal launching

- `show`/`tick` open the UI in a new terminal.
- Launcher: `xdg-terminal-exec` (recommended; respects your configured default terminal and supports app-id/title)

Recommended setup:

- Setup `xdg-terminal-exec` and let `hypr-breaktimer` use it.

Hyprland

- When launched via `xdg-terminal-exec`, `tick` passes `--app-id=hypr-breaktimer` so you can match that in your window rules.

```
windowrule = opacity 0.8 0.8, match:title ^(hypr-breaktimer-popup)$
windowrule = maximize on, match:title ^(hypr-breaktimer-popup)$
```
