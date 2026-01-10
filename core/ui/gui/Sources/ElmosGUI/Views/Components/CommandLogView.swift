import SwiftUI

struct CommandLogView: View {
    @Binding var isPresented: Bool
    let title: String
    let output: String
    let isRunning: Bool

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            // Header
            HStack {
                Text(title)
                    .font(.headline)

                if isRunning {
                    ProgressView()
                        .controlSize(.small)
                        .padding(.leading, 8)
                }

                Spacer()

                Button(action: { isPresented = false }) {
                    Image(systemName: "xmark.circle.fill")
                        .foregroundStyle(.secondary)
                        .font(.title2)
                }
                .buttonStyle(.plain)
            }
            .padding()
            .background(Color(nsColor: .windowBackgroundColor))

            Divider()

            // Console Output
            ScrollView {
                Text(output)
                    .font(.monospaced(.caption)())
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding()
                    .textSelection(.enabled)
            }
            .background(Color.black)
            .foregroundStyle(.white)
        }
        .frame(minWidth: 600, minHeight: 400)
    }
}
