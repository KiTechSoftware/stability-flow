# Maintainer Notes

## 1. Overview

This document provides guidance for maintaining the Stability Flow repository.

It exists to help keep the project consistent as it evolves.

Stability Flow is structured as:

1. a **public branching strategy specification**
2. supporting **public documentation**
3. optional **reference tooling and integrations**

Maintainers should preserve that separation.

---

## 2. Project Structure Philosophy

A core principle of this repository is:

> the specification is the primary artifact; tooling is secondary.

This means:

- the branching model is the product
- tooling exists to help people adopt or enforce it
- tooling must not redefine the specification

This distinction should remain clear in both documentation and implementation.

---

## 3. Documentation Boundaries

The documentation is intentionally split into two areas.

### 3.1 `docs/`

The root of `docs/` is for **specification-level and public concept documentation**.

These documents describe the Stability Flow model itself.

Examples:

- `index.md`
- `spec.md`
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
- how it can be enforced in principle

They should **not** become tool manuals.

---

### 3.2 `docs/tools/`

The `docs/tools/` directory is for **tooling and implementation documentation**.

These documents describe tooling built to support the Stability Flow specification.

Examples:

- CLI validator documentation
- GitHub Actions documentation
- reusable workflow documentation
- future integration docs

These documents may describe:

- commands
- flags
- outputs
- usage examples
- CI integration
- automation behavior

They should not redefine the spec.

---

## 4. Keep the Spec Tool-Neutral

The specification must remain independent of any single implementation.

That means the following should **not** appear as normative language in the spec-level docs:

- CLI command names
- GitHub Action names
- reusable workflow assumptions
- CI vendor-specific requirements
- tool-specific flags or invocation examples

### Bad pattern

> “Use `validate-merge` to check whether a release branch may merge into main.”

### Better pattern

> “Merge eligibility should be validated before protected branch integration.”

The CLI or workflow docs can then explain how to do that with specific tooling.

This distinction is important.

---

## 5. Where New Documentation Should Go

When adding new documentation, maintainers should decide whether it belongs to the:

- **specification**
- **design rationale**
- **enforcement guidance**
- **tooling / implementation docs**

### Use `docs/` when the content answers:
- what is the rule?
- why does the rule exist?
- how does the flow work?
- how can this be enforced in principle?

### Use `docs/tools/` when the content answers:
- how does this tool work?
- how do I run this validator?
- how do I use this GitHub Action?
- how do I integrate this workflow?

This is the most important documentation boundary in the project.

---

## 6. Specification Changes vs Tool Changes

Not every tooling change is a specification change.

Not every specification change requires a tooling change.

Maintainers should treat these as separate concerns.

---

### 6.1 A Specification Change Usually Means

A specification change typically affects one or more of:

- allowed branch roles
- allowed branch origins
- allowed merge targets
- branch naming rules
- release behavior
- hotfix behavior
- reintegration behavior
- normative wording (`MUST`, `SHOULD`, etc.)

These changes should usually update:

- `docs/spec.md`
- possibly `docs/design.md`
- possibly `docs/release-flow.md`
- possibly `docs/enforcement.md`

They may also require tooling changes, but tooling should follow the spec, not lead it.

---

### 6.2 A Tooling Change Usually Means

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

## 7. Keep Reference Implementations Honest

Tooling in this repository should be presented as:

# reference implementations

That means maintainers should avoid language that implies:

- the tool is the definition of the standard
- the workflow requires this repository’s tooling
- Stability Flow is inseparable from one CLI or CI implementation

### Preferred framing

- “reference validator”
- “reference tooling”
- “example GitHub Actions integration”
- “example reusable workflow”

This keeps the project open and implementation-friendly.

---

## 8. Repository Hygiene Expectations

Maintainers should aim to keep the repository:

- easy to navigate
- explicit in purpose
- low in duplication
- clear in ownership of concepts

### Recommended hygiene

- avoid duplicate explanations across docs
- keep examples aligned with the current spec
- avoid stale workflow diagrams
- keep implementation docs separate from standard docs
- prefer one canonical explanation of each concept

---

## 9. Scripts and Tooling Boundaries

Scripts and tooling should follow the same separation principles as the docs.

### Good separation

- build scripts build artifacts
- run scripts run artifacts
- test scripts test artifacts
- flow demo scripts demonstrate the branching model
- validators validate policy

### Avoid

- mixing documentation generation concerns into validator tooling
- making example scripts appear normative
- embedding implementation assumptions into the spec

This keeps the project easier to maintain.

---

## 10. Documentation Maintenance Guidance

When updating docs:

### Prefer
- short, explicit sections
- consistent terminology
- examples that reflect the current spec
- diagrams that match real workflow behavior

### Avoid
- over-explaining obvious Git concepts
- mixing policy with implementation
- conversational drift from earlier drafting
- “temporary” notes left in public docs

If a concept already has a clear home, update that document instead of repeating it elsewhere.

---

## 11. Versioning and Maturity

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

## 12. Recommended Maintainer Checklist

When making changes, ask:

### Specification questions
- does this change alter the rules of Stability Flow?
- does it change branch behavior or release behavior?
- does it need normative wording updates?

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

## 13. Summary

The most important maintainer rule for this repository is:

> keep the specification separate from the tooling.

If that boundary remains clear, the project stays coherent.

If that boundary blurs, the project becomes harder to trust, harder to adopt, and harder to maintain.

That separation is one of the most important design decisions in the repository.
