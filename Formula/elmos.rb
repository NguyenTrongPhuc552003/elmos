# typed: false
# frozen_string_literal: true

# ELMOS - Embedded Linux on Multi-OS
# Homebrew formula for installing elmos CLI tool
#
# Installation (this repository is the tap; there is no homebrew-elmos
# repository, so the tap needs its URL):
#   brew tap NguyenTrongPhuc552003/elmos https://github.com/NguyenTrongPhuc552003/elmos
#   brew install elmos
#
# Development install (latest from git):
#   brew install --HEAD elmos

class Elmos < Formula
  desc "Embedded Linux on Multi-OS: kernels, modules, root file systems and QEMU"
  homepage "https://github.com/NguyenTrongPhuc552003/elmos"
  license "MIT"

  head "https://github.com/NguyenTrongPhuc552003/elmos.git", branch: "main"
  url "https://github.com/NguyenTrongPhuc552003/elmos/archive/v3.3.0.tar.gz"
  sha256 "1f44ad75fc3f8524981d37a36c6d4dc16f705ebc6b4cd105dca2c13472f12a6b"
  version "3.3.0"

  depends_on "go" => :build
  depends_on "go-task" => :build

  # Runtime dependencies
  depends_on "llvm"
  depends_on "lld"
  depends_on "gnu-sed"
  depends_on "make"
  depends_on "libelf"
  depends_on "qemu"
  depends_on "e2fsprogs"
  depends_on "coreutils"
  depends_on "fakeroot"

  def install
    # A release tarball has no .git: the version comes from the formula
    # (a --HEAD build from git describe).
    version_str = if build.head?
      Utils.safe_popen_read("git", "describe", "--tags", "--always", "--dirty").strip
    else
      "v#{version}"
    end
    build_date = Time.now.utc.iso8601

    ldflags = %W[
      -s -w
      -X github.com/NguyenTrongPhuc552003/elmos/core/app/version.Version=#{version_str}
      -X github.com/NguyenTrongPhuc552003/elmos/core/app/version.BuildDate=#{build_date}
    ]

    # The main package is cmd/elmos (the repository root has no Go files).
    system "go", "build", *std_go_args(ldflags: ldflags), "./cmd/elmos"

    # Generate and install shell completions
    generate_completions_from_executable(bin/"elmos", "completion")

    # Install supporting files
    pkgshare.install "libraries" if File.directory?("libraries")
    pkgshare.install "patches" if File.directory?("patches")
  end

  def caveats
    <<~EOS
      ELMOS has been installed!

      Quick start:
        elmos doctor              # Check dependencies
        elmos ui                  # Launch interactive TUI
        elmos init                # Initialize workspace
        elmos config set arch arm64
        elmos kernel config
        elmos build
        elmos qemu run

      Shell completions have been installed for bash, zsh, and fish.
    EOS
  end

  test do
    assert_match "ELMOS", shell_output("#{bin}/elmos --help")
    assert_match "Version", shell_output("#{bin}/elmos version")
  end
end
