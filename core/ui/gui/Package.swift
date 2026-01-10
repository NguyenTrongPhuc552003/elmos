// swift-tools-version: 6.0
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "ElmosGUI",
    platforms: [.macOS(.v14)],
    products: [
        // Binary executable named 'elmos' (no hyphens or underscores)
        .executable(
            name: "elmos",
            targets: ["ElmosGUI"]
        )
    ],
    dependencies: [],
    targets: [
        .executableTarget(
            name: "ElmosGUI",
            path: "Sources/ElmosGUI"
        )
    ]
)
