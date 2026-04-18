#!/usr/bin/env bash
# fireflies-cli installer — interactive wizard.
#   curl -fsSL https://raw.githubusercontent.com/fvdm-otinga/fireflies-cli/main/scripts/install.sh | bash
# Non-interactive (CI): pipe with --yes, or set FIREFLIES_INSTALL_YES=1.
set -euo pipefail

REPO="fvdm-otinga/fireflies-cli"
BIN="fireflies"

YES=${FIREFLIES_INSTALL_YES:-0}
VERSION=${FIREFLIES_INSTALL_VERSION:-}
PREFIX=${FIREFLIES_INSTALL_PREFIX:-}
for arg in "$@"; do
  case "$arg" in
    --yes|-y) YES=1 ;;
    --version=*) VERSION="${arg#*=}" ;;
    --prefix=*) PREFIX="${arg#*=}" ;;
    -h|--help)
      cat <<EOF
fireflies-cli installer
Usage: install.sh [--yes] [--version=v1.0.2] [--prefix=/usr/local/bin]
  --yes, -y        non-interactive; accept defaults
  --version=X      install a specific tag (default: latest)
  --prefix=PATH    install dir (default: ask, else ~/.local/bin)
EOF
      exit 0 ;;
  esac
done

bold() { printf "\033[1m%s\033[0m\n" "$*"; }
dim()  { printf "\033[2m%s\033[0m\n" "$*"; }
err()  { printf "\033[31merror:\033[0m %s\n" "$*" >&2; exit 1; }
ok()   { printf "\033[32m✓\033[0m %s\n" "$*"; }

ask() {
  # ask "prompt" "default"
  local prompt="$1" default="${2-}" reply
  if [ "$YES" = "1" ] || [ ! -t 0 ]; then
    echo "$default"
    return
  fi
  if [ -n "$default" ]; then
    printf "%s [%s]: " "$prompt" "$default" >&2
  else
    printf "%s: " "$prompt" >&2
  fi
  read -r reply </dev/tty || reply=""
  echo "${reply:-$default}"
}

confirm() {
  local prompt="$1" default="${2:-y}" reply
  if [ "$YES" = "1" ]; then return 0; fi
  reply=$(ask "$prompt (y/n)" "$default")
  case "$reply" in y|Y|yes|YES) return 0;; *) return 1;; esac
}

bold "fireflies-cli installer"
echo

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in x86_64|amd64) ARCH=amd64;; arm64|aarch64) ARCH=arm64;; *) err "unsupported arch: $ARCH";; esac
case "$OS" in linux|darwin) ;; *) err "unsupported OS: $OS (for Windows, download the .zip from releases)";; esac
ok "detected: ${OS}/${ARCH}"

for cmd in curl tar shasum; do
  command -v "$cmd" >/dev/null 2>&1 || err "missing required tool: $cmd"
done

if [ -z "$VERSION" ]; then
  dim "resolving latest release..."
  VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' | head -n1 | sed -E 's/.*"([^"]+)".*/\1/')
  [ -n "$VERSION" ] || err "could not resolve latest version"
fi
ok "version: $VERSION"

V=${VERSION#v}
ARCHIVE="${BIN}_${V}_${OS}_${ARCH}.tar.gz"
URL_BASE="https://github.com/${REPO}/releases/download/${VERSION}"

if [ -z "$PREFIX" ]; then
  default_prefix="$HOME/.local/bin"
  if [ -w /usr/local/bin ] 2>/dev/null; then default_prefix="/usr/local/bin"; fi
  PREFIX=$(ask "install directory" "$default_prefix")
fi
mkdir -p "$PREFIX" || err "cannot create $PREFIX"
[ -w "$PREFIX" ] || err "no write permission on $PREFIX (try a different --prefix, or run with sudo)"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

dim "downloading $ARCHIVE ..."
curl -fsSL "$URL_BASE/$ARCHIVE" -o "$tmp/$ARCHIVE" || err "download failed"
curl -fsSL "$URL_BASE/checksums.txt" -o "$tmp/checksums.txt" || err "checksum download failed"

( cd "$tmp" && grep " $ARCHIVE$" checksums.txt | shasum -a 256 -c - >/dev/null ) \
  || err "checksum mismatch for $ARCHIVE"
ok "checksum verified"

tar -xzf "$tmp/$ARCHIVE" -C "$tmp"
install -m 755 "$tmp/$BIN" "$PREFIX/$BIN"
ok "installed: $PREFIX/$BIN"

if ! echo ":$PATH:" | grep -q ":$PREFIX:"; then
  echo
  dim "note: $PREFIX is not on your PATH."
  dim "add this to ~/.zshrc or ~/.bashrc:"
  printf '    export PATH="%s:$PATH"\n' "$PREFIX"
fi

echo
if confirm "run auth setup now?" "y"; then
  if [ -n "${FIREFLIES_API_KEY:-}" ]; then
    ok "FIREFLIES_API_KEY already set — skipping login"
  elif [ -t 0 ]; then
    "$PREFIX/$BIN" auth login
  else
    dim "stdin is not a tty; run \`$BIN auth login\` yourself"
  fi
fi

if confirm "install shell completions?" "n"; then
  SHELL_NAME=$(basename "${SHELL:-/bin/zsh}")
  case "$SHELL_NAME" in
    zsh)
      target="${HOME}/.zsh/completions"
      mkdir -p "$target"
      "$PREFIX/$BIN" completion zsh >"$target/_${BIN}"
      ok "zsh completion → $target/_${BIN}"
      dim "ensure your .zshrc has: fpath+=(\"$target\"); autoload -U compinit && compinit"
      ;;
    bash)
      target="${HOME}/.local/share/bash-completion/completions"
      mkdir -p "$target"
      "$PREFIX/$BIN" completion bash >"$target/$BIN"
      ok "bash completion → $target/$BIN"
      ;;
    fish)
      target="${HOME}/.config/fish/completions"
      mkdir -p "$target"
      "$PREFIX/$BIN" completion fish >"$target/$BIN.fish"
      ok "fish completion → $target/$BIN.fish"
      ;;
    *) dim "unknown shell: $SHELL_NAME — skipping" ;;
  esac
fi

echo
bold "done — try: $BIN --help"
