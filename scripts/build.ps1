# Build script for JSSON V0.0.6
# Builds binaries for Windows, Linux, and macOS

Set-Location (Split-Path -Parent $PSScriptRoot)

$VERSION = "v0.0.6"
$BIN_DIR = "bin"
$DIST_DIR = "dist/$VERSION"

# Create output directories
New-Item -ItemType Directory -Force -Path $BIN_DIR | Out-Null
New-Item -ItemType Directory -Force -Path $DIST_DIR | Out-Null

Write-Host "Building JSSON $VERSION..." -ForegroundColor Green

# Local Development Builds (bin/)

Write-Host "`n📦 Building local development binaries to '$BIN_DIR'..." -ForegroundColor Cyan

# CLI (native)
Write-Host "  Building jsson CLI..." -ForegroundColor White
go build -o "$BIN_DIR/jsson" ./cmd/jsson

# LSP Server
Write-Host "  Building jsson-lsp..." -ForegroundColor White
go build -o "$BIN_DIR/jsson-lsp" ./cmd/lsp

# WASM (for playground)
Write-Host "  Building jsson.wasm..." -ForegroundColor White
$env:GOOS = "js"
$env:GOARCH = "wasm"
go build -o "$BIN_DIR/jsson.wasm" ./cmd/wasm

# Reset GOOS/GOARCH
$env:GOOS = ""
$env:GOARCH = ""

Write-Host "✅ Local binaries ready in '$BIN_DIR'" -ForegroundColor Green

# Distribution Builds (dist/)

Write-Host "`n📦 Building distribution binaries to '$DIST_DIR'..." -ForegroundColor Cyan

# Windows AMD64
Write-Host "  Building for Windows (amd64)..." -ForegroundColor White
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o "$DIST_DIR/jsson-$VERSION-windows-amd64.exe" ./cmd/jsson

# Linux AMD64
Write-Host "  Building for Linux (amd64)..." -ForegroundColor White
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o "$DIST_DIR/jsson-$VERSION-linux-amd64" ./cmd/jsson

# macOS AMD64 (Intel)
Write-Host "  Building for macOS (amd64)..." -ForegroundColor White
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o "$DIST_DIR/jsson-$VERSION-darwin-amd64" ./cmd/jsson

# macOS ARM64 (Apple Silicon)
Write-Host "  Building for macOS (arm64)..." -ForegroundColor White
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o "$DIST_DIR/jsson-$VERSION-darwin-arm64" ./cmd/jsson

# Reset GOOS/GOARCH
$env:GOOS = ""
$env:GOARCH = ""

Write-Host "`n Build complete!" -ForegroundColor Green

Write-Host "`nLocal binaries ($BIN_DIR):" -ForegroundColor Yellow
Get-ChildItem $BIN_DIR | ForEach-Object {
    $size = [math]::Round($_.Length / 1MB, 2)
    Write-Host "  - $($_.Name) ($size MB)"
}

Write-Host "`nDistribution binaries ($DIST_DIR):" -ForegroundColor Yellow
Get-ChildItem $DIST_DIR | ForEach-Object {
    $size = [math]::Round($_.Length / 1MB, 2)
    Write-Host "  - $($_.Name) ($size MB)"
}
