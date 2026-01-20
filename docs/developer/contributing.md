# Contributing

Guidelines for contributing to ELMOS.

---

## Quick Start

```bash
# Clone
git clone https://github.com/NguyenTrongPhuc552003/elmos.git
cd elmos

# Setup
task dev:setup

# Develop
task dev:check   # Format + lint
task test        # Run tests
task build       # Build binary
```

---

## Commit Convention

Use the project's commit message format:

```
<scope>: <Title>

- Change 1
- Change 2

Signed-off-by: Your Name <email@example.com>
```

**Scope examples:**

| Scope                   | Example                |
| ----------------------- | ---------------------- |
| `core: domain: builder` | Kernel builder changes |
| `core: app: commands`   | CLI command changes    |
| `docs: user: qemu`      | User doc updates       |
| `.github: workflows`    | CI changes             |
| `Taskfile`              | Build task changes     |

**Examples:**

```bash
git commit -sm "core: domain: emulator: Add machine validation

- Add ValidateMachine function
- Improve error messages"

git commit -sm "docs: user: kernel-building: Add troubleshooting table

- Add common issues and solutions"
```

---

## Pull Request Process

1. **Branch from `main`**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make changes with tests**
   ```bash
   task dev:check
   task test
   ```

3. **Push and create PR**
   ```bash
   git push origin feature/my-feature
   ```

4. **PR requirements:**
   - Descriptive title
   - Link related issues
   - CI must pass
   - One maintainer approval

---

## Code Standards

### Go Style

- Run `gofmt` and `goimports`
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use interfaces for testability
- Constructor pattern: `NewXxx(deps...) *Xxx`

### Testing

- Table-driven tests
- Mock infra interfaces
- Aim for 80%+ coverage on new code

```bash
task test:cover
```

---

## Documentation

- Update user docs for CLI changes
- Update developer docs for API changes
- Add code comments for exported types/functions

---

## Issue Guidelines

### Bug Reports

Include:

- ELMOS version (`elmos version`)
- macOS version
- Steps to reproduce
- Expected vs actual behavior

### Feature Requests

Include:

- Use case description
- Proposed solution
- Alternatives considered

---

## License

All contributions are under MIT license.