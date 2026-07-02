# Project Name

One-line description of what this project does.

## Overview

<!-- Describe the project purpose and key capabilities -->

## Getting Started

```bash
make setup   # install dependencies
make test    # run tests
make build   # compile / bundle
make ci      # full pipeline: lint → test → build
```

## Project Structure

```
domain/          # core types and port interfaces
usecase/         # business logic
adapter/         # format translation and port implementations
infrastructure/  # external integrations (databases, APIs, queues)
fake/            # test doubles — for use in tests only
internal/        # private implementation details
examples/        # runnable examples
```

See [CLAUDE.md](CLAUDE.md) for architecture rules and code style guidelines.

## Configuration

Copy `.env.example` to `.env` and fill in your values. Never commit `.env`.
