import SwiftUI

struct ToolchainsView: View {
    @StateObject private var manager = ToolchainManager()
    @State private var selectedArch: String?

    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            Text("Toolchains")
                .font(.largeTitle)
                .bold()

            HStack {
                Text("Manage cross-compiler toolchains")
                    .foregroundStyle(.secondary)
                Spacer()
                Button(action: {
                    Task { await manager.loadToolchains() }
                }) {
                    Label("Refresh", systemImage: "arrow.clockwise")
                }
                .disabled(manager.isLoading)
            }

            if manager.isLoading {
                ProgressView("Loading toolchains...")
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if let error = manager.errorMessage {
                Text(error)
                    .foregroundStyle(.red)
            } else if manager.toolchains.isEmpty {
                VStack {
                    Image(systemName: "hammer.slash")
                        .font(.system(size: 48))
                        .foregroundStyle(.secondary)
                    Text("No toolchains found")
                        .font(.title3)
                    Text("Run 'elmos toolchains list' to see available toolchains")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                List(manager.toolchains) { toolchain in
                    ToolchainRow(toolchain: toolchain, manager: manager)
                }
            }
        }
        .padding()
        .task {
            await manager.loadToolchains()
        }
    }
}

struct ToolchainRow: View {
    let toolchain: Toolchain
    @ObservedObject var manager: ToolchainManager
    @State private var isBuilding = false

    var body: some View {
        HStack {
            VStack(alignment: .leading, spacing: 4) {
                Text(toolchain.arch)
                    .font(.headline)
                Text(toolchain.path)
                    .font(.caption)
                    .foregroundStyle(.secondary)
            }

            Spacer()

            if toolchain.isInstalled {
                Image(systemName: "checkmark.circle.fill")
                    .foregroundStyle(.green)
            } else {
                Button(action: {
                    Task {
                        isBuilding = true
                        _ = await manager.buildToolchain(arch: toolchain.arch)
                        isBuilding = false
                    }
                }) {
                    if isBuilding {
                        ProgressView()
                            .scaleEffect(0.7)
                    } else {
                        Text("Build")
                    }
                }
                .disabled(isBuilding || manager.isLoading)
            }
        }
        .padding(.vertical, 4)
    }
}
