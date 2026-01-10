import Foundation

struct Toolchain: Identifiable, Codable {
    var id = UUID()
    let arch: String
    let status: String
    let path: String

    var isInstalled: Bool {
        status.lowercased().contains("installed") || status.contains("✓")
    }

    private enum CodingKeys: String, CodingKey {
        case arch, status, path
    }
}

class ToolchainManager: ObservableObject {
    @Published var toolchains: [Toolchain] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let executor = CommandExecutor()

    @MainActor
    func loadToolchains() async {
        isLoading = true
        errorMessage = nil

        let result = await executor.executeAsync(args: ["toolchains", "list"])

        switch result {
        case .success(let output):
            parseToolchains(output)
        case .failure(let error):
            errorMessage = error.localizedDescription
        }

        isLoading = false
    }

    private func parseToolchains(_ output: String) {
        var parsed: [Toolchain] = []
        let lines = output.split(separator: "\n")

        for line in lines {
            let text = String(line).trimmingCharacters(in: .whitespaces)
            if text.isEmpty || text.hasPrefix("Architecture") || text.hasPrefix("---") {
                continue
            }

            let parts = text.split(separator: "|").map { $0.trimmingCharacters(in: .whitespaces) }
            if parts.count >= 3 {
                parsed.append(
                    Toolchain(
                        arch: parts[0],
                        status: parts[1],
                        path: parts[2]
                    ))
            }
        }

        toolchains = parsed
    }

    @MainActor
    func buildToolchain(arch: String) async -> Bool {
        isLoading = true
        let result = await executor.executeAsync(args: ["toolchains", "build", arch])
        isLoading = false

        if case .success = result {
            await loadToolchains()
            return true
        }
        return false
    }
}
