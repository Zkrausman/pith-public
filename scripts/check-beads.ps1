# [bd-3qwe self-doc] RATIONALE: This script ensures all closed beads in the last 24h
# have mandatory RATIONALE and EVIDENCE anchors, as per the Beads Closure Mandate.

$checkDuration = (Get-Date).AddHours(-24)
$beads = bd ready --all | ConvertFrom-Json

# Heuristic: Check recently closed beads (in the last 24h)
$recentClosed = $beads | Where-Object { $_.status -eq 'closed' -and (Get-Date $_.updated_at) -gt $checkDuration }

if ($recentClosed.Count -eq 0) {
    Write-Host "No beads closed in the last 24h. Skipping validation." -ForegroundColor Gray
    exit 0
}

$failCount = 0

foreach ($bead in $recentClosed) {
    Write-Host "Validating bead $($bead.id): $($bead.title)" -ForegroundColor Cyan
    
    # Get bead notes
    $notes = bd show $($bead.id) --notes
    
    $hasRationale = $notes -match "RATIONALE"
    $hasEvidence = $notes -match "EVIDENCE"
    
    if (-not $hasRationale) {
        Write-Host "  [FAIL] Missing 'RATIONALE' anchor in notes." -ForegroundColor Red
        $failCount++
    }
    
    if (-not $hasEvidence) {
        Write-Host "  [FAIL] Missing 'EVIDENCE' anchor in notes." -ForegroundColor Red
        $failCount++
    }
    
    if ($hasRationale -and $hasEvidence) {
        Write-Host "  [PASS] Compliant." -ForegroundColor Green
    }
}

if ($failCount -gt 0) {
    Write-Host "`nFound $failCount non-compliant beads. Documentation is mandatory for closure." -ForegroundColor Red
    exit 1
}

Write-Host "`nAll recently closed beads are compliant." -ForegroundColor Green
exit 0
