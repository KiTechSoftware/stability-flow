# CLI Validator

## 1. Overview

The Stability Flow CLI Validator is a reference tool for validating whether a repository workflow follows the Stability Flow specification.

It is intended to help teams validate:

- branch names
- branch origins
- merge eligibility
- commit messages

This tool is a **reference implementation** of parts of the Stability Flow specification.

It is **not** the specification itself.

For the normative rules, see the specification.

---

## 2. What It Validates

The CLI validator is designed to validate high-value, machine-checkable workflow rules.

These include:

- branch naming conventions
- allowed branch origins
- allowed merge targets
- commit message structure

This makes it useful for:

- local development checks
- pull request validation
- CI enforcement
- reusable workflow integration
- custom policy automation

---

## 3. Scope

The CLI validator is intentionally focused on validation only.

It does **not**:

- create branches
- modify Git history
- perform merges
- rewrite commits
- automate releases

Its purpose is to tell you whether a proposed action is valid according to the Stability Flow rules it implements.

---

## 4. Current Commands

The CLI currently supports the following validation commands.

### Branch Name Validation

Validates whether a branch name matches the expected Stability Flow naming model.

```bash id="uw1a3q"
stability-flow-validator validate-branch-name --branch feat/add-authentication
````

---

### Branch Origin Validation

Validates whether a branch was created from an allowed base branch.

```bash id="whb5h8"
stability-flow-validator validate-origin --branch hotfix/1.2.4 --base main
```

---

### Merge Validation

Validates whether a source branch is allowed to merge into a target branch.

```bash id="5rtuzl"
stability-flow-validator validate-merge --source release/1.2.3 --target main
```

---

### Commit Validation

Validates whether a commit message matches the configured commit rules.

```bash id="r52bdb"
stability-flow-validator validate-commit --mode squash --message "feat: complete validator v1"
```

---

## 5. Branch Name Validation

### Purpose

Branch name validation checks whether a branch follows the expected Stability Flow branch naming conventions.

### Example valid names

```text id="5vp4dy"
feat/add-authentication
fix/race-on-authentication
docs/update-release-policy
ci/add-validator-check
refactor/simplify-branch-rules
chore/update-dependencies
hotfix/1.2.4
release/1.3.0
sync/main-into-develop
```

### Example

```bash id="if8xyj"
stability-flow-validator validate-branch-name --branch feat/add-authentication
```

### Example result

```text id="m5h3f2"
PASS: branch name allowed: feat/add-authentication
reason: valid branch type: feat
```

---

## 6. Branch Origin Validation

### Purpose

Branch origin validation checks whether a branch was created from an allowed base branch.

Examples:

* regular work branches should come from `develop`
* `hotfix/*` should come from `main`
* `release/*` should come from `develop` or `hotfix/*`

### Example

```bash id="sagqxh"
stability-flow-validator validate-origin --branch hotfix/1.2.4 --base main
```

### Example result

```text id="st0y9n"
PASS: branch origin allowed: hotfix/1.2.4 from main
reason: hotfix/* must be created from main
```

---

## 7. Merge Validation

### Purpose

Merge validation checks whether a source branch is allowed to merge into a given target branch.

Examples:

* `feat/*` may merge into `develop`
* `release/*` may merge into `main`
* `feat/*` may not merge into `main`

### Example valid merge

```bash id="2xv47u"
stability-flow-validator validate-merge --source release/1.2.3 --target main
```

### Example result

```text id="87h9rf"
PASS: merge allowed: release/1.2.3 -> main
reason: only release/* may merge into main, using fast-forward only
```

### Example invalid merge

```bash id="bq9r1k"
stability-flow-validator validate-merge --source feat/add-authentication --target main
```

### Example result

```text id="ajblw9"
FAIL: merge not allowed: feat/add-authentication -> main
reason: merge not allowed by Stability Flow: feat -> main
```

---

## 8. Commit Validation

### Purpose

Commit validation checks whether a commit message matches the expected commit format.

This is especially useful for:

* squash merge commit messages
* release preparation commits
* CI validation

### Example valid commit

```bash id="b7o6mo"
stability-flow-validator validate-commit --mode squash --message "feat: complete validator v1"
```

### Example result

```text id="hvtg8y"
PASS: commit message allowed (squash): feat: complete validator v1
reason: valid commit type: feat
```

### Example breaking change commit

```bash id="w3w0q0"
stability-flow-validator validate-commit --mode squash --message "feat!: remove legacy auth flow

BREAKING CHANGE: legacy auth removed"
```

---

## 9. Commit Modes

The validator supports commit validation modes to reflect different workflow contexts.

### Example modes

* squash
* branch
* release

The exact supported modes may evolve as the validator grows.

The purpose of modes is to allow different commit expectations depending on the workflow event being validated.

Examples:

* a squash merge commit
* a release preparation commit
* a branch-local commit

---

## 10. Output Formats

The validator supports multiple output formats so it can be used by both humans and automation.

### Supported formats

* `text`
* `json`
* `jsonl`
* `markdown`

### Example

```bash id="zvf4e3"
stability-flow-validator validate-merge --source release/1.2.3 --target main --format json
```

### Example JSON output

```json id="jsnnfk"
{
  "ok": true,
  "command": "validate-merge",
  "reason": "only release/* may merge into main, using fast-forward only",
  "fields": {
    "source": "release/1.2.3",
    "target": "main"
  }
}
```

### Example JSONL output

```bash id="nwlxmr"
stability-flow-validator validate-origin --branch hotfix/1.2.4 --base main --format jsonl
```

```json id="flr77a"
{"ok":true,"command":"validate-origin","reason":"hotfix/* must be created from main","fields":{"base":"main","branch":"hotfix/1.2.4"}}
```

### Example Markdown output

```bash id="rj9wyy"
stability-flow-validator validate-merge --source release/1.2.3 --target main --format markdown
```

```md id="74xyxq"
## Stability Flow Validation Result

- **Command:** `validate-merge`
- **Status:** ✅ Passed
- **Source:** `release/1.2.3`
- **Target:** `main`
- **Reason:** only release/* may merge into main, using fast-forward only
```

---

## 11. Exit Codes

The validator is designed to return a non-zero exit code when validation fails.

This makes it suitable for:

* shell scripting
* CI pipelines
* pre-push hooks
* pull request validation

### General behavior

* success → exit code `0`
* validation failure → non-zero exit code

---

## 12. Local Usage

The validator can be used locally as a lightweight workflow safety check.

Examples:

* validate a branch before opening a pull request
* validate a squash commit before merging
* validate a hotfix branch origin before release work begins

### Example

```bash id="4gw3w2"
stability-flow-validator validate-origin --branch release/1.2.4 --base hotfix/1.2.4
```

This is useful when teams want fast feedback before CI.

---

## 13. CI and Workflow Usage

The validator is also intended to work well in automation.

Common uses include:

* validating branch names in pull requests
* validating merge direction rules
* validating commit messages
* producing structured output for job summaries or machine parsing

The CLI is especially useful in CI because it supports machine-readable output formats.

Examples:

* `json`
* `jsonl`
* `markdown`

This allows teams to integrate it into:

* GitHub Actions
* reusable workflows
* custom CI systems
* release policy checks

---

## 14. Reference Tooling Role

The CLI validator is one implementation of Stability Flow enforcement.

It exists to help teams adopt and validate the specification more easily.

Teams may:

* use it directly
* wrap it in their own scripts
* integrate it into CI
* adapt its rules
* build their own validator instead

The Stability Flow specification does not require this tool.

---

## 15. Summary

The CLI validator provides a practical way to validate parts of the Stability Flow specification.

It is useful when teams want machine-checkable validation for:

* branch names
* branch origins
* merge eligibility
* commit messages

It is designed to support both:

* human workflows
* automation workflows
