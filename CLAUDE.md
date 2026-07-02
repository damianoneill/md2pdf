# Project Name

One-line description of what this project does.

## Technology Stack

- Language: <!-- e.g. Go 1.22+, Python 3.12+, Node 20+ -->
- Key dependencies: <!-- list them here -->
- Testing: <!-- e.g. make test, pytest, go test -->
- Linting: <!-- e.g. make lint, ruff, golangci-lint -->

## Build & Test Commands

```bash
# Fill in language-specific commands
make setup      # install dependencies and tools
make ci         # full pipeline: lint → test → build
make test       # run tests
make coverage   # tests with coverage report
make build      # compile / bundle
make fmt        # format source
make lint       # run linters
make clean      # remove build artifacts
```

## Architecture

This project follows **Clean Architecture**. Dependencies flow inward only — outer layers depend on inner layers, never the reverse.

```
domain/
├── <concept>/   # Port interfaces and core types — no external dependencies
usecase/
└── <feature>/   # Business logic — orchestrates domain interfaces only
adapter/         # Translates between external formats and domain types
infrastructure/  # External integrations (databases, APIs, queues)
fake/            # Test doubles for all domain ports — test use only
internal/        # Private implementation details
examples/        # Runnable examples demonstrating usage
```

### Dependency Rules

| Layer | May depend on | Must not depend on |
|-------|--------------|-------------------|
| `domain/` | other `domain/`, stdlib | `usecase`, `adapter`, `infrastructure`, `fake` |
| `usecase/` | `domain/`, stdlib | `adapter`, `infrastructure`, `fake` |
| `adapter/` | `domain/`, `usecase/` | `infrastructure`, `fake` |
| `infrastructure/` | `domain/`, `usecase/`, `adapter/` | `fake` |
| `fake/` | anything | — |

- **No framework or third-party imports** in `domain/` or `usecase/`
- **`fake/` is for testing only** — never imported in production code

## Code Style Guidelines

- Follow the conventions of the language in use
- Prefer explicit error handling
- Keep interfaces small and focused (prefer 1–3 methods)
- Add comments only where logic is non-obvious
- Table-driven / parameterised tests for any non-trivial logic

## KISS — Keep It Simple

This is the most important principle. The design must be proportionate to the problem.

- Solve the stated requirement only — do not design for hypothetical future requirements
- Choose the simplest design that satisfies the behaviour; introduce patterns only when justified by a current need
- Do not add layers, indirection, or extensibility points speculatively
- If the approach feels more complex than the problem, stop and reconsider
- **Before proposing a new port or abstraction, ask: does something already in the project cover this?**

## Important Notes

- Never commit API keys or secrets — use environment variables
- Run the full CI pipeline before committing
- Use conventional commits (feat:, fix:, docs:, chore:)
