# Getting Started

This guide shows the easiest ways to adopt **Stability Flow** in an existing repository.

You do **not** need to adopt everything at once.

The recommended path is:

1. start with PR validation
2. optionally add branch validation
3. optionally add GitHub rulesets
4. optionally use helper workflows to create release/hotfix/sync branches

---

## Branch Model

Stability Flow uses these branch types:

- `main`
- `develop`
- `release/*`
- `hotfix/*`
- `sync/*`
- `wip/*` (optional)

If you are not yet familiar with the model, read the [Specification](spec.md) first.

---

## Option A — Use the full reusable workflow

The easiest way to try Stability Flow is to call the included workflow from your repository.

Create:

```yaml
# .github/workflows/stability-flow.yml
name: Stability Flow

on:
  pull_request:
  push:

jobs:
  stability-flow:
    uses: KiTechSoftware/stability-flow/.github/workflows/stability-flow.yml@main
```

This is the fastest way to begin validating Stability Flow rules.

> Note: `@main` is used during the current pre-1.0 refinement phase. Versioned tags will be recommended once Stability Flow reaches `v1`.

---

## Option B — Copy the workflow into your repo

If you prefer to own and customize the workflow directly, copy the provided workflow files from this repository into your own:

### Recommended starting files

- `stability-flow.yml`
- `pr-to-main.yml`
- `pr-to-develop.yml`

### Optional branch validation files

- `release-branch.yml`
- `hotfix-branch.yml`
- `sync-branch.yml`
- `wip-branch.yml`

This approach is useful if you want to adapt the behavior before Stability Flow reaches `v1`.

---

## Option C — Roll out incrementally (recommended for existing teams)

If you are introducing Stability Flow into an active repository, use phased rollout.

### Phase 1 — Pull request validation only

Start with:

- `pr-to-main.yml`
- `pr-to-develop.yml`

This gives you immediate value with minimal disruption.

### Phase 2 — Branch validation

Add:

- `release-branch.yml`
- `hotfix-branch.yml`
- `sync-branch.yml`

This helps enforce branch naming and branch-specific rules.

### Phase 3 — Optional work branch validation

Add:

- `wip-branch.yml`

Use this only if your team wants to formalize temporary work branches.

---

## Helper Workflows (Optional)

To reduce friction and naming mistakes, Stability Flow can also include helper workflows for creating standard branches.

### Available helper workflows

- `create-release-branch.yml`
- `create-hotfix-branch.yml`
- `create-sync-branch.yml`

These are convenience workflows and are **not required** to use Stability Flow.

They are most useful for teams that want:

- consistent branch naming
- lower onboarding friction
- fewer manual branch creation mistakes

---

## GitHub Rulesets

Stability Flow includes optional GitHub rulesets that can be pasted or adapted for your repository.

### Included rulesets

- `main.ruleset.json`
- `develop.ruleset.json`
- `release.ruleset.json`
- `hotfix.ruleset.json`
- `sync.ruleset.json`
- `wip.ruleset.json`
- `tag.ruleset.json`

These rulesets are designed to align with the included workflows and required status checks.

---

## Recommended Ruleset Usage

### Required / strongly recommended

- `main.ruleset.json`
- `develop.ruleset.json`

These protect the two most important branches.

### Recommended if you use release/hotfix/sync branches

- `release.ruleset.json`
- `hotfix.ruleset.json`
- `sync.ruleset.json`

### Optional

- `wip.ruleset.json`
- `tag.ruleset.json`

---

## Suggested Adoption Path

If you want the simplest practical rollout:

### Minimal

- use `pr-to-main.yml`
- use `pr-to-develop.yml`

### Strong

- use `stability-flow.yml`
- add `main.ruleset.json`
- add `develop.ruleset.json`

### Strict

- use `stability-flow.yml`
- add all branch rulesets
- use helper branch-creation workflows

---

## What to Do Next

Once your repository is wired up:

1. create or confirm `main` and `develop`
2. protect `main`
3. test a `release/*` pull request
4. test a `hotfix/*` pull request
5. verify reconciliation using a `sync/*` branch

Then read:

- [Specification](spec.md)
- [Release Flow](release-flow.md)
- [Enforcement](enforcement.md)
