#!/usr/bin/env bash
set -euo pipefail

log()  { printf '%s\n' "[hypr-breaktimer] $*"; }
warn() { printf '%s\n' "[hypr-breaktimer] WARN: $*" >&2; }

BIN_NAME="hypr-breaktimer"
BIN_DIR="${XDG_BIN_HOME:-$HOME/.local/bin}"
BIN_PATH="$BIN_DIR/$BIN_NAME"

CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/hypr-breaktimer"
STATE_DIR="${XDG_STATE_HOME:-$HOME/.local/state}/hypr-breaktimer"

SYSTEMD_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/systemd/user"
SERVICE_NAME="hypr-breaktimer.service"
TIMER_NAME="hypr-breaktimer.timer"

PURGE=0
for arg in "$@"; do
  case "$arg" in
    --purge) PURGE=1 ;;
    -h|--help)
      cat <<EOF
Usage: ./uninstall.sh [--purge]

Removes:
- $BIN_PATH
- $SYSTEMD_DIR/$SERVICE_NAME
- $SYSTEMD_DIR/$TIMER_NAME

With --purge also removes:
- $CONFIG_DIR
- $STATE_DIR
EOF
      exit 0
      ;;
    *) printf '%s\n' "[hypr-breaktimer] ERROR: Unknown argument: $arg" >&2; exit 1 ;;
  esac
done

log "Disabling systemd timer (if installed)"
if command -v systemctl >/dev/null 2>&1; then
  systemctl --user disable --now "$TIMER_NAME" >/dev/null 2>&1 || true
  systemctl --user daemon-reload >/dev/null 2>&1 || true
fi

log "Removing systemd unit files"
rm -f "$SYSTEMD_DIR/$SERVICE_NAME" "$SYSTEMD_DIR/$TIMER_NAME" || true

log "Removing binary"
rm -f "$BIN_PATH" || true

if [ "$PURGE" -eq 1 ]; then
  log "Purging config/state"
  rm -rf "$CONFIG_DIR" "$STATE_DIR" || true
else
  warn "Leaving config/state in place. Use --purge to remove:"
  warn "  $CONFIG_DIR"
  warn "  $STATE_DIR"
fi

log "Done"
