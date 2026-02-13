# QR Code Generator Installer for Windows
# Usage: irm https://raw.githubusercontent.com/DalyChouikh/qr-code-generator/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "DalyChouikh/qr-code-generator"
$Binary = "qrgen"
$InstallDir = "$env:LOCALAPPDATA\Programs\qrgen"

Write-Host ""
Write-Host "  QR Code Generator Installer" -ForegroundColor Cyan
Write-Host "  ============================" -ForegroundColor Cyan
Write-Host ""

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Host "x 32-bit systems are not supported." -ForegroundColor Red
    exit 1
}

Write-Host "  > Detected: windows/$Arch" -ForegroundColor Gray

# Get latest release version
Write-Host "  > Fetching latest release..." -ForegroundColor Gray
try {
    $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $Release.tag_name
} catch {
    Write-Host "x Failed to fetch latest release. Check your internet connection." -ForegroundColor Red
    exit 1
}

$VersionNum = $Version.TrimStart("v")
$Archive = "${Binary}_${VersionNum}_windows_${Arch}.zip"
$Url = "https://github.com/$Repo/releases/download/$Version/$Archive"

Write-Host "  > Installing $Binary $Version..." -ForegroundColor Gray

# Create temp directory
$TmpDir = Join-Path $env:TEMP "qrgen-install-$(Get-Random)"
New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null

try {
    # Download
    Write-Host "  > Downloading $Archive..." -ForegroundColor Gray
    $ZipPath = Join-Path $TmpDir $Archive
    Invoke-WebRequest -Uri $Url -OutFile $ZipPath -UseBasicParsing

    # Extract
    Write-Host "  > Extracting..." -ForegroundColor Gray
    Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

    # Create install directory
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    # Move binary
    $BinaryPath = Join-Path $TmpDir "$Binary.exe"
    Copy-Item -Path $BinaryPath -Destination (Join-Path $InstallDir "$Binary.exe") -Force

    # Add to PATH if not already there
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        Write-Host "  > Adding to PATH..." -ForegroundColor Gray
        [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
        $env:Path = "$env:Path;$InstallDir"
    }

    Write-Host ""
    Write-Host "  OK $Binary $Version installed successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "  Installed to: $InstallDir\$Binary.exe" -ForegroundColor Gray
    Write-Host ""
    Write-Host "  Restart your terminal, then run '$Binary' to get started!" -ForegroundColor Cyan
    Write-Host ""

} catch {
    Write-Host "x Installation failed: $_" -ForegroundColor Red
    exit 1
} finally {
    # Cleanup
    Remove-Item -Path $TmpDir -Recurse -Force -ErrorAction SilentlyContinue
}
