# Stability Flow Design

## 1. Purpose

This document explains the design rationale behind **Stability Flow**.

The specification defines the branching model.

This document explains:

- why the model exists
- what operational problems it is intended to solve
- why the branch roles are shaped the way they are
- what tradeoffs the model makes

It is intentionally explanatory rather than normative.

For the normative contract, see:

- [Specification](spec.md)

---

## 2. Problem Statement

Many teams need a workflow that can support all of the following at the same time:

- ongoing planned development
- controlled production releases
- urgent production hotfixes
- clear recovery after production and development diverge

This is where many branching models start to break down in practice.

The failure usually does not happen during normal development.

It happens when teams are under pressure and need to answer questions such as:

- what is actually safe to release?
- where should a hotfix start?
- how should a production fix return to future development?
- how do we keep production promotion auditable?

Stability Flow exists to make those transitions structurally explicit.

---

## 3. Core Design Goals

Stability Flow is designed around a small set of operational goals.

### 3.1 Keep Production Stable

Teams often need a branch that clearly represents the production line.

That branch should be easy to trust during:

- deployments
- incident response
- release approval
- rollback decisions
- audits

Stability Flow treats production as a protected destination rather than a normal development lane.

---

### 3.2 Make Promotion Explicit

A change should not reach production accidentally.

Production promotion should be visible and intentional.

That is why Stability Flow separates:

- **integration of planned work**
- **promotion of releasable work**

This keeps “work in progress” and “production candidate” from collapsing into the same path.

---

### 3.3 Support Safe Hotfixes

A production hotfix should start from production reality, not from whatever planned work happens to be in progress.

That sounds obvious, but many teams still end up pulling unreleased work into emergency fixes because the workflow does not make the correct path obvious enough.

Stability Flow is designed to make the safe hotfix path structurally clear.

---

### 3.4 Treat Divergence as Normal

In real repositories, production and planned development often diverge.

For example:

- `develop` continues moving
- `main` receives a hotfix
- the branches temporarily represent different truths

This is not a workflow failure.

It is normal release-management behavior.

Stability Flow is designed to make that divergence understandable rather than pretending it should never happen.

---

### 3.5 Make Reconciliation Explicit

When production changes and planned work diverge, bringing those realities back together is not bookkeeping.

It is a meaningful operational event.

That is why Stability Flow treats reconciliation as a first-class part of the model rather than an informal cleanup step.

---

### 3.6 Keep the Model Enforceable

A branching strategy is much easier to follow consistently when its most important rules are structurally visible and easy to validate.

That is why Stability Flow is designed around:

- explicit branch roles
- constrained branch origins
- constrained merge targets
- explicit promotion paths
- explicit reconciliation paths

The model is intentionally shaped so that high-value rules can be checked by policy or tooling.

---

## 4. Why the Branch Roles Exist

The Stability Flow model is built around branch roles with deliberately different responsibilities.

Each branch exists because it serves a distinct operational purpose.

---

### 4.1 Why `main` Exists

`main` exists to represent production truth.

Many teams still need a branch that answers a simple operational question:

> what is production right now?

That matters for:

- release confidence
- incident handling
- auditability
- rollback reasoning
- production support

Stability Flow keeps `main` narrow on purpose.

Its value comes from being trustworthy.

---

### 4.2 Why `develop` Exists

`develop` exists to carry planned work forward without forcing production promotion.

This is useful when:

- multiple changes are in progress at once
- not all work is intended for the same release
- planned work should continue while production remains protected
- a hotfix may need to ship before current planned work is ready

The purpose of `develop` is not just “another branch.”

It is the integration line for future release work.

---

### 4.3 Why Regular Work Stays Off `main`

Day-to-day work such as:

- features
- fixes
- docs
- CI changes
- refactors
- chores

is not automatically production-ready just because it is approved for integration.

Stability Flow keeps that work off the production line until it is intentionally promoted.

This creates a cleaner boundary between:

- **work that is being integrated**
- **work that is being released**

That separation is one of the most important design choices in the model.

---

### 4.4 Why `release/*` Exists

Production promotion is not the same thing as normal development.

Even when code is “done,” teams often still need a release lane for activities such as:

- release validation
- version preparation
- changelog preparation
- release metadata updates
- final release checks

A dedicated `release/*` branch keeps those activities separate from ordinary development work.

That is why Stability Flow treats release preparation as a distinct operational phase.

---

### 4.5 Why `release/*` Is Not a Work Branch

One of the fastest ways for a release process to collapse is for the release branch to quietly become “another place where coding continues.”

That makes it much harder to reason about:

- what is actually being released
- whether the release candidate is still trustworthy
- whether production promotion is still auditable

Stability Flow avoids that by treating `release/*` as a promotion lane, not a development lane.

That keeps the release boundary explicit.

---

### 4.6 Why Hotfixes Start From `main`

A hotfix is specifically a repair to production reality.

That means it should start from the production line, not from future planned work.

If a hotfix starts from an integration branch, it can accidentally absorb unrelated or unreleased changes.

That undermines the entire point of isolating an urgent production fix.

Stability Flow makes the hotfix starting point explicit so the safe path is structurally obvious.

---

### 4.7 Why Hotfixes Still Promote Through `release/*`

Even urgent production fixes still benefit from a consistent promotion path.

Allowing one special “direct-to-production” path for hotfixes tends to create a second, looser set of rules around the most sensitive production changes.

Stability Flow intentionally avoids that.

The model keeps one production promotion shape:

- production changes are promoted intentionally
- production stays protected
- promotion remains auditable

This keeps emergency fixes from becoming an exception that weakens the whole model.

---

### 4.8 Why `sync/*` Exists

Reconciliation is one of the most important parts of the model.

A common scenario looks like this:

- `develop` continues moving
- production receives a hotfix
- the production fix now needs to return to future development

That reconciliation is meaningful enough that it should be visible in repository history.

A dedicated `sync/*` branch makes that visible.

It gives teams a clear place to:

- merge production truth back into planned work
- resolve conflicts
- review the reconciliation
- preserve the operational meaning of what happened

That is why Stability Flow treats reconciliation as a distinct branch role rather than an informal merge step.

---

### 4.9 Why Direct `main` → `develop` Merges Are Outside the Model

A direct merge from `main` into `develop` can technically work in Git.

But Stability Flow is not designed around what is merely possible.

It is designed around what keeps intent visible.

A dedicated reconciliation branch preserves important context such as:

- when production changed
- why reconciliation happened
- where conflicts were resolved
- what release or hotfix caused the divergence

That is operationally valuable.

So Stability Flow prefers explicit reconciliation over convenience.

---

### 4.10 Why `wip/*` Exists

Some teams benefit from a branch type that supports:

- exploration
- debugging spikes
- investigations
- rough experiments
- dead-end discovery

That kind of work is real and often useful.

The problem is not that exploratory work exists.

The problem is when exploratory history is accidentally treated as delivery history.

Stability Flow allows `wip/*` so teams have a place for temporary exploration without pretending it is already integration-ready.

---

## 5. Why the History Shape Matters

Stability Flow is not only about which branches exist.

It also cares about what repository history means.

The merge strategy choices are intentional because they preserve different kinds of meaning.

---

### 5.1 Why Regular Work Is Squash Merged

Regular work branches usually contain branch-local history such as:

- checkpoints
- partial refactors
- local reversions
- temporary commits
- incremental experimentation

That history is useful while work is in progress, but it is not always useful as long-lived integration history.

Squash merge keeps the long-lived integration history more intentional.

This helps with:

- release review
- repository readability
- auditability
- automation

---

### 5.2 Why Production Promotion Stays Linear

A production branch is easier to trust when its promotion history is straightforward.

A linear production line makes it easier to answer questions such as:

- what shipped?
- in what order?
- what release introduced this state?

That is why Stability Flow treats production promotion as a clean, auditable path rather than as an arbitrary merge history.

---

### 5.3 Why Reconciliation Preserves Merge History

Reconciliation is not just “another merge.”

It is a visible moment where production truth is brought back into future development.

Preserving that merge history helps maintainers and reviewers understand:

- when divergence occurred
- when it was resolved
- where conflict resolution happened

That makes the repository easier to reason about over time.

---

## 6. Why Stability Flow Is More Structured Than Simpler Workflows

Stability Flow deliberately introduces more structure than looser workflows such as direct trunk-only development.

That is not because simpler workflows are inherently wrong.

It is because some teams need explicit support for:

- planned releases
- protected production promotion
- hotfix isolation
- post-release reconciliation
- machine-checkable branch behavior

For teams that do not need those things, Stability Flow may be unnecessary.

For teams that do, the extra structure is usually a net benefit.

---

## 7. Tradeoffs

Stability Flow is intentionally opinionated.

That means it makes tradeoffs.

### It adds process shape

Compared to lighter workflows, it introduces more explicit branch roles such as:

- `develop`
- `release/*`
- `hotfix/*`
- `sync/*`

### In return, it gives teams

- clearer production promotion
- safer hotfix handling
- explicit divergence handling
- clearer repository intent
- stronger auditability
- stronger enforcement opportunities

This tradeoff is deliberate.

The model prefers operational clarity over minimal branching.

---

## 8. When Stability Flow Is a Good Fit

Stability Flow is a good fit for teams that need one or more of the following:

- a clearly protected production branch
- planned release boundaries
- occasional urgent production hotfixes
- explicit release promotion
- explicit reconciliation after production divergence
- a branching model that can be validated by policy or tooling

It is especially useful when teams need both:

- **ongoing planned work**
- **production safety under pressure**

---

## 9. When Stability Flow May Be a Poor Fit

Stability Flow may be unnecessarily heavy for teams that:

- ship continuously from a single trunk without release staging
- do not need a long-lived integration branch
- do not need explicit hotfix isolation
- do not need explicit reconciliation history
- prefer minimal process over stronger branch semantics

That does not make Stability Flow wrong.

It simply means it is designed for a particular operational shape, not every possible team.

---

## 10. Summary

At a high level, Stability Flow is designed so that:

- production remains trustworthy
- promotion is explicit
- hotfixes are isolated from unreleased work
- divergence is treated as normal
- reconciliation is explicit and reviewable
- branch movement is structurally meaningful
- the most important rules are enforceable

That is the core design logic behind the model.
