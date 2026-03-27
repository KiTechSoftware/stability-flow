# Design

## 1. Overview

Stability Flow is designed for teams that need a branching strategy centered around **production stability**, **planned release promotion**, and **safe hotfix reintegration**.

It is not intended to be the most minimal possible Git workflow.

It is intended to make release movement and production safety explicit.

This document explains the reasoning behind the design of the Stability Flow specification.

For normative rules, see the specification.

---

## 2. Why Stability Flow Exists

Many teams eventually run into the same set of problems:

- production needs to remain stable
- planned work continues in parallel
- urgent hotfixes must sometimes ship before the next planned release
- development and production can diverge temporarily
- branch history becomes hard to reason about under pressure

Many common Git workflows address parts of this problem, but often leave important operational behavior implicit.

Stability Flow exists to make those behaviors explicit and easier to reason about.

---

## 3. Core Design Priorities

Stability Flow is built around the following priorities:

1. **Protect production**
2. **Allow urgent intervention**
3. **Make promotion explicit**
4. **Make reintegration explicit**
5. **Keep regular work isolated from production**
6. **Remain practical to enforce**

These priorities drive the shape of the branch model.

---

## 4. Why `main` Is Protected

In Stability Flow, `main` represents the stable production line.

That means it should behave differently from an ordinary integration branch.

The design intent is that `main` should answer a simple question:

> “What is the current stable production-ready state?”

That becomes much harder if `main` is used for:

- regular feature development
- direct maintenance work
- ad hoc release preparation
- mixed-purpose merges

Restricting `main` to promotion through `release/*` makes production movement easier to audit and reason about.

---

## 5. Why `develop` Exists

Stability Flow intentionally keeps a separate development integration line.

That line is `develop`.

This is a deliberate choice.

Some workflows attempt to eliminate long-lived integration branches entirely, but Stability Flow keeps `develop` because it provides a clear answer to another important question:

> “What is currently being prepared for the next planned release?”

This matters especially for teams that:

- release on a cadence
- have multiple in-flight work branches
- need to stage work before production promotion
- may need to hotfix production while planned work continues

`develop` is the place where regular work accumulates before being promoted.

---

## 6. Why Regular Work Must Stay Off `main`

A core design decision in Stability Flow is that regular work should never flow directly into `main`.

This is not only about “branch naming”.

It is about **promotion discipline**.

If regular work can reach `main` directly, then:

- production movement becomes harder to audit
- release boundaries become less clear
- hotfix behavior becomes less predictable
- enforcement becomes weaker

The specification therefore isolates regular work on the development line and only allows promotion through `release/*`.

This is one of the most important design choices in the model.

---

## 7. Why Work Branches Should Be Squash Merged into `develop`

Regular work branches are recommended to be squash merged into `develop`.

This is a deliberate tradeoff.

### Benefits

Squash merging regular work:

- keeps integration history cleaner
- reduces noisy intermediate commits
- makes the intent of merged work clearer
- aligns well with pull-request-based review

### Tradeoff

Squash merging does discard the exact internal branch commit structure from the main integration history.

That is acceptable in Stability Flow because:

- the work branch itself carries that history while active
- pull requests preserve review and discussion context
- the integration line benefits from clearer intent

This design is especially useful when teams want `develop` to reflect **meaningful integrated changes**, not every intermediate work-in-progress commit.

---

## 8. Why `release/*` Exists

A central design choice in Stability Flow is the use of an explicit promotion branch:

# `release/*`

This branch exists to make one thing clear:

> “This change is being prepared for production promotion.”

That explicit transition is valuable.

It gives teams a place to:

- prepare release metadata
- validate the release candidate
- run release checks
- stage release approval
- keep promotion behavior separate from ordinary development

This is why `release/*` is not treated as a normal work branch.

It has a very specific role.

---

## 9. Why `release/*` Must Be Disposable

Another deliberate design decision is that release branches should be easy to discard.

If a release candidate becomes invalid, teams should prefer:

- deleting the release branch
- creating a fresh one

rather than repeatedly mutating and repairing a stale release branch.

### Why this matters

Long-lived, repeatedly repaired release branches tend to accumulate confusion:

- what changed for the release?
- what was carried over?
- what was abandoned?
- what is still intended?

Disposable release branches keep release intent cleaner and reduce branch drift.

That is a major operational benefit.

---

## 10. Why Hotfixes Start from `main`

Stability Flow requires hotfixes to start from `main`.

This is a critical design rule.

A hotfix exists to patch the current production line.

That means it should be created from the actual stable production state, not from whatever is currently in `develop`.

If hotfixes start from `develop`, they risk inheriting:

- unreleased features
- unfinished fixes
- integration-only changes
- unrelated branch drift

That undermines the entire purpose of a hotfix.

Starting from `main` keeps the urgent fix isolated and production-relevant.

---

## 11. Why Hotfixes Still Promote Through `release/*`

It may be tempting to allow `hotfix/*` to merge directly into `main`.

Stability Flow intentionally does **not** do that.

Instead, hotfixes still promote through `release/*`.

### Why?

Because it preserves a consistent production promotion rule:

> only `release/*` branches promote into `main`

This consistency matters.

It means:

- production promotion is always explicit
- release checks can stay consistent
- tooling can validate one promotion model
- the workflow is easier to reason about under pressure

That consistency is worth the extra step.

---

## 12. Why Reintegration Is Explicit

One of the most important design goals of Stability Flow is to treat production divergence as **normal**, not exceptional.

A common real-world scenario is:

- `develop` is ahead with planned work
- a production hotfix must ship immediately
- the production line and development line diverge

This is not a mistake. It is an expected operational reality.

What matters is how the workflow handles it.

Stability Flow therefore requires explicit reintegration from `main` back into `develop`.

This is one of the defining ideas in the model.

---

## 13. Why `sync/*` Exists

Some teams are comfortable merging `main` back into `develop` directly.

Others want that reintegration to be explicit, reviewable, and operationally consistent.

That is why Stability Flow includes the optional `sync/*` branch role.

A `sync/*` branch provides a dedicated place to:

- merge production changes back into development
- resolve conflicts safely
- review reintegration intentionally
- build muscle memory around production synchronization

This is especially valuable after hotfixes.

### Important note

`sync/*` is not required because Git needs it.

It exists because teams often benefit from making reintegration behavior explicit and repeatable.

That is a design choice in favor of operational clarity.

---

## 14. Why This Is Not Just Gitflow Renamed

Stability Flow is related to Gitflow, but it is not simply Gitflow with different names.

The differences are meaningful.

### Key differences include:

- regular work is expected to squash into `develop`
- production promotion is explicitly centered on `release/*`
- hotfix reintegration is treated as a first-class concern
- `sync/*` exists to make reintegration explicit where desired
- the model is intentionally designed with enforceability in mind

Gitflow is historically useful, but many teams experience it as either:

- too heavy in some places
- too implicit in others

Stability Flow attempts to keep the useful release separation while making operational expectations clearer.

---

## 15. Why This Is Not a Trunk-Based Workflow

Stability Flow is also intentionally not a pure trunk-based model.

That is not because trunk-based development is “wrong”.

It is because Stability Flow optimizes for a different set of needs.

Trunk-based workflows are often strongest when teams have:

- very high deployment confidence
- short-lived changes
- strong continuous delivery discipline
- low operational friction for rollback and recovery

Stability Flow is designed for teams that instead want:

- an explicit development line
- explicit release promotion
- explicit hotfix isolation
- explicit reintegration after divergence

That is a different optimization target.

---

## 16. Why Enforceability Matters

A workflow that only exists as tribal knowledge is fragile.

A major design goal of Stability Flow is that it should be understandable enough to follow manually and structured enough to enforce automatically.

That means the workflow intentionally favors things like:

- explicit branch roles
- clear allowed origins
- clear allowed merge targets
- clear promotion boundaries
- machine-checkable naming patterns

This is not because tooling should define the workflow.

It is because good workflows should be clear enough that tooling can support them.

---

## 17. Tradeoffs

Stability Flow makes deliberate tradeoffs.

### It gives you:

- strong production discipline
- explicit release promotion
- safe hotfix handling
- clear reintegration paths
- a practical enforcement surface

### It costs you:

- more structure than ad hoc Git usage
- more branch roles than a minimal workflow
- one extra promotion step for hotfixes
- more intentional branch hygiene

These tradeoffs are intentional.

Stability Flow is not trying to be the smallest possible Git workflow.

It is trying to be a stable and operationally sane one.

---

## 18. Design Summary

Stability Flow is built around a simple idea:

> keep production safe, make promotion explicit, and treat reintegration as a first-class part of the workflow.

That leads to a design where:

- `main` stays protected
- `develop` carries the next planned release line
- regular work stays off production
- `release/*` controls promotion
- `hotfix/*` isolates urgent production fixes
- `sync/*` makes reintegration explicit when needed

That is the reasoning behind the model.
