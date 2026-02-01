#!/usr/bin/env bash
set -euo pipefail

log()  { printf '%s\n' "[hypr-breaktimer] $*"; }
warn() { printf '%s\n' "[hypr-breaktimer] WARN: $*" >&2; }
die()  { printf '%s\n' "[hypr-breaktimer] ERROR: $*" >&2; exit 1; }

need_cmd() { command -v "$1" >/dev/null 2>&1 || die "Missing dependency: $1"; }

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$SCRIPT_DIR"

BIN_NAME="hypr-breaktimer"
BIN_DIR="${XDG_BIN_HOME:-$HOME/.local/bin}"
BIN_PATH="$BIN_DIR/$BIN_NAME"

CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/hypr-breaktimer"
CONFIG_PATH="$CONFIG_DIR/config.toml"

STATE_DIR="${XDG_STATE_HOME:-$HOME/.local/state}/hypr-breaktimer"
STATE_PATH="$STATE_DIR/state.json"

SYSTEMD_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/systemd/user"
SERVICE_NAME="hypr-breaktimer.service"
TIMER_NAME="hypr-breaktimer.timer"

ENABLE_TIMER=1

for arg in "$@"; do
  case "$arg" in
    --no-timer) ENABLE_TIMER=0 ;;
    -h|--help)
      cat <<EOF
Usage: ./install.sh [--no-timer]

Installs $BIN_NAME to:
  $BIN_PATH

Creates a default config if missing:
  $CONFIG_PATH

Installs and enables a systemd user timer by default.
Pass --no-timer to skip systemd integration.
EOF
      exit 0
      ;;
    *) die "Unknown argument: $arg" ;;
  esac
done

log "Starting per-user install"

need_cmd go
need_cmd install
need_cmd mkdir
need_cmd mktemp

if ! command -v xdg-terminal-exec >/dev/null 2>&1; then
  warn "xdg-terminal-exec not found. Popups will not launch until it's installed."
  warn "On Arch: pacman -S xdg-terminal-exec"
fi

log "Building binary"
tmp="$(mktemp -d)"
cleanup() { rm -rf "$tmp"; }
trap cleanup EXIT

OUT="$tmp/$BIN_NAME"
( cd "$REPO_ROOT" && go build -trimpath -ldflags "-s -w" -o "$OUT" "./cmd/$BIN_NAME" )

log "Installing binary to $BIN_PATH"
mkdir -p "$BIN_DIR"
install -m 0755 "$OUT" "$BIN_PATH"

if [ ! -f "$CONFIG_PATH" ]; then
  log "Creating default config at $CONFIG_PATH"
  mkdir -p "$CONFIG_DIR"
  if [ -f "$REPO_ROOT/config/config.toml.example" ]; then
    install -m 0644 "$REPO_ROOT/config/config.toml.example" "$CONFIG_PATH"
  else
    die "Missing sample config: $REPO_ROOT/config/config.toml.example"
  fi
else
  log "Config already exists, leaving it unchanged: $CONFIG_PATH"
fi

if [ ! -f "$STATE_PATH" ]; then
  log "Initializing state at $STATE_PATH"
  mkdir -p "$STATE_DIR"

  now="$(date +"%Y-%m-%dT%H:%M:%S%:z")"
  tmp_state="$STATE_PATH.tmp.$$"
  {
    printf '{\n'
    printf '  "last_popup_shown_at": "%s"\n' "$now"
    printf '}\n'
  } >"$tmp_state"
  chmod 0644 "$tmp_state" 2>/dev/null || true
  mv -f "$tmp_state" "$STATE_PATH"
else
  log "State already exists, leaving it unchanged: $STATE_PATH"
fi

if [ "$ENABLE_TIMER" -eq 1 ]; then
  if command -v systemctl >/dev/null 2>&1; then
    log "Installing systemd user units"
    mkdir -p "$SYSTEMD_DIR"

    [ -f "$REPO_ROOT/systemd/$SERVICE_NAME" ] || die "Missing $REPO_ROOT/systemd/$SERVICE_NAME"
    [ -f "$REPO_ROOT/systemd/$TIMER_NAME" ] || die "Missing $REPO_ROOT/systemd/$TIMER_NAME"

    install -m 0644 "$REPO_ROOT/systemd/$SERVICE_NAME" "$SYSTEMD_DIR/$SERVICE_NAME"
    install -m 0644 "$REPO_ROOT/systemd/$TIMER_NAME" "$SYSTEMD_DIR/$TIMER_NAME"

    log "Enabling timer: $TIMER_NAME"
    if systemctl --user daemon-reload && systemctl --user enable --now "$TIMER_NAME"; then
      log "Timer enabled"
    else
      warn "Could not enable systemd user timer automatically."
      warn "Try: systemctl --user daemon-reload && systemctl --user enable --now $TIMER_NAME"
    fi
  else
    warn "systemctl not found; skipping systemd integration"
  fi
else
  log "Skipping systemd timer install (--no-timer)"
fi

log "Done"
log "Try:"
log "  $BIN_NAME --help"
log "  $BIN_NAME status"
if [ "$ENABLE_TIMER" -eq 1 ]; then
  log "  systemctl --user status $TIMER_NAME"
fi
