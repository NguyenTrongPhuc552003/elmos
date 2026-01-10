import SwiftUI

struct SettingsView: View {
    @AppStorage("elmosPath") private var elmosPath: String =
        "/Users/trongphucnguyen/Documents/kernel-dev/linux/build/elmos"
    @AppStorage("theme") private var theme: String = "System"

    var body: some View {
        Form {
            Section(header: Text("General")) {
                TextField("ELMOS CLI Path", text: $elmosPath)

                Picker("Theme", selection: $theme) {
                    Text("System").tag("System")
                    Text("Light").tag("Light")
                    Text("Dark").tag("Dark")
                }
            }

            Section(header: Text("About")) {
                HStack {
                    Text("Version")
                    Spacer()
                    Text("1.0.0")
                        .foregroundStyle(.secondary)
                }

                Link(
                    "Documentation",
                    destination: URL(string: "https://github.com/NguyenTrongPhuc552003/elmos")!)
            }
        }
        .formStyle(.grouped)
        .navigationTitle("Settings")
    }
}
