package builder

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/NguyenTrongPhuc552003/elmos/assets"
)

// CreateModule creates a new module from template.
func (m *ModuleBuilder) CreateModule(name string) error {
	modPath := filepath.Join(m.cfg.Paths.ModulesDir, name)

	// Check if already exists
	if m.fs.Exists(modPath) {
		return fmt.Errorf("module already exists: %s", name)
	}

	// Create directory
	if err := m.fs.MkdirAll(modPath, 0755); err != nil {
		return err
	}

	// Convert name to valid C identifier (replace dashes with underscores)
	cName := strings.ReplaceAll(name, "-", "_")

	// Template data
	data := struct {
		Name        string
		CName       string
		Description string
	}{
		Name:        name,
		CName:       cName,
		Description: "A simple kernel module",
	}

	// Load and execute source template
	srcTmpl, err := assets.GetModuleTemplate()
	if err != nil {
		return fmt.Errorf("failed to load module template: %w", err)
	}

	srcContent, err := executeModuleTemplate("module.c", string(srcTmpl), data)
	if err != nil {
		return fmt.Errorf("failed to execute module template: %w", err)
	}

	srcPath := filepath.Join(modPath, cName+".c")
	if err := m.fs.WriteFile(srcPath, []byte(srcContent), 0644); err != nil {
		return err
	}

	// Load and execute Makefile template
	makeTmpl, err := assets.GetModuleMakefile()
	if err != nil {
		return fmt.Errorf("failed to load module makefile template: %w", err)
	}

	makeContent, err := executeModuleTemplate("Makefile", string(makeTmpl), data)
	if err != nil {
		return fmt.Errorf("failed to execute module makefile template: %w", err)
	}

	makePath := filepath.Join(modPath, "Makefile")
	return m.fs.WriteFile(makePath, []byte(makeContent), 0644)
}

// executeModuleTemplate executes a Go template with the given data.
func executeModuleTemplate(name, tmplContent string, data interface{}) (string, error) {
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
