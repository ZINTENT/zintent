# ZINTENT Build Script (v2)
# Features: Container-First, Tokens, Animations, Antigravity

Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "               ZINTENT v2.0 - Build System                  " -ForegroundColor Cyan
Write-Host "  Container-First | AI Tokens | Animations | Antigravity   " -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

# Check for Go
$hasGo = $null -ne (Get-Command go -ErrorAction SilentlyContinue)

if (-not $hasGo) {
    Write-Host "WARNING: Go not found. Please install Go 1.21+ to use the native compiler." -ForegroundColor Yellow
    Write-Host "   Falling back to Node.js prototype compiler." -ForegroundColor Gray
}

# Determine which compiler to use
$useV2 = Test-Path ".\compiler\main-v2.go"
$sourceFile = "main-v2.go"
$targetFile = "zintent-v2.exe"

if (-not $useV2) {
    $sourceFile = "main.go"
    $targetFile = "zintent.exe"
}

# Compile Go Engine (entire compiler package: main-v2.go + scanner_parser.go, etc.)
if ($hasGo -and (Test-Path ".\compiler\$sourceFile")) {
    Write-Host "Compiling ZINTENT Engine (v2.0)..." -ForegroundColor Yellow
    go build -o $targetFile ./compiler
    if (Test-Path ".\$targetFile") {
        Write-Host "Engine Built: $targetFile" -ForegroundColor Green
    } else {
        Write-Host "Build failed" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "Using JavaScript compiler" -ForegroundColor Yellow
}

# Determine input/output files
$inputFile = "src\phase1-demo.html"
$outputFile = "dist\phase1-styles.css"

if ($args[0]) {
    $inputFile = $args[0]
}
if ($args[1]) {
    $outputFile = $args[1]
}

# Ensure dist directory exists
if (-not (Test-Path "dist")) {
    New-Item -ItemType Directory -Name "dist" | Out-Null
    Write-Host "Created dist/ directory" -ForegroundColor Gray
}

# Build process
Write-Host "Building: $inputFile -> $outputFile" -ForegroundColor Cyan

$exePath = Join-Path "." $targetFile

if (Test-Path $exePath) {
    & $exePath $inputFile -o $outputFile
} elseif ($hasGo) {
    go run .\compiler\$sourceFile $inputFile -o $outputFile
} else {
    node .\compiler\zintent-compiler.js $inputFile -o $outputFile
}

if (Test-Path $outputFile) {
    $fileSize = (Get-Item $outputFile).Length
    $fileSizeKB = [math]::Round($fileSize / 1024, 2)
    Write-Host ""
    Write-Host "Build Complete!" -ForegroundColor Green
    Write-Host "   Output: $outputFile ($fileSizeKB KB)" -ForegroundColor White
    Write-Host ""
    Write-Host "Features Included:" -ForegroundColor Cyan
    Write-Host "   - Container-First Responsive (@container queries)" -ForegroundColor Gray
    Write-Host "   - AI Design Token System" -ForegroundColor Gray
    Write-Host "   - Intent-Driven Animations (physics-based)" -ForegroundColor Gray
    Write-Host "   - Antigravity Layout Engine" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "   npm run dev     - Start watch mode" -ForegroundColor White
    Write-Host "   npm run tokens:generate  - Generate AI design tokens" -ForegroundColor White
} else {
    Write-Host "Build failed - output file not created" -ForegroundColor Red
}
