# Direcom Vision

Cross-platform filesystem tree monitoring with a Go CLI/server and an embedded web dashboard. Dashboard rendering is driven by `requestAnimationFrame` (up to 60 FPS on a 60 Hz display); the filesystem is scanned every 250 ms by default. Rendering at 60 FPS does not make the filesystem scanner run 60 times per second.

## Features

- Recursive file and directory inventory with create, modify, and delete activity history.
- Browser dashboard with tree filtering, counts, aggregate file size, and recent events.
- JSON snapshot endpoint at `/api/snapshot`.
- PowerShell and Bash launchers; Go implementation has no third-party dependencies.
- Optional Rust inventory helper using only the Rust standard library.

## Requirements

- Go 1.22 or later to run/build the monitor.
- Optional: Rust stable and Cargo for the inventory helper.
- Windows PowerShell 5.1+ or PowerShell 7+, or Bash on Linux/macOS.

## Run on PowerShell

From the cloned repository:

```powershell
.\scripts\run.ps1 -Path . -Address ":8080" -Interval "250ms"
```

Open <http://localhost:8080>. Press Ctrl+C to stop. You can also run the Go command directly:

```powershell
go run ./cmd/direcom-vision -root . -addr :8080 -interval 250ms
```

## Run on Bash

```bash
chmod +x scripts/run.sh
./scripts/run.sh . :8080 250ms
```

Or invoke Go directly:

```bash
go run ./cmd/direcom-vision -root . -addr :8080 -interval 250ms
```

Arguments are monitored-directory, listen-address, and scan interval for the Bash launcher. The server binds to all interfaces by default (`:8080`); use `127.0.0.1:8080` to restrict access to the local machine. This is a lightweight development dashboard and has no authentication; do not expose it to untrusted networks.

## Refresh and scan rates

The web UI paints continuously with `requestAnimationFrame`, typically up to 60 frames per second on a 60 Hz display (browser/device scheduling can cap or reduce this). The UI requests fresh monitor snapshots every 250 ms. The Go scanner also runs every 250 ms by default; set `-interval 1s` to reduce disk activity or choose another duration no shorter than 100 ms. A scan-based monitor can detect changes with a delay of up to approximately one scan interval and may not observe very short-lived files.

## Optional Rust inventory helper

```bash
cargo run --manifest-path rust/direcom-scan/Cargo.toml -- .
```

PowerShell equivalent:

```powershell
cargo run --manifest-path rust/direcom-scan/Cargo.toml -- .
```

It prints a JSON summary of file count, directory count, and total file bytes. Rust is an optional companion utility; the live dashboard and monitor run entirely in Go.

## Build and test

```bash
go test ./...
go build -o direcom-vision ./cmd/direcom-vision
```

On Windows, omit the output extension or use `-o direcom-vision.exe`.
