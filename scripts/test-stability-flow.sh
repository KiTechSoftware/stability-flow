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

assert_work_branch_name_valid() {
  local branch="$1"
  if [[ ! "$branch" =~ ^(feat|fix|docs|ci|refactor|wip|chore)/.+$ ]]; then
    fail "Invalid work branch prefix: ${branch}"
  fi
}

assert_release_branch_name_valid() {
  local branch="$1"
  if [[ ! "$branch" =~ ^release/[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    fail "Invalid release branch name: ${branch}"
  fi
}

assert_hotfix_branch_name_valid() {
  local branch="$1"
  if [[ ! "$branch" =~ ^hotfix/[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    fail "Invalid hotfix branch name: ${branch}"
  fi
}

assert_sync_branch_name_valid() {
  local branch="$1"
  if [[ ! "$branch" =~ ^sync/main-into-develop-[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    fail "Invalid sync branch name: ${branch}"
  fi
}

assert_develop_contains_main_before_planned_release() {
  git merge-base --is-ancestor main develop \
    || fail "Cannot begin planned release: develop does not contain the latest reconciled main"
}

develop_contains_main() {
  git merge-base --is-ancestor main develop
}

assert_demo_release_branch_contains_only_allowed_release_files_since_base() {
  local branch="$1"
  local base="$2"

  # Demo-repo policy:
  # Stability Flow requires release-safe preparation changes only,
  # but the exact allowlist is repository-specific.
  local changed
  changed="$(git diff --name-only "${base}..${branch}")"

  while IFS= read -r file; do
    [[ -z "$file" ]] && continue
    case "$file" in
      VERSION|CHANGELOG.md|release.json)
        ;;
      *)
        fail "Release branch '${branch}' changed disallowed demo file '${file}' after base '${base}'"
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
  git commit -m "init: repository bootstrap"

  git checkout -b develop
  pass "Initialized main and develop"
}

create_and_squash_merge_work_branch() {
  local work_branch="fix/race-on-authentication"
  assert_work_branch_name_valid "${work_branch}"

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

  git branch -D "${work_branch}"
  assert_branch_not_exists "${work_branch}"
  pass "Deleted ${work_branch}"
}

reconcile_main_into_develop() {
  local version="$1"
  local sync_branch="sync/main-into-develop-${version}"

  assert_sync_branch_name_valid "${sync_branch}"

  log "Reconciling main back into develop through ${sync_branch}"

  git checkout develop
  git checkout -b "${sync_branch}"
  assert_branch_exists "${sync_branch}"

  git merge --no-ff main -m "chore: reconcile main into develop after v${version}"
  assert_ancestor main "${sync_branch}"

  git checkout develop
  git merge --no-ff "${sync_branch}" -m "chore: merge ${sync_branch} into develop"

  assert_ancestor main develop
  pass "Reconciled main back into develop through ${sync_branch}"

  git branch -D "${sync_branch}"
  assert_branch_not_exists "${sync_branch}"
  pass "Deleted ${sync_branch}"
}

planned_release() {
  local version="$1"
  local release_branch="release/${version}"
  local develop_tip_before_release

  assert_release_branch_name_valid "${release_branch}"

  log "Planned release ${version}"

  git checkout develop
  assert_develop_contains_main_before_planned_release

  develop_tip_before_release="$(git rev-parse develop)"
  git checkout -b "${release_branch}"
  assert_branch_exists "${release_branch}"

  [[ "$(git rev-parse HEAD)" == "${develop_tip_before_release}" ]] \
    || fail "${release_branch} was not created from develop"

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

  assert_demo_release_branch_contains_only_allowed_release_files_since_base \
    "${release_branch}" "${release_base}"

  git checkout main
  git merge --ff-only "${release_branch}"
  assert_equal_commit main "${release_branch}"
  pass "main fast-forwarded from ${release_branch}"

  git tag "v${version}"
  assert_tag_exists "v${version}"
  pass "Tagged v${version}"

  reconcile_main_into_develop "${version}"

  git branch -D "${release_branch}"
  assert_branch_not_exists "${release_branch}"
  pass "Deleted ${release_branch}"
}

hotfix_release() {
  local version="$1"
  local hotfix_branch="hotfix/${version}"
  local release_branch="release/${version}"
  local main_tip_before_hotfix
  local hotfix_tip_before_release

  assert_hotfix_branch_name_valid "${hotfix_branch}"
  assert_release_branch_name_valid "${release_branch}"

  log "Hotfix release ${version}"

  git checkout main
  main_tip_before_hotfix="$(git rev-parse main)"
  git checkout -b "${hotfix_branch}"
  assert_branch_exists "${hotfix_branch}"

  [[ "$(git rev-parse HEAD)" == "${main_tip_before_hotfix}" ]] \
    || fail "${hotfix_branch} was not created from main"

  mkdir -p src
  echo "urgent production fix" > src/hotfix.txt
  git add src/hotfix.txt
  git commit -m "fix: patch production issue"

  hotfix_tip_before_release="$(git rev-parse "${hotfix_branch}")"
  git checkout -b "${release_branch}"
  assert_branch_exists "${release_branch}"

  [[ "$(git rev-parse HEAD)" == "${hotfix_tip_before_release}" ]] \
    || fail "${release_branch} was not created from ${hotfix_branch}"

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

  assert_demo_release_branch_contains_only_allowed_release_files_since_base \
    "${release_branch}" "${release_base}"

  git checkout main
  git merge --ff-only "${release_branch}"
  assert_equal_commit main "${release_branch}"
  pass "main fast-forwarded from ${release_branch}"

  git tag "v${version}"
  assert_tag_exists "v${version}"
  pass "Tagged v${version}"

  reconcile_main_into_develop "${version}"

  git branch -D "${hotfix_branch}"
  git branch -D "${release_branch}"
  assert_branch_not_exists "${hotfix_branch}"
  assert_branch_not_exists "${release_branch}"
  pass "Deleted ${hotfix_branch} and ${release_branch}"
}

stale_develop_release_block_test() {
  local version="2.0.1"
  local hotfix_branch="hotfix/${version}"
  local release_branch="release/${version}"

  assert_hotfix_branch_name_valid "${hotfix_branch}"
  assert_release_branch_name_valid "${release_branch}"

  log "Testing planned release is blocked when develop is stale"

  git checkout main
  git checkout -b "${hotfix_branch}"
  echo "critical prod fix" > src/stale-check-hotfix.txt
  git add src/stale-check-hotfix.txt
  git commit -m "fix: critical production issue"

  git checkout -b "${release_branch}"
  git rebase main

  echo "${version}" > VERSION
  git add VERSION
  git commit -m "chore: prepare hotfix release ${version}"

  git checkout main
  git merge --ff-only "${release_branch}"
  assert_equal_commit main "${release_branch}"

  git checkout develop
  if develop_contains_main; then
    fail "Expected develop to be stale before reconciliation, but it already contains main"
  fi

  pass "Planned release correctly blocked when develop is stale"

  git checkout main
  git branch -D "${hotfix_branch}"
  git branch -D "${release_branch}"
  pass "Planned release correctly blocked when develop is stale"

  reconcile_main_into_develop "${version}"

  git tag "v${version}"
  assert_tag_exists "v${version}"
  pass "Tagged v${version}"
}

failed_release_disposal_test() {
  local release_branch="release/9.9.9"

  assert_release_branch_name_valid "${release_branch}"

  log "Testing failed release disposal and recreation from correct source"

  git checkout develop
  git checkout -b "${release_branch}"
  assert_branch_exists "${release_branch}"

  echo "9.9.9" > VERSION
  git add VERSION
  git commit -m "chore: prepare release 9.9.9"

  git checkout develop
  mkdir -p src
  echo "release source fix" > src/release-source-fix.txt
  git add src/release-source-fix.txt
  git commit -m "fix: correct source branch before recreating release"

  git branch -D "${release_branch}"
  assert_branch_not_exists "${release_branch}"

  local develop_tip_before_recreate
  develop_tip_before_recreate="$(git rev-parse develop)"

  git checkout -b "${release_branch}"
  assert_branch_exists "${release_branch}"

  [[ "$(git rev-parse HEAD)" == "${develop_tip_before_recreate}" ]] \
    || fail "${release_branch} was not recreated from the correct source branch"

  git checkout develop
  git branch -D "${release_branch}"
  assert_branch_not_exists "${release_branch}"
  pass "Failed release branch deleted and recreated from correct source branch"
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
  stale_develop_release_block_test
  failed_release_disposal_test
  summary

  printf "\nAll Stability Flow tests passed.\n"
  printf "Temporary repo: %s\n" "${REPO_DIR}"
}

main "$@"