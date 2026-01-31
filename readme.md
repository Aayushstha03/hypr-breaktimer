## hypr-breaktimer

Tiny Go + Bubble Tea break reminder for Hyprland.

### Commands

- `hypr-breaktimer show`: force open popup in a new terminal
- `hypr-breaktimer tick`: headless scheduler entrypoint (spawns popup when due)
- `hypr-breaktimer status`: print current config/state and next due time

### Config

- Path: `~/.config/hypr-breaktimer/config.toml`
- Defaults: work interval 30m, break 5m, snooze 10m

### State

- Path: `~/.local/state/hypr-breaktimer/state.json`
- Lock: `~/.local/state/hypr-breaktimer/lock`

### Terminal launching

- `show`/`tick` open the UI in a new terminal.
- Launcher: `xdg-terminal-exec` (recommended; respects your configured default terminal and supports app-id/title)

### Install (from source)

Dependencies (Arch examples):

- Go (build): `pacman -S go`
- Terminal launcher (runtime): `pacman -S xdg-terminal-exec`

Install (per-user, to `~/.local/bin`):

```bash
git clone <repo-url>
cd hypr-breaktimer
./install.sh
```

What `install.sh` does:

- Builds `hypr-breaktimer` from source.
- Installs it to `~/.local/bin/hypr-breaktimer`.
- Creates `~/.config/hypr-breaktimer/config.toml` if it does not already exist (from `config/config.toml.example`).
- Installs and enables a systemd user timer by default.

Skip systemd integration:

```bash
./install.sh --no-timer
```

### systemd (user)

These are the exact unit files installed by `install.sh`:

`~/.config/systemd/user/hypr-breaktimer.service`

```ini
[Unit]
Description=hypr-breaktimer (check if break is due)

[Service]
Type=oneshot
ExecStart=%h/.local/bin/hypr-breaktimer tick
```

`~/.config/systemd/user/hypr-breaktimer.timer`

```ini
[Unit]
Description=Run hypr-breaktimer periodically

[Timer]
OnBootSec=1m
OnUnitActiveSec=60s
Persistent=true

[Install]
WantedBy=timers.target
```

Manage the timer:

```bash
systemctl --user daemon-reload
systemctl --user enable --now hypr-breaktimer.timer
systemctl --user status hypr-breaktimer.timer
```

### Uninstall

```bash
./uninstall.sh
```

Remove config/state too:

```bash
./uninstall.sh --purge
```

### Hyprland

- When launched via `xdg-terminal-exec`, `tick` passes `--app-id=hypr-breaktimer` so you can match it in your window rules.

Example rules:

```ini
windowrule = opacity 0.8 0.8, match:app-id ^(hypr-breaktimer)$
windowrule = maximize on, match:app-id ^(hypr-breaktimer)$
```
