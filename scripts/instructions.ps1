# JSSON - Quick Build Instructions
# All binaries are generated to the bin/ directory

Write-Host "JSSON - Quick Build Instructions" -ForegroundColor Cyan
Write-Host ""

Write-Host "All binaries are built to the 'bin/' directory" -ForegroundColor Yellow
Write-Host ""

Write-Host "Build CLI (native):" -ForegroundColor Yellow
Write-Host "  go build -o bin/jsson ./cmd/jsson"
Write-Host ""

Write-Host "Build LSP Server:" -ForegroundColor Yellow
Write-Host "  go build -o bin/jsson-lsp ./cmd/lsp"
Write-Host ""

Write-Host "Build WASM (for playground):" -ForegroundColor Yellow
Write-Host '  $env:GOOS="js"; $env:GOARCH="wasm"; go build -o bin/jsson.wasm ./cmd/wasm'
Write-Host ""

Write-Host "Build all (run build script):" -ForegroundColor Yellow
Write-Host "  ./scripts/build.ps1"
Write-Host ""

Write-Host "Package VSIX (VS Code extension):" -ForegroundColor Yellow
Write-Host "  npx --yes vsce package"
Write-Host ""

Write-Host "Done." -ForegroundColor Green