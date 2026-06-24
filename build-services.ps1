# Build all unsia services
$ErrorActionPreference = "Stop"

$servicePath = "d:\Superman\Superman\Coding\New folder\Dockument\candi\unsia-docs-md\services"

Write-Host "Building unsia-core-service..." -ForegroundColor Cyan
Set-Location "$servicePath\unsia-core-service"
go build -o bin\core-service.exe .\cmd\core-service\
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build unsia-core-service" -ForegroundColor Red
    exit 1
}
Write-Host "unsia-core-service built successfully" -ForegroundColor Green

Write-Host "Building unsia-reference-service..." -ForegroundColor Cyan
Set-Location "$servicePath\unsia-reference-service"
go build -o bin\reference-service.exe .\cmd\reference-service\
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build unsia-reference-service" -ForegroundColor Red
    exit 1
}
Write-Host "unsia-reference-service built successfully" -ForegroundColor Green

Write-Host "Building unsia-finance-service..." -ForegroundColor Cyan
Set-Location "$servicePath\unsia-finance-service"
go build -o bin\finance-service.exe .\cmd\finance-service\
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build unsia-finance-service" -ForegroundColor Red
    exit 1
}
Write-Host "unsia-finance-service built successfully" -ForegroundColor Green

Write-Host "Building unsia-pmb-service..." -ForegroundColor Cyan
Set-Location "$servicePath\unsia-pmb-service"
go build -o bin\pmb-service.exe .\cmd\pmb-service\
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build unsia-pmb-service" -ForegroundColor Red
    exit 1
}
Write-Host "unsia-pmb-service built successfully" -ForegroundColor Green

Write-Host "Building unsia-academic-service..." -ForegroundColor Cyan
Set-Location "$servicePath\unsia-academic-service"
go build -o bin\academic-service.exe .\cmd\academic-service\
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build unsia-academic-service" -ForegroundColor Red
    exit 1
}
Write-Host "unsia-academic-service built successfully" -ForegroundColor Green

Write-Host "`nAll services built successfully!" -ForegroundColor Green
