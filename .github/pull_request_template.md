## What

A short description of the change.

## Why

The problem it solves. Link the issue if there is one (for example
`Fixes #123`).

## How it was tested

Commands you ran and what you observed (for a fix, the test that fails on
the old code).

## Checklist

- [ ] `gofmt -l $(git ls-files '*.go')` prints nothing
- [ ] `go vet ./...` and `go test ./...` pass
- [ ] `go build -o build/ ./cmd/elmos` succeeds
- [ ] Docs updated (`docs/user`, `docs/developer`, README) if behavior changed
- [ ] Commit messages use a conventional prefix (`feat:`, `fix:`, `docs:`, `chore:`)
