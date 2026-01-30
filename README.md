# GitHub Actions Watcher

A terminal UI for monitoring GitHub Actions workflow runs. Auto-detects the repo from your current directory, shows recent runs, and lets you drill into jobs and steps -- all without leaving the terminal.

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) for flicker-free rendering.

## Install

### Homebrew

```
brew install dzoba/tap/ghaw
```

### Go

```
go install github.com/dzoba/github-actions-watcher/cmd/ghaw@latest
```

### Binary

Download from [GitHub Releases](https://github.com/dzoba/github-actions-watcher/releases).

### Prerequisites

Requires the [GitHub CLI](https://cli.github.com/) (`gh`) to be installed and authenticated.

## Usage

```bash
# Run from any directory with a GitHub remote
ghaw

# Custom polling interval (default: 10s)
ghaw --interval 5
ghaw -i 30
```

## Features

- **Auto-detects repo** from git remote (SSH or HTTPS)
- **Live countdown timer** showing seconds until next refresh (flicker-free)
- **Drill into runs** to see individual jobs and steps with durations
- **Switch repos** on the fly with `s`
- **Open in browser** with `o` from the detail view
- **Responsive layout** -- columns adapt to terminal width
- **Single binary** -- no Node.js runtime required

## Keybindings

### List view

| Key | Action |
|-----|--------|
| Up/Down | Navigate runs |
| Enter | View jobs and steps |
| s | Switch repository |
| r | Refresh now |
| q | Quit |

### Detail view

| Key | Action |
|-----|--------|
| Up/Down | Scroll |
| Esc | Back to list |
| o | Open run in browser |
| r | Refresh |
| q | Quit |

## Development

```bash
git clone git@github.com:dzoba/github-actions-watcher.git
cd github-actions-watcher
go run ./cmd/ghaw
```

## License

MIT
