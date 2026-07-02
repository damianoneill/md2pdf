---
name: new-domain-concept
description: Scaffold a new domain concept (port interface + test double) following the established Clean Architecture pattern. Use when adding a new concern to the domain layer.
argument-hint: "<package-name> <description>"
allowed-tools: Read, Write, Glob
---

# Scaffold New Domain Concept: $ARGUMENTS

## Existing domain packages (for reference)
Use the Glob tool to list existing packages under `domain/`.

## Existing port + fake pattern (for reference)
Read existing domain interfaces and their corresponding fakes before scaffolding.

---

The first argument is the package/module name, the rest describes the concept.

Create the following files following the established pattern in this project:

### `domain/<name>/<name>.<ext>`
- Module/package comment describing the concept
- The port interface with 1–3 focused methods
- Any supporting types (structs, enums) needed by the interface
- Zero imports outside the standard library and other `domain/` packages

### `fake/<name>.<ext>`
- A test double implementing the port interface
- Thread-safe where applicable
- Fields to capture calls for assertion in tests (e.g. `call_count`, `last_input`)
- Configurable behaviour via function/callable fields

### `domain/<name>/test_<name>.<ext>` (or equivalent test file)
- Tests for any non-trivial logic on the types themselves (if any)

### Update `README.md`
- Add the package to the structure listing
- Add an interface definition and a one-line usage note
- Add the fake to the test doubles listing

Keep interfaces small. If unsure whether a method belongs, leave it out — it can be added when there is a concrete use case.
