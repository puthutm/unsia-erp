# Reorganize Workspace Script
# Run this script in local PowerShell to flatten repository layout.

$root = $PSScriptRoot
$sourceDir = Join-Path $root "unsia-docs-md"
$docsTarget = Join-Path $root "docs"

if (-not (Test-Path $sourceDir)) {
    Write-Host "Directory 'unsia-docs-md' not found. Restructuring might already be completed." -ForegroundColor Yellow
    exit
}

# 1. Create docs directory if not exists
if (-not (Test-Path $docsTarget)) {
    New-Item -ItemType Directory -Path $docsTarget | Out-Null
    Write-Host "Created 'docs' directory." -ForegroundColor Green
}

# 2. Define docs items to be moved to docs/
$docsItems = @(
    "01-prd", "02-brd", "03-fsd", "04-api-contract", "05-event-contract", 
    "06-erd-dbml", "07-uat", "08-developer", "09-workplan", "10-repo-structure", 
    "11-srs", "UI", "backend_plan.md", "frontend_plan.md", "frontend_todo.md", 
    "FOLDER_TREE.md", "README-PROJECT.md", "README.md", "TODO.md", ".kiro"
)

foreach ($item in $docsItems) {
    $srcPath = Join-Path $sourceDir $item
    if (Test-Path $srcPath) {
        $destPath = Join-Path $docsTarget $item
        if (Test-Path $destPath) {
            Remove-Item -Recururse -Force $destPath
        }
        Move-Item -Path $srcPath -Destination $destPath -Force
        Write-Host "Moved docs item: $item" -ForegroundColor Cyan
    }
}

# 3. Move code folders to root workspace
$codeFolders = @("services", "packages", "frontend", "infra")
foreach ($folder in $codeFolders) {
    $srcPath = Join-Path $sourceDir $folder
    if (Test-Path $srcPath) {
        $destPath = Join-Path $root $folder
        if (Test-Path $destPath) {
            Write-Host "Destination folder '$folder' already exists in root. Merging contents..." -ForegroundColor Yellow
            Copy-Item -Path "$srcPath\*" -Destination $destPath -Recurse -Force
            Remove-Item -Path $srcPath -Recurse -Force
        } else {
            Move-Item -Path $srcPath -Destination $destPath -Force
        }
        Write-Host "Moved code folder: $folder" -ForegroundColor Green
    }
}

# 4. Move configuration files to root workspace
$configFiles = @("go.work", "go.work.sum", "docker-compose.yml", "run-wsl.sh", ".env.example")
foreach ($file in $configFiles) {
    $srcPath = Join-Path $sourceDir $file
    if (Test-Path $srcPath) {
        $destPath = Join-Path $root $file
        if (Test-Path $destPath) {
            Remove-Item -Path $destPath -Force
        }
        Move-Item -Path $srcPath -Destination $destPath -Force
        Write-Host "Moved configuration file: $file" -ForegroundColor Green
    }
}

# 5. Remove original unsia-docs-md directory if empty
$remaining = Get-ChildItem -Path $sourceDir
if ($remaining.Count -eq 0) {
    Remove-Item -Path $sourceDir -Force
    Write-Host "Cleaned up empty 'unsia-docs-md' directory." -ForegroundColor Green
} else {
    Write-Host "Some files remained in 'unsia-docs-md'. Please check manually." -ForegroundColor Yellow
}

Write-Host "Workspace reorganization successfully completed!" -ForegroundColor Green
