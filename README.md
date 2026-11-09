# Direcom Vision

Cross-platform file system monitoring with a Go engine, terminal CLI and desktop-style web dashboard. Dashboard animation targets 60 FPS on a typical 60 Hz display; filesystem scanning defaults to every 250 ms.

## Features

- Recursive tree inventory with create, modify and delete tracking.
- Detailed event history with timestamps, sequence numbers, file sizes and previous sizes.
- Advanced tree filters: search by path, extension, kind, min size and sorting.
- Terminal mode for live CLI output.
- Desktop-style dashboard with dark UI and fast animation.
- JSON and CSV export endpoints for history logs.
- PowerShell and Bash launchers plus optional Rust helper.

## Requirements

- Go 1.22 or later
- Optional: Rust stable + Cargo
- PowerShell 5.1+, PowerShell 7+, or Bash on Linux/macOS

## Run in PowerShell

Dashboard:

```powershell
.\scripts\run.ps1 -Path . -Address ":8080" -Interval "250ms"
```

CLI mode:

```powershell
go run ./cmd/direcom-vision -root . -mode cli -interval 250ms
```

## Run in Bash

```bash
chmod +x scripts/run.sh
./scripts/run.sh . :8080 250ms
```

CLI mode:

```bash
go run ./cmd/direcom-vision -root . -mode cli -interval 250ms
```

## Dashboard filters and exports

The filter box supports path name search, extension filtering, directory/file-only filtering, minimum size, and sorting by path, size, or modified time. Use the export buttons or call `/api/export?format=json` or `/api/export?format=csv` to download a recent history snapshot.

## Refresh rate and scan tuning

The browser UI paints via `requestAnimationFrame`, typically up to 60 FPS on a 60 Hz display. The monitor scans the file tree at default 250 ms. Use a larger interval to reduce CPU and disk activity; values under 100 ms are clamped.

## Optional Rust helper

```bash
cargo run --manifest-path rust/direcom-scan/Cargo.toml -- .
```

PowerShell:

```powershell
cargo run --manifest-path rust/direcom-scan/Cargo.toml -- .
```

## Build and test

```bash
go test ./...
go build -o direcom-vision ./cmd/direcom-vision
```
