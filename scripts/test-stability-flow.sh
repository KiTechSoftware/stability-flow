#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(mktemp -d)"
REPO_DIR="${ROOT_DIR}/stability-flow-repo"

cleanup() {
  rm -rf "${ROOT_DIR}"
}
trap cleanup EXIT

log() {
  printf "\n==> %s\n" "$*"
}

pass() {
  printf "[PASS] %s\n" "$*"
}

fail() {
  printf "[FAIL] %s\n" "$*" >&2
  exit 1
}

assert_branch_exists() {
  local branch="$1"
  git show-ref --verify --quiet "refs/heads/${branch}" \
    || fail "Expected branch '${branch}' to exist"
}

assert_branch_not_exists() {
  local branch="$1"
  if git show-ref --verify --quiet "refs/heads/${branch}"; then
    fail "Expected branch '${branch}' to be deleted"
  fi
}

assert_equal_commit() {
  local a="$1"
  local b="$2"
  local sha_a sha_b
  sha_a="$(git rev-parse "$a")"
  sha_b="$(git rev-parse "$b")"
  [[ "$sha_a" == "$sha_b" ]] || fail "Expected '${a}' and '${b}' to point to same commit"
}

assert_ancestor() {
  local ancestor="$1"
  local descendant="$2"
  git merge-base --is-ancestor "$ancestor" "$descendant" \
    || fail "Expected '${ancestor}' to be an ancestor of '${descendant}'"
}

assert_tag_exists() {
  local tag="$1"
  git rev-parse "$tag^{commit}" >/dev/null 2>&1 \
    || fail "Expected tag '${tag}' to exist"
}

assert_branch_name_valid() {
  local branch="$1"
  if [[ ! "$branch" =~ ^(feat|fix|docs|ci|refactor|wip|chore)/.+$ ]]; then
    fail "Invalid work branch prefix: ${branch}"
  fi
}

assert_no_direct_work_commits_on_main() {
  local bad
  bad="$(git log main --first-parent --pretty=%s | grep -E '^(feat|fix|docs|ci|refactor|chore):' || true)"
  [[ -z "$bad" ]] || fail "Detected direct work-style commits on main first-parent history"
}

assert_release_branch_contains_only_release_metadata_since_base() {
  local branch="$1"
  local base="$2"

  local changed
  changed="$(git diff --name-only "${base}..${branch}")"

  while IFS= read -r file; do
    [[ -z "$file" ]] && continue
    case "$file" in
      VERSION|CHANGELOG.md|release.json)
        ;;
      *)
        fail "Release branch '${branch}' changed disallowed file '${file}' after base '${base}'"
        ;;
    esac
  done <<< "$changed"
}

create_repo() {
  log "Creating temp repo"
  mkdir -p "${REPO_DIR}"
  cd "${REPO_DIR}"

  git init -b main
  git config user.name "Stability Flow Tester"
  git config user.email "tester@example.com"

  echo "0.0.0" > VERSION
  echo "# Demo" > README.md
  git add VERSION README.md
  git commit -m "chore: initial commit"

  git checkout -b develop
  pass "Initialized main and develop"
}

create_and_squash_merge_work_branch() {
  local work_branch="fix/race-on-authentication"
  assert_branch_name_valid "${work_branch}"

  log "Creating work branch ${work_branch} from develop"
  git checkout develop
  git checkout -b "${work_branch}"
  assert_branch_exists "${work_branch}"

  mkdir -p src
  cat > src/auth.txt <<'EOF'
authentication feature fixed
EOF
  git add src/auth.txt
  git commit -m "fix: race condition in authentication"

  mkdir -p tests
  cat > tests/auth.txt <<'EOF'
authentication feature fixed
EOF
  git add tests/auth.txt
  git commit -m "chore: improve test coverage for auth service"

  log "Squash merging ${work_branch} into develop"
  git checkout develop
  git merge --squash "${work_branch}"
  git commit -m "fix: race condition in authentication"

  pass "Work branch squash merged into develop"

  git branch -d "${work_branch}"
  assert_branch_not_exists "${work_branch}"
  pass "Deleted ${work_branch}"
}

planned_release() {
  local version="$1"
  local release_branch="release/${version}"

  log "Planned release ${version}"

  git checkout develop
  git checkout -b "${release_branch}"
  assert_branch_exists "${release_branch}"

  git rebase main

  local release_base
  release_base="$(git rev-parse HEAD)"

  echo "${version}" > VERSION
  cat > CHANGELOG.md <<EOF
## v${version}
- planned release
EOF
  git add VERSION CHANGELOG.md
  git commit -m "chore: prepare release ${version}"

  assert_release_branch_contains_only_release_metadata_since_base "${release_branch}" "${release_base}"

  git checkout main
  git merge --ff-only "${release_branch}"
  assert_equal_commit main "${release_branch}"
  pass "main fast-forwarded from ${release_branch}"

  git tag "v${version}"
  assert_tag_exists "v${version}"
  pass "Tagged v${version}"

  git checkout develop
  git merge --no-ff main -m "chore: merge main back into develop after v${version}"
  assert_ancestor main develop
  pass "Merged main back into develop"

  git branch -d "${release_branch}"
  assert_branch_not_exists "${release_branch}"
  pass "Deleted ${release_branch}"
}

hotfix_release() {
  local version="$1"
  local hotfix_branch="hotfix/${version}"
  local release_branch="release/${version}"

  log "Hotfix release ${version}"

  git checkout main
  git checkout -b "${hotfix_branch}"
  assert_branch_exists "${hotfix_branch}"

  mkdir -p src
  echo "urgent production fix" > src/hotfix.txt
  git add src/hotfix.txt
  git commit -m "fix: patch production issue"

  git checkout -b "${release_branch}"
  assert_branch_exists "${release_branch}"

  git rebase main

  local release_base
  release_base="$(git rev-parse HEAD)"

  echo "${version}" > VERSION
  {
    echo ""
    echo "## v${version}"
    echo "- emergency hotfix release"
  } >> CHANGELOG.md
  git add VERSION CHANGELOG.md
  git commit -m "chore: prepare hotfix release ${version}"

  assert_release_branch_contains_only_release_metadata_since_base "${release_branch}" "${release_base}"

  git checkout main
  git merge --ff-only "${release_branch}"
  assert_equal_commit main "${release_branch}"
  pass "main fast-forwarded from ${release_branch}"

  git tag "v${version}"
  assert_tag_exists "v${version}"
  pass "Tagged v${version}"

  git checkout develop
  git merge --no-ff main -m "chore: merge main back into develop after v${version}"
  assert_ancestor main develop
  pass "Merged main back into develop"

  git branch -d "${hotfix_branch}"
  git branch -d "${release_branch}"
  assert_branch_not_exists "${hotfix_branch}"
  assert_branch_not_exists "${release_branch}"
  pass "Deleted ${hotfix_branch} and ${release_branch}"
}

failed_release_disposal_test() {
  log "Testing failed release disposal"

  git checkout develop
  git checkout -b "release/9.9.9"
  assert_branch_exists "release/9.9.9"

  echo "9.9.9" > VERSION
  git add VERSION
  git commit -m "chore: prepare release 9.9.9"

  git checkout develop
  git branch -D "release/9.9.9"
  assert_branch_not_exists "release/9.9.9"
  pass "Failed release branch deleted instead of repaired"
}

summary() {
  log "Git graph"
  git log --graph --decorate --oneline --all
}

main() {
  create_repo
  create_and_squash_merge_work_branch
  planned_release "1.0.0"
  hotfix_release "1.0.1"
  failed_release_disposal_test
  assert_no_direct_work_commits_on_main
  pass "No direct work commits on main first-parent history"
  summary

  printf "\nAll Stability Flow tests passed.\n"
  printf "Temporary repo: %s\n" "${REPO_DIR}"
}

main "$@"
