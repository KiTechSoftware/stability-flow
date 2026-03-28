# Stability Flow Specification

## Version

- Version: v1
- Status: Stable

## 1. Overview

**Stability Flow** is a branching strategy specification for teams that require:

- a **stable production branch**
- **planned releases**
- **safe hotfixes**
- **explicit reconciliation after production divergence**

It defines a structured branching model intended to keep production promotion explicit, preserve a stable production line, and make release and reconciliation behavior machine-checkable.

---

## 2. Goals

Stability Flow is designed to achieve the following goals:

1. **Protect production stability**
2. **Support planned releases**
3. **Allow urgent production hotfixes**
4. **Make branch intent explicit**
5. **Keep production promotion auditable**
6. **Make enforcement practical through policy and automation**

---

## 3. Normative Language

The key words **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT**, **MAY**, **ONLY**, and **RECOMMENDED** in this specification are to be interpreted as follows:

| Term            | Meaning                                                                                                                   |
| --------------- | ------------------------------------------------------------------------------------------------------------------------- |
| **MUST**        | Mandatory. No exceptions unless explicitly allowed by the flow.                                                           |
| **MUST NOT**    | Prohibited. This action is not allowed.                                                                                   |
| **SHOULD**      | Strong recommendation. Expected in normal operation, but may be bypassed only with deliberate justification.              |
| **SHOULD NOT**  | Strong recommendation against. Acceptable only with deliberate justification.                                             |
| **MAY**         | Optional. Allowed, but not required.                                                                                      |
| **ONLY**        | Exclusive constraint. No other path, branch source, or merge target is permitted.                                         |
| **RECOMMENDED** | Preferred approach when multiple valid options exist.                                                                     |

### Interpretation Rule

If a rule uses **MUST**, **MUST NOT**, or **ONLY**, it is part of the enforceable Stability Flow contract.

If a rule uses **SHOULD**, **SHOULD NOT**, or **RECOMMENDED**, it is operational guidance unless elevated elsewhere.

If a rule uses **MAY**, it is optional and **MUST NOT** be interpreted as required.

---

## 4. Branch Roles

### 4.1 `main`

Production branch.

#### Main Branch Requirements

- `main` **MUST** remain stable
- `main` **MUST NOT** receive direct regular work
- `main` **MUST ONLY** receive changes from `release/*`
- `main` **SHOULD** remain linear where practical

---

### 4.2 `develop`

Integration branch for planned work.

#### Develop Branch Requirements

- regular work branches **MUST** branch from `develop`
- regular work branches **MUST** merge into `develop`
- `develop` **MAY** be ahead of `main`
- `develop` **MUST** contain the latest reconciled production state before a planned release begins

---

### 4.3 Regular Work Branches

Regular work branches include:

- `feat/*`
- `fix/*`
- `docs/*`
- `ci/*`
- `refactor/*`
- `chore/*`

#### Regular Work Branch Requirements

- regular work branches **MUST** branch from `develop`
- regular work branches **MUST NOT** target `main`
- regular work branches **MUST NOT** target `release/*`
- regular work branches **MUST** be squash merged into `develop`
- regular work branches **SHOULD** be short-lived
- regular work branches **SHOULD** be deleted after merge

---

### 4.4 `release/*`

Promotion branch.

#### Release Branch Requirements

- `release/*` **MUST** be created from:
  - `develop`, or
  - `hotfix/*`
- `release/*` **MUST ONLY** merge into `main`
- `release/*` **MUST** be the only branch type that merges into `main`
- `release/*` **MUST NOT** receive feature or bug-fix commits
- `release/*` **MUST NOT** be used for ongoing development
- `release/*` **MUST** be rebased onto the latest `main` before promotion
- `release/*` **MUST ONLY** contain release-safe preparation changes such as:
  - version bumps
  - changelog updates
  - release metadata
- `release/*` **SHOULD** be deleted after merge or disposal

---

### 4.5 `hotfix/*`

Emergency production repair branch.

#### Hotfix Branch Requirements

- `hotfix/*` **MUST** branch from `main`
- `hotfix/*` **MUST NOT** branch from `develop`
- `hotfix/*` **MUST NOT** merge directly into `main`
- `hotfix/*` **MUST** be promoted through `release/*`
- `hotfix/*` **SHOULD** be short-lived
- `hotfix/*` **SHOULD** be deleted after release

---

### 4.6 `sync/*`

Reconciliation branch.

#### Sync Branch Requirements

- `sync/*` **MUST** branch from `develop`
- `sync/*` **MUST** be used to reconcile `main` back into `develop`
- direct `main` → `develop` merges **MUST NOT** be used
- `sync/*` **MUST NOT** contain unrelated product work
- `sync/*` **SHOULD** be deleted after merge

---

### 4.7 `wip/*`

Exploratory branch.

#### WIP Branch Requirements

- `wip/*` **MAY** be used for temporary exploration
- `wip/*` **MUST NOT** be merged
- `wip/*` **MUST NOT** be promoted into `develop`, `main`, `release/*`, `hotfix/*`, or `sync/*`
- accepted work discovered in `wip/*` **MUST** be recreated through the correct branch type before integration

---

## 5. Core Rules

1. **Nothing goes directly to `main`**
2. **ONLY `release/*` may merge into `main`**
3. **All releases MUST go through `release/*`**
4. **All regular work MUST branch from `develop`**
5. **All approved regular work MUST be squash merged into `develop`**
6. **Hotfixes MUST originate from `main`**
7. **`main` is the single source of production truth**
8. **Production reconciliation back into `develop` MUST happen through `sync/*`**
9. **No planned release may begin unless `develop` already contains the latest reconciled `main`**
10. **If reconciliation conflicts are unresolved, planned releases are blocked until reconciliation is complete**
11. **`release/*` branches MUST be treated as disposable promotion artifacts**
12. **If a release candidate becomes invalid, the `release/*` branch SHOULD be discarded and recreated from the correct source branch**
13. **`wip/*` branches are exploratory only and MUST NEVER be promoted into the delivery flow**

---

## 6. Required Branch Origins

The following branch origins are required.

| Branch Type                            | MUST Branch From        |
| -------------------------------------- | ----------------------- |
| regular work (`feat/*`, `fix/*`, etc.) | `develop`               |
| `release/*`                            | `develop` or `hotfix/*` |
| `hotfix/*`                             | `main`                  |
| `sync/*`                               | `develop`               |

Additional rules:

- `wip/*` **MAY** branch from `develop`
- `wip/*` **MAY** branch from `main`
- no delivery branch **MAY** branch from `wip/*`

---

## 7. Required Merge Targets

The following merge targets are required.

| Source Branch                          | Required Target |
| -------------------------------------- | --------------- |
| regular work (`feat/*`, `fix/*`, etc.) | `develop`       |
| `release/*`                            | `main`          |
| `sync/*`                               | `develop`       |

Additional rules:

- `wip/*` **MUST NOT** merge into any branch
- `hotfix/*` **MUST NOT** merge directly into `main`
- direct `main` → `develop` merges **MUST NOT** be used

---

## 8. Merge Strategies

The following merge strategies are required.

| Flow                     | Strategy                 |
| ------------------------ | ------------------------ |
| regular work → `develop` | **squash merge**         |
| `release/*` → `main`     | **fast-forward only**    |
| `main` → `sync/*`        | **regular merge commit** |
| `sync/*` → `develop`     | **regular merge commit** |

### Additional Rule

Before promotion into `main`, `release/*` **MUST** be rebased onto the latest `main`.

---

## 9. Release Lifecycle

### 9.1 Planned Release

A planned release follows this shape:

```text
develop → release/X.Y.Z → main → sync/main-into-develop-X.Y.Z → develop
```

#### Planned Release Requirements

- planned releases **MUST** begin from `develop`
- planned releases **MUST NOT** begin unless `develop` already contains the latest reconciled `main`
- `release/*` **MUST** be rebased onto the latest `main` before promotion
- `release/*` **MUST** be fast-forward promoted into `main`
- after promotion, `main` **MUST** be reconciled back into `develop` through `sync/*`

---

### 9.2 Hotfix Release

A hotfix release follows this shape:

```text
main → hotfix/X.Y.Z → release/X.Y.Z → main → sync/main-into-develop-X.Y.Z → develop
```

#### Hotfix Release Requirements

- hotfix releases **MUST** begin from `main`
- hotfixes **MUST** be isolated from unreleased work
- hotfixes **MUST** still promote through `release/*`
- after promotion, `main` **MUST** be reconciled back into `develop` through `sync/*`

---

## 10. Conventions

Branch naming and commit conventions are defined separately in:

- [Conventions](conventions.md)

Conventions support Stability Flow, but are not required to understand the core branching model.

---

## 11. Compliance

A repository is considered compliant with Stability Flow if it satisfies:

- branch origin rules
- merge target rules
- required merge strategies
- promotion constraints (`release/*` → `main`)
- reconciliation constraints (`main` → `develop` through `sync/*`)
- squash merge requirements for regular work

Additional naming and commit conventions **MAY** be adopted separately.

Compliance **MAY** be validated through tooling or enforced through policy.

---

## 12. Non-Goals

Stability Flow does **not** define:

- how teams must version software
- how release approval must work
- how deployments must be performed
- how CI/CD must be implemented
- a required tooling stack
- a required hosting platform

These concerns are intentionally left to implementation and organizational policy.

---

## 13. Summary

At a high level:

- regular work happens from `develop`
- `wip/*` is exploratory and never mergeable
- production stays protected on `main`
- only `release/*` may promote into `main`
- hotfixes start from `main`
- production changes must be reconciled back into `develop` through `sync/*`

This is the core Stability Flow model.
