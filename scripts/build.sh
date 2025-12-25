#!/bin/bash
# Build script for JSSON V0.0.6
# Builds binaries for local development and distribution

cd "$(dirname "$0")/.." || exit 1

VERSION="v0.0.6"
BIN_DIR="bin"
DIST_DIR="dist/$VERSION"

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Create output directories
mkdir -p "$BIN_DIR"
mkdir -p "$DIST_DIR"

echo -e "${GREEN}Building JSSON $VERSION...${NC}"

# Local Development Builds (bin/)

echo -e "\n${CYAN}📦 Building local development binaries to '$BIN_DIR'...${NC}"

# CLI (native)
echo -e "${WHITE}  Building jsson CLI...${NC}"
go build -o "$BIN_DIR/jsson" ./cmd/jsson

# LSP Server
echo -e "${WHITE}  Building jsson-lsp...${NC}"
go build -o "$BIN_DIR/jsson-lsp" ./cmd/lsp

# WASM (for playground)
echo -e "${WHITE}  Building jsson.wasm...${NC}"
GOOS=js GOARCH=wasm go build -o "$BIN_DIR/jsson.wasm" ./cmd/wasm

echo -e "${GREEN}✅ Local binaries ready in '$BIN_DIR'${NC}"

# Distribution Builds (dist/)

echo -e "\n${CYAN}📦 Building distribution binaries to '$DIST_DIR'...${NC}"

# Windows AMD64
echo -e "${WHITE}  Building for Windows (amd64)...${NC}"
GOOS=windows GOARCH=amd64 go build -o "$DIST_DIR/jsson-$VERSION-windows-amd64.exe" ./cmd/jsson

# Linux AMD64
echo -e "${WHITE}  Building for Linux (amd64)...${NC}"
GOOS=linux GOARCH=amd64 go build -o "$DIST_DIR/jsson-$VERSION-linux-amd64" ./cmd/jsson

# macOS AMD64 (Intel)
echo -e "${WHITE}  Building for macOS (amd64)...${NC}"
GOOS=darwin GOARCH=amd64 go build -o "$DIST_DIR/jsson-$VERSION-darwin-amd64" ./cmd/jsson

# macOS ARM64 (Apple Silicon)
echo -e "${WHITE}  Building for macOS (arm64)...${NC}"
GOOS=darwin GOARCH=arm64 go build -o "$DIST_DIR/jsson-$VERSION-darwin-arm64" ./cmd/jsson

echo -e "\n${GREEN}Build complete!${NC}"

echo -e "\n${YELLOW}Local binaries ($BIN_DIR):${NC}"
ls -lh "$BIN_DIR" | tail -n +2 | awk '{printf "  - %s (%s)\n", $9, $5}'

echo -e "\n${YELLOW}Distribution binaries ($DIST_DIR):${NC}"
ls -lh "$DIST_DIR" | tail -n +2 | awk '{printf "  - %s (%s)\n", $9, $5}'
