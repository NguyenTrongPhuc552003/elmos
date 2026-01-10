import Foundation

/// WorkspaceManager manages workspace state and configuration
class WorkspaceManager: ObservableObject {
    @Published var state: WorkspaceState = WorkspaceState()
    @Published var isLoading: Bool = false
    @Published var errorMessage: String?

    private let executor: CommandExecutor

    init(executor: CommandExecutor = CommandExecutor()) {
        self.executor = executor
    }

    /// Refresh workspace status from CLI
    @MainActor
    func refreshStatus() async {
        isLoading = true
        errorMessage = nil

        // 1. Check Status
        let statusResult = await executor.executeAsync(args: ["status"])
        switch statusResult {
        case .success(let output):
            parseStatusOutput(output)
        case .failure(let error):
            errorMessage = "Status check failed: \(error.localizedDescription)"
        }

        // 2. Check Architecture
        if let arch = await getArchitecture() {
            state.architecture = arch
        }

        isLoading = false
    }

    /// Parse status command output
    private func parseStatusOutput(_ output: String) {
        // Logic based on actual 'elmos status' output: "✓ Workspace mounted at /Volumes/..."
        if output.contains("Workspace mounted") {
            state.isInitialized = true
            state.isMounted = true
        } else {
            // Reset if not valid
            state.isInitialized = false
            state.isMounted = false
        }

        // Extract mount path if present
        if let match = output.firstMatch(of: /Workspace mounted at (.+)/) {
            state.mountPath = String(match.1).trimmingCharacters(in: .whitespacesAndNewlines)
        }
    }

    /// Get architecture info
    @MainActor
    func getArchitecture() async -> String? {
        let result = await executor.executeAsync(args: ["arch"])
        if case .success(let output) = result {
            // Output format: "Architecture: arm64"
            let components = output.split(separator: ":")
            if components.count > 1 {
                return String(components[1]).trimmingCharacters(in: .whitespacesAndNewlines)
            }
            return output.trimmingCharacters(in: .whitespacesAndNewlines)
        }
        return nil
    }
}
