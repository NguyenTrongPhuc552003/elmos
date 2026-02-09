package builder

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/NguyenTrongPhuc552003/elmos/assets"
)

// CreateApp creates a new application from template.
func (a *AppBuilder) CreateApp(name string) error {
	appPath := filepath.Join(a.cfg.Paths.AppsDir, name)

	// Ensure apps directory exists
	if err := a.fs.MkdirAll(a.cfg.Paths.AppsDir, 0755); err != nil {
		return err
	}

	// Check if already exists
	if a.fs.Exists(appPath) {
		return fmt.Errorf("app already exists: %s", name)
	}

	// Create directory
	if err := a.fs.MkdirAll(appPath, 0755); err != nil {
		return err
	}

	// Convert name to valid C identifier (replace dashes with underscores)
	cName := strings.ReplaceAll(name, "-", "_")

	// Template data
	data := struct {
		Name  string
		CName string
	}{
		Name:  name,
		CName: cName,
	}

	// Load and execute source template
	srcTmpl, err := assets.GetAppTemplate()
	if err != nil {
		return fmt.Errorf("failed to load app template: %w", err)
	}

	srcContent, err := executeAppTemplate("app.c", string(srcTmpl), data)
	if err != nil {
		return fmt.Errorf("failed to execute app template: %w", err)
	}

	srcPath := filepath.Join(appPath, cName+".c")
	if err := a.fs.WriteFile(srcPath, []byte(srcContent), 0644); err != nil {
		return err
	}

	// Load and execute Makefile template
	makeTmpl, err := assets.GetAppMakefile()
	if err != nil {
		return fmt.Errorf("failed to load makefile template: %w", err)
	}

	makeContent, err := executeAppTemplate("Makefile", string(makeTmpl), data)
	if err != nil {
		return fmt.Errorf("failed to execute makefile template: %w", err)
	}

	makePath := filepath.Join(appPath, "Makefile")
	return a.fs.WriteFile(makePath, []byte(makeContent), 0644)
}

// executeAppTemplate executes a Go template with the given data.
func executeAppTemplate(name, tmplContent string, data interface{}) (string, error) {
	tmpl, err := template.New(name).Parse(tmplContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
