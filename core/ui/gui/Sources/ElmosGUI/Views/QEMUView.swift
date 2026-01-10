import SwiftUI

struct QEMUView: View {
    @StateObject private var manager = QEMUManager()

    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            Text("QEMU Emulator")
                .font(.largeTitle)
                .bold()

            Text("Run and debug the built kernel in QEMU.")
                .foregroundStyle(.secondary)

            HStack(spacing: 20) {
                Button(action: {
                    Task { await manager.runQEMU() }
                }) {
                    Label(manager.isRunning ? "Running..." : "Run QEMU", systemImage: "play.fill")
                        .frame(minWidth: 150)
                        .padding()
                }
                .disabled(manager.isRunning)
                .buttonStyle(.borderedProminent)
                .controlSize(.large)
                .tint(.green)

                if manager.isRunning {
                    ProgressView()
                        .controlSize(.small)
                }

                Spacer()
            }
            .padding(.vertical)

            Divider()

            // Output Log
            VStack(alignment: .leading) {
                HStack {
                    Text("Session Log")
                        .font(.headline)
                    Spacer()
                    Button("Clear") {
                        manager.clearOutput()
                    }
                    .font(.caption)
                }

                ScrollView {
                    Text(manager.output.isEmpty ? "Ready to launch..." : manager.output)
                        .font(.monospaced(.caption)())
                        .frame(maxWidth: .infinity, alignment: .leading)
                        .padding()
                }
                .background(Color.black.opacity(0.8))
                .cornerRadius(8)
                .foregroundStyle(.white)
            }
        }
        .padding()
    }
}
