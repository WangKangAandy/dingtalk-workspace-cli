#!/bin/sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)"

cd "$ROOT"
go test ./test/mock_mcp/... \
  ./test/integration/discovery/... \
  ./test/integration/extensions/... \
  ./test/integration/recovery/...
