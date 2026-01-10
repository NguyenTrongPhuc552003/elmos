import Foundation

class QEMUManager: ObservableObject {
    @Published var isRunning = false
    @Published var output: String = ""
    @Published var errorMessage: String?

    private let executor = CommandExecutor()

    @MainActor
    func runQEMU() async {
        isRunning = true
        errorMessage = nil
        appendOutput("Starting QEMU...\n")

        // This will block until QEMU exits
        let result = await executor.executeAsync(args: ["qemu", "run"])

        switch result {
        case .success(let out):
            appendOutput(out)
            appendOutput("\nQEMU exited successfully.\n")
        case .failure(let error):
            errorMessage = error.localizedDescription
            appendOutput("\nQEMU exited with error: \(error.localizedDescription)\n")
        }

        isRunning = false
    }

    private func appendOutput(_ text: String) {
        output += text
    }

    func clearOutput() {
        output = ""
    }
}
