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
- `internal/` — private implementation details that must never be part of the public API

## Dependency rules

| Layer             | May import                        | Must never import                      |
| ----------------- | --------------------------------- | -------------------------------------- |
| `domain/`         | other `domain/`, stdlib           | `usecase`, `adapter`, `infrastructure` |
| `usecase/`        | `domain/`, stdlib                 | `adapter`, `infrastructure`            |
| `adapter/`        | `domain/`, `usecase/`             | `infrastructure`                       |
| `infrastructure/` | `domain/`, `usecase/`, `adapter/` | —                                      |

## When creating new code

- Identify the layer before writing — ask: is this a type/interface (domain), business logic (usecase), translation (adapter), or external integration (infrastructure)?
- New port interfaces always belong in `domain/`; their implementations belong in `adapter/` or `infrastructure/`
- No framework or third-party imports in `domain/` or `usecase/` — they become transitive dependencies for all consumers
- Test doubles for domain ports live in `adapter/` and are only imported by test files

## When reviewing or modifying existing code

- Flag any import that violates the dependency rules above
- Flag any third-party framework import in `domain/` or `usecase/`
- Flag adapter test doubles imported outside of test files
