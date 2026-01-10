import SwiftUI

struct ModulesView: View {
    @StateObject private var manager = ModulesManager()

    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            Text("Kernel Modules")
                .font(.largeTitle)
                .bold()

            HStack {
                Text("Manage loaded kernel modules")
                    .foregroundStyle(.secondary)
                Spacer()
                Button(action: {
                    Task { await manager.loadModules() }
                }) {
                    Label("Refresh", systemImage: "arrow.clockwise")
                }
                .disabled(manager.isLoading)
            }

            if manager.isLoading {
                ProgressView("Loading modules...")
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else if let error = manager.errorMessage {
                Text("Error: \(error)")
                    .foregroundStyle(.red)
                    .padding()
            } else if manager.modules.isEmpty {
                VStack {
                    Image(systemName: "memorychip")
                        .font(.system(size: 48))
                        .foregroundStyle(.secondary)
                    Text("No modules loaded")
                        .font(.title3)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                Table(manager.modules) {
                    TableColumn("Name", value: \.name)
                    TableColumn("Size", value: \.size)
                    TableColumn("Used By", value: \.usedBy)
                }
            }
        }
        .padding()
        .task {
            await manager.loadModules()
        }
    }
}
