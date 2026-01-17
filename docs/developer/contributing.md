# Contributing

Guidelines for contributing to ELMOS.

## Workflow

1. Fork the repo
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make changes, add tests
4. Run checks: `task dev:check`
5. Commit: `git commit -m "feat: add my feature"`
6. Push and PR

## Code Style

- **Go**: Follow `gofmt`, `goimports`
- **Commits**: Conventional commits (e.g., `feat:`, `fix:`, `docs:`)
- **PRs**: Descriptive, link issues
- **Tests**: Required for new code

## Development Setup

```bash
git clone https://github.com/NguyenTrongPhuc552003/elmos.git
cd elmos
task dev:setup  # Install deps, build
```

## Testing

- `task test` - Run tests
- `task test:cover` - With coverage
- Add mocks for new infra interfaces

## Documentation

- Update docs for user-facing changes
- Add API docs for new interfaces

## Reviews

- At least one maintainer review
- CI must pass
- Squash merges

## Issues

- Bug reports with reproduction steps
- Feature requests with use cases

## License

Contributions under MIT license.