# Contributing to Stability Flow

Thanks for your interest in contributing to Stability Flow.

This project is structured around a simple idea:

> Stability Flow is a **branching strategy specification first**, with optional **reference tooling and integrations** built around it.

That distinction matters.

Please try to preserve it when contributing.

---

## Ways to Contribute

Contributions are welcome in areas such as:

- specification clarity
- examples and diagrams
- design rationale
- enforcement guidance
- reference tooling
- CI integrations
- GitHub Actions
- reusable workflows
- documentation quality
- bug fixes and polish

You do not need to contribute code to contribute meaningfully.

Clear docs and good examples are high-value contributions here.

---

## Project Structure

The repository is intentionally structured to separate the **specification** from the **tooling**.

### Specification and public docs

```text
docs/
```

This area contains the public Stability Flow documentation, including:

- the specification
- design rationale
- release examples
- enforcement guidance

This is where content about **the standard itself** belongs.

---

### Tooling and implementation docs

```text
docs/tools/
```

This area contains documentation for tooling and integrations built around the specification.

Examples include:

- CLI validator docs
- GitHub Actions docs
- reusable workflow docs

This is where content about **how a specific implementation works** belongs.

---

### Reference implementations

```text
tools/
```

This area contains tooling built to support the specification.

Examples may include:

- validators
- CI support tooling
- future reference implementations

---

### Supporting project files

Other top-level areas may include:

- `scripts/` for local helper and demo scripts
- `docker/` for publishable container artifacts
- root docs such as `README.md`, `LICENSE`, and project metadata

---

## The Most Important Contribution Rule

When making changes, ask yourself:

> “Am I changing the specification, or am I changing a tool?”

This is the most important distinction in the repository.

---

## If You Are Changing the Specification

Specification changes affect the Stability Flow model itself.

Examples include changes to:

- branch roles
- allowed branch origins
- allowed merge targets
- release behavior
- hotfix behavior
- reintegration behavior
- normative wording such as `MUST`, `SHOULD`, and `MAY`

These changes usually belong in:

- `docs/spec.md`
- and sometimes also:

  - `docs/design.md`
  - `docs/release-flow.md`
  - `docs/enforcement.md`

Specification changes should be made carefully and intentionally.

---

## If You Are Changing Tooling

Tooling changes affect one implementation of the spec, not the spec itself.

Examples include changes to:

- CLI commands
- output formats
- GitHub Actions
- reusable workflows
- scripts
- container images
- CI behavior

These changes usually belong in:

- `docs/tools/`
- `tools/`
- `scripts/`
- `docker/`

Tooling should support the specification, not redefine it.

---

## Keep the Spec Tool-Neutral

Please avoid turning the specification into tool-specific documentation.

### Bad pattern

> “Use `validate-merge` to ensure a release branch can merge into main.”

### Better pattern

> “Merge eligibility should be validated before protected branch integration.”

The first belongs in tooling docs.
The second belongs in the specification or enforcement guidance.

This distinction is important to the project.

---

## Writing Guidelines

Please prefer writing that is:

- clear
- direct
- explicit
- practical
- easy to validate against

Please avoid writing that is:

- overly conversational
- redundant across multiple docs
- implementation-specific in spec-level docs
- unnecessarily abstract

Good docs here should be understandable to both:

- humans reading the standard
- people building tooling around it

---

## Diagrams and Examples

Examples and diagrams are encouraged.

They are especially useful in:

- `docs/release-flow.md`
- `docs/design.md`
- tooling and integration docs

When contributing examples:

- prefer realistic branch names
- prefer realistic release examples
- keep examples aligned with the current specification
- avoid stale or contradictory diagrams

If the spec changes, examples should be updated too.

---

## Tooling Contributions

Tooling contributions are welcome, especially when they help make Stability Flow easier to adopt or enforce.

Good tooling contributions usually:

- align clearly with the specification
- do one thing well
- avoid overreaching into “workflow orchestration”
- remain understandable to users

Reference tooling should help teams:

- validate the flow
- adopt the flow
- automate around the flow

It should not become the definition of the flow.

---

## Pull Request Guidance

When opening a pull request, it helps to make the intent explicit.

A good pull request should make it clear whether it is primarily:

- a spec change
- a docs clarification
- a tooling change
- a bug fix
- an example / diagram improvement

This helps reviewers evaluate the change correctly.

---

## Suggested Pull Request Checklist

Before opening a PR, ask:

### Scope

- Does this change affect the specification or only an implementation?
- Is the change in the right part of the repo?

### Docs

- Did I update any affected documentation?
- Did I accidentally duplicate content that already exists elsewhere?

### Consistency

- Does this still align with the current Stability Flow model?
- Do any examples or diagrams need updating?

### Tooling

- If I changed tooling, does it still match the current spec?
- If I changed the spec, do any tools now need updating?

---

## Local Validation and Testing

If you are changing tooling or scripts, please run the relevant local checks before opening a PR.

Examples may include:

- validator tests
- script smoke tests
- docs build checks
- container build checks

The exact commands may evolve over time as tooling grows.

If a command or workflow is unclear, improving the docs for it is a valid contribution too.

---

## Small Contributions Are Welcome

Not every useful contribution needs to be a major change.

Helpful small contributions include:

- fixing unclear wording
- correcting diagrams
- improving examples
- clarifying edge cases
- fixing typos
- tightening docs structure

These often improve the project more than large but unfocused changes.

---

## Questions and Discussions

If you are unsure where a change belongs, the safest rule is:

> spec-level behavior belongs in `docs/`, implementation behavior belongs in `docs/tools/` and `tools/`.

That single rule will prevent most structural drift.

---

## Thank You

Thanks for helping improve Stability Flow.

The project gets stronger when contributors help keep it:

- clear
- explicit
- enforceable
- implementation-friendly
