# Direcom Vision

Cross-platform filesystem tree monitoring with a Go engine, terminal CLI and embedded desktop-style web dashboard. Dashboard animation uses `requestAnimationFrame` (up to 60 FPS on a 60 Hz display); filesystem scanning defaults to every 250 ms.

## Features

- Recursive inventory with detailed create/modify/delete history, timestamps, sequence numbers and file sizes.
- Advanced tree filters: path/name search, file extension, file/directory kind, minimum size and sort order.
- Terminal mode with initial inventory and live change stream.
- Desktop-style responsive web dashboard, activity timeline and JSON/CSV export.
- JSON snapshot at `/api/snapshot`; history exports at `/api/export?format=json` or `?format=csv`.
- PowerShell and Bash launchers; Go implementation uses no third-party dependencies. Optional Rust inventory helper is included.

## Requirements

Go 1.22 or later. Optional Rust stable/Cargo. Windows PowerShell 5.1+ or PowerShell 7+, or Bash on Linux/macOS.

## PowerShell

Web dashboard:

```powershell
.\scripts\run.ps1 -Path . -Address ":8080" -Interval "250ms"
```

Open <http://localhost:8080>. Terminal mode:

```powershell
go run ./cmd/direcom-vision -root . -mode cli -interval 250ms
```

## Bash

```bash
chmod +x scripts/run.sh
./scripts/run.sh . :8080 250ms
```

Terminal mode:

```bash
go run ./cmd/direcom-vision -root . -mode cli -interval 250ms
```

Bash launcher arguments are directory, address, and scan interval. The server binds to all interfaces by default; use `127.0.0.1:8080` to restrict access to localhost. The dashboard has no authentication; do not expose it to untrusted networks.

## Dashboard filters and exports

Use the search field to match any path substring. Extension accepts values such as `.go` or `go`; kind switches between all/files/directories; minimum size applies to files; sort supports path, size, and modified time. Export buttons download the retained event history as JSON or CSV. History is kept in memory and limited to the most recent 5,000 events; it is not persisted between runs.

## Refresh rate

UI paints with `requestAnimationFrame`, typically up to 60 FPS on a 60 Hz display. Dashboard data is polled every 250 ms. Go scans every 250 ms by default; `-interval 1s` reduces disk activity. Intervals below 100 ms are clamped. Polling scans can miss short-lived files and have detection latency up to roughly one scan interval.

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
