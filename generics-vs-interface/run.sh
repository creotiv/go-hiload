#!/usr/bin/env bash
set -euo pipefail

GOCACHE="${GOCACHE:-$(pwd)/.gocache}" go test -bench . -benchmem
