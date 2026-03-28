# Stability Flow Conventions

## Version

- Version: v1
- Status: Stable

## 1. Overview

This document defines the **branch naming** and **commit message** conventions used alongside the Stability Flow specification.

These conventions support:

- clearer branch intent
- more consistent repository history
- easier automation and validation
- better release and audit workflows

This document complements the Stability Flow branching model, but does not redefine it.

For the branching specification itself, see:

- [Specification](spec.md)

---

## 2. Scope

This document defines conventions for:

- branch prefixes
- branch naming
- final squash commit formatting
- breaking change signaling
- revert categorization

It does **not** define:

- the branching model itself
- release workflow behavior
- enforcement mechanisms
- CI/CD implementation details

Those concerns are defined elsewhere.

---

## 3. Normative Language

The key words **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT**, **MAY**, **ONLY**, and **RECOMMENDED** in this document are to be interpreted as defined in the Stability Flow specification.

---

## 4. Branch Naming Conventions

### 4.1 Branch Prefixes

The following branch prefixes are valid.

| Prefix     | Purpose                                      |
| ---------- | -------------------------------------------- |
| `feat`     | New feature                                  |
| `fix`      | Bug fix                                      |
| `docs`     | Documentation                                |
| `ci`       | CI/CD changes                                |
| `refactor` | Code restructuring with no behavior change   |
| `chore`    | General maintenance or uncategorized work    |
| `wip`      | Temporary exploration only                   |
| `hotfix`   | Emergency production fix                     |
| `release`  | Production promotion                         |
| `sync`     | Reconciliation of `main` back into `develop` |

Branch prefixes **MUST** align with the branch roles defined in the Stability Flow specification.

### 4.2 Branch Name Format

Branch names **MUST**:

- be explicit enough to communicate intent
- be machine-checkable
- use a valid branch prefix
- use lowercase
- use slash-separated format

Recommended format:

```text
<prefix>/<description>
```

Examples:

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

### 4.3 Naming Guidance

Branch descriptions **SHOULD** be:

- concise
- descriptive
- stable enough to remain understandable in history

Branch names **SHOULD NOT**:

- rely on internal ticket numbers alone
- use vague descriptions such as `misc-fixes` or `stuff`
- encode implementation details that are likely to change rapidly

---

## 5. Final Squash Commit Conventions

Stability Flow focuses commit conventions on the **final squash commit** that enters long-lived history.

Branch-local history may be noisy during development, but the final squash commit **MUST** be intentional and machine-checkable.

### 5.1 Required Format

The final squash commit **MUST** use this format:

```text
<type>: <description>
```

Examples:

```text
feat: add authentication flow
fix: patch session race condition
docs: clarify hotfix reconciliation
ci: add branch validation workflow
refactor: simplify release branching rules
chore: update dependency maintenance tasks
```

### 5.2 Allowed Final Commit Types

The following final squash commit types are valid.

| Type       | Purpose                                     |
| ---------- | ------------------------------------------- |
| `feat`     | New user-facing or system-facing capability |
| `fix`      | Bug fix or behavior correction              |
| `docs`     | Documentation change                        |
| `ci`       | CI/CD workflow or automation change         |
| `refactor` | Code restructuring with no behavior change  |
| `chore`    | Maintenance or uncategorized change         |
| `test`     | Test-only work                              |
| `perf`     | Performance improvement                     |
| `build`    | Build system or packaging change            |
| `style`    | Formatting or style-only change             |

### 5.3 Type Mapping Guidance

Some commit types do not have a corresponding branch prefix.

When that happens, the branch **SHOULD** map to the closest valid branch category.

Recommended mapping:

| Commit Type              | Recommended Branch Prefix |
| ------------------------ | ------------------------- |
| `feat`                   | `feat`                    |
| `fix`                    | `fix`                     |
| `docs`                   | `docs`                    |
| `ci`                     | `ci`                      |
| `refactor`               | `refactor`                |
| `test`                   | `chore`                   |
| `perf`                   | `chore`                   |
| `build`                  | `chore`                   |
| `style`                  | `chore`                   |
| other uncategorized work | `chore`                   |

This mapping is intended to keep branch categories small and predictable while still allowing expressive final commit types.

---

## 6. Breaking Changes

Breaking changes **MUST** be indicated by `!` in the type header.

Required format:

```text
<type>!: <description>
```

Example:

```text
feat!: remove legacy authentication flow
```

A `BREAKING CHANGE:` footer **MAY** also be used, but does not replace the required `!`.

Example:

```text
feat!: remove legacy authentication flow

BREAKING CHANGE: legacy authentication endpoints were removed
```

---

## 7. Reverts

Reverts **MUST NOT** use `revert:` as the final squash commit categorization.

The final squash commit **MUST** be categorized by the **net effect** of the merged change.

### 7.1 Revert Categorization

Recommended categorization:

| Scenario                  | Recommended Type |
| ------------------------- | ---------------- |
| Restores correct behavior | `fix`            |
| Reverts documentation     | `docs`           |
| Reverts CI/CD changes     | `ci`             |
| Reverts refactor work     | `refactor`       |
| Everything else           | `chore`          |

### 7.2 Rationale

A revert is an action, not a long-lived category.

The important thing in repository history is the effect of the final integrated change.

Examples:

```text
fix: restore correct login redirect behavior
docs: revert outdated release notes guidance
chore: revert experimental dependency update
```

---

## 8. Work-in-Progress Commit Guidance

Temporary `wip:` commits **MAY** appear in branch-local history on valid work branches.

Examples:

```text
wip: refactor payment module
wip: debug authentication issue
wip: sketch release validation changes
```

This is acceptable during local or branch-local development.

However:

- `wip:` commits **MUST NOT** appear as the final squash commit
- the final squash commit **MUST** use a valid final type
- incomplete work **MUST NOT** be merged into long-lived history

---

## 9. Optional Versioning Guidance

Teams **MAY** use commit types as versioning signals.

A common mapping is:

| Pattern | Typical Version Impact |
| ------- | ---------------------- |
| `!`     | major                  |
| `feat`  | minor                  |
| `fix`   | patch                  |
| others  | none                   |

This mapping is optional and is not required by Stability Flow.

---

## 10. Summary

At a practical level, Stability Flow conventions are designed to keep branch and commit history:

- explicit
- consistent
- machine-checkable
- easier to review
- easier to automate

The most important expectations are:

- branch names use valid prefixes
- branch names clearly communicate intent
- final squash commits use valid final types
- breaking changes use `!`
- final squash commits categorize the net effect of the change

These conventions support the Stability Flow specification without redefining it.
