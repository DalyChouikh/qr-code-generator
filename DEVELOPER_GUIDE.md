# Developer & Publishing Guide

A complete guide for developing, building, and publishing the QR Code Generator CLI.

---

## Table of Contents

- [Project Overview](#project-overview)
- [Prerequisites](#prerequisites)
- [Local Development](#local-development)
- [Project Structure](#project-structure)
- [Git Workflow](#git-workflow)
- [Publishing a Release](#publishing-a-release)
- [How the CI/CD Pipeline Works](#how-the-cicd-pipeline-works)
- [Makefile Commands](#makefile-commands)
- [GoReleaser Explained](#goreleaser-explained)
- [Homebrew Tap Setup](#homebrew-tap-setup)
- [Install Script](#install-script)
- [FAQ](#faq)

---

## Project Overview

This project is a terminal CLI application written in Go. It uses:

| Component | Library | Purpose |
|-----------|---------|---------|
| TUI Framework | [Bubbletea](https://github.com/charmbracelet/bubbletea) | Interactive terminal UI (The Elm Architecture) |
| TUI Components | [Bubbles](https://github.com/charmbracelet/bubbles) | Text inputs, spinners, etc. |
| TUI Styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) | Colors, borders, layout |
| QR Generation | [go-qrcode](https://github.com/skip2/go-qrcode) | Encode content into QR codes |
| Release Automation | [GoReleaser](https://goreleaser.com) | Cross-compile + publish binaries |
| CI/CD | [GitHub Actions](https://github.com/features/actions) | Automated build & release on tag push |

---

## Prerequisites

### Required

- **Go 1.21+** — [Install Go](https://go.dev/doc/install)
- **Git** — For version control

### Optional (not required for development)

- **GoReleaser** — Only needed if you want to test releases locally. Users get binaries from GitHub Releases.
- **golangci-lint** — Optional linter for deeper code analysis. `make lint` falls back to `go vet` if not installed.

---

## Local Development

### First-time setup

```bash
# Clone the repo
git clone https://github.com/DalyChouikh/qr-code-generator.git
cd qr-code-generator

# Download dependencies
go mod tidy

# Build
make build

# Run
./qrgen
```

### Day-to-day workflow

```bash
# Build and run in one step
make run

# Just build
make build

# Run linter (uses go vet; golangci-lint if installed)
make lint

# Run tests
make test

# Clean build artifacts
make clean
```

---

## Project Structure

```
qr-code-generator/
├── .github/
│   └── workflows/
│       └── release.yml           # GitHub Actions: auto-release on tag push
├── cmd/
│   └── qrgen/
│       └── main.go               # Entry point, --version flag, starts TUI
├── internal/
│   ├── config/
│   │   └── config.go             # QRConfig struct, color parsing, validation
│   ├── generator/
│   │   └── generator.go          # QR code generation (PNG & SVG)
│   └── ui/
│       ├── model.go              # TUI model: steps, input handling, rendering
│       └── styles.go             # Lipgloss styles, color palette, progress bar
├── .gitignore
├── .goreleaser.yaml              # GoReleaser cross-compilation config
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
├── install.sh                    # One-line installer for end users
├── Makefile                      # Build/run/lint/release shortcuts
├── README.md                     # User-facing documentation
├── DEVELOPER_GUIDE.md            # This file
└── LICENSE
```

### How the code is organized

| Package | Responsibility |
|---------|---------------|
| `cmd/qrgen` | Application entry point. Parses `--version` flag, creates and runs the Bubbletea program. |
| `internal/config` | Defines `QRConfig` (content, format, size, colors, output path). Validates input, parses hex colors, provides predefined color palette. |
| `internal/generator` | Takes a `QRConfig` and produces the actual QR code file. Handles both PNG (using `image/png`) and SVG (custom builder). |
| `internal/ui` | The interactive TUI. `styles.go` defines the visual theme. `model.go` implements the Bubbletea Model (Init/Update/View) with a 6-step wizard. |

### TUI Architecture (Bubbletea / The Elm Architecture)

The TUI follows the **Elm Architecture** pattern:

1. **Model** — The application state (`Model` struct in `model.go`)
2. **Update** — Handles messages (keyboard input, window resize) and returns updated state
3. **View** — Renders the current state as a string for the terminal

The wizard has 7 steps: `StepURL → StepFormat → StepColor → StepSize → StepOutput → StepConfirm → StepComplete`

---

## Git Workflow

### Feature Development

```bash
# Create a feature branch
git checkout -b feature/my-new-feature

# Make changes, then commit
git add .
git commit -m "feat: add my new feature"

# Push the branch
git push origin feature/my-new-feature

# Open a Pull Request on GitHub, merge to main
```

### Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/) for clean changelogs:

| Prefix | When to use | Example |
|--------|------------|---------|
| `feat:` | New feature | `feat: add background color selection` |
| `fix:` | Bug fix | `fix: text input not accepting keystrokes` |
| `docs:` | Documentation only | `docs: update README install instructions` |
| `chore:` | Maintenance (deps, CI) | `chore: update bubbletea to v1.4` |
| `refactor:` | Code restructure (no behavior change) | `refactor: extract color picker into own step` |
| `test:` | Adding tests | `test: add config validation tests` |
| `ci:` | CI/CD changes | `ci: add lint job to workflow` |

---

## Publishing a Release

### Step-by-step

```bash
# 1. Make sure you're on main with everything committed
git checkout main
git pull origin main

# 2. Create a version tag (semantic versioning)
git tag v1.0.1

# 3. Push the tag — this triggers the release pipeline
git push origin v1.0.1
```

That's it. GitHub Actions will automatically:
1. Cross-compile binaries for 6 platform/arch combinations
2. Create a GitHub Release with downloadable archives
3. Generate a changelog from commit messages
4. Publish to Homebrew tap (if configured)

### Version Numbering (Semantic Versioning)

Format: `vMAJOR.MINOR.PATCH`

| Bump | When | Example |
|------|------|---------|
| PATCH | Bug fixes, small tweaks | `v1.0.0` → `v1.0.1` |
| MINOR | New features (backward compatible) | `v1.0.1` → `v1.1.0` |
| MAJOR | Breaking changes | `v1.1.0` → `v2.0.0` |

### If a release fails

```bash
# Delete the tag locally and remotely
git tag -d v1.0.1
git push origin --delete v1.0.1

# Fix the issue, commit, then re-tag
git add .
git commit -m "fix: resolve release build issue"
git tag v1.0.1
git push origin main
git push origin v1.0.1
```

---

## How the CI/CD Pipeline Works

### File: `.github/workflows/release.yml`

```
Push tag (v*)  →  GitHub Actions triggers  →  GoReleaser runs  →  Release published
```

**Trigger:** Any tag matching `v*` (e.g., `v1.0.0`, `v2.3.1`)

**What happens:**
1. **Checkout** — Pulls the full repo (with history for changelog)
2. **Setup Go** — Installs the Go version from `go.mod`
3. **GoReleaser** — Cross-compiles, packages, and publishes

**Permissions:** Uses the built-in `GITHUB_TOKEN` (no secrets to configure).

### File: `.goreleaser.yaml`

This tells GoReleaser:
- **What to build** — `./cmd/qrgen` compiled for linux/darwin/windows × amd64/arm64
- **How to package** — `.tar.gz` for macOS/Linux, `.zip` for Windows
- **What to include** — README.md and LICENSE in each archive
- **Where to publish** — GitHub Releases + Homebrew tap
- **Build flags** — `-s -w` (strip debug info for smaller binaries) + version/commit/date injection

---

## Makefile Commands

| Command | What it does |
|---------|-------------|
| `make build` | Compile the binary with version info embedded via ldflags |
| `make run` | Build then immediately run the app |
| `make clean` | Remove the binary and `dist/` directory |
| `make test` | Run all Go tests |
| `make lint` | Run `go vet` (and `golangci-lint` if installed) |
| `make snapshot` | Build a local release snapshot (no publish) — **requires goreleaser** |
| `make release-dry` | Validate goreleaser config and do a dry-run — **requires goreleaser** |
| `make help` | Show all available commands |

### Build with version info

When you run `make build`, the Makefile injects version metadata:

```bash
go build -ldflags "-s -w -X main.version=v1.0.0 -X main.commit=abc1234 -X main.date=2026-02-13T12:00:00Z"
```

This is why `./qrgen --version` shows the correct version, commit, and build date.

---

## GoReleaser Explained

[GoReleaser](https://goreleaser.com) is a tool that automates Go binary releases. You **don't need it installed locally** — it runs in GitHub Actions.

### What it produces for each release

```
GitHub Release "v1.0.0"
├── qrgen_1.0.0_darwin_amd64.tar.gz    (macOS Intel)
├── qrgen_1.0.0_darwin_arm64.tar.gz    (macOS Apple Silicon)
├── qrgen_1.0.0_linux_amd64.tar.gz     (Linux x86_64)
├── qrgen_1.0.0_linux_arm64.tar.gz     (Linux ARM64)
├── qrgen_1.0.0_windows_amd64.zip      (Windows x86_64)
├── qrgen_1.0.0_windows_arm64.zip      (Windows ARM64)
└── checksums.txt                       (SHA256 checksums)
```

### Key config sections

| Section | Purpose |
|---------|---------|
| `builds` | What to compile and for which platforms |
| `archives` | How to package (tar.gz/zip), what extra files to include |
| `checksum` | Generate SHA256 checksums for verification |
| `changelog` | Auto-generate changelog from git commits |
| `release` | GitHub Release settings (draft, prerelease, naming) |
| `brews` | Auto-publish Homebrew formula to your tap repo |

---

## Homebrew Tap Setup

Homebrew lets macOS/Linux users install with `brew install DalyChouikh/tap/qrgen`.

### One-time setup

1. Create a **separate** public GitHub repo called `homebrew-tap`:
   ```
   https://github.com/DalyChouikh/homebrew-tap
   ```

2. The repo can be empty — GoReleaser will automatically push a formula file to it on each release.

3. For GoReleaser to push to the tap repo, you need a **Personal Access Token** with `repo` scope:
   - Go to GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Create a token with `repo` scope
   - In your `qr-code-generator` repo, go to Settings → Secrets and variables → Actions
   - Add a secret named `HOMEBREW_TAP_TOKEN` with the token value

4. Update `.github/workflows/release.yml` to use the token:
   ```yaml
   env:
     GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
     HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
   ```

> **Note:** If you skip Homebrew setup, the release will still work — it will just skip the Homebrew publish step. Users can still download binaries directly.

---

## Install Script

The `install.sh` script lets users install with one command:

```bash
curl -fsSL https://raw.githubusercontent.com/DalyChouikh/qr-code-generator/main/install.sh | sh
```

**What it does:**
1. Detects OS (Linux/macOS/Windows) and architecture (amd64/arm64)
2. Fetches the latest release version from GitHub API
3. Downloads the correct archive
4. Extracts the binary
5. Moves it to `/usr/local/bin` (with `sudo` if needed)

---

## FAQ

### Do I need GoReleaser installed locally?

**No.** GoReleaser runs in GitHub Actions. You only need it locally if you want to test the release process with `make snapshot` or `make release-dry`.

### Do I need golangci-lint?

**No.** `make lint` uses `go vet` by default (which is built into Go). If you install `golangci-lint`, it will also run that for deeper analysis. It's optional.

### How do users without Go install this?

Three ways (none require Go):
1. **Quick install:** `curl -fsSL https://raw.githubusercontent.com/DalyChouikh/qr-code-generator/main/install.sh | sh`
2. **Homebrew:** `brew install DalyChouikh/tap/qrgen` (after tap setup)
3. **Manual:** Download binary from [GitHub Releases](https://github.com/DalyChouikh/qr-code-generator/releases)

### What if I push a tag and the release fails?

Delete the tag, fix the issue, and re-tag:
```bash
git tag -d v1.0.1
git push origin --delete v1.0.1
# fix and commit
git tag v1.0.1
git push origin main && git push origin v1.0.1
```

### How do I add a new feature?

1. Create a branch: `git checkout -b feature/my-feature`
2. Make changes in the relevant `internal/` package
3. Build and test: `make run`
4. Commit: `git commit -m "feat: description"`
5. Push and open a PR
6. After merge, tag a new version if ready to release

### What's the difference between `go build` and `make build`?

`make build` wraps `go build` but adds **ldflags** that inject version, commit hash, and build date into the binary. This is why `./qrgen --version` shows useful info instead of "dev".

### Why is the code in `internal/`?

The `internal/` directory is a Go convention that prevents other Go modules from importing these packages. It keeps the public API surface clean — only the `cmd/qrgen` entry point is "public".
