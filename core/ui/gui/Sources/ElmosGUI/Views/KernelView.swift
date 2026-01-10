import SwiftUI

struct KernelView: View {
    @StateObject private var manager = KernelManager()
    @State private var showingLog = false
    @State private var searchText = ""

    var filteredOptions: [KernelManager.ConfigOption] {
        if searchText.isEmpty {
            return manager.configOptions
        }
        return manager.configOptions.filter { $0.key.localizedCaseInsensitiveContains(searchText) }
    }

    var body: some View {
        HSplitView {
            // SIDEBAR: Configuration Editor
            VStack(alignment: .leading, spacing: 0) {
                HStack {
                    Text("Configuration")
                        .font(.headline)
                    Spacer()
                    Button(action: { Task { await manager.saveConfig() } }) {
                        Label("Save", systemImage: "square.and.arrow.down")
                    }
                }
                .padding()
                .background(Color(nsColor: .controlBackgroundColor))

                Divider()

                // Search Bar
                HStack {
                    Image(systemName: "magnifyingglass")
                        .foregroundStyle(.secondary)
                    TextField("Search config...", text: $searchText)
                        .textFieldStyle(.plain)
                }
                .padding(8)
                .background(Color(nsColor: .controlBackgroundColor))

                Divider()

                // Config List
                List {
                    if manager.configOptions.isEmpty {
                        ContentUnavailableView(
                            "No Configuration Loaded",
                            systemImage: "gearshape.2",
                            description: Text("Select an architecture to load config"))
                    } else {
                        ForEach(filteredOptions) { option in
                            HStack {
                                VStack(alignment: .leading) {
                                    Text(option.key.replacingOccurrences(of: "CONFIG_", with: ""))
                                        .font(.system(.body, design: .monospaced))
                                        .fontWeight(.medium)
                                    Text(option.key)
                                        .font(.caption2)
                                        .foregroundStyle(.tertiary)
                                }

                                Spacer()

                                if option.isBoolean {
                                    Toggle(
                                        "",
                                        isOn: Binding(
                                            get: { option.isEnabled },
                                            set: { _ in manager.toggleOption(option.id) }
                                        )
                                    )
                                    .toggleStyle(.switch)
                                    .controlSize(.mini)
                                } else {
                                    TextField(
                                        "",
                                        text: Binding(
                                            get: { option.value },
                                            set: {
                                                manager.updateOption(option.id, stringValue: $0)
                                            }
                                        )
                                    )
                                    .textFieldStyle(.roundedBorder)
                                    .frame(width: 100)
                                }
                            }
                            .padding(.vertical, 2)
                        }
                    }
                }
            }
            .frame(minWidth: 300)

            // MAIN AREA: Actions & Logs
            VStack {
                // Header / Actions
                HStack(spacing: 20) {
                    Button(action: { Task { await manager.loadConfig(for: "arm64") } }) {
                        Label("Load Config", systemImage: "arrow.triangle.2.circlepath")
                    }

                    Spacer()

                    Button(action: {
                        Task {
                            await manager.buildKernel()
                            showingLog = true
                        }
                    }) {
                        Label(
                            manager.isBuilding ? "Building..." : "Build Kernel",
                            systemImage: "hammer.fill"
                        )
                        .frame(minWidth: 120)
                        .padding()
                    }
                    .disabled(manager.isBuilding)
                    .buttonStyle(.borderedProminent)
                    .controlSize(.large)
                }
                .padding()

                // Status Box
                if let status = manager.lastBuildStatus {
                    HStack {
                        Image(
                            systemName: status == "Success"
                                ? "checkmark.circle.fill" : "exclamationmark.triangle.fill")
                        Text("Last Build: \(status)")
                    }
                    .foregroundStyle(status == "Success" ? .green : .red)
                    .padding(.bottom)
                }

                Divider()

                // Embedded Log View (Not sheet)
                VStack(alignment: .leading, spacing: 0) {
                    HStack {
                        Text("Console Output")
                            .font(.caption)
                            .fontWeight(.bold)
                            .foregroundStyle(.secondary)
                        Spacer()
                        Button("Clear") { manager.clearOutput() }
                            .buttonStyle(.plain)
                            .font(.caption)
                    }
                    .padding(8)
                    .background(Color.black)

                    ScrollView {
                        Text(manager.output.isEmpty ? "Ready..." : manager.output)
                            .font(.monospaced(.caption)())
                            .frame(maxWidth: .infinity, alignment: .leading)
                            .padding(8)
                    }
                    .background(Color.black)
                    .foregroundStyle(.white)
                }
            }
        }
        .task {
            // Auto load defaults
            await manager.loadConfig(for: "arm64")
        }
    }
}
