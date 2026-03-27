# Stability Flow Specification

## 1. Overview

**Stability Flow** is a branching strategy specification for teams that need:

- a **stable production branch**
- **planned releases**
- **safe hotfixes**
- **explicit reintegration after production divergence**

It is designed as an alternative to Gitflow for teams that want stronger control over production stability, clearer release intent, and a simpler enforcement surface.

Stability Flow is especially suited to teams that:

- release on a planned cadence
- occasionally need urgent production hotfixes
- want to keep `main` stable and auditable
- want regular work isolated from production promotion

---

## 2. Goals

Stability Flow is designed to achieve the following goals:

1. **Protect production stability**
2. **Allow urgent hotfix releases**
3. **Support planned releases from ongoing development**
4. **Make branch intent explicit**
5. **Keep production promotion auditable**
6. **Make enforcement practical through automation and policy**

---

## 3. Branch Roles

Stability Flow defines the following branch roles.

### 3.1 `main`

`main` is the **stable production branch**.

It represents the current production-ready state of the system.

#### Requirements

- `main` **MUST** remain stable
- `main` **MUST NOT** receive direct regular work
- `main` **MUST ONLY** receive changes from `release/*`
- `main` **SHOULD** remain linear where practical

---

### 3.2 `develop`

`develop` is the **integration branch for regular work**.

It represents the next planned release line.

#### Requirements

- regular work branches **MUST** branch from `develop`
- regular work branches **MUST ONLY** merge into `develop`
- `develop` **MAY** move ahead of `main`
- `develop` **MUST** eventually receive production changes from `main` after releases and hotfixes

---

### 3.3 Regular Work Branches

Regular work branches are short-lived branches used for day-to-day work.

Examples include:

- `feat/*`
- `fix/*`
- `docs/*`
- `ci/*`
- `refactor/*`
- `chore/*`

#### Requirements

- regular work branches **MUST** branch from `develop`
- regular work branches **MUST ONLY** merge into `develop`
- regular work branches **SHOULD** be short-lived
- regular work branches **SHOULD** be deleted after merge

#### Merge Strategy

Regular work branches **SHOULD** be **squash merged** into `develop`.

This keeps the integration history clean while preserving work context in pull requests and branch history.

---

### 3.4 `release/*`

`release/*` branches are **promotion branches** used to move code into `main`.

A `release/*` branch represents a candidate production release.

Examples:

- `release/1.3.0`
- `release/2.0.0`
- `release/1.2.4`

#### Requirements

- `release/*` **MUST** be created from:
  - `develop`, or
  - `hotfix/*`
- `release/*` **MUST ONLY** merge into `main`
- `release/*` **SHOULD** contain only release preparation and promotion-safe changes
- `release/*` **SHOULD** be deleted after merge or disposal

#### Intended Use

A `release/*` branch exists to:

- prepare a release
- validate a release
- promote a release into `main`

A `release/*` branch is **not** a general-purpose work branch.

---

### 3.5 `hotfix/*`

`hotfix/*` branches are used for urgent production fixes.

Examples:

- `hotfix/1.2.4`
- `hotfix/2.0.1`

#### Requirements

- `hotfix/*` **MUST** branch from `main`
- `hotfix/*` **MUST NOT** merge directly into `main`
- `hotfix/*` **MUST** be promoted through `release/*`
- `hotfix/*` **SHOULD** be short-lived
- `hotfix/*` **SHOULD** be deleted after release

#### Intended Use

A `hotfix/*` branch exists only to isolate an urgent production fix from the current stable line.

It is not a general-purpose maintenance branch.

---

### 3.6 `sync/*`

`sync/*` branches are optional short-lived branches used to reintegrate production changes back into `develop`.

Examples:

- `sync/main-into-develop`
- `sync/1.2.4-backport`

#### Requirements

- `sync/*` **SHOULD** branch from `develop`
- `sync/*` **SHOULD** merge `main` into `develop`
- `sync/*` **SHOULD** be used when explicit reintegration review is desired
- `sync/*` **SHOULD** be deleted after merge

#### Intended Use

A `sync/*` branch exists to make reintegration explicit and reviewable.

It is especially useful after:

- hotfix releases
- production divergence
- any release where `develop` has moved ahead of `main`

---

## 4. Core Rules

### 4.1 Production Promotion Rule

Only `release/*` branches may promote changes into `main`.

#### Therefore:

- `feat/*` **MUST NOT** merge into `main`
- `fix/*` **MUST NOT** merge into `main`
- `docs/*` **MUST NOT** merge into `main`
- `ci/*` **MUST NOT** merge into `main`
- `refactor/*` **MUST NOT** merge into `main`
- `chore/*` **MUST NOT** merge into `main`
- `hotfix/*` **MUST NOT** merge directly into `main`

---

### 4.2 Regular Work Isolation Rule

Regular work must stay on the development line until intentionally promoted.

#### Therefore:

- regular work branches **MUST** start from `develop`
- regular work branches **MUST ONLY** merge into `develop`

---

### 4.3 Hotfix Isolation Rule

Urgent production fixes must start from the current production state.

#### Therefore:

- `hotfix/*` **MUST** branch from `main`

This ensures the fix is made against the actual stable production line.

---

### 4.4 Explicit Reintegration Rule

Production changes must be reintegrated into the development line after release.

#### Therefore:

- after a release to `main`, production changes **MUST** be brought back into `develop`
- this **MAY** happen:
  - directly, or
  - through a `sync/*` branch

Teams that want stronger consistency and clearer review **SHOULD** use a `sync/*` branch for reintegration.

---

## 5. Allowed Branch Origins

The following origins are allowed.

| Branch Type | Allowed Origin |
|---|---|
| `feat/*` | `develop` |
| `fix/*` | `develop` |
| `docs/*` | `develop` |
| `ci/*` | `develop` |
| `refactor/*` | `develop` |
| `chore/*` | `develop` |
| `hotfix/*` | `main` |
| `release/*` | `develop`, `hotfix/*` |
| `sync/*` | `develop` |

---

## 6. Allowed Merge Targets

The following merge targets are allowed.

| Source Branch | Allowed Target |
|---|---|
| regular work branches | `develop` |
| `release/*` | `main` |
| `sync/*` | `develop` |
| `main` | `develop` (directly or via `sync/*`) |

---

## 7. Recommended Merge Strategies

Stability Flow defines branch roles and merge intent. Teams **MAY** implement the exact merge mechanics differently if the branch role and history semantics are preserved.

The following strategies are recommended.

| Flow | Recommended Strategy |
|---|---|
| regular work → `develop` | squash merge |
| `release/*` → `main` | fast-forward only |
| `main` → `develop` | merge commit |
| `sync/*` → `develop` | merge commit |

### Rationale

#### Squash into `develop`
Keeps regular work history concise and intentional.

#### Fast-forward into `main`
Keeps production promotion linear and auditable.

#### Merge `main` back into `develop`
Preserves the fact that production changes were reintegrated.

---

## 8. Release Lifecycle

### 8.1 Planned Release

A planned release follows this shape:

1. regular work is integrated into `develop`
2. a `release/*` branch is created from `develop`
3. release preparation and validation occur on `release/*`
4. `release/*` is merged into `main`
5. production changes are reintegrated into `develop`
6. temporary branches are deleted

---

### 8.2 Hotfix Release

A hotfix release follows this shape:

1. a `hotfix/*` branch is created from `main`
2. the urgent fix is implemented on `hotfix/*`
3. a `release/*` branch is created from `hotfix/*`
4. release preparation and validation occur on `release/*`
5. `release/*` is merged into `main`
6. production changes are reintegrated into `develop`
7. temporary branches are deleted

---

## 9. Divergence and Reintegration

A core Stability Flow scenario is:

- `develop` is ahead with planned work
- `main` receives an urgent hotfix release
- the hotfix must be safely brought back into `develop`

This is expected behavior, not an edge case.

Stability Flow handles this by requiring explicit reintegration from `main` back into `develop`.

This reintegration **MUST** happen before future planned releases continue normally.

Teams **SHOULD** prefer a `sync/*` branch when they want:

- reviewable reintegration
- explicit conflict handling
- consistent operational muscle memory

---

## 10. Release Branch Disposal

Not all release branches should be repaired or reused.

If a release candidate becomes invalid, teams **SHOULD** prefer disposing of the branch and creating a fresh release branch rather than repeatedly mutating a stale release branch.

This keeps release intent clear and reduces branch drift.

---

## 11. Branch Naming

Branch names **SHOULD** be explicit and machine-checkable.

Recommended formats:

```text
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

Branch names **SHOULD** be:

* lowercase
* slash-prefixed by role
* descriptive enough to express intent

---

## 12. Commit Message Guidance

Stability Flow does not require a specific commit convention, but teams **SHOULD** use a consistent, machine-checkable commit format.

A Conventional Commits style is recommended.

Examples:

```text
feat: add authentication flow
fix: patch production issue
docs: clarify release promotion rules
chore: prepare release 1.3.0
```

### Squash Commits

When regular work branches are squash merged into `develop`, the resulting squash commit **SHOULD** clearly describe the merged change.

Examples:

```text
feat: add authentication flow
fix: patch session race condition
docs: clarify hotfix reconciliation
```

---

## 13. Normative Language

The key words **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT**, and **MAY** in this specification are to be interpreted as described in RFC 2119 and RFC 8174 when, and only when, they appear in all capitals.

### Interpretation

* **MUST / MUST NOT**
  Absolute requirement.

* **SHOULD / SHOULD NOT**
  Strong recommendation. There may be valid reasons to deviate, but the tradeoff should be understood.

* **MAY**
  Optional behavior.

---

## 14. Non-Goals

Stability Flow does **not** define:

* how teams must version software
* how release approval must work
* how deployments must be performed
* how CI/CD must be implemented
* a required tooling stack
* a required hosting platform

These concerns are intentionally left to implementation and organizational policy.

---

## 15. Summary

At a high level:

* regular work happens from `develop`
* production stays protected on `main`
* only `release/*` may promote into `main`
* hotfixes start from `main`
* production changes must come back into `develop`
* `sync/*` provides explicit reintegration when desired

This is the core Stability Flow model.
