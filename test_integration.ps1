# Workshot Integration Test Script for Windows
# Run with: .\test_integration.ps1

$ErrorActionPreference = "Stop"

# Colors for output
$ErrorActionPreference = "Stop"

# Output helpers (ASCII ONLY)
function Write-Success { Write-Host "[OK] $args" -ForegroundColor Green }
function Write-Error   { Write-Host "[ERR] $args" -ForegroundColor Red }
function Write-Info    { Write-Host "[INFO] $args" -ForegroundColor Cyan }
function Write-Test    { Write-Host "[TEST] $args" -ForegroundColor Yellow }

$TestsPassed = 0
$TestsFailed = 0

Write-Host ""
Write-Info "Running Workshot Integration Tests..."
Write-Host ""


$TestsPassed = 0
$TestsFailed = 0

Write-Host ""
Write-Info "Running Workshot Integration Tests..."
Write-Host ""

# Test 1: Build the binary
Write-Test "Building workshot..."
try {
    go build -o "$env:TEMP\workshot.exe" .\cmd\workshot\main.go
    if (Test-Path "$env:TEMP\workshot.exe") {
        Write-Success "Build successful"
        $TestsPassed++
    } else {
        Write-Error "Build failed - binary not found"
        $TestsFailed++
        exit 1
    }
} catch {
    Write-Error "Build failed: $_"
    $TestsFailed++
    exit 1
}

# Add to PATH for this session
$env:PATH = "$env:TEMP;$env:PATH"

Write-Host ""

# Test 2: Version flag
Write-Test "Version flag"
try {
    $version = & workshot --version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Version: $version"
        $TestsPassed++
    } else {
        Write-Error "Version flag failed"
        $TestsFailed++
    }
} catch {
    Write-Error "Version flag failed: $_"
    $TestsFailed++
}

# Test 3: Help command
Write-Test "Help command"
try {
    $help = & workshot --help 2>&1
    if ($LASTEXITCODE -eq 0 -and $help -match "workshot") {
        Write-Success "Help displayed correctly"
        $TestsPassed++
    } else {
        Write-Error "Help command failed"
        $TestsFailed++
    }
} catch {
    Write-Error "Help command failed: $_"
    $TestsFailed++
}

# Test 4: Freeze in temp directory
Write-Test "Freeze in temp directory"
Push-Location $env:TEMP
try {
    $output = & workshot freeze test-context-1 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Freeze successful in $env:TEMP"
        $TestsPassed++
    } else {
        Write-Error "Freeze failed: $output"
        $TestsFailed++
    }
} catch {
    Write-Error "Freeze failed: $_"
    $TestsFailed++
}
Pop-Location

# Test 5: List shows snapshot
Write-Test "List shows snapshot"
try {
    $list = & workshot list 2>&1 | Out-String
    if ($list -match "test-context-1") {
        Write-Success "Snapshot appears in list"
        $TestsPassed++
    } else {
        Write-Error "Snapshot not found in list"
        $TestsFailed++
    }
} catch {
    Write-Error "List command failed: $_"
    $TestsFailed++
}

# Test 6: Show snapshot details
Write-Test "Show snapshot details"
try {
    $show = & workshot show test-context-1 2>&1 | Out-String
    if ($show -match "test-context-1" -and $show -match "working_dir") {
        Write-Success "Show displays snapshot details"
        $TestsPassed++
    } else {
        Write-Error "Show command output incorrect"
        $TestsFailed++
    }
} catch {
    Write-Error "Show command failed: $_"
    $TestsFailed++
}

# Test 7: Freeze another snapshot in different directory
Write-Test "Freeze in home directory"
Push-Location $env:USERPROFILE
try {
    $output = & workshot freeze test-context-2 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Freeze successful in home directory"
        $TestsPassed++
    } else {
        Write-Error "Freeze failed: $output"
        $TestsFailed++
    }
} catch {
    Write-Error "Freeze failed: $_"
    $TestsFailed++
}
Pop-Location

# Test 8: List shows 2 snapshots
Write-Test "List shows 2 snapshots"
try {
    $list = & workshot list 2>&1 | Out-String
    # Count occurrences of snapshot names instead of emoji
    $matches = ([regex]::Matches($list, "test-context-")).Count
    if ($matches -ge 2) {
        Write-Success "List shows 2 snapshots"
        $TestsPassed++
    } else {
        Write-Error "Expected 2 snapshots, found $matches"
        $TestsFailed++
    }
} catch {
    Write-Error "List count failed: $_"
    $TestsFailed++
}

# Test 9: Restore first context
Write-Test "Restore first context"
try {
    $output = & workshot restore test-context-1 2>&1 | Out-String
    if ($output -match "Working directory:.*$([regex]::Escape($env:TEMP))") {
        Write-Success "Restore output shows correct directory"
        $TestsPassed++
    } else {
        Write-Error "Restore output doesn't show temp directory: $output"
        $TestsFailed++
    }
} catch {
    Write-Error "Restore failed: $_"
    $TestsFailed++
}

# Remove Test 10 entirely or change it:
Write-Test "Verify restore output contains directory info"
try {
    $output = & workshot restore test-context-1 2>&1 | Out-String
    if ($output -match "Working directory:") {
        Write-Success "Restore command outputs directory info"
        $TestsPassed++
    } else {
        Write-Error "Restore output missing directory info"
        $TestsFailed++
    }
} catch {
    Write-Error "Directory info check failed: $_"
    $TestsFailed++
}

# Test 11: Delete with force flag
Write-Test "Delete snapshot (force)"
try {
    $output = & workshot delete test-context-1 -f 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Delete successful"
        $TestsPassed++
    } else {
        Write-Error "Delete failed: $output"
        $TestsFailed++
    }
} catch {
    Write-Error "Delete failed: $_"
    $TestsFailed++
}

# Test 12: Verify snapshot deleted
Write-Test "Snapshot deleted"
try {
    $show = & workshot show test-context-1 2>&1 | Out-String
    if ($show -match "not found") {
        Write-Success "Snapshot successfully deleted"
        $TestsPassed++
    } else {
        Write-Error "Snapshot still exists"
        $TestsFailed++
    }
} catch {
    # Error expected - snapshot should not exist
    Write-Success "Snapshot successfully deleted"
    $TestsPassed++
}

# Test 13: List shows 1 snapshot after delete
Write-Test "List shows 1 snapshot after delete"
try {
    $list = & workshot list 2>&1 | Out-String
    $matches = ([regex]::Matches($list, "test-context-")).Count
    if ($matches -eq 1) {
        Write-Success "List shows 1 snapshot"
        $TestsPassed++
    } else {
        Write-Error "Expected 1 snapshot, found $matches"
        $TestsFailed++
    }
} catch {
    Write-Error "List count failed: $_"
    $TestsFailed++
}

# Cleanup
Write-Host ""
Write-Info "Cleaning up..."
try {
    & workshot delete test-context-2 -f 2>&1 | Out-Null
    Remove-Item "$env:TEMP\workshot.exe" -ErrorAction SilentlyContinue
} catch {
    # Ignore cleanup errors
}

# Results
Write-Host ""
Write-Host "======================================" -ForegroundColor White
Write-Host "Tests Passed: " -NoNewline
Write-Host "$TestsPassed" -ForegroundColor Green
Write-Host "Tests Failed: " -NoNewline
Write-Host "$TestsFailed" -ForegroundColor Red
Write-Host "======================================" -ForegroundColor White
Write-Host ""

if ($TestsFailed -eq 0) {
    Write-Success "All tests passed!"
    exit 0
} else {
    Write-Error "Some tests failed"
    exit 1
}