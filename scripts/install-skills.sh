#!/bin/sh
set -eu

# Install DWS agent skills from GitHub Releases into one agent skill directory.
# Usage:
#   curl -fsSL .../install-skills.sh | sh
#   sh scripts/install-skills.sh --agent openclaw
#   DWS_SKILLS_ROOT=$HOME DWS_AGENT=openclaw sh scripts/install-skills.sh
#
# Default installs only to <root>/.agents/skills/dws.
# Pass --agent / DWS_AGENT to choose another home, or --agent all for the old
# multi-home copy (only when those parent dirs already exist).
#
# Environment variables (optional):
#   DWS_VERSION        — release tag (default: latest)
#   DWS_SKILLS_ROOT    — base path for agent dirs (default: $PWD)
#   DWS_AGENT          — agents | openclaw | claude | cursor | ... | all
#                        (default: agents; overridden by --agent)

REPO="DingTalk-Real-AI/dingtalk-workspace-cli"
VERSION="${DWS_VERSION:-latest}"
SKILL_NAME="dws"
ROOT="${DWS_SKILLS_ROOT:-$PWD}"
AGENT="${DWS_AGENT:-agents}"

while [ $# -gt 0 ]; do
  case "$1" in
    --agent)
      [ $# -ge 2 ] || { printf '❌ --agent needs a value\n' >&2; exit 1; }
      AGENT="$2"
      shift 2
      ;;
    --agent=*)
      AGENT="${1#*=}"
      shift
      ;;
    -h|--help)
      printf 'Usage: install-skills.sh [--agent agents|openclaw|claude|cursor|gemini|codex|github|windsurf|augment|cline|amp|kiro|trae|hermes|all]\n'
      printf 'Env: DWS_VERSION DWS_SKILLS_ROOT DWS_AGENT\n'
      exit 0
      ;;
    *)
      printf '❌ Unknown argument: %s\n' "$1" >&2
      exit 1
      ;;
  esac
done

# ── Helpers ──────────────────────────────────────────────────────────────────

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    printf '❌ Missing required command: %s\n' "$1" >&2
    exit 1
  fi
}

resolve_version() {
  if [ "$VERSION" = "latest" ]; then
    VERSION="$(curl -fsSI "https://github.com/${REPO}/releases/latest" 2>/dev/null \
      | grep -i '^location:' | sed 's|.*/tag/||;s/[[:space:]]*$//')"
    if [ -z "$VERSION" ]; then
      printf '❌ Could not determine the latest version. Set DWS_VERSION explicitly.\n' >&2
      exit 1
    fi
  fi
}

extract_zip() {
  archive="$1"
  dest="$2"
  if command -v unzip >/dev/null 2>&1; then
    unzip -q "$archive" -d "$dest"
    return 0
  fi
  if command -v tar >/dev/null 2>&1 && tar -xf "$archive" -C "$dest" >/dev/null 2>&1; then
    return 0
  fi
  printf '❌ Missing required command: unzip (or tar with zip support)\n' >&2
  exit 1
}

agent_skill_relpath() {
  case "$1" in
    agents) printf '%s\n' ".agents/skills" ;;
    openclaw) printf '%s\n' ".openclaw/skills" ;;
    claude) printf '%s\n' ".claude/skills" ;;
    cursor) printf '%s\n' ".cursor/skills" ;;
    gemini) printf '%s\n' ".gemini/skills" ;;
    codex) printf '%s\n' ".codex/skills" ;;
    github) printf '%s\n' ".github/skills" ;;
    windsurf) printf '%s\n' ".windsurf/skills" ;;
    augment) printf '%s\n' ".augment/skills" ;;
    cline) printf '%s\n' ".cline/skills" ;;
    amp) printf '%s\n' ".amp/skills" ;;
    kiro) printf '%s\n' ".kiro/skills" ;;
    trae) printf '%s\n' ".trae/skills" ;;
    hermes) printf '%s\n' ".hermes/skills" ;;
    all) printf '%s\n' "all" ;;
    *)
      printf '❌ Unknown agent %s (use agents|openclaw|claude|cursor|...|all)\n' "$1" >&2
      exit 1
      ;;
  esac
}

# One-line summary copy (2nd+ targets).
_copy_skill_summary() {
  _src="$1"
  _dest="$2"
  _label="$3"

  if [ -d "$_dest" ]; then
    rm -rf "$_dest"
  fi

  mkdir -p "$_dest"
  cp -R "$_src/"* "$_dest/" 2>/dev/null || cp -r "$_src/"* "$_dest/"
  file_count="$(find "$_dest" -type f | wc -l | tr -d ' ')"

  printf '  ✅ Skills → %s (%s files)\n' "$_label" "$file_count"
}

# Full copy with top-level listing (1st target).
_copy_skill() {
  _src="$1"
  _dest="$2"
  _label="$3"

  if [ -d "$_dest" ]; then
    rm -rf "$_dest"
  fi

  mkdir -p "$_dest"
  cp -R "$_src/"* "$_dest/" 2>/dev/null || cp -r "$_src/"* "$_dest/"
  file_count="$(find "$_dest" -type f | wc -l | tr -d ' ')"

  printf '  ✅ Skills → %s (%s files)\n' "$_label" "$file_count"

  for entry in "$_dest"/*; do
    entry_name="$(basename "$entry")"
    if [ -d "$entry" ]; then
      sub_count="$(find "$entry" -type f | wc -l | tr -d ' ')"
      printf '     📁 %s/ (%s files)\n' "$entry_name" "$sub_count"
    else
      printf '     📄 %s\n' "$entry_name"
    fi
  done
}

_copy_one() {
  skill_src="$1"
  root="$2"
  agent_dir="$3"
  first="$4"
  dest="$root/$agent_dir/$SKILL_NAME"
  if [ "$root" = "$HOME" ]; then
    label="~/$agent_dir/$SKILL_NAME"
  else
    label="$root/$agent_dir/$SKILL_NAME"
  fi
  if [ "$first" = "1" ]; then
    _copy_skill "$skill_src" "$dest" "$label"
  else
    _copy_skill_summary "$skill_src" "$dest" "$label"
  fi
}

# Default: one agent dir. --agent all keeps the old multi-home copy.
install_skills_to_root() {
  skill_src="$1"
  root="$2"
  agent_key="$3"
  relpath="$(agent_skill_relpath "$agent_key")"

  if [ "$relpath" != "all" ]; then
    _copy_one "$skill_src" "$root" "$relpath" 1
    return 0
  fi

  installed=0
  idx=0
  for agent_dir in \
    ".agents/skills" \
    ".claude/skills" \
    ".cursor/skills" \
    ".gemini/skills" \
    ".codex/skills" \
    ".github/skills" \
    ".windsurf/skills" \
    ".augment/skills" \
    ".cline/skills" \
    ".amp/skills" \
    ".kiro/skills" \
    ".trae/skills" \
    ".openclaw/skills" \
    ".hermes/skills"
  do
    base_dir="$root/$agent_dir"
    parent_gate="$(dirname "$base_dir")"
    if [ "$idx" -gt 0 ] && [ ! -e "$parent_gate" ]; then
      idx=$((idx + 1))
      continue
    fi
    if [ "$installed" -eq 0 ]; then
      _copy_one "$skill_src" "$root" "$agent_dir" 1
    else
      _copy_one "$skill_src" "$root" "$agent_dir" 0
    fi
    installed=$((installed + 1))
    idx=$((idx + 1))
  done
  if [ "$installed" -eq 0 ]; then
    _copy_one "$skill_src" "$root" ".agents/skills" 1
  fi
}

# ── Main ─────────────────────────────────────────────────────────────────────

main() {
  need_cmd curl
  resolve_version
  # Validate agent early.
  agent_skill_relpath "$AGENT" >/dev/null

  printf '\n'
  printf '  ┌──────────────────────────────────────┐\n'
  printf '  │     DWS Skill Installer              │\n'
  printf '  │     DingTalk Workspace CLI            │\n'
  printf '  └──────────────────────────────────────┘\n'
  printf '\n'

  TMPDIR_WORK="$(mktemp -d)"
  trap 'rm -rf "$TMPDIR_WORK"' EXIT INT TERM

  ASSET_URL="https://github.com/${REPO}/releases/download/${VERSION}/dws-skills.zip"
  printf '  ⬇  Downloading skills from GitHub Releases: %s (%s)\n' "$REPO" "$VERSION"
  curl -fsSL "$ASSET_URL" -o "$TMPDIR_WORK/dws-skills.zip"
  extract_zip "$TMPDIR_WORK/dws-skills.zip" "$TMPDIR_WORK/extracted"

  SKILL_SRC="$TMPDIR_WORK/extracted"
  if [ -f "$TMPDIR_WORK/extracted/${SKILL_NAME}/SKILL.md" ]; then
    SKILL_SRC="$TMPDIR_WORK/extracted/${SKILL_NAME}"
  fi

  if [ ! -f "$SKILL_SRC/SKILL.md" ]; then
    printf '  ❌ Skill source not found in release asset\n' >&2
    exit 1
  fi

  printf '\n'
  printf '  Installing under root: %s (agent=%s)\n' "$ROOT" "$AGENT"
  install_skills_to_root "$SKILL_SRC" "$ROOT" "$AGENT"

  printf '\n'
  printf '  📖 Skill includes:\n'
  printf '     • SKILL.md — Main skill with product overview and intent routing\n'
  printf '     • references/ — Detailed product command references\n'
  printf '     • scripts/ — Batch operation scripts for all products\n'
  printf '\n'
  printf '  ⚡ Requires: dws CLI installed and on $PATH\n'
  printf '     Install: go install github.com/%s/cmd@latest\n' "$REPO"
  printf '\n'
}

main
