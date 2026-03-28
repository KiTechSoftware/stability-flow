# Maintainer Guide

## 1. Purpose

This document provides guidance for maintaining the Stability Flow repository.

Its purpose is to help preserve:

- clear documentation boundaries
- a stable public specification
- clean separation between the model and its implementations

Stability Flow is structured as:

1. a **public branching strategy specification**
2. supporting **public documentation**
3. optional **reference tooling and integrations**

Maintainers should preserve that separation as the project evolves.

---

## 2. Repository Philosophy

A core principle of this repository is:

> the specification is the primary artifact; tooling is secondary.

This means:

- the branching model is the product
- documentation exists to explain the model
- tooling exists to help teams adopt or validate the model
- tooling must not redefine the model

This distinction should remain clear in both documentation and implementation.

---

## 3. Documentation Structure

The documentation is intentionally split by responsibility.

### 3.1 Public Documentation

The root of `docs/` is for **specification-level and public concept documentation**.

These documents describe Stability Flow itself.

Examples:

- `index.md`
- `spec.md`
- `conventions.md`
- `design.md`
- `release-flow.md`
- `enforcement.md`

These documents should remain:

- tool-neutral
- implementation-neutral
- vendor-neutral where practical

They should describe:

- what Stability Flow is
- how it works
- why it is shaped this way
- how it can be validated in principle

They should **not** become tool manuals.

---

### 3.2 Tooling Documentation

The `docs/tools/` directory is for **tooling and implementation documentation**.

These documents describe tools built to support Stability Flow.

Examples:

- CLI validator documentation
- GitHub Actions documentation
- reusable workflow documentation
- integration docs

These documents may describe:

- commands
- flags
- outputs
- usage examples
- CI integration
- automation behavior

They should not redefine the spec.

---

## 4. Documentation Boundaries

Every important concept in the project should have **one canonical home**.

That prevents drift and duplicate explanations.

Use the following ownership model:

| Concept                              | Canonical Home    |
| ------------------------------------ | ----------------- |
| branching rules                      | `spec.md`         |
| naming and commit conventions        | `conventions.md`  |
| rationale and tradeoffs              | `design.md`       |
| validation and enforcement surfaces  | `enforcement.md`  |
| worked examples and branch histories | `release-flow.md` |
| tooling usage and integration        | `docs/tools/*`    |

If a concept already has a clear home, update that document instead of repeating it elsewhere.

---

## 5. Keep the Specification Tool-Neutral

The specification must remain independent of any single implementation.

That means the following should **not** appear as normative language in spec-level docs:

- CLI command names
- GitHub Action names
- reusable workflow assumptions
- CI vendor-specific requirements
- tool-specific flags or invocation examples

### Bad pattern

> “Use `validate-merge` to check whether a release branch may merge into `main`.”

### Better pattern

> “Merge eligibility should be validated before protected branch integration.”

Implementation-specific docs can then explain how to do that with actual tooling.

This distinction is important.

---

## 6. Where New Content Should Go

When adding new content, maintainers should first decide what kind of content it is.

### Add to `spec.md` if it answers

- what is the rule?
- what branch movement is allowed?
- what merge behavior is required?

### Add to `conventions.md` if it answers

- how should branches be named?
- how should final squash commits be formatted?
- how should breaking changes or reverts be represented?

### Add to `design.md` if it answers

- why does this rule exist?
- what tradeoff is the model making?
- why is the workflow shaped this way?

### Add to `enforcement.md` if it answers

- what should be validated?
- where can this be enforced?
- how enforceable is this rule in practice?

### Add to `release-flow.md` if it answers

- what does this workflow look like in practice?
- how does a release or hotfix actually move?

### Add to `docs/tools/` if it answers

- how does this tool work?
- how do I run this validator?
- how do I integrate this workflow?

This is the most important documentation discipline in the repository.

---

## 7. Specification Changes vs Tool Changes

Not every tooling change is a specification change.

Not every specification change requires a tooling change.

Maintainers should treat these as separate concerns.

### 7.1 A Specification Change Usually Means

A specification change typically affects one or more of:

- branch roles
- branch origins
- merge targets
- merge strategies
- release behavior
- hotfix behavior
- reconciliation behavior
- normative wording (`MUST`, `SHOULD`, etc.)

These changes should usually update one or more of:

- `docs/spec.md`
- `docs/conventions.md`
- `docs/design.md`
- `docs/release-flow.md`
- `docs/enforcement.md`

Tooling should follow the specification, not lead it.

---

### 7.2 A Tooling Change Usually Means

A tooling change typically affects one or more of:

- command behavior
- output formats
- CI integration
- GitHub Actions support
- reusable workflow support
- CLI flags
- implementation details

These changes should usually update:

- `docs/tools/*`
- tool source code
- examples and scripts

They do not automatically imply a spec change.

---

## 8. Keep Reference Implementations Honest

Tooling in this repository should be presented as:

> reference implementations

That means maintainers should avoid language that implies:

- the tool is the definition of the standard
- the workflow requires this repository’s tooling
- Stability Flow is inseparable from one CLI or CI implementation

Preferred framing includes:

- reference validator
- reference tooling
- example GitHub Actions integration
- example reusable workflow

This keeps the project open and implementation-friendly.

---

## 9. Repository Hygiene

Maintainers should aim to keep the repository:

- easy to navigate
- explicit in purpose
- low in duplication
- clear in concept ownership

### Recommended hygiene

- avoid duplicate explanations across docs
- keep examples aligned with the current spec
- avoid stale workflow diagrams
- keep implementation docs separate from standard docs
- prefer one canonical explanation of each concept
- ensure workflow diagrams and examples reflect the canonical release and reconciliation path

---

## 10. Scripts and Tooling Boundaries

Scripts and tooling should follow the same separation principles as the docs.

### Good separation

- build scripts build artifacts
- run scripts run artifacts
- test scripts test artifacts
- flow demo scripts demonstrate the branching model
- validators validate policy

### Avoid in Tooling

- mixing documentation generation concerns into validator tooling
- making example scripts appear normative
- embedding implementation assumptions into the spec

This keeps the project easier to maintain.

---

## 11. Documentation Maintenance Guidance

When updating docs:

### Prefer

- short, explicit sections
- consistent terminology
- examples that reflect the current spec
- diagrams that match real workflow behavior

### Avoid

- over-explaining obvious Git concepts
- mixing policy with implementation
- conversational drafting leftovers
- “temporary” notes left in public docs

If a concept already has a clear home, update that document instead of repeating it elsewhere.

---

## 12. Versioning and Maturity

As the project evolves, maintainers should distinguish clearly between:

- specification maturity
- tooling maturity
- integration maturity

Example:

- the spec may be stable at `v1`
- the CLI may still be evolving
- GitHub Actions support may still be growing

That is acceptable.

The maturity of one part of the project does not need to be artificially tied to the others.

---

## 13. Maintainer Checklist

When making changes, ask:

### Specification questions

- does this change alter the rules of Stability Flow?
- does it change branch behavior or release behavior?
- does it require normative wording updates?

### Documentation questions

- is this the right document for this content?
- am I duplicating an explanation that already exists?
- is this still tool-neutral where it should be?

### Tooling questions

- is this a reference implementation concern or a spec concern?
- does this belong in `docs/tools/` instead?
- does the implementation still match the spec?

This simple discipline will prevent most structural drift.

---

## 14. Summary

The most important maintainer rule for this repository is:

> keep the specification separate from the tooling.

If that boundary remains clear, the project stays coherent.

If that boundary blurs, the project becomes harder to trust, harder to adopt, and harder to maintain.

That separation is one of the most important design decisions in the repository.
