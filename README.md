# QR Code Generator

A beautiful terminal-based QR code generator built with Go, featuring an interactive TUI (Terminal User Interface).

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)
[![Release](https://img.shields.io/github/v/release/DalyChouikh/qr-code-generator?style=flat)](https://github.com/DalyChouikh/qr-code-generator/releases/latest)

## Features

- 🎨 **Interactive TUI** — Step-by-step wizard for creating QR codes
- 📋 **Smart Content Templates** — Guided forms for WiFi, Contact (vCard), Email, SMS, URL, and plain text
- 🖼️ **Multiple Formats** — Generate PNG or SVG output
- 🎭 **Custom Colors** — Foreground & background color pickers with predefined palette or custom hex
- 📐 **Flexible Dimensions** — Set size from 64 to 4096 pixels
- 📂 **File Picker** — Built-in file browser for choosing output location
- 📱 **Terminal Preview** — Scan the QR code directly in your terminal after generation
- 📜 **Generation History** — Automatically saves your last 50 generations for quick re-use
- 🔄 **Self-Updater** — Update to the latest version with a single command
- 🚀 **Cross-Platform** — Works on macOS, Linux, and Windows

## Installation

### Quick Install (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/DalyChouikh/qr-code-generator/main/install.sh | sh
```

### Quick Install (Windows)

Open PowerShell and run:

```powershell
irm https://raw.githubusercontent.com/DalyChouikh/qr-code-generator/main/install.ps1 | iex
```

### Homebrew (macOS / Linux)

```bash
brew install DalyChouikh/tap/qrgen
```

### Download Binary

Pre-built binaries for all platforms are available on the [Releases](https://github.com/DalyChouikh/qr-code-generator/releases/latest) page.

| Platform | Architecture | Download |
|----------|-------------|----------|
| macOS    | Apple Silicon (M1/M2/M3) | `qrgen_*_darwin_arm64.tar.gz` |
| macOS    | Intel | `qrgen_*_darwin_amd64.tar.gz` |
| Linux    | x86_64 | `qrgen_*_linux_amd64.tar.gz` |
| Linux    | ARM64 | `qrgen_*_linux_arm64.tar.gz` |
| Windows  | x86_64 | `qrgen_*_windows_amd64.zip` |
| Windows  | ARM64 | `qrgen_*_windows_arm64.zip` |

Download, extract, and move the binary to your PATH:

```bash
# Example for Linux x86_64
tar -xzf qrgen_*_linux_amd64.tar.gz
sudo mv qrgen /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/DalyChouikh/cmd/qrgen@latest
```

### From Source

```bash
git clone https://github.com/DalyChouikh/qr-code-generator.git
cd qr-code-generator
make build
./qrgen
```

## Usage

Simply run the application and follow the interactive prompts:

```bash
qrgen
```

### CLI Commands

| Command | Description |
|---------|-------------|
| `qrgen` | Launch the interactive QR code generator |
| `qrgen history` | Show your generation history |
| `qrgen regen <id>` | Re-generate a QR code from history |
| `qrgen update` | Update qrgen to the latest version |
| `qrgen check-update` | Check if a newer version is available |
| `qrgen --version` | Print version information |
| `qrgen --help` | Show help message |

> **Note:** On Linux/macOS, if qrgen is installed in `/usr/local/bin`, updating requires elevated privileges: `sudo qrgen update`. On Windows, run the update from an Administrator terminal.

### Wizard Steps

1. **Content Type** — Choose what to encode: URL, WiFi, Contact, Email, SMS, or plain text
2. **Content Details** — Enter the content (guided form for WiFi/Contact/Email/SMS, or free text for URL/Text)
3. **Output Format** — Select PNG (raster) or SVG (vector)
4. **Foreground Color** — Pick the QR code color from a palette or enter a custom hex value
5. **Background Color** — Pick the background color
6. **Dimensions** — Set the output size (64–4096 pixels)
7. **Output Location** — Type a path or browse with the built-in file picker
8. **Review & Generate** — Confirm settings and generate your QR code

### Content Templates

| Type | Description | Example Output |
|------|-------------|---------------|
| 🔗 URL | Website link | Opens browser on scan |
| 📶 WiFi | Network credentials (SSID, password, encryption) | Auto-connects to network |
| 👤 Contact | vCard with name, phone, email, org, title, URL | Saves contact to phone |
| ✉️ Email | Pre-filled email with address, subject, body | Opens email compose |
| 💬 SMS | Pre-filled text message with phone and message | Opens messaging app |
| 📝 Text | Plain text | Displays text |

### Keyboard Navigation

| Key | Action |
|-----|--------|
| `Enter` | Confirm selection |
| `↑/↓` or `j/k` | Navigate lists |
| `←/→` or `h/l` | Switch options |
| `Tab` / `Shift+Tab` | Next / previous field (in template forms) |
| `Space` | Select option / toggle |
| `Esc` | Go back |
| `Ctrl+C` | Quit |
| `c` | Enter custom color (in color steps) |
| `Tab` | Toggle file browser (in output step) |
| `r` | Create another (after completion) |

## Project Structure

```
qr-code-generator/
├── cmd/
│   └── qrgen/
│       └── main.go              # Application entry point & CLI commands
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration types & color utilities
│   ├── generator/
│   │   ├── generator.go         # PNG & SVG QR code generation
│   │   └── terminal.go          # Terminal QR preview renderer
│   ├── history/
│   │   └── history.go           # Generation history storage
│   ├── templates/
│   │   └── templates.go         # Content templates (WiFi, vCard, Email, SMS)
│   ├── ui/
│   │   ├── model.go             # Main TUI model & wizard logic
│   │   ├── styles.go            # UI styling
│   │   ├── filepicker.go        # Built-in file/directory browser
│   │   └── template_wizard.go   # Template form UI component
│   └── updater/
│       └── updater.go           # Self-update via GitHub Releases
├── .goreleaser.yaml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Website

`website/` is a Next.js static site that generates QR codes in the browser. It runs the CLI's own Go generator compiled to WebAssembly (`cmd/qrgen-wasm`), so web and terminal output are identical and nothing is uploaded.

```bash
make website-dev   # build qrgen.wasm and start the dev server
make website       # build qrgen.wasm and export the static site to website/out
```

Requires Go and Node.js 20+. Run `npm install` in `website/` once before `make website-dev`.

When building for deployment, set `NEXT_PUBLIC_SITE_URL` to the site's public URL (for example `NEXT_PUBLIC_SITE_URL=https://qrgen.example.com make website`). It is used for the canonical link, the sitemap and social-media previews.

## Dependencies

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [go-qrcode](https://github.com/skip2/go-qrcode) - QR code generation

## Examples

### Generate a URL QR code
```bash
qrgen
# Select: URL
# Enter: https://github.com
# Select: PNG format
# Select: Black foreground, White background
# Size: 256 (default)
# Output: qrcode (saves as qrcode.png)
```

### Generate a WiFi QR code
```bash
qrgen
# Select: WiFi
# Enter SSID: MyNetwork
# Enter Password: secret123
# Select encryption: WPA/WPA2/WPA3
# Hidden: No
# Select: PNG format, colors, size, and output
```

### Generate a Contact card QR code
```bash
qrgen
# Select: Contact
# Fill in: First Name, Last Name, Phone, Email, etc.
# Generates a vCard QR — scanning saves the contact to your phone
```

### View generation history
```bash
qrgen history
```

### Re-generate a previous QR code
```bash
qrgen regen 3   # Re-generate entry #3 from history
```

### Update to the latest version
```bash
qrgen check-update     # Check if update is available
sudo qrgen update      # Update (use sudo on Linux/macOS if needed)
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
