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

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/doctor"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
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

func TestDoctorFixDownloadsElfHeader(t *testing.T) {
	librariesDir := prepareDoctorLibraries(t)
	const header = "/* test elf.h */\n"
	requests := stubHeaderDownload(t, http.StatusOK, header)

	output, err := executeDoctor(t, newDoctorCommand(librariesDir), "--fix")
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(librariesDir, "elf.h"))
	if err != nil || string(content) != header {
		t.Fatalf("written elf.h = %q, err=%v", content, err)
	}
	if *requests != 1 || !strings.Contains(output, "elf.h downloaded") {
		t.Fatalf("fix did not report one successful download; requests=%d, output=%s", *requests, output)
	}
}

func TestDoctorFixReportsDownloadFailure(t *testing.T) {
	librariesDir := prepareDoctorLibraries(t)
	requests := stubHeaderDownload(t, http.StatusServiceUnavailable, "")

	output, err := executeDoctor(t, newDoctorCommand(librariesDir), "--fix")
	if err == nil || !strings.Contains(err.Error(), "HTTP 503") {
		t.Fatalf("download error = %v, want HTTP 503", err)
	}
	if !strings.Contains(output, "elf.h (missing)") {
		t.Fatalf("doctor did not report the missing header after download failure: %s", output)
	}
	if *requests != 1 {
		t.Fatalf("fix made %d HTTP requests, want 1", *requests)
	}
	if _, err := os.Stat(filepath.Join(librariesDir, "elf.h")); !os.IsNotExist(err) {
		t.Fatalf("elf.h should be absent after failure, stat error: %v", err)
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
