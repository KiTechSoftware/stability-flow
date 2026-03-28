# CLI Validator

## 1. Overview

The Stability Flow CLI Validator is a **reference implementation** for validating whether a repository workflow follows the Stability Flow specification.

It is intended to help teams validate high-value, machine-checkable rules such as:

- branch naming
- branch origin
- merge eligibility
- final squash commit messages

This tool is an **implementation of the Stability Flow specification**.

It is **not** the specification itself.

For the normative rules, see the Stability Flow specification.

---

## 2. What It Validates

The CLI validator is designed to validate the parts of Stability Flow that are most valuable to enforce automatically.

These include:

- branch naming conventions
- allowed branch origins
- allowed merge targets
- final squash commit message structure
- breaking change notation in final squash commits

This makes it useful for:

- local workflow checks
- pull request validation
- CI enforcement
- reusable workflow integration
- custom policy automation

---

## 3. Scope

The CLI validator is intentionally focused on **validation only**.

It does **not**:

- create branches
- modify Git history
- perform merges
- rewrite commits
- automate releases
- decide whether a change is “worthy” of a release or hotfix

Its purpose is to tell you whether a proposed branch, merge, or final squash commit is valid according to the Stability Flow rules it implements.

---

## 4. Current Commands

The CLI currently supports the following validation commands.

### 4.1. Branch Name Validation

Validates whether a branch name matches the expected Stability Flow naming model.

```bash
stability-flow-validator validate-branch-name --branch feat/add-authentication
```

---

### 4.2. Branch Origin Validation

Validates whether a branch was created from an allowed base branch.

```bash
stability-flow-validator validate-origin --branch hotfix/1.2.4 --base main
```

---

### 4.3. Merge Validation

Validates whether a source branch is allowed to merge into a target branch.

```bash
stability-flow-validator validate-merge --source release/1.2.4 --target main
```

---

### 4.4. Commit Validation

Validates whether a commit message matches the Stability Flow commit rules.

```bash
stability-flow-validator validate-commit --mode squash --message "feat: complete validator v1"
```

---

## 5. Branch Name Validation

### 5.1. Purpose

Branch name validation checks whether a branch follows the expected Stability Flow branch naming conventions.

### 5.2. Example valid names

```text
feat/add-authentication
fix/race-on-authentication
docs/update-release-policy
ci/add-validator-check
refactor/simplify-branch-rules
chore/update-dependencies
wip/auth-investigation
hotfix/1.2.4
release/1.3.0
sync/main-into-develop-1.2.4
```

### 5.3. Example

```bash
stability-flow-validator validate-branch-name --branch feat/add-authentication
```

### 5.4. Example result

```text
PASS: branch name allowed: feat/add-authentication
reason: valid branch type: feat
```

---

## 6. Branch Origin Validation

### 6.1. Purpose

Branch origin validation checks whether a branch was created from an allowed base branch.

### 6.2. Expected origin behavior

- regular work branches must come from `develop`
- `wip/*` may come from `develop`
- `wip/*` may come from `main` for hotfix troubleshooting
- `wip/*` is never mergeable
- `hotfix/*` must come from `main`
- `release/*` must come from `develop` or `hotfix/*`
- `sync/*` must come from `develop`

### 6.3. Example valid origin

```bash
stability-flow-validator validate-origin --branch hotfix/1.2.4 --base main
```

### 6.4. Example result

```text
PASS: branch origin allowed: hotfix/1.2.4 from main
reason: hotfix/* must be created from main
```

### 6.5. Example invalid origin

```bash
stability-flow-validator validate-origin --branch feat/add-authentication --base main
```

### 6.6. Example result

```text
FAIL: branch origin not allowed: feat/add-authentication from main
reason: regular work branches must be created from develop
```

### 6.7. Important note on `wip/*`

`wip/*` branches may be created from `develop` or `main` for temporary exploration, but they are not integration branches and are not mergeable under Stability Flow.
`wip/*` branches derived from `main` MUST be used ONLY for hotfix troubleshooting/exploration

---

## 7. Merge Validation

### 7.1. Purpose

Merge validation checks whether a source branch is allowed to merge into a given target branch.

### 7.2. Expected merge behavior

- regular work branches may merge into `develop`
- `release/*` MUST merge into `main`
- `sync/*` MUST merge into `develop`
- `hotfix/*` MUST NOT merge directly into `main`
- `wip/*` MUST NOT merge into any branch
- direct `main` → `develop` merges are not allowed

### 7.3. Example valid merge

```bash
stability-flow-validator validate-merge --source release/1.2.4 --target main
```

### 7.4. Example result

```text
PASS: merge allowed: release/1.2.4 -> main
reason: only release/* may merge into main
```

### 7.5. Example invalid merge

```bash
stability-flow-validator validate-merge --source feat/add-authentication --target main
```

### 7.6. Example result

```text
FAIL: merge not allowed: feat/add-authentication -> main
reason: merge not allowed by Stability Flow: regular work branches must not merge into main
```

### 7.7. Example invalid exploratory merge

```bash
stability-flow-validator validate-merge --source wip/auth-investigation --target develop
```

### 7.8. Example result

```text
FAIL: merge not allowed: wip/auth-investigation -> develop
reason: wip/* branches are exploratory only and must never be merged
```

---

## 8. Commit Validation

### 8.1. Purpose

Commit validation checks whether a commit message matches the Stability Flow rules for **final squash commits**.

This is especially useful for:

- pull request squash merges
- protected branch merge validation
- CI enforcement
- release review workflows

The validator is primarily concerned with the **final squash commit that enters `develop`**.

It is not intended to police every intermediate commit in branch history.

---

### 8.2. Allowed Final Squash Commit Types

The validator accepts the following commit types for final squash commits:

- `feat`
- `fix`
- `docs`
- `ci`
- `refactor`
- `chore`
- `test`
- `perf`
- `build`
- `style`

#### 8.2.1 Example valid commits

```text
feat: add authentication flow
fix: patch session race condition
docs: clarify hotfix reconciliation
chore: prepare release metadata
test: add branch origin coverage
perf: improve validator startup
```

---

### 8.3 Breaking Changes

Breaking changes must be indicated by:

- `!` in the type header

A `BREAKING CHANGE:` footer may also be used, but does not replace the required `!`.

#### 8.3.1 Example valid breaking change commit

```text
feat!: remove legacy auth flow

BREAKING CHANGE: legacy auth flow removed
```

---

### 8.4. Reverts

The validator treats `revert:` differently depending on context.

#### 8.4.1 Allowed

- `revert:` may appear in branch-local history

#### 8.4.2 Not allowed

- `revert:` must not be used as the **final squash commit** entering `develop`

#### 8.4.3 Example invalid final squash commit

```text
revert: undo previous change
```

---

## 9. Commit Validation Modes

The validator may support different commit validation modes depending on workflow context.

### 9.1 Common modes

- `squash` — validates the final squash commit that will enter `develop`
- `work` — validates ordinary work-branch commits, if enabled by implementation policy

### 9.2 Important note

The Stability Flow specification defines rules for **final squash commits**.

Additional validation of ordinary branch-local commits is an implementation choice.

That means:

- `squash` mode reflects spec-level validation
- other modes may exist for convenience or stricter workflow policy

---

## 10. Output Formats

The validator supports multiple output formats so it can be used by both humans and automation.

### 10.1. Supported formats

- `text`
- `json`
- `jsonl`
- `markdown`

#### 10.1.1 Example JSON output

```bash
stability-flow-validator validate-merge --source release/1.2.4 --target main --format json
```

```json
{
  "ok": true,
  "command": "validate-merge",
  "reason": "only release/* may merge into main",
  "fields": {
    "source": "release/1.2.4",
    "target": "main"
  }
}
```

#### 10.1.2 Example JSONL output

```bash
stability-flow-validator validate-origin --branch hotfix/1.2.4 --base main --format jsonl
```

```json
{"ok":true,"command":"validate-origin","reason":"hotfix/* must be created from main","fields":{"base":"main","branch":"hotfix/1.2.4"}}
```

#### 10.1.3 Example Markdown output

```bash
stability-flow-validator validate-merge --source release/1.2.4 --target main --format markdown
```

```md
## Stability Flow Validation Result

- **Command:** `validate-merge`
- **Status:** ✅ Passed
- **Source:** `release/1.2.4`
- **Target:** `main`
- **Reason:** only release/* may merge into main
```

---

## 11. Exit Codes

The validator is designed to return a non-zero exit code when validation fails.

This makes it suitable for:

- shell scripting
- CI pipelines
- pre-push hooks
- pull request validation

### 11.1 General behavior

- success → exit code `0`
- validation failure → non-zero exit code

---

## 12. Local Usage

The validator can be used locally as a lightweight workflow safety check.

Examples:

- validate a branch before opening a pull request
- validate a squash commit before merging
- validate a hotfix branch origin before release work begins

### 12.1 Example

```bash
stability-flow-validator validate-origin --branch release/1.2.4 --base hotfix/1.2.4
```

This is useful when teams want fast feedback before CI.

---

## 13. CI and Workflow Usage

The validator is also intended to work well in automation.

Common uses include:

- validating branch names in pull requests
- validating merge direction rules
- validating final squash commit messages
- producing structured output for job summaries or machine parsing

The CLI is especially useful in CI because it supports machine-readable output formats.

Examples:

- `json`
- `jsonl`
- `markdown`

This allows teams to integrate it into:

- CI pipelines
- reusable workflows
- repository policy checks
- merge automation

---

## 14. Relationship to the Specification

The CLI validator is a **reference implementation** of the Stability Flow specification.

That means:

- it should follow the spec contract
- it should not redefine the spec
- it should not become the only way to adopt Stability Flow

Teams may:

- use it directly
- wrap it in their own scripts
- integrate it into CI
- extend it with internal policy
- build their own validator instead

The Stability Flow specification does **not** require this tool.

---

## 15. Summary

The CLI validator provides a practical way to validate the machine-checkable parts of Stability Flow.

It is especially useful for validating:

- branch names
- branch origins
- merge eligibility
- final squash commit rules

It is designed to support both:

- human workflows
- automation workflows
