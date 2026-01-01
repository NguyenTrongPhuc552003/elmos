// Package cmd implements the Cobra CLI commands for elmos.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/internal/core"
)

// doctorCmd - check environment and dependencies
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment and dependencies",
	Long:  `Verify that all required tools and dependencies are installed correctly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDoctor()
	},
}

// RequiredPackage represents a Homebrew package dependency
type RequiredPackage struct {
	Name        string
	Description string
	Required    bool
}

// requiredPackages lists all Homebrew dependencies
var requiredPackages = []RequiredPackage{
	{"llvm", "LLVM/Clang toolchain", true},
	{"lld", "LLVM linker", true},
	{"gnu-sed", "GNU sed (kernel requires it)", true},
	{"make", "GNU make 4.0+", true},
	{"libelf", "ELF library", true},
	{"git", "Git version control", true},
	{"qemu", "QEMU emulator", true},
	{"fakeroot", "Fake root for packaging", true},
	{"e2fsprogs", "ext4 filesystem tools", true},
	{"wget", "File downloader", false},
	{"coreutils", "GNU core utilities", true},
	{"go", "Go programming language", true},
	{"go-task", "Go task runner", true},
}

// requiredTaps lists required Homebrew taps
var requiredTaps = []string{
	"messense/macos-cross-toolchains",
}

// requiredHeaders lists header files that should exist
var requiredHeaders = []string{
	"elf.h",
	"byteswap.h",
}

func runDoctor() error {
	fmt.Println()
	printInfo("ELMOS Doctor - Environment Check")
	fmt.Println()

	issuesFound := 0

	// Check Homebrew
	printStep("Checking Homebrew...")
	if !checkCommandExists("brew") {
		printError("Homebrew not found")
		printInfo("Install from: https://brew.sh")
		issuesFound++
	} else {
		printSuccess("Homebrew found")
	}
	fmt.Println()

	// Check taps
	printStep("Checking Homebrew taps...")
	issuesFound += checkTaps()
	fmt.Println()

	// Check packages
	printStep("Checking required packages...")
	issuesFound += checkPackages()
	fmt.Println()

	// Check headers
	printStep("Checking custom headers...")
	issuesFound += checkHeaders()
	fmt.Println()

	// Offer to fix missing elf.h - if fixed, decrement issue count
	if fixMissingElfH() {
		issuesFound--
	}

	// Check architecture-specific GDB
	printStep("Checking cross-debuggers...")
	checkCrossGDB()
	fmt.Println()

	// Check architecture-specific GCC
	printStep("Checking cross-compilers...")
	checkCrossGCC()
	fmt.Println()

	// Summary
	if issuesFound == 0 {
		printSuccess("All checks passed! Environment ready for kernel development.")
		return nil
	}

	printWarn("Found %d issue(s). Please fix before proceeding.", issuesFound)
	return fmt.Errorf("doctor found %d issues", issuesFound)
}

func checkCommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func checkTaps() int {
	issues := 0
	out, err := exec.Command("brew", "tap").Output()
	if err != nil {
		printError("Failed to list Homebrew taps")
		return 1
	}

	installedTaps := strings.Split(strings.TrimSpace(string(out)), "\n")

	for _, tap := range requiredTaps {
		found := false
		for _, installed := range installedTaps {
			if installed == tap {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("  ✓ %s\n", tap)
		} else {
			fmt.Printf("  ✗ %s (missing)\n", tap)
			printInfo("  Fix: brew tap %s", tap)
			issues++
		}
	}

	return issues
}

func checkPackages() int {
	issues := 0
	out, err := exec.Command("brew", "list", "--formulae").Output()
	if err != nil {
		printError("Failed to list Homebrew packages")
		return 1
	}

	installedPkgs := strings.Split(strings.TrimSpace(string(out)), "\n")
	installedSet := make(map[string]bool)
	for _, pkg := range installedPkgs {
		installedSet[pkg] = true
	}

	var missing []string

	for _, pkg := range requiredPackages {
		if installedSet[pkg.Name] {
			fmt.Printf("  ✓ %s\n", pkg.Name)
		} else {
			status := "missing"
			if !pkg.Required {
				status = "optional, missing"
			}
			fmt.Printf("  ✗ %s (%s)\n", pkg.Name, status)
			if pkg.Required {
				missing = append(missing, pkg.Name)
				issues++
			}
		}
	}

	if len(missing) > 0 {
		printInfo("Fix: brew install %s", strings.Join(missing, " "))
	}

	return issues
}

func checkHeaders() int {
	cfg := core.GetConfig()
	headersDir := cfg.Paths.LibrariesDir
	issues := 0

	// Check if headers directory exists
	if _, err := os.Stat(headersDir); os.IsNotExist(err) {
		printError("Headers directory not found: %s", headersDir)
		return 1
	}

	for _, header := range requiredHeaders {
		headerPath := filepath.Join(headersDir, header)
		if _, err := os.Stat(headerPath); err == nil {
			fmt.Printf("  ✓ %s\n", header)
		} else {
			fmt.Printf("  ✗ %s (missing)\n", header)
			issues++
		}
	}

	// Check asm symlinks
	asmDir := filepath.Join(headersDir, "asm")
	if info, err := os.Stat(asmDir); err == nil && info.IsDir() {
		fmt.Printf("  ✓ asm/ (directory)\n")
	} else {
		fmt.Printf("  ✗ asm/ (missing)\n")
		issues++
	}

	return issues
}

// fixMissingElfH downloads elf.h from glibc if it's missing
// Returns true if it successfully fixed a missing elf.h
func fixMissingElfH() bool {
	cfg := core.GetConfig()
	headersDir := cfg.Paths.LibrariesDir
	elfPath := filepath.Join(headersDir, "elf.h")

	// Check if elf.h already exists
	if _, err := os.Stat(elfPath); err == nil {
		return false
	}

	// Prompt user
	fmt.Print("elf.h missing — download from glibc? (Y/n): ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	if choice == "n" || choice == "no" {
		return false
	}

	// Download from glibc
	glibcVersion := "2.42"
	url := fmt.Sprintf("https://raw.githubusercontent.com/bminor/glibc/glibc-%s/elf/elf.h", glibcVersion)

	// Ensure directory exists
	if err := os.MkdirAll(headersDir, 0755); err != nil {
		printError("Failed to create directory: %v", err)
		return false
	}

	printStep("Downloading elf.h from glibc %s...", glibcVersion)

	resp, err := http.Get(url)
	if err != nil {
		printError("Download failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		printError("Download failed: HTTP %d", resp.StatusCode)
		return false
	}

	out, err := os.Create(elfPath)
	if err != nil {
		printError("Failed to create file: %v", err)
		return false
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		printError("Failed to write file: %v", err)
		return false
	}

	printSuccess("elf.h downloaded successfully")
	fmt.Println()
	return true
}

func checkCrossGDB() {
	gdbs := []struct {
		arch string
		bin  string
	}{
		{"RISC-V", "riscv64-elf-gdb"},
		{"ARM64", "aarch64-unknown-linux-gnu-gdb"},
		{"ARM32", "arm-unknown-linux-gnueabihf-gdb"},
	}

	for _, gdb := range gdbs {
		if checkCommandExists(gdb.bin) {
			fmt.Printf("  ✓ %s (%s)\n", gdb.arch, gdb.bin)
		} else {
			fmt.Printf("  ○ %s (%s) - optional\n", gdb.arch, gdb.bin)
		}
	}
}

func checkCrossGCC() {
	gccs := []struct {
		arch string
		bin  string
	}{
		{"RISC-V", "riscv64-elf-gcc"},
		{"ARM64", "aarch64-unknown-linux-gnu-gcc"},
		{"ARM32", "arm-unknown-linux-gnueabihf-gcc"},
	}

	for _, gcc := range gccs {
		if checkCommandExists(gcc.bin) {
			fmt.Printf("  ✓ %s (%s)\n", gcc.arch, gcc.bin)
		} else {
			fmt.Printf("  ○ %s (%s) - optional\n", gcc.arch, gcc.bin)
		}
	}
}
