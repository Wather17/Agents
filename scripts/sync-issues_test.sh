#!/usr/bin/env bash

set -Eeuo pipefail

root_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
script_path="$root_dir/scripts/sync-issues.sh"
template_path="$root_dir/agent-init/internal/templates/files/sync-issues.sh"
test_path="$(mktemp -d)"
trap 'rm -rf -- "$test_path"' EXIT

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  exit 1
}

assert_file_contains() {
  local file="$1"
  local expected="$2"
  grep -F -- "$expected" "$file" >/dev/null \
    || fail "expected $file to contain: $expected"
}

run_success_case() {
  local script="$1"
  local name="$2"
  local target="$test_path/$name"
  local log="$target/gh.log"

  mkdir -p -- "$target/issues"
  printf 'stale cache\n' > "$target/issues/stale.md"

  if ! (
    cd -- "$target"
    GH_CALL_LOG="$log" \
      PATH="$root_dir/scripts/testdata:$PATH" \
      bash "$script"
  ); then
    fail "$name synchronization failed"
  fi

  [[ ! -e "$target/issues/stale.md" ]] \
    || fail "$name kept a stale issue file after a successful sync"
  [[ -f "$target/issues/42-first-issue.md" ]] \
    || fail "$name did not create the titled issue file"
  [[ -f "$target/issues/7-issue.md" ]] \
    || fail "$name did not create the fallback-slug issue file"
  [[ "$(wc -l < "$log")" -eq 1 ]] \
    || fail "$name made more than one GitHub CLI call"
  ! grep -F -- 'issue view' "$log" >/dev/null \
    || fail "$name still called gh issue view"

  assert_file_contains "$target/issues/42-first-issue.md" '**Labels**: bug, status:ready'
  assert_file_contains "$target/issues/42-first-issue.md" $'line two\twith tab'
  assert_file_contains "$target/issues/42-first-issue.md" '### Comentário por @alice:'
  assert_file_contains "$target/issues/42-first-issue.md" 'comentário 2'
}

run_failure_case() {
  local target="$test_path/failure"
  local log="$target/gh.log"

  mkdir -p -- "$target/issues"
  printf 'existing cache\n' > "$target/issues/existing.md"

  if (
    cd -- "$target"
      GH_MODE=failure \
      GH_CALL_LOG="$log" \
      PATH="$root_dir/scripts/testdata:$PATH" \
      bash "$script_path"
  ); then
    fail 'a failed GitHub query returned success'
  fi

  assert_file_contains "$target/issues/existing.md" 'existing cache'
}

run_empty_case() {
  local target="$test_path/empty"
  local log="$target/gh.log"

  mkdir -p -- "$target/issues"
  printf 'stale cache\n' > "$target/issues/stale.md"

  if ! (
    cd -- "$target"
      GH_MODE=empty \
      GH_CALL_LOG="$log" \
      PATH="$root_dir/scripts/testdata:$PATH" \
      bash "$script_path"
  ); then
    fail 'empty issue synchronization failed'
  fi

  [[ ! -e "$target/issues/stale.md" ]] \
    || fail 'empty synchronization kept a stale issue file'
  [[ "$(find "$target/issues" -maxdepth 1 -type f -name '*.md' -print | wc -l)" -eq 0 ]] \
    || fail 'empty synchronization created issue files'
}

cmp -s "$script_path" "$template_path" \
  || fail 'installed script and embedded template differ'

run_success_case "$script_path" root
run_success_case "$template_path" template
run_failure_case
run_empty_case

printf 'PASS: sync-issues.sh aggregation and cache safety checks\n'
