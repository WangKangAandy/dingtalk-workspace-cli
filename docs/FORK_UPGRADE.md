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

An optional expected-identity login guard is a separate change on this baseline.
It must compare both organization and user identities after authoritative
identity enrichment and before persisting the new credential, for OAuth and
device flows. Missing identity must fail closed; historical Profile names are
not proof of the scanning user. Legacy MCP login skipped this check.

This baseline PR does not deploy a binary, change a connector, migrate tokens,
or advance a deployment repository's gitlink. Those are separate integration
and release steps. Until the connector adopts these contracts, this branch
must not be treated as a drop-in replacement for the old sender-id CLI.

## Verification

The baseline's tracked files, excluding this document, exactly match upstream
`v1.0.66`. Use `git diff v1.0.66 -- . ':!docs/FORK_UPGRADE.md'` with a locally
available tag (or its pinned commit) to verify. This proves merge resolution,
not end-to-end integration or platform support.
