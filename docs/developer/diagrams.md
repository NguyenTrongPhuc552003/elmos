# Architecture Diagrams

Diagrams of ELMOS as the code stands, drawn as Mermaid text: GitHub draws
them here and the [documentation site](https://docs-haven.vercel.app/elmos/)
draws the same text in its own colours. Everything shown is in the code and
runs in ELMOS's CI. The evidence key shared with HAVEN and KDAL is on the
site's [Reading the diagrams](https://docs-haven.vercel.app/start/diagrams/)
page.

The component diagram is checked against the code:
`go test ./docs/developer/` fails when a package or an import between the
groups below is missing from it, or when it shows one the code does not
have.

---

## Components

Package groups and the imports between them. An arrow means "imports".

```mermaid
flowchart TB
  entry["Entry point<br/>cmd/ elmos"]
  sdk["SDK build helpers (CI)<br/>sdk/ demoboot, rootfs"]
  app["Commands<br/>core/app/ app, commands, version"]
  ui["Terminal output<br/>core/ui/ ui, tui"]
  domain["What ELMOS does<br/>core/domain/ builder, demo, doctor, emulator,<br/>patch, rootfs, toolchain"]
  workspace["Workspace and settings<br/>core/ context, config"]
  infra["The host<br/>core/infra/ executor, filesystem, homebrew"]
  assets["Embedded files<br/>assets"]
  thirdparty["Vendored code<br/>third_party/ext4/ compactext4, format"]
  entry --> app
  entry --> workspace
  entry --> infra
  app --> ui
  app --> domain
  app --> workspace
  app --> infra
  sdk --> domain
  sdk --> workspace
  sdk --> infra
  domain --> ui
  domain --> workspace
  domain --> infra
  domain --> assets
  domain --> thirdparty
  workspace --> infra
  classDef observed stroke:#0ca30c,stroke-width:2px
  class entry,sdk,app,ui,domain,workspace,infra,assets,thirdparty observed
```

`core/domain/toolchain` imports `core/ui` (its manager holds a printer),
the one place where the domain reaches the terminal layer.

---

## Kernel build

`elmos kernel build` (`core/app/commands/kernel.go`,
`core/domain/builder/kernel.go`, `core/domain/builder/env.go`).

```mermaid
sequenceDiagram
  actor U as User
  participant C as kernel build command<br/>core/app/commands
  participant B as KernelBuilder<br/>core/domain/builder
  participant T as toolchain.Manager<br/>core/domain/toolchain
  participant X as Executor<br/>core/infra/executor
  participant M as make (Linux kbuild)
  U->>C: elmos kernel build [targets]
  C->>B: Build(ctx, BuildOptions{Jobs, Targets})
  B->>B: check each target against config.ValidBuildTargets
  B->>T: getToolchainEnv(arch)
  T-->>B: PATH and CROSS_COMPILE of the installed crosstool-ng toolchain
  B->>X: RunWithEnv(env, make, ...)
  X->>M: make -C kernel_dir -jN ARCH=arch LLVM=1 CROSS_COMPILE=prefix targets
  M-->>X: output, streamed to the terminal
  X-->>B: exit status
  B-->>C: error or nil
  C-->>U: result
```

---

## The HAVEN demo

`elmos demo haven` (`core/app/commands/demo.go`, `core/domain/demo/demo.go`)
fetches HAVEN's try bundle and runs its `try.sh`, described in HAVEN's
[`scripts/release/try/README.md`](https://github.com/NguyenTrongPhuc552003/haven/blob/main/scripts/release/try/README.md).

```mermaid
sequenceDiagram
  actor U as User
  participant E as elmos demo haven
  participant F as demo.Fetcher<br/>cache ~/.elmos/demo
  participant GH as GitHub release<br/>haven, demo.HavenTag
  participant T as try.sh<br/>in the bundle
  participant R as readback.py<br/>in the bundle
  participant Q as qemu-system-aarch64
  U->>E: elmos demo haven
  E->>F: Release("haven", HavenTag)
  F->>GH: SHA256SUMS, then each file it lists
  GH-->>F: files (a cached file with the right SHA-256 is reused)
  F-->>E: release directory
  E->>E: unpack the bundle, check that try.sh exists
  E->>T: sh ./try.sh
  T->>Q: partition-b-check.py: boot the good images (two partitions)
  T->>R: read back the good images against grant.json
  R->>Q: boot, stop, read stage-2, SMMU and GIC state through QMP
  R-->>T: no difference
  T->>R: read back the image with a seeded fault
  R->>Q: boot, stop, read the enforced state
  R-->>T: finding: partition memory mapped at EL2
  T-->>U: both results
```

---

## Main types

The domain types the commands use (`core/domain/*`, `core/config/types.go`).
Fields are unexported unless marked `+`.

```mermaid
classDiagram
  class KernelBuilder {
    -exec Executor
    -fs FileSystem
    -cfg *Config
    -ctx *Context
    -tm *toolchain.Manager
    +Build(ctx, BuildOptions) error
    +Configure(ctx, configType) error
    +EnableKVMConfig(ctx) error
    +Clean(ctx) error
  }
  class BuildOptions {
    +Jobs int
    +Targets []string
  }
  class QEMURunner {
    -exec Executor
    -fs FileSystem
    -cfg *Config
    -ctx *Context
    +Run(ctx, RunOptions) error
    +Debug(ctx, graphical) error
    +ConnectGDB() error
  }
  class RunOptions {
    +Debug bool
    +Graphical bool
  }
  class Creator {
    -exec Executor
    -fs FileSystem
    -cfg *Config
    +Create(ctx, CreateOptions) error
    +Status() RootfsInfo
    +Clean(ctx) error
  }
  class CreateOptions {
    +Size string
    +Engine string
    +Image string
  }
  class Fetcher {
    +CacheDir string
    +Release(ctx, repo, tag) string
  }
  class Config {
    +Image ImageConfig
    +Build BuildConfig
    +QEMU QEMUConfig
    +Paths PathsConfig
    +Profiles map
  }
  KernelBuilder ..> BuildOptions
  KernelBuilder --> Config
  QEMURunner ..> RunOptions
  QEMURunner --> Config
  Creator ..> CreateOptions
  Creator --> Config
```

`Creator.Create` uses an OCI image by default (`alpine:3`); debootstrap is
the other engine.

---

## Workspace states

The case-sensitive volume that holds the kernel source, root file system and
toolchains (`core/app/commands/init.go`, `core/context/context.go`).

```mermaid
stateDiagram-v2
  [*] --> NoImage
  NoImage --> Mounted: elmos init (creates the image)
  Unmounted --> Mounted: elmos init (attaches it)
  Mounted --> Unmounted: elmos exit
  note right of Mounted
    kernel, module, patch, qemu, rootfs,
    toolchains and arch commands run here
  end note
  note right of Unmounted
    those commands stop with
    "kernel volume not mounted"
  end note
```

`elmos init [name] [size]` creates a sparse, case-sensitive APFS image
(`hdiutil create`) if there is none, then attaches it; `elmos exit [--force]`
detaches it. `demo`, `doctor`, `status` and `version` work in any state.

---

## Deployment

Where things live on the Mac, and what QEMU gets (`core/config/loader.go`,
`core/domain/emulator/qemu.go`). Paths are the defaults; `elmos.yaml`
can change them. The first `elmos.yaml` found is used.

```mermaid
flowchart LR
  subgraph host["macOS host"]
    bin["elmos binary"]
    yaml["elmos.yaml<br/>first found in ., build/,<br/>~/.config/elmos, /etc/elmos"]
    image[("data/NAME.sparseimage<br/>in the project root")]
    cache[("~/.elmos/demo<br/>demo release cache")]
    subgraph vol["/Volumes/NAME (the mounted image)"]
      linux["linux/<br/>kernel source and build"]
      tools["toolchains/<br/>crosstool-ng"]
      disk[("disk.img<br/>root file system")]
    end
    mods["examples/modules/<br/>kernel modules"]
  end
  subgraph vm["QEMU virt machine"]
    kernel["Linux kernel"]
    rootfs["root file system"]
    share["9p share<br/>mount tag modules_mount"]
  end
  bin --> yaml
  image -->|hdiutil attach| vol
  linux -->|-kernel| kernel
  disk -->|virtio drive| rootfs
  mods -->|virtio-9p| share
  bin -->|launches| vm
```
