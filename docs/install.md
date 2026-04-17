# Installing fireflies

## Homebrew (macOS / Linux)

```sh
brew install fvdm-otinga/tap/fireflies
```

This installs a prebuilt binary and automatically registers shell completions.

## Go install

Requires Go 1.21+ and places the binary in `$GOPATH/bin` (typically `~/go/bin`):

```sh
go install github.com/fvdm-otinga/fireflies-cli@latest
```

## Download a prebuilt binary

Binaries for macOS (arm64 / amd64), Linux (amd64 / arm64), and Windows (amd64) are
published on every [GitHub Release](https://github.com/fvdm-otinga/fireflies-cli/releases).

### macOS (arm64)

```sh
VERSION=v1.0.0
curl -sSL "https://github.com/fvdm-otinga/fireflies-cli/releases/download/${VERSION}/fireflies_${VERSION#v}_darwin_arm64.tar.gz" \
  | tar xz -C /usr/local/bin fireflies
```

### macOS (amd64 / Intel)

```sh
VERSION=v1.0.0
curl -sSL "https://github.com/fvdm-otinga/fireflies-cli/releases/download/${VERSION}/fireflies_${VERSION#v}_darwin_amd64.tar.gz" \
  | tar xz -C /usr/local/bin fireflies
```

### Linux (amd64)

```sh
VERSION=v1.0.0
curl -sSL "https://github.com/fvdm-otinga/fireflies-cli/releases/download/${VERSION}/fireflies_${VERSION#v}_linux_amd64.tar.gz" \
  | tar xz -C /usr/local/bin fireflies
```

### Linux (arm64)

```sh
VERSION=v1.0.0
curl -sSL "https://github.com/fvdm-otinga/fireflies-cli/releases/download/${VERSION}/fireflies_${VERSION#v}_linux_arm64.tar.gz" \
  | tar xz -C /usr/local/bin fireflies
```

### Windows (amd64)

Download `fireflies_<version>_windows_amd64.zip` from the [releases page](https://github.com/fvdm-otinga/fireflies-cli/releases),
unzip it, and move `fireflies.exe` to a directory on your `PATH`.

### Verify the checksum

Each release includes a `checksums.txt` (SHA-256). After downloading:

```sh
sha256sum --check checksums.txt
```

---

## Auth quickstart

### Option A — environment variable (recommended for CI/scripts)

```sh
export FIREFLIES_API_KEY=ff_xxxxxxxxxxxxxxxxxxxxxxxx
fireflies users whoami
```

Generate your key at [app.fireflies.ai/integrations/custom/fireflies](https://app.fireflies.ai/integrations/custom/fireflies).

### Option B — interactive login (persisted in `~/.config/fireflies/config.toml`)

```sh
fireflies auth login
# Paste your API key when prompted.
fireflies auth status   # verify
```

### Multi-profile usage

```sh
# Save a second profile
fireflies auth login --profile work

# Use it for a single command
fireflies meetings list --profile work
```

---

## Shell completions

### Bash

```sh
source <(fireflies completion bash)
# or persist:
fireflies completion bash > /usr/local/etc/bash_completion.d/fireflies
```

### Zsh

```sh
source <(fireflies completion zsh)
# or install to fpath:
fireflies completion zsh > "${fpath[1]}/_fireflies"
```

### Fish

```sh
fireflies completion fish | source
# or persist:
fireflies completion fish > ~/.config/fish/completions/fireflies.fish
```

### PowerShell

```powershell
fireflies completion powershell | Out-String | Invoke-Expression
```
