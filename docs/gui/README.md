# ELMOS GUI - Native macOS Application

## Overview

ELMOS GUI is a native macOS application built with Swift and SwiftUI, providing a graphical interface for the ELMOS (Embedded Linux on MacOS) toolchain.

## Architecture

### Technology Stack
- **Language**: Swift 6.2
- **Framework**: SwiftUI (macOS 14+)
- **Build**: Swift Package Manager
- **Integration**: Go CLI spawns Swift binary

### Structure
```
core/ui/gui/
├── Package.swift                  # SPM manifest
├── Sources/
│   └── ElmosGUI/
│       └── ElmosGUI.swift        # Main app + all views
└── .build/                        # SPM build output
```

## Building

### Build Swift GUI
```bash
task gui:build
```

This will:
1. Build Swift app in release mode
2. Copy binary to `build/gui/elmos`
3. Run verification checks

### Clean Build
```bash
task gui:clean
```

## Running

### Launch from CLI
```bash
./build/elmos gui
```

### Direct Launch (development)
```bash
task gui:run
```

## Development

### How It Works

1. **Go CLI Integration**: The `gui` subcommand is registered in `core/app/commands/registry.go`
2. **Command Execution**: `gui.go` spawns the Swift binary at `build/gui/elmos`
3. **Swift App**: Native macOS app built with SwiftUI
4. **CLI Integration**: Swift app can execute `./build/elmos` commands

### Communication Flow
```
User → ./build/elmos gui → build/gui/elmos (Swift)
                                 ↓
                        SwiftUI GUI launches
                                 ↓
                    Executes CLI commands ← ./build/elmos [cmd]
```

### Adding New Views

1. Edit `core/ui/gui/Sources/ElmosGUI/ElmosGUI.swift`
2. Add new view structs
3. Rebuild: `task gui:build`

### Calling CLI from Swift

```swift
import Foundation

func runElmosCommand(_ args: [String]) -> String? {
    let process = Process()
    process.executableURL = URL(fileURLWithPath: "./build/elmos")
    process.arguments = args
    
    let pipe = Pipe()
    process.standardOutput = pipe
    
    try? process.run()
    process.waitUntilExit()
    
    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    return String(data: data, encoding: .utf8)
}

// Usage
let status = runElmosCommand(["status"])
```

## Features

### Current Views
- **Dashboard**: Workspace overview
- **Toolchains**: Manage cross-compilers
- **Kernel**: Build configuration
- **Modules**: Kernel modules
- **QEMU**: Emulator control
- **Settings**: App preferences

### Planned Features
- [ ] Real-time command execution
- [ ] Build progress indicators
- [ ] Log viewing
- [ ] Configuration editor
- [ ] YAML file management

## Troubleshooting

### GUI binary not found
```bash
# Rebuild GUI
task gui:build

# Check binary exists
ls -lh build/gui/elmos
```

### Swift build fails
```bash
# Clean and rebuild
task gui:clean
task gui:build
```

### Command not showing
```bash
# Rebuild Go CLI to register command
task build

# Verify command exists
./build/elmos --help | grep gui
```

## Design Patterns

### ELMOS Conventions Followed
- ✅ Command in `core/app/commands/`
- ✅ UI in `core/ui/gui/`
- ✅ Binary naming: `elmos` (no hyphens)
- ✅ Taskfile integration with `gui:*` tasks
- ✅ Follows existing patterns (similar to `tui`)

### Swift Conventions
- SwiftUI for declarative UI
- macOS 14+ minimum target
- No external dependencies
- Single executable output

## Benefits

### Why Swift instead of Qt/C++?
1. **Native**: First-class macOS support
2. **Modern**: SwiftUI is declarative and intuitive
3. **No Dependencies**: No Homebrew Qt issues
4. **Performance**: Compiled, native execution
5. **Maintenance**: Easier to maintain than C++/Qt

## Future Enhancements

- [ ] Parse `elmos.yaml` directly from Swift
- [ ] Real-time workspace monitoring
- [ ] Kernel build visualization
- [ ] QEMU integration with VNC
- [ ] Code signing and notarization
- [ ] DMG distribution
