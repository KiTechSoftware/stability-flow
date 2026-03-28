# Stability Flow Enforcement

## 1. Purpose

This document explains how the Stability Flow specification can be **validated and enforced in practice**.

The specification defines the branching model.

This document explains:

- what parts of the model are practical to validate
- where enforcement can happen
- what kinds of checks are most valuable
- where automation helps and where human review still matters

It is intentionally **implementation-neutral**.

For the normative branching model, see:

- [Specification](spec.md)

For branch naming and commit conventions, see:

- [Conventions](conventions.md)

---

## 2. Enforcement Philosophy

The purpose of enforcement is not to make Git more complicated.

The purpose is to make the workflow:

- safer
- more predictable
- more auditable
- easier to follow under pressure

Good enforcement should:

- prevent invalid branch movement
- preserve branch role meaning
- reduce accidental policy drift
- support review and automation
- remain understandable to humans

Good enforcement should **not**:

- hide the workflow behind opaque automation
- turn normal development into ceremony
- rely entirely on tribal knowledge
- depend on one vendor, platform, or CI system

The goal is not maximum restriction.

The goal is preserving the integrity of the model.

---

## 3. What Enforcement Should Protect

Not every rule has equal operational value.

The most useful enforcement protects the parts of Stability Flow that preserve its branching contract.

High-value enforcement usually focuses on:

- branch role integrity
- promotion safety
- reconciliation correctness
- history shape
- protected branch behavior

These are the areas where invalid workflow behavior causes the most confusion or operational risk.

---

## 4. High-Value Enforcement Surfaces

### 4.1. Branch Naming

Branch naming is one of the easiest and highest-value checks.

Branch names often communicate:

- branch role
- intended workflow path
- likely merge behavior
- expected automation behavior

#### 4.1.1 Why it matters

A branch with an invalid or ambiguous name makes the rest of the workflow harder to reason about.

Naming is often the first signal used by:

- reviewers
- CI checks
- repository policy
- tooling

#### 4.1.2 Common validation targets

Typical branch naming validation checks include:

- valid prefix
- valid format
- lowercase naming
- explicit branch intent

#### 4.1.3 Typical enforcement points

- local hooks
- CI checks
- pull request validation
- server-side Git hooks

---

### 4.2. Branch Origin

Branch origin is one of the most important structural checks in the model.

A correctly named branch can still violate the workflow if it was created from the wrong source.

Examples of invalid branch origin include:

- a hotfix branch created from `develop`
- a release branch created from an exploratory branch
- a regular work branch created from `main`

#### 4.2.1. Why it matters

Branch origin is what preserves the meaning of branch roles.

If branch origins drift, the workflow may still *look* compliant while behaving incorrectly.

#### 4.2.2. Common validation targets

Typical branch origin checks include:

- whether the branch ancestry includes the required source branch
- whether the branch diverged from the correct branch family
- whether disallowed origins are present

#### 4.2.3. Typical enforcement points

- local tooling
- CI validation
- pull request checks
- server-side Git hooks

#### Practical note

Branch origin validation usually requires repository graph inspection rather than simple naming checks.

That makes it slightly harder than branch naming, but still high-value.

---

### 4.3. Merge Eligibility

Merge eligibility determines whether a source branch is allowed to merge into a given target branch.

This is one of the most important controls in Stability Flow.

#### 4.3.1 Why it matters

If invalid merge targets are allowed, the workflow can collapse even if branch names and origins are correct.

This is what prevents:

- regular work from bypassing `develop`
- hotfixes from bypassing promotion rules
- invalid reconciliation paths
- exploratory branches from leaking into delivery flow

#### 4.3.2 Common validation targets

Typical merge eligibility checks include:

- source branch type
- target branch type
- whether the transition is valid under the model
- whether the merge path bypasses promotion or reconciliation rules

#### 4.3.3 Typical enforcement points

- pull request checks
- merge queue validation
- CI validation
- protected branch policy
- custom merge automation

---

### 4.4. Merge Strategy

Merge strategy is part of the workflow contract because Stability Flow uses history shape intentionally.

A valid merge target with the wrong merge strategy can still weaken the model.

#### 4.4.1 Why it matters

Examples of invalid merge behavior include:

- merge commits from regular work into `develop` when squash is expected
- non-fast-forward promotion into `main`
- reconciliation that loses its explicit merge history

#### 4.4.2 Common validation targets

Typical merge strategy checks include:

- squash requirement for regular work
- fast-forward requirement for promotion
- merge-commit preservation for reconciliation

#### 4.4.3 Typical enforcement points

- repository branch protection
- platform merge settings
- merge queue rules
- custom merge automation

---

### 4.5. Protected Branch Writes

Some branches should be protected from direct mutation.

In practice, the most important protected branches are usually:

- `main`
- `develop`

#### 4.5.1 Why it matters

Direct pushes can bypass review, validation, and merge policy.

That weakens the workflow even if the written code is otherwise valid.

#### 4.5.2 Common validation targets

Typical protected branch controls include:

- direct push restrictions
- pull request requirements
- required review gates
- required passing checks before merge

#### 4.5.3 Typical enforcement points

- repository branch protection
- platform policy settings
- server-side Git controls

---

### 4.6. Promotion Safety

Promotion safety focuses on whether code reaches production through the correct path.

This is one of the most important areas to protect.

#### 4.6.1 Why it matters

The most important production question is not just:

> did this code work?

It is also:

> did this code reach production through the intended path?

Promotion safety is what preserves the meaning of `main` as the production line.

#### 4.6.2 Common validation targets

Typical promotion safety checks include:

- only promotion branches may target `main`
- release candidates are based on current production truth
- invalid direct production paths are rejected

#### 4.6.3 Typical enforcement points

- pull request checks
- protected branch policy
- release automation
- merge queue validation

---

### 4.7. Reconciliation Safety

Reconciliation safety focuses on whether production truth returns to future development through the correct path.

#### 4.7.1 Why it matters

This is where many workflows quietly drift.

Teams often protect promotion into production but under-protect the path that brings production changes back into future work.

That creates confusion later when planned releases are cut from stale or partially reconciled branches.

#### 4.7.2 Common validation targets

Typical reconciliation checks include:

- reconciliation occurs through `sync/*`
- direct `main` → `develop` merges are rejected
- planned release branches are not created from stale `develop`
- unresolved reconciliation state blocks future planned promotion

#### 4.7.3 Typical enforcement points

- CI validation
- pull request checks
- merge queue validation
- release gating automation

---

### 4.8. Final Squash Commit Conventions

If the repository adopts Stability Flow conventions, the final squash commit is a useful enforcement surface.

#### 4.8.1 Why it matters

Final squash commits are part of long-lived repository history.

Consistent final commit formatting improves:

- readability
- release note generation
- auditability
- automation compatibility

#### 4.8.2 Common validation targets

Typical checks include:

- valid final commit type
- breaking change formatting
- prohibited final forms such as `revert:` if conventions disallow them

#### 4.8.3 Typical enforcement points

- local hooks
- CI validation
- merge-time checks

This enforcement surface is optional if a team chooses not to adopt the conventions document.

---

## 5. What Is Easier vs Harder to Enforce

Not every rule is equally automatable.

Some parts of Stability Flow are structurally easy to validate.

Others require more context or judgment.

---

### 5.1. Easier to Enforce Reliably

The following are usually strong automation candidates:

- branch naming
- branch origin
- merge target validity
- merge strategy
- protected branch writes
- final squash commit format
- direct promotion path restrictions
- direct reconciliation path restrictions

These rules are usually well-suited to validation because they are structurally observable.

---

### 5.2. Harder to Enforce Reliably

The following are often harder to validate perfectly:

- whether a release branch contains only release-safe changes
- whether a branch is truly a hotfix rather than just “important work”
- whether a reconciliation conflict was resolved correctly
- whether a release candidate should be discarded rather than repaired
- whether exploratory work should be recreated instead of reused

These rules are still important.

They are just more likely to require policy, review, or judgment in addition to automation.

---

## 6. Enforcement Layers

Stability Flow can be enforced at multiple layers.

Strong adoption usually combines more than one.

---

### 6.1. Local Validation

Local validation helps developers catch issues before code is pushed.

Examples include:

- branch naming checks
- commit message validation
- branch origin checks
- pre-push merge eligibility checks

#### 6.1.1 Benefits

- fast feedback
- lower CI noise
- easier developer correction

#### 6.1.2 Limitations

- can be bypassed
- depends on local setup
- should not be the only enforcement layer for protected workflows

---

### 6.2. Pull Request and CI Validation

Pull request and CI validation are often the most practical centralized enforcement layer.

Examples include:

- reject invalid branch names
- reject invalid merge targets
- reject invalid branch origins
- validate merge strategy expectations
- validate promotion and reconciliation paths
- validate final squash commit formatting

#### 6.2.1 Benefits

- centralized
- visible
- auditable
- works across teams

#### 6.2.2 Limitations

- catches issues only after code is pushed
- depends on repository and platform integration quality

---

### 6.3. Branch Protection and Repository Policy

Branch protection should reinforce the most important branch guarantees.

Examples include:

- prevent direct pushes to `main`
- prevent direct pushes to `develop`
- require pull requests into protected branches
- restrict merge methods
- require passing checks before merge

#### 6.3.1 Benefits

- strong protection for the most sensitive branches
- difficult to bypass accidentally
- aligns well with the Stability Flow promotion model

#### 6.3.2 Limitations

- platform capabilities vary
- not every rule can be expressed natively

---

### 6.4. Merge Queue and Release Gating

Teams with stricter workflow requirements may also validate Stability Flow at merge queue or release promotion time.

Examples include:

- promotion path validation before merge queue admission
- stale reconciliation detection before release creation
- release ancestry validation
- release gating based on policy compliance

#### 6.4.1 Benefits

- protects the highest-risk workflow transitions
- supports stronger release discipline
- helps prevent policy bypass late in the flow

#### 6.4.2 Limitations

- more workflow complexity
- usually requires more custom integration

---

### 6.5. Server-Side Git Enforcement

Teams with stronger governance needs may enforce parts of Stability Flow at the Git server layer.

Examples include:

- pre-receive hooks
- update hooks
- repository policy gateways
- internal Git control services

#### 6.5.1 Benefits

- difficult to bypass
- strong organizational control
- platform-independent in principle

#### 6.5.2 Limitations

- more operational overhead
- may require custom infrastructure

---

## 7. Recommended Enforcement Layering

The exact implementation may vary by team and platform.

A useful enforcement model usually grows in layers.

### 7.1. Recommended minimum

- branch naming validation
- branch origin validation
- merge eligibility validation
- protected `main`
- protected `develop`

### 7.2. Recommended stronger model

- local validation
- CI validation
- protected branches
- merge strategy enforcement
- promotion path validation
- reconciliation path validation

### 7.3. Recommended mature model

- local validation
- CI validation
- branch protection
- merge strategy enforcement
- promotion and reconciliation gating
- optional server-side enforcement

The important thing is not using every possible mechanism.

The important thing is preserving the branching contract consistently.

---

## 8. Human Review Still Matters

Stability Flow is designed to be enforceable, but not every important decision should be reduced to automation.

Examples that still benefit from human review include:

- whether a hotfix is appropriate
- whether a release branch should be discarded
- whether a release candidate is trustworthy
- whether a reconciliation conflict was resolved correctly
- whether a branch still reflects its intended purpose

Automation should reduce avoidable mistakes.

It should not replace judgment.

---

## 9. Reference Tooling

Reference tooling may be provided to help teams adopt or validate Stability Flow.

Examples may include:

- local validators
- CLI tools
- reusable workflows
- CI integration examples
- policy templates

Reference tooling is an implementation of the model.

It is not the definition of the model.

Teams may use reference tooling, adapt it, or build their own.

---

## 10. Summary

At a practical level, Stability Flow is best enforced by validating:

- branch naming
- branch origin
- merge eligibility
- merge strategy
- protected branch writes
- promotion safety
- reconciliation safety
- optional naming and commit conventions

The exact tooling is flexible.

The important thing is that the workflow remains:

- explicit
- reviewable
- stable
- auditable
- enforceable
