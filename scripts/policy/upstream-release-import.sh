#!/bin/sh

# One-time community-fork import, not a general release-policy bypass.
# Require the reviewed base and upstream commit, preserve their ancestry,
# and accept only byte-for-byte upstream release metadata (including modes).
# Later PRs use the ordinary release rules because their base will differ.
is_verified_upstream_release_import() {
  _uri_root="$1"
  _uri_base="$2"
  _uri_head="$3"
  _uri_expected_base=dd34e2241c479bde9d70d29c4847594b2544a60a
  _uri_upstream=f215043e676d6889e1ddce3464828d8b1ec40c31
  [ "$_uri_base" = "$_uri_expected_base" ] || return 1
  git -C "$_uri_root" merge-base --is-ancestor "$_uri_base" "$_uri_head" || return 1
  git -C "$_uri_root" merge-base --is-ancestor "$_uri_upstream" "$_uri_head" 2>/dev/null || return 1
  if git -C "$_uri_root" merge-base --is-ancestor "$_uri_upstream" "$_uri_base"; then
    return 1
  fi
  git -C "$_uri_root" diff --no-ext-diff --quiet "$_uri_upstream" "$_uri_head" -- CHANGELOG.md .changes
}
