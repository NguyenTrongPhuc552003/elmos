import Foundation

/// CommandExecutor handles running elmos CLI commands from Swift
final class CommandExecutor: ObservableObject, @unchecked Sendable {
    private var elmosPath: String {
        UserDefaults.standard.string(forKey: "elmosPath")
            ?? "/Users/trongphucnguyen/Documents/kernel-dev/linux/build/elmos"
    }

    init() {}

    /// Execute a command synchronously and return output
    func execute(args: [String]) -> Result<String, Error> {
        let process = Process()
        let pipe = Pipe()

        // Resolve absolute path if local
        let executablePath = (elmosPath as NSString).expandingTildeInPath
        process.executableURL = URL(fileURLWithPath: executablePath)
        process.arguments = args
        process.standardOutput = pipe
        process.standardError = pipe

        do {
            try process.run()
            process.waitUntilExit()

            let data = pipe.fileHandleForReading.readDataToEndOfFile()
            let output = String(data: data, encoding: .utf8) ?? ""

            if process.terminationStatus == 0 {
                return .success(output)
            } else {
                return .failure(
                    NSError(
                        domain: "CommandExecutor", code: Int(process.terminationStatus),
                        userInfo: [NSLocalizedDescriptionKey: output]))
            }
        } catch {
            return .failure(error)
        }
    }

    /// Execute command asynchronously
    func executeAsync(args: [String]) async -> Result<String, Error> {
        await Task.detached {
            self.execute(args: args)
        }.value
    }
}
