import Foundation

struct KernelModule: Identifiable, Codable {
    var id = UUID()
    let name: String
    let size: String
    let usedBy: String
    let status: String

    // Custom coding keys to exclude id
    private enum CodingKeys: String, CodingKey {
        case name, size, usedBy, status
    }
}

class ModulesManager: ObservableObject {
    @Published var modules: [KernelModule] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let executor = CommandExecutor()

    @MainActor
    func loadModules() async {
        isLoading = true
        errorMessage = nil

        // Mock output for now or implement 'elmos module list'
        // Let's assume 'elmos module list' returns parseable text similar to lsmod
        let result = await executor.executeAsync(args: ["module", "list"])

        switch result {
        case .success(let output):
            parseModules(output)
        case .failure(let error):
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    private func parseModules(_ output: String) {
        var parsed: [KernelModule] = []
        let lines = output.split(separator: "\n")

        // Skip header
        for line in lines.dropFirst() {
            let parts = line.split(separator: " ", omittingEmptySubsequences: true)
            if parts.count >= 3 {
                parsed.append(
                    KernelModule(
                        name: String(parts[0]),
                        size: String(parts[1]),
                        usedBy: parts.count > 3 ? String(parts[3]) : "-",
                        status: "Loaded"
                    ))
            }
        }

        modules = parsed
    }
}
