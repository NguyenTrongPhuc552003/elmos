# Architecture

ELMOS follows a layered, domain-driven architecture in Go.

## Overview

```
┌─────────────────┐
│   App Layer     │  CLI commands, TUI
├─────────────────┤
│  Domain Layer   │  Business logic
├─────────────────┤
│  Infra Layer    │  External interfaces
└─────────────────┘
```

## Layers

### App Layer (`core/app/`)

- **Purpose**: Command-line interface and user interactions
- **Components**:
  - `commands/`: Individual CLI commands (e.g., `kernel.go`, `toolchain.go`)
  - `app.go`: App wiring, dependency injection
  - `helpers.go`: Shared utilities
  - `version/`: Version command
- **Dependencies**: Domain layer

### Domain Layer (`core/domain/`)

- **Purpose**: Core business logic, independent of frameworks
- **Modules**:
  - `builder/`: Kernel, module, app builders
  - `doctor/`: Environment checks and fixes
  - `emulator/`: QEMU options and execution
  - `patch/`: Kernel patch management
  - `rootfs/`: RootFS creation (debootstrap)
  - `toolchain/`: Toolchain management
- **Principles**: Pure functions, interfaces for testing

### Infra Layer (`core/infra/`)

- **Purpose**: External dependencies and system interactions
- **Components**:
  - `executor/`: Command execution (shell, mock)
  - `filesystem/`: File operations (os, mock)
  - `homebrew/`: Homebrew resolver
- **Interfaces**: Allow mocking for tests

## Key Patterns

- **Dependency Injection**: App layer injects infra into domain
- **Interfaces**: Domain defines interfaces, infra implements
- **Embedded Assets**: Templates in `assets/` for code generation
- **Configuration**: YAML-based config in `core/config/`
- **Errors**: Custom errors in `core/context/errors.go`

## Data Flow

1. CLI parses args → App layer
2. App calls domain logic with config
3. Domain uses infra interfaces for I/O
4. Results returned via channels or callbacks

## Testing

- Unit tests in each package
- Mocks for infra interfaces
- Integration tests via `task test`

## Extensibility

- Add commands in `core/app/commands/`
- Extend domain with new modules
- Implement infra interfaces for new platforms