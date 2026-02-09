# Architecture Diagrams

This page contains visual diagrams of the ELMOS architecture.

---

## GUI Architecture

![GUI Architecture](../diagrams/gui_architecture.puml)

## GUI Build Sequence

![GUI Build Sequence](../diagrams/gui_build_sequence.puml)

## Platform Abstraction

![Platform Abstraction](../diagrams/platform_abstraction.puml)

## Component

High-level component architecture showing package relationships.

```puml
@startuml
!theme plain
skinparam componentStyle rectangle

title ELMOS Component Architecture

package "cmd/elmos" {
    [main.go] as main
}

package "core/app" {
    [app.go] as app
    [commands/] as commands
}

package "core/config" {
    [types.go] as types
    [loader.go] as loader
    [arch.go] as arch
}

package "core/context" {
    [context.go] as context
}

package "core/domain" {
    [builder/] as builder
    [emulator/] as emulator
    [rootfs/] as rootfs
    [patch/] as patch
    [module/] as module
    [toolchain/] as toolchain
}

package "core/infra" {
    [executor/] as executor
    [filesystem/] as filesystem
    [printer/] as printer
}

package "core/ui" {
    [tui/] as tui
}

package "assets" {
    [templates/] as templates
}

main --> app
app --> commands
commands --> context
context --> loader
context --> executor
context --> filesystem
commands --> builder
commands --> emulator
commands --> rootfs
builder --> executor
builder --> templates
emulator --> executor
rootfs --> filesystem
tui --> commands
@enduml
```

---

## Sequence

Kernel build workflow sequence.

```puml
@startuml
!theme plain

title Kernel Build Workflow

actor User
participant "elmos CLI" as CLI
participant "KernelBuilder" as Builder
participant "Executor" as Exec
participant "FileSystem" as FS
participant "Make" as Make

User -> CLI: elmos kernel build
CLI -> Builder: Build(targets)
Builder -> FS: Check kernel dir exists
FS --> Builder: OK

Builder -> Builder: getToolchainEnv()
Builder -> Exec: Run make with env
Exec -> Make: make ARCH=arm64 LLVM=1 Image

Make --> Exec: Output stream
Exec --> Builder: Success
Builder --> CLI: Build complete
CLI --> User: ✓ Kernel built
@enduml
```

---

## Class

Core domain class relationships.

```puml
@startuml
!theme plain

title Core Domain Classes

package "domain/builder" {
    class KernelBuilder {
        -cfg: *Config
        -ctx: *Context
        -exec: Executor
        -fs: FileSystem
        +Build(targets []string) error
        +Configure(configType string) error
        +Clean() error
    }
}

package "domain/emulator" {
    class QEMURunner {
        -cfg: *Config
        -ctx: *Context
        -exec: Executor
        +Run(opts RunOptions) error
        +GetListMachines() []MachineInfo
    }
    
    class RunOptions {
        +Debug: bool
        +Run: bool
        +Graphical: bool
        +Targets: []Target
        +Machine: string
    }
}

package "domain/rootfs" {
    class RootfsManager {
        -cfg: *Config
        -exec: Executor
        -fs: FileSystem
        +Create() error
        +UpdateDisk() error
    }
}

package "config" {
    class Config {
        +Image: ImageConfig
        +Build: BuildConfig
        +QEMU: QEMUConfig
        +Paths: PathsConfig
    }
}

KernelBuilder --> Config
QEMURunner --> Config
QEMURunner --> RunOptions
RootfsManager --> Config
@enduml
```

---

## State

Workspace state machine.

```puml
@startuml
!theme plain

title ELMOS Workspace State Machine

[*] --> Unmounted : Initial

Unmounted --> Mounted : elmos init
Mounted --> Unmounted : elmos exit

state Mounted {
    [*] --> Ready
    
    Ready --> Configuring : elmos kernel config
    Configuring --> Ready : config complete
    
    Ready --> Building : elmos kernel build
    Building --> Ready : build complete
    Building --> Error : build failed
    
    Ready --> Running : elmos qemu -r
    Running --> Ready : QEMU exit
    
    Ready --> Debugging : elmos qemu -d
    Debugging --> Ready : Debug session end
    
    Error --> Ready : fix & retry
}
@enduml
```

---

## Deployment

Runtime deployment architecture.

```puml
@startuml
!theme plain

title ELMOS Deployment View

node "macOS Host" {
    package "elmos CLI" {
        [elmos binary]
    }
    
    folder "/Volumes/elmos" as volume {
        folder "linux/" {
            [Kernel Source]
        }
        folder "rootfs/" {
            [Debian Root]
        }
        file "disk.img" as disk
    }
    
    database "~/.config/elmos" {
        [elmos.yaml]
    }
}

cloud "QEMU VM" {
    [Linux Kernel]
    [Root Filesystem]
    folder "/mnt/modules" {
        [9p mount]
    }
}

[elmos binary] --> volume : mounts
[elmos binary] --> [QEMU VM] : launches
disk --> [Root Filesystem] : virtio
[Kernel Source] --> [Linux Kernel] : builds
@enduml
```
