# Enforcement

## 1. Overview

Stability Flow is a specification.

It defines:

- branch roles
- allowed branch origins
- allowed merge targets
- promotion rules
- reintegration expectations
- naming and workflow conventions

This document describes **how those rules can be enforced**.

It is intentionally **tool-neutral**.

Teams may enforce Stability Flow through:

- local developer tooling
- Git hooks
- pull request checks
- CI pipelines
- branch protection rules
- reusable workflows
- release automation
- custom internal tooling
- reference tooling provided by this project

---

## 2. Enforcement Philosophy

The purpose of enforcement is not to make Git “more complicated”.

The purpose is to make the workflow:

- predictable
- auditable
- safer under pressure
- easier to follow consistently

Good enforcement should:

- prevent invalid branch movement
- make intended branch roles explicit
- reduce accidental policy drift
- support review and automation
- remain understandable to humans

Good enforcement should **not**:

- hide the workflow behind opaque automation
- make normal development painful
- rely entirely on tribal knowledge
- depend on one specific vendor or platform

---

## 3. What Can Be Enforced

Some Stability Flow rules are easy to enforce automatically.

Others are only partially enforceable and still require human review.

---

## 4. High-Value Enforcement Targets

The following rules are especially valuable to enforce.

### 4.1 Branch Naming

Branch naming is one of the easiest and highest-value checks.

Examples:

```text id="c0yd3g"
feat/add-authentication
fix/race-on-authentication
docs/update-release-policy
hotfix/1.2.4
release/1.3.0
sync/main-into-develop
```

#### Why enforce it

Branch naming makes branch role explicit.

It allows automation and reviewers to infer intent from branch identity.

#### Typical enforcement points

* local hooks
* CI checks
* pull request validation
* server-side Git hooks

---

### 4.2 Branch Origin

Some branch types must be created from specific base branches.

Examples:

* regular work branches should start from `develop`
* `hotfix/*` should start from `main`
* `release/*` should start from `develop` or `hotfix/*`

#### Why enforce it

A correctly named branch created from the wrong origin can still violate the flow.

Example:

* `hotfix/1.2.4` created from `develop` is structurally wrong even if the name looks valid

#### Typical enforcement points

* local tooling
* CI validation
* pull request checks
* server-side Git hooks

#### Important note

Branch origin is often harder to validate than branch naming because it may require repository graph inspection rather than simple string matching.

---

### 4.3 Merge Eligibility

Merge eligibility determines whether a source branch may merge into a given target branch.

Examples:

* `feat/*` may merge into `develop`
* `release/*` may merge into `main`
* `hotfix/*` should not merge directly into `main`

#### Why enforce it

This is one of the most important controls in Stability Flow.

It prevents invalid promotion paths.

#### Typical enforcement points

* pull request checks
* merge queue validation
* CI validation
* protected branch policies
* custom merge automation

---

### 4.4 Merge Strategy

The merge strategy matters because Stability Flow relies on branch history semantics.

Examples:

* regular work into `develop` should usually be squash merged
* `release/*` into `main` should usually be fast-forward only
* reintegration into `develop` should usually preserve an explicit merge

#### Why enforce it

The wrong merge strategy can preserve the branch target while still weakening the workflow.

Example:

* a merge commit from a feature branch into `develop` may technically “work”, but violates the intended history shape if squash merge is required by policy

#### Typical enforcement points

* repository branch protection
* platform merge settings
* PR workflow rules
* custom merge automation

---

### 4.5 Protected Branch Writes

Some branches should be protected from direct mutation.

Examples:

* `main`
* `develop`

#### Why enforce it

Direct pushes can bypass review and invalidate the flow.

#### Typical enforcement points

* repository branch protection
* server-side Git hooks
* platform policy settings

---

### 4.6 Commit Message Conventions

Commit messages are not the core of Stability Flow, but they are often useful to enforce alongside it.

Examples:

```text id="h75rnt"
feat: add authentication flow
fix: patch production issue
docs: clarify hotfix reconciliation
chore: prepare release 1.3.0
```

#### Why enforce it

Consistent commit messages improve:

* traceability
* release notes
* auditability
* tooling compatibility

#### Typical enforcement points

* local hooks
* CI validation
* merge-time checks

---

## 5. What Is Harder to Enforce Reliably

Some rules are harder to validate fully through automation.

These should still be reviewed, but enforcement may be partial or heuristic.

---

### 5.1 Release Branch Content Purity

A `release/*` branch should usually contain only release-safe changes.

That is easy to describe, but harder to enforce reliably.

For example:

* version bumps are easy to detect
* changelog changes are easy to detect
* “this file change is acceptable release preparation” is often context-dependent

#### Practical approach

Teams may choose to enforce:

* allowed file patterns
* allowed commit types
* release PR review requirements

But some judgment usually remains human.

---

### 5.2 “Was This Really a Hotfix?”

A branch may be named `hotfix/*`, but that does not prove it is truly urgent production work.

This is partly a policy and review concern, not only a technical one.

#### Practical approach

Use:

* PR review
* incident review
* approval gates
* release sign-off

Automation can support this, but cannot fully decide intent.

---

### 5.3 Sync Branch Necessity

A team may prefer using `sync/*` for reintegration, but determining when it is “required” may be contextual.

Example considerations:

* size of divergence
* merge conflict likelihood
* operational risk
* release timing

This is often better handled by team policy than strict automation.

---

## 6. Where Enforcement Can Happen

Stability Flow can be enforced at multiple layers.

Strong adoption usually combines more than one.

---

## 7. Local Enforcement

Local enforcement helps catch issues before code is pushed.

Examples:

* branch name checks
* commit message validation
* branch origin checks
* pre-push merge eligibility checks

### Benefits

* fast feedback
* lower CI noise
* easier developer correction

### Limitations

* can be bypassed
* depends on local setup
* not sufficient on its own for protected workflows

---

## 8. Pull Request and CI Enforcement

CI and pull request validation are among the most practical enforcement layers.

Examples:

* reject invalid branch names
* reject invalid merge targets
* reject invalid branch origins
* validate commit messages
* validate release branch rules

### Benefits

* centralized
* visible
* auditable
* works across teams

### Limitations

* only catches issues once code is pushed
* depends on platform integration quality

---

## 9. Branch Protection and Repository Policy

Repository policy should reinforce the most important branch protections.

Examples:

* prevent direct pushes to `main`
* require pull requests into `main`
* restrict merge methods
* require passing checks before merge

### Benefits

* strong control over protected branches
* difficult to bypass accidentally
* aligns well with Stability Flow’s promotion model

### Limitations

* platform capabilities vary
* not every rule can be expressed in native repository settings

---

## 10. Server-Side Enforcement

Teams with stronger governance requirements may enforce Stability Flow through server-side Git controls.

Examples:

* pre-receive hooks
* update hooks
* protected repository services
* internal Git policy gateways

### Benefits

* difficult to bypass
* platform-independent in principle
* strong organizational control

### Limitations

* more operational overhead
* may require custom infrastructure

---

## 11. Recommended Enforcement Layering

The following layered model is recommended.

### Recommended minimum

* branch naming validation
* merge eligibility validation
* protected `main`
* protected `develop`

### Recommended stronger model

* local validation
* CI validation
* protected branches
* merge strategy enforcement
* release PR review requirements

### Recommended mature model

* local validation
* CI validation
* branch protection
* reusable workflows
* release automation
* optional server-side enforcement

The exact implementation may vary by team and platform.

---

## 12. Human Review Still Matters

Stability Flow should be enforceable, but not every important decision can be automated.

Examples that still require human review include:

* whether a hotfix is appropriate
* whether a release branch should be discarded
* whether a reintegration path is safe
* whether a release is ready for promotion

Automation should reduce mistakes, not replace judgment.

---

## 13. Reference Tooling

This project may provide reference tooling and integrations to help enforce Stability Flow.

These are implementations of the specification, not the definition of the specification.

Examples may include:

* local validators
* CLI tools
* GitHub Actions
* reusable workflows
* CI integration examples
* branch policy templates

Teams may use these directly, adapt them, or build their own tooling.

---

## 14. Enforcement Summary

At a practical level, Stability Flow is best enforced by validating:

* branch names
* branch origins
* merge targets
* merge strategy
* protected branch writes
* release promotion paths

The exact tooling is flexible.

The important thing is that the workflow remains:

* explicit
* stable
* reviewable
* auditable
