# QR Code Generator

A beautiful terminal-based QR code generator built with Go, featuring an interactive TUI (Terminal User Interface).

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)
[![Release](https://img.shields.io/github/v/release/DalyChouikh/qr-code-generator?style=flat)](https://github.com/DalyChouikh/qr-code-generator/releases/latest)

## Features

- 🎨 **Interactive TUI** - Step-by-step wizard for creating QR codes
- 🖼️ **Multiple Formats** - Generate PNG or SVG output
- 🎭 **Custom Colors** - Choose from predefined colors or enter custom hex values
- 📐 **Flexible Dimensions** - Set size from 64 to 4096 pixels
- 💾 **Flexible Output** - Save anywhere on your system
- 🚀 **Cross-platform** - Works on macOS, Linux, and Windows

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
./qrgen
```

### Steps

1. **Enter URL/Text** - The content to encode in the QR code
2. **Choose Format** - Select PNG (raster) or SVG (vector)
3. **Select Color** - Pick a predefined color or enter a custom hex value
4. **Set Dimensions** - Specify the output size (64-4096 pixels)
5. **Output Location** - Choose where to save the file
6. **Review & Generate** - Confirm and create your QR code

### Keyboard Navigation

| Key | Action |
|-----|--------|
| `Enter` | Confirm selection |
| `↑/↓` or `j/k` | Navigate lists |
| `←/→` or `h/l` | Switch options |
| `Space` | Select option |
| `Esc` | Go back |
| `Ctrl+C` | Quit |
| `c` | Enter custom color (in color step) |
| `r` | Create another (after completion) |

## Project Structure

```
qr-code-generator/
├── cmd/
│   └── qrgen/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/

│       └── styles.go        # UI styling
├── go.mod
├── go.sum
└── README.md
```

## Dependencies

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [go-qrcode](https://github.com/skip2/go-qrcode) - QR code generation

## Examples

### Generate a QR code with default settings
```bash
./qrgen
# Enter: https://github.com
# Select: PNG format
# Select: Black color
# Size: 256 (default)
# Output: qrcode (saves as qrcode.png)
```

### Generate a colored SVG QR code
```bash
./qrgen
# Enter: https://example.com
# Select: SVG format
# Select: Blue color (or press 'c' and enter #0066CC)
# Size: 512
# Output: ~/Downloads/my-qr (saves as ~/Downloads/my-qr.svg)
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
