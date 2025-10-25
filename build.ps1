# PowerShell build script for GoInGo
Write-Host "Running tests..." -ForegroundColor Green
go test ./pkg/...

if ($LASTEXITCODE -ne 0) {
    Write-Host "Tests failed. Exiting build process." -ForegroundColor Red
    exit 1
}

Write-Host "Building main executable..." -ForegroundColor Green
$env:CGO_ENABLED = "1"
go build -o goengine.exe ./cmd/main.go

Write-Host "Checking CGO availability..." -ForegroundColor Yellow

$cgoEnabled = (go env CGO_ENABLED)
if ($cgoEnabled -eq "1") {
    Write-Host "Building shared library..." -ForegroundColor Green
    go build -buildmode=c-shared -o libgoengine.so ./export/export.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Build completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "Shared library build failed, but main executable was built successfully." -ForegroundColor Yellow
    }
} else {
    Write-Host "CGO is disabled. Skipping shared library build." -ForegroundColor Yellow
    Write-Host "To enable CGO, install a C compiler (like MinGW-w64) and set CGO_ENABLED=1" -ForegroundColor Yellow
}

Write-Host "Build process completed!" -ForegroundColor Green
