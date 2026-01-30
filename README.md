# GitHub Actions Watcher

A terminal UI for monitoring GitHub Actions workflow runs. Auto-detects the repo from your current directory, shows recent runs, and lets you drill into jobs and steps -- all without leaving the terminal.

## Install

```bash
npm install -g github-actions-watcher
```

Requires the [GitHub CLI](https://cli.github.com/) (`gh`) to be installed and authenticated.

## Usage

```bash
# Run from any directory with a GitHub remote
ghaw

# Custom polling interval (default: 10s)
ghaw --interval 5
ghaw -i 30

# Or use without installing
npx github-actions-watcher
```

## Features

- **Auto-detects repo** from git remote (SSH or HTTPS)
- **Live polling** with countdown timer showing next refresh
- **Drill into runs** to see individual jobs and steps with durations
- **Switch repos** on the fly with `s`
- **Open in browser** with `o` from the detail view
- **Responsive layout** -- columns adapt to terminal width

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
npm install
npx tsx src/cli.tsx
```

## License

MIT
