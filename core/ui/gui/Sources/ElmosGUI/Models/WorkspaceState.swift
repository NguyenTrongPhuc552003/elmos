import Foundation

/// WorkspaceState holds the current elmos workspace status
struct WorkspaceState: Codable {
    var isInitialized: Bool = false
    var architecture: String = "unknown"
    var kernelVersion: String = "unknown"
    var mountPath: String = ""
    var isMounted: Bool = false
}
