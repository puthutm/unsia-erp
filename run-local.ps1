# UNSIA ERP - Run Locally (Tanpa Docker!)
# Hanya perlu: Go, Node.js, dan PostgreSQL lokal

$ErrorActionPreference = "Stop"
$projectRoot = "d:\Superman\Superman\Coding\New folder\Dockument\candi\unsia-docs-md"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "UNSIA ERP - Jalankan Lokal (Tanpa Docker)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check Go
Write-Host "1. Go:" -ForegroundColor Yellow
$goVersion = go version
Write-Host "   $goVersion" -ForegroundColor Green

# Check Node
Write-Host "2. Node.js:" -ForegroundColor Yellow
$nodeVersion = node --version
Write-Host "   $nodeVersion" -ForegroundColor Green

# Check PostgreSQL
Write-Host "3. PostgreSQL:" -ForegroundColor Yellow
try {
    $pgConn = Test-NetConnection -ComputerName localhost -Port 5432 -WarningAction SilentlyContinue
    if ($pgConn.TcpTestSucceeded) {
        Write-Host "   Running di port 5432" -ForegroundColor Green
    }
} catch {
    Write-Host "   Pastikan PostgreSQL sudah running!" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "LANGKAH MENJALANKAN:" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Step 1 - Databases (SUDAH DIBUAT via Docker init script)
Write-Host "Langkah 1: Database (Sudah dibuat otomatis oleh Docker init script)" -ForegroundColor Green
Write-Host "---------------------------------------------------------------" -ForegroundColor Green
Write-Host "  core_db, reference_db, crm_db, pmb_db, finance_db," -ForegroundColor White
Write-Host "  academic_db, hris_db, lms_db, assessment_db, portal_db" -ForegroundColor White
Write-Host ""

# Step 2
Write-Host "Langkah 2: Setup Environment Variable" -ForegroundColor Cyan
Write-Host "--------------------------------" -ForegroundColor Yellow
Write-Host "Buat file .env di setiap service, contoh:" -ForegroundColor White
Write-Host "  DATABASE_URL=postgres://postgres:@localhost:5432/core_db?sslmode=disable" -ForegroundColor Gray
Write-Host "  PORT=8001" -ForegroundColor Gray
Write-Host "  JWKS_URL=http://localhost:8001/.well-known/jwks.json" -ForegroundColor Gray
Write-Host ""

# Step 3 - Backend Services (needs .env files first)
Write-Host "Langkah 3: Jalankan Backend (needs .env setup)" -ForegroundColor Cyan
Write-Host "---------------------------------------" -ForegroundColor Yellow
Write-Host "  # Note: Services need DATABASE_URL in .env" -ForegroundColor White
Write-Host "  # Each service requires port and database config" -ForegroundColor White
Write-Host ""

# Step 4
Write-Host "Langkah 4: Frontend" -ForegroundColor Cyan
Write-Host "-----------------" -ForegroundColor Yellow
Write-Host "  cd $projectRoot\frontend\unsia-portal-web" -ForegroundColor Gray
Write-Host "  npm install" -ForegroundColor Gray
Write-Host "  npm run dev" -ForegroundColor Gray
Write-Host ""

Write-Host "PORT:" -ForegroundColor Cyan
Write-Host "----" -ForegroundColor Yellow
Write-Host "  Frontend: http://localhost:3000" -ForegroundColor White
Write-Host "  Core API: http://localhost:8001" -ForegroundColor White
Write-Host ""
