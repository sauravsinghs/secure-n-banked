# Run after starting Docker Desktop.
# Rotates the Postgres password and clears DB-backed sessions.

$ErrorActionPreference = "Stop"

$appEnvPath = Join-Path $PSScriptRoot "..\app.env"
if (-not (Test-Path $appEnvPath)) {
    throw "app.env not found. Copy app.env.example to app.env and set POSTGRES_PASSWORD first."
}

$newPasswordLine = Get-Content -LiteralPath $appEnvPath | Where-Object { $_ -match '^POSTGRES_PASSWORD=' } | Select-Object -First 1
if (-not $newPasswordLine) {
    throw "POSTGRES_PASSWORD is missing from app.env."
}

$newPassword = $newPasswordLine.Substring("POSTGRES_PASSWORD=".Length)

Write-Host "Starting postgres..."
docker compose up -d postgres

Write-Host "Waiting for postgres..."
$ready = $false
for ($i = 0; $i -lt 30; $i++) {
    docker compose exec -T postgres pg_isready -U root 2>$null
    if ($LASTEXITCODE -eq 0) {
        $ready = $true
        break
    }
    Start-Sleep -Seconds 2
}

if (-not $ready) {
    throw "Postgres did not become ready in time."
}

Write-Host "Rotating DB password and truncating sessions..."
docker compose exec -T postgres psql -U root -d simple_bank -c "ALTER USER root WITH PASSWORD '$newPassword'; TRUNCATE TABLE sessions;"

Write-Host "Verifying sessions table is empty..."
docker compose exec -T postgres psql -U root -d simple_bank -c "SELECT COUNT(*) AS session_count FROM sessions;"

Write-Host "Done. docker-compose.yaml now reads POSTGRES_PASSWORD from app.env."
