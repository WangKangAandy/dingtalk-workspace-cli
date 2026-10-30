# Community fork baseline and integration contract

This fork tracks upstream stable `v1.0.66`
(`f215043e676d6889e1ddce3464828d8b1ec40c31`). The baseline upgrade merges
upstream history into the existing fork history. It does not reset `main`.

## Retired patches

- `--sender-id`, `DWS_AUTH_IDENTITY`, `users/<senderId>` token directories and
  the fork's process token cache are replaced by upstream exact Profiles.
- Business calls must receive `--profile <corpId>:<userId>` from trusted
  channel context. Do not select an organization-current account on behalf of
  a message sender, and do not omit the selector when it is unavailable.
- The custom `IDENTITY_NOT_AUTHENTICATED` / `IDENTITY_MISMATCH` JSON protocol
  is retired. Inspect upstream structured errors. Exit 5 means internal error,
  not unauthenticated; a missing profile can fail validation before a request.
- Shell-only `DWS_AGENT` installation patches are replaced by upstream
  `dws skill setup --mode multi --target openclaw --dry-run`. Upstream supports
  a canonical shared skill tree and target links, with compatibility fallbacks;
  this is not a promise to write exactly one directory.

## Requirements still owned by the integration

The fork's allowed runtime delta is only the optional pre-save expected-login
identity guard. Profile selection/storage, token refresh/cache, default login,
error rendering, product APIs and skill installation remain upstream behavior.
CI import checks, whitespace fixes, tests and maintainer documentation support
the upgrade; they do not introduce additional runtime features.

| Required behavior | Implementation owner and boundary |
|---|---|
| Trusted message sender selects an exact account | Connector derives `corpId:userId`; DWS uses upstream `--profile` |
| Missing authorization, link in current conversation, first login | Agent/connector coordinates the flow; DWS login uses `--expect-identity`, without `--profile` |
| Actual authorizer matches the requester | This fork's optional guard rejects mismatched or unverified identities before saving |
| Robot and user OAuth credentials stay separate | Connector configuration and child-process environment |
| Downloads use `media/inbound` | Deployment-managed Agent rules and explicit output paths |
| Skill installation | Upstream setup; no custom single-directory installer unless separately required |

An integration that injects `--profile` into every DWS command must distinguish
business calls from `auth login`. For login, inject the trusted
`--expect-identity` instead: a new user's Profile does not exist yet, and the
two flags are deliberately incompatible. The DWS PRs alone do not implement
this command routing or conversation delivery.

The connector owns the trusted sender-to-profile mapping, authorization-link
delivery to the originating conversation, login concurrency and expiry, and
matching the completed authorization to that sender. Profile selection alone
is not an expected-login-identity check. For first login, an absent exact
Profile is not a valid `auth login --profile` target.

Keep robot application credentials separate from user OAuth credentials.
Do not inject the robot AppKey/secret into business DWS invocations. Permission
or scope denial is not automatically a reason to restart OAuth.

Keep the following OpenClaw rule in the deployment's managed Agent rules,
rather than modifying upstream product skills:

> For local document/drive downloads, pass an explicit, expanded absolute
> output path below `$HOME/.openclaw/media/inbound/`. Do not write downloads
> into the agent workspace or the command's current working directory.

## Follow-up and release boundary

The follow-up change provides an optional
[`auth login --expect-identity corpId:userId`](expected-login-identity.md) guard
for first and subsequent authorization. It compares organization and user
identities after authoritative identity enrichment and before persisting the
new credential, for OAuth and device flows. Missing identity fails closed;
historical Profile names are not proof of the scanning user. Legacy MCP login
skipped this check.

This baseline PR does not deploy a binary, change a connector, migrate tokens,
or advance a deployment repository's gitlink. Those are separate integration
and release steps. Until the connector adopts these contracts, this branch
must not be treated as a drop-in replacement for the old sender-id CLI.

## Verification

The original baseline merge (`2e4649a6`) matches upstream `v1.0.66`, except
for this document. A subsequent CI repair normalizes whitespace in eight
upstream documentation/skill/script files and adapts the release-policy checks
for this one-time import. Runtime Go code in the baseline PR matches upstream;
the follow-up PR adds the expected-login-identity guard described above.

The import policy pins both the old fork base and the upstream commit, requires
both histories to remain ancestors, and checks that `CHANGELOG.md` and the
entire `.changes` tree are identical to upstream, including file modes. It
does not grant the changelog-only fast path. Later PRs and any changed release
metadata use the ordinary release rules. Policy tests cover rejected imports.

Use `git diff f215043e676d6889e1ddce3464828d8b1ec40c31` to inspect all fork
differences. These checks do not prove end-to-end integration or platform support.
