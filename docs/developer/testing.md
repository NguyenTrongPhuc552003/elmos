# Testing

ELMOS uses Go's testing framework with mocks for infra.

## Structure

- Unit tests in `*_test.go` files
- Mocks in `core/infra/*/mock.go`
- Table-driven tests for commands

## Running Tests

```bash
task test          # Run all tests
task test:cover    # With coverage report
```

## Mocking

Infra interfaces mocked for domain testing:

```go
type MockExecutor struct {
    RunFunc func(cmd string, args ...string) error
    // ...
}

func (m *MockExecutor) Run(cmd string, args ...string) error {
    return m.RunFunc(cmd, args...)
}
```

Used in domain tests to verify calls.

## Coverage

Target: 80%+ coverage. Generated via `go test -cover`.

## Integration Tests

Limited; focus on unit tests. Manual testing for full workflows.

## CI

Tests run on PRs and pushes.

## Best Practices

- Test public APIs
- Mock external deps
- Use `assert` for checks
- Table-driven for multiple cases