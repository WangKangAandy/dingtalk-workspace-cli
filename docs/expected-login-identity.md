# Expected identity for a fresh authorization

This community-fork option preserves the expected-user requirement without
restoring the old per-sender token store:

```sh
dws auth login --device --no-browser --format json --expect-identity ding_example:user_123
```

The argument is one exact `corpId:userId`, supplied by the integration from
trusted message context. It is not a profile name, list, or identity supplied
by an LLM. It works before the user has a local Profile. Omit `--device` for
the ordinary OAuth flow. Device authorization still requires the user to open
the authorization URL; this option validates the result, not access to the URL.

The existing OAuth and device identity-enrichment hooks run while the new
credential is in memory. The guard then checks both IDs before either provider
persists the token. It does not populate identity from the expected value or
from a historical Profile. This check also applies to MCP-managed device
tokens: unavailable user identity is a failure, not a reason to skip checking.

| Result | Behavior |
|---|---|
| Both IDs match | Continue normal upstream Profile persistence and login output |
| Either ID differs | Authentication error, reason `identity_mismatch`, exit 2 |
| Either authenticated ID is unavailable | Authentication error, reason `identity_unverified`, exit 2 |
| Invalid selector or combined `--profile` / `--token` | Validation error, exit 3; no authorization starts |
| Option omitted | Preserve upstream login behavior, including unresolved identities |

The guard does not cover `auth exchange`, `auth provider-login`, importing an
auth bundle, or manual token storage. Those are different entry points; this
option belongs only to `auth login` and rejects its manual `--token` mode.
It does not revoke a token that has already been issued remotely to the
scanning account. It prevents that new credential from being saved locally
through this guarded login. Existing pre-login storage maintenance remains
upstream behavior; this is not a guarantee that no local metadata is touched.

The host must deliver the link to the originating conversation, keep login
attempts associated with their trusted expected identity, handle failure, and
invoke later business operations with the same exact `--profile`. Successful
login retains upstream current-profile behavior, so never depend on the global
current profile in a shared bot. Permission/scope errors are not login errors.

`auth login` remains an existing reviewed CLI-management Schema exclusion;
this change does not introduce an Agent business command or hand-written JSON
envelope. It preserves classified provider errors so the standard renderer can
publish the reason. Error messages do not include the credential or actual
other user's identity. The old `IDENTITY_*` protocol and exit 5 are not restored.

Tests exercise malformed options, both login-flow hooks, first login, wrong
organization/user, missing identity, historical identity hints, unchanged
unguarded behavior, and the actual provider pre-persistence boundaries using
isolated storage and fake network seams. They do not perform real account
authorization.
