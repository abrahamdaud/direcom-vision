param(
    [string]$Path = ".",
    [string]$Address = ":8080",
    [string]$Interval = "250ms"
)
$ErrorActionPreference = "Stop"
$repoRoot = Split-Path -Parent $PSScriptRoot
Push-Location $repoRoot
try {
    go run ./cmd/direcom-vision -root $Path -addr $Address -interval $Interval
} finally {
    Pop-Location
}
