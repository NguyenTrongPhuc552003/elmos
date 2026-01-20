# Package Version API

Version information management.

## Info Struct

```go
type Info struct {
    Version   string `json:"version"`
    Commit    string `json:"commit"`
    BuildDate string `json:"buildDate"`
    GoVersion string `json:"goVersion"`
    OS        string `json:"os"`
    Arch      string `json:"arch"`
}
```

Holds build and runtime info.

## Functions

### Get

```go
func Get() Info
```

Returns current version info, populated via ldflags.

### String / Short

```go
func (i Info) String() string
func (i Info) Short() string
```

Formatted output for CLI.

## Build Process

Version info set at build time:

```bash
go build -ldflags "-X 'github.com/NguyenTrongPhuc552003/elmos/core/app/version.Version=1.0.0' -X 'github.com/NguyenTrongPhuc552003/elmos/core/app/version.Commit=$(git rev-parse HEAD)' -X 'github.com/NguyenTrongPhuc552003/elmos/core/app/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" -o build/elmos ./cmd/elmos
```

Used in `Taskfile.yml`.