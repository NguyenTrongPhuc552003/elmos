import SwiftUI

struct ContentView: View {
    @State private var selectedTab = 0

    var body: some View {
        NavigationSplitView {
            List(selection: $selectedTab) {
                Label("Dashboard", systemImage: "house.fill")
                    .tag(0)
                Label("Toolchains", systemImage: "hammer.fill")
                    .tag(1)
                Label("Kernel", systemImage: "cpu.fill")
                    .tag(2)
                Label("Modules", systemImage: "puzzlepiece.extension.fill")
                    .tag(3)
                Label("QEMU", systemImage: "desktopcomputer")
                    .tag(4)
                Label("Settings", systemImage: "gearshape.fill")
                    .tag(5)
            }
            .navigationTitle("ELMOS")
        } detail: {
            Group {
                switch selectedTab {
                case 0: DashboardView()
                case 1: ToolchainsView()
                case 2: KernelView()
                case 3: ModulesView()
                case 4: QEMUView()
                case 5: SettingsView()
                default: Text("Select a section")
                }
            }
        }
    }
}
