---
paths:
  - "**/*"
---

# Architecture Rules

This project follows **Clean Architecture**. Dependencies flow inward only — outer layers depend on inner layers, never the reverse.

## Layer responsibilities

- `domain/` — core types and port interfaces; zero external dependencies; the contract everything else builds against
- `usecase/` — business logic; orchestrates domain interfaces; no framework imports
- `adapter/` — translates between external formats and domain types; implements domain port interfaces; added as needed
- `infrastructure/` — external integrations (databases, APIs, queues); implements domain port interfaces; added as needed
- `fake/` — test doubles for all domain ports; for use in tests only; never imported in production code
- `internal/` — private implementation details that must never be part of the public API

## Dependency rules

| Layer | May import | Must never import |
|-------|-----------|-------------------|
| `domain/` | other `domain/`, stdlib | `usecase`, `adapter`, `infrastructure`, `fake` |
| `usecase/` | `domain/`, stdlib | `adapter`, `infrastructure`, `fake` |
| `adapter/` | `domain/`, `usecase/` | `infrastructure`, `fake` |
| `infrastructure/` | `domain/`, `usecase/`, `adapter/` | `fake` |
| `fake/` | anything | — |

## When creating new code

- Identify the layer before writing — ask: is this a type/interface (domain), business logic (usecase), translation (adapter), or external integration (infrastructure)?
- New port interfaces always belong in `domain/`; their implementations belong in `adapter/` or `infrastructure/`
- No framework or third-party imports in `domain/` or `usecase/` — they become transitive dependencies for all consumers
- `fake/` implementations live in `fake/`, not alongside their interfaces

## When reviewing or modifying existing code

- Flag any import that violates the dependency rules above
- Flag any third-party framework import in `domain/` or `usecase/`
- Flag `fake/` imported outside of test files or `examples/`
