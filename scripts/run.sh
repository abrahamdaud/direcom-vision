#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PATH_TO_WATCH="${1:-.}"
ADDRESS="${2:-:8080}"
INTERVAL="${3:-250ms}"
cd "$ROOT"
exec go run ./cmd/direcom-vision -root "$PATH_TO_WATCH" -addr "$ADDRESS" -interval "$INTERVAL"
