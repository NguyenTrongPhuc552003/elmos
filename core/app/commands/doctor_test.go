package commands

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"codeberg.org/NguyenTrongPhuc552003/elmos/core/config"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/domain/doctor"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/ui"
)

func TestDoctorReportsMissingElfHeaderWithoutDownload(t *testing.T) {
	librariesDir := prepareDoctorLibraries(t)
	before := libraryEntries(t, librariesDir)
	requests := stubHeaderDownload(t, http.StatusOK, "header contents")

	output, err := executeDoctor(t, newDoctorCommand(librariesDir))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "elf.h (missing)") || !strings.Contains(output, "doctor --fix") {
		t.Fatalf("doctor did not report missing elf.h and fix option: %s", output)
	}
	if *requests != 0 {
		t.Fatalf("ordinary doctor made %d HTTP requests", *requests)
	}
	if after := libraryEntries(t, librariesDir); !reflect.DeepEqual(before, after) {
		t.Fatalf("libraries directory changed: before=%v, after=%v", before, after)
	}
}

// Audit E14: the headers ship inside elmos, so --fix writes them without
// the network and restores every missing one, not only elf.h.
func TestDoctorFixWritesEmbeddedHeaders(t *testing.T) {
	librariesDir := prepareDoctorLibraries(t)
	requests := stubHeaderDownload(t, http.StatusServiceUnavailable, "")

	output, err := executeDoctor(t, newDoctorCommand(librariesDir), "--fix")
	if err != nil {
		t.Fatal(err)
	}
	if *requests != 0 {
		t.Fatalf("fix made %d HTTP requests, want none", *requests)
	}
	for _, name := range []string{"elf.h", "endian.h", "asm/types.h"} {
		if _, err := os.Stat(filepath.Join(librariesDir, name)); err != nil {
			t.Fatalf("%s not written: %v\n%s", name, err, output)
		}
	}
	content, err := os.ReadFile(filepath.Join(librariesDir, "byteswap.h"))
	if err != nil || string(content) != "existing" {
		t.Fatalf("existing byteswap.h overwritten: %q %v", content, err)
	}
	if !strings.Contains(output, "host headers written") || strings.Contains(output, "elf.h (missing)") {
		t.Fatalf("fix output: %s", output)
	}
}

func prepareDoctorLibraries(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "asm"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "byteswap.h"), []byte("existing"), 0644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func libraryEntries(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}
	return names
}

type doctorRoundTrip func(*http.Request) (*http.Response, error)

func (roundTrip doctorRoundTrip) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTrip(request)
}

func stubHeaderDownload(t *testing.T, status int, body string) *int {
	t.Helper()
	requests := new(int)
	previous := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: doctorRoundTrip(func(request *http.Request) (*http.Response, error) {
		*requests++
		if request.URL.Host != "raw.githubusercontent.com" {
			t.Errorf("unexpected download host: %s", request.URL.Host)
		}
		return &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
			Request:    request,
		}, nil
	})}
	t.Cleanup(func() { http.DefaultClient = previous })
	return requests
}

func newDoctorCommand(librariesDir string) *cobra.Command {
	cfg := &config.Config{Paths: config.PathsConfig{LibrariesDir: librariesDir}}
	fs := filesystem.NewOSFileSystem()
	exec := executor.NewMockExecutor()
	printer := ui.NewPrinter()
	tm := toolchain.NewManager(exec, fs, cfg, printer)
	ctx := &Context{
		HealthChecker: doctor.NewHealthChecker(exec, fs, cfg, tm),
		AutoFixer:     doctor.NewAutoFixer(fs, cfg),
		Printer:       printer,
	}
	cmd := BuildDoctor(ctx)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	return cmd
}

func executeDoctor(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	previous := os.Stdout
	os.Stdout = writer
	defer func() {
		os.Stdout = previous
		_ = reader.Close()
		_ = writer.Close()
	}()
	cmd.SetArgs(args)
	runErr := cmd.Execute()
	_ = writer.Close()
	os.Stdout = previous
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	return string(output), runErr
}
