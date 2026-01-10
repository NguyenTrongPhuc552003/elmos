import AppKit
import Foundation

/// Helper to launch commands in external MAC Terminal
enum TerminalLauncher {
    static func launch(command: String, args: [String]) {
        // Construct the full command string
        // We use 'elmos' as the binary name, assuming it's in the build path
        // We need to resolve the absolute path to elmos binary based on Settings
        let elmosPath =
            UserDefaults.standard.string(forKey: "elmosPath")
            ?? "/Users/trongphucnguyen/Documents/kernel-dev/linux/build/elmos"
        let fullCommand = "\(elmosPath) \(command) \(args.joined(separator: " "))"

        // AppleScript to activate Terminal and run command
        let scriptSource = """
            tell application "Terminal"
                activate
                do script "\(fullCommand)"
            end tell
            """

        if let script = NSAppleScript(source: scriptSource) {
            var error: NSDictionary?
            script.executeAndReturnError(&error)
            if let error = error {
                print("Terminal launch error: \(error)")
            }
        }
    }
}
