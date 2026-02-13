#!/bin/sh
# QR Code Generator Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/DalyChouikh/qr-code-generator/main/install.sh | sh

set -e

REPO="DalyChouikh/qr-code-generator"
BINARY="qrgen"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

info() {
    printf "${CYAN}▸${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}✓${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}⚠${NC} %s\n" "$1"
}

error() {
    printf "${RED}✗${NC} %s\n" "$1" >&2
    exit 1
}

# Detect OS
detect_os() {
    OS="$(uname -s)"
    case "$OS" in
        Linux*)   OS="linux" ;;
        Darwin*)  OS="darwin" ;;
        MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
        *)        error "Unsupported operating system: $OS" ;;
    esac
    echo "$OS"
}

# Detect architecture
detect_arch() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64)  ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *)             error "Unsupported architecture: $ARCH" ;;
    esac
    echo "$ARCH"
}

# Get latest release version
get_latest_version() {
    if command -v curl > /dev/null 2>&1; then
        curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget > /dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Download and install
install() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    info "Detecting system... ${OS}/${ARCH}"

    info "Fetching latest release..."
    VERSION=$(get_latest_version)

    if [ -z "$VERSION" ]; then
        error "Could not determine the latest version. Check your internet connection."
    fi

    # Strip the 'v' prefix for the archive name
    VERSION_NUM="${VERSION#v}"

    info "Installing ${BINARY} ${VERSION}..."

    # Construct download URL
    if [ "$OS" = "windows" ]; then
        ARCHIVE="${BINARY}_${VERSION_NUM}_${OS}_${ARCH}.zip"
    else
        ARCHIVE="${BINARY}_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
    fi

    URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    info "Downloading ${URL}..."

    # Download
    if command -v curl > /dev/null 2>&1; then
        curl -fsSL "$URL" -o "${TMP_DIR}/${ARCHIVE}"
    elif command -v wget > /dev/null 2>&1; then
        wget -q "$URL" -O "${TMP_DIR}/${ARCHIVE}"
    fi

    # Extract
    info "Extracting..."
    cd "$TMP_DIR"

    if [ "$OS" = "windows" ]; then
        unzip -q "$ARCHIVE"
    else
        tar -xzf "$ARCHIVE"
    fi

    # Install binary
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY" "$INSTALL_DIR/"
    else
        warn "Elevated permissions required to install to ${INSTALL_DIR}"
        sudo mv "$BINARY" "$INSTALL_DIR/"
    fi

    chmod +x "${INSTALL_DIR}/${BINARY}"

    success "${BINARY} ${VERSION} installed to ${INSTALL_DIR}/${BINARY}"
    echo ""
    info "Run '${BINARY}' to get started!"
}

# Run installer
install
