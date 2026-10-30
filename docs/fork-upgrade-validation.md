# Fork upgrade validation — 2026-10-30

Baseline: upstream `v1.0.66`, `f215043e676d6889e1ddce3464828d8b1ec40c31`.
Environment: Windows amd64, Go 1.26.7, `CGO_ENABLED=0`,
`DWS_PACKAGE_VERSION=0.0.0-test`. No real account authorization or deployment.

## Passed

- Baseline merge tree equals upstream except `docs/FORK_UPGRADE.md`.
- CLI build (`go build -o <temporary executable> ./cmd`). This is a stub build
  for optional CGO functionality, not validation of SafeChat.
- Focused auth/app tests: expected-identity validation, both OAuth/device
  hooks, first login, different user/organization, missing identity, history
  not proving identity, unguarded compatibility, and actual providers rejecting
  before token save. Existing same-organization account storage, exact-profile
  logout, TokenManager and command-tree completeness tests also passed.
- Real CLI smoke checks with temporary HOME/config: help exposes the flag;
  invalid selector and manual-token conflict both return validation/exit 3.
- `scripts/policy/check-generated-drift.sh`, including two byte-identical
  Schema assemblies. Windows `core.autocrlf=true` initially made the byte
  comparison fail; the two generated working files were normalized to the
  upstream LF bytes for this check, without changing their Git content.
  The optional protobuf regeneration check was skipped because libprotoc
  35.1 is not installed.

## Full-suite limits

`go test ./...` was attempted but not passed. The app package hit its ten-minute
timeout; auth tests also reported failures. The broader run was stopped after
these failures, so this is not a complete full-suite result.

The following failures were reproduced separately on an untouched worktree of
the upstream commit with `go test -p 1 ... -count=1`:

- `TestCrossPlatformCoveragePortableAuthBundleCoverageEdges`: Windows denies
  creation of the test's symbolic link (`A required privilege is not held`).
- `TestAuthStatusJSONReportsKeychainUnavailable` and
  `TestAuthStatusJSONReportsDEKMissing`: TempDir cleanup cannot remove an open
  `config/logs/dws.log` on Windows.

The full run additionally encountered local HTTP test timeouts and a readonly
status timing assertion. Those failures and the package timeout were not
separately reproduced on upstream and are not classified as proven baseline
defects. No unrelated production-code changes were made to silence them.

Keep both PRs in draft until CI/platform checks and connector integration are
reviewed. These results do not establish production readiness, Linux/CGO
compatibility, or the success of an actual DingTalk authorization exchange.
