import SwiftUI

struct DashboardView: View {
    @StateObject private var workspaceManager = WorkspaceManager()

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 20) {
                Text("Dashboard")
                    .font(.largeTitle)
                    .bold()

                WorkspaceStatusCard(workspaceManager: workspaceManager)
                QuickActionsCard(workspaceManager: workspaceManager)

                Spacer()
            }
            .padding()
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .topLeading)
        .task {
            await workspaceManager.refreshStatus()
        }
    }
}

// MARK: - Workspace Status Card
private struct WorkspaceStatusCard: View {
    @ObservedObject var workspaceManager: WorkspaceManager

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Text("Workspace Status")
                    .font(.title2)
                    .bold()
                Spacer()
                Button(action: {
                    Task {
                        await workspaceManager.refreshStatus()
                    }
                }) {
                    Label("Refresh", systemImage: "arrow.clockwise")
                }
                .disabled(workspaceManager.isLoading)
            }

            if workspaceManager.isLoading {
                ProgressView("Loading workspace status...")
                    .padding()
            } else {
                StatusRow(
                    label: "Initialized",
                    value: workspaceManager.state.isInitialized ? "✓ Yes" : "✗ No",
                    isGood: workspaceManager.state.isInitialized)
                StatusRow(
                    label: "Volume Mounted",
                    value: workspaceManager.state.isMounted ? "✓ Yes" : "✗ No",
                    isGood: workspaceManager.state.isMounted)
                StatusRow(
                    label: "Architecture",
                    value: workspaceManager.state.architecture)
                StatusRow(
                    label: "Kernel Version",
                    value: workspaceManager.state.kernelVersion)

                if let error = workspaceManager.errorMessage {
                    Text(error)
                        .foregroundStyle(.red)
                        .font(.caption)
                }
            }
        }
        .padding()
        .background(Color(.windowBackgroundColor))
        .cornerRadius(12)
        .shadow(radius: 2)
    }
}

// MARK: - Quick Actions Card
struct QuickActionsCard: View {
    @ObservedObject var workspaceManager: WorkspaceManager

    @State private var showingLog = false
    @State private var logTitle = ""
    @State private var logOutput = ""
    @State private var isCommandRunning = false

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Quick Actions")
                .font(.title2)
                .bold()

            HStack(spacing: 12) {
                QuickActionButton(title: "Initialize", icon: "play.circle") {
                    runCommand("Initialize", ["init"])
                }
                QuickActionButton(title: "Doctor", icon: "stethoscope") {
                    runCommand("Doctor", ["doctor"])
                }
                QuickActionButton(title: "Build", icon: "hammer") {
                    runCommand("Build", ["build"])
                }
                QuickActionButton(title: "Run QEMU", icon: "desktopcomputer") {
                    runCommand("QEMU", ["qemu", "run"])
                }
            }
        }
        .padding()
        .background(Color(nsColor: .windowBackgroundColor))
        .cornerRadius(12)
        .shadow(radius: 2)
        .sheet(isPresented: $showingLog) {
            CommandLogView(
                isPresented: $showingLog,
                title: logTitle,
                output: logOutput,
                isRunning: isCommandRunning
            )
        }
    }

    private func runCommand(_ title: String, _ args: [String]) {
        logTitle = title
        logOutput = "Starting command: elmos \(args.joined(separator: " "))\n\n"
        isCommandRunning = true
        showingLog = true

        Task {
            let executor = CommandExecutor()
            let result = await executor.executeAsync(args: args)

            switch result {
            case .success(let output):
                logOutput += output
                logOutput += "\n\n✅ Command completed successfully."
            case .failure(let error):
                logOutput += "\n\n❌ Error: \(error.localizedDescription)"
            }

            isCommandRunning = false

            if args.contains("init") {
                await workspaceManager.refreshStatus()
            }
        }
    }
}
