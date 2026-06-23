# Maintainer Automation Notes

This document keeps agent- and maintainer-specific workflow notes out of the
repository root while preserving repo-local guidance for automation.

## Read Order

1. `README.md`
2. `CONTRIBUTING.md`
3. `docs/architecture.md`
4. This document

## Project Snapshot

- `dws` is a Go-based DingTalk Workspace CLI and MCP runtime bridge.
- One internal Tool IR drives canonical CLI, schema, docs, skills, and snapshots.
- Compatibility and helper surfaces are overlays, not the canonical truth.

## Repository Map

- `cmd`: public CLI entrypoint
- `internal/app`: root command wiring and command tree mount points
- `internal/discovery`, `internal/market`, `internal/transport`: runtime discovery and MCP transport
- `internal/generator`: CLI/schema/docs/skills generation pipeline
- `internal/compat`, `internal/helpers`: legacy-compatible aliases and helper commands
- `docs/`: public architecture and reference docs
- `hack/`: developer-only helper commands not shipped as public binaries
- `scripts/`: build, test, lint, packaging, and policy checks
- `test/`: integration, contract, compatibility, and script validation suites

## Task Routing

- Add or fix a command path: start from `internal/app` and the related module under `internal/*`
- Discovery or protocol issues: inspect `internal/discovery`, `internal/market`, `internal/transport`
- Generated output drift: inspect `internal/generator` and run drift checks
- Legacy behavior mismatch: inspect `internal/compat` and `test/cli_compat`
- Failure or degraded mode: inspect `internal/discovery`, `internal/errors`

## Generated Artifacts

Prefer editing source logic instead of generated files directly.

- Generated-heavy paths:
  - `docs/generated/`
  - `skills/generated/`
  - `test/golden/generated_outputs/`
- When generator or command surface changes, run:
  - `./scripts/policy/check-generated-drift.sh`
  - `./scripts/policy/check-command-surface.sh --strict`

## Common Commands

```bash
make build
make test
make lint
./scripts/dev/ci-local.sh
./scripts/policy/check-generated-drift.sh
./scripts/policy/check-command-surface.sh --strict
./scripts/policy/check-open-source-assets.sh
git diff --check
```

## Handoff Checklist

Before handoff, include:

1. Changed files and why
2. Verification commands run and outcomes
3. Known risks or follow-up work
