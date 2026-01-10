import SwiftUI

/// Reusable status row component
struct StatusRow: View {
    let label: String
    let value: String
    var isGood: Bool? = nil

    var body: some View {
        HStack {
            Text(label)
                .foregroundStyle(.secondary)
            Spacer()
            Text(value)
                .bold()
                .foregroundStyle(statusColor)
        }
        .padding(.vertical, 4)
    }

    private var statusColor: Color {
        guard let isGood = isGood else { return .primary }
        return isGood ? .green : .red
    }
}

/// Reusable quick action button component
struct QuickActionButton: View {
    let title: String
    let icon: String
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            VStack {
                Image(systemName: icon)
                    .font(.title)
                Text(title)
                    .font(.caption)
            }
            .frame(maxWidth: .infinity)
            .padding()
            .background(Color.accentColor.opacity(0.1))
            .cornerRadius(8)
        }
        .buttonStyle(.plain)
    }
}
