// Package assets provides embedded template files for elmos.
package assets

import "embed"

//go:embed templates/*
var Templates embed.FS

// Patches holds the built-in kernel patches, laid out as
// patches/<kernel ref>/<arch or generic>/NNNN-name.patch.
//
//go:embed patches
var Patches embed.FS

// Libraries holds the Linux headers that macOS lacks (elf.h, byteswap.h,
// endian.h, asm/*), needed by the kernel host tools. They are written into
// the project's libraries directory, which the kernel build adds to
// HOSTCFLAGS.
//
//go:embed libraries
var Libraries embed.FS

// GetModuleTemplate returns the module source template.
func GetModuleTemplate() ([]byte, error) {
	return Templates.ReadFile("templates/module/module.c.tmpl")
}

// GetModuleMakefile returns the module Makefile template.
func GetModuleMakefile() ([]byte, error) {
	return Templates.ReadFile("templates/module/Makefile.tmpl")
}

// GetAppTemplate returns the app source template.
func GetAppTemplate() ([]byte, error) {
	return Templates.ReadFile("templates/app/main.c.tmpl")
}

// GetAppMakefile returns the app Makefile template.
func GetAppMakefile() ([]byte, error) {
	return Templates.ReadFile("templates/app/Makefile.tmpl")
}

// GetInitScript returns the init script template.
func GetInitScript() ([]byte, error) {
	return Templates.ReadFile("templates/init/init.sh.tmpl")
}

// GetGuestSync returns the guest sync script template.
func GetGuestSync() ([]byte, error) {
	return Templates.ReadFile("templates/init/guesync.sh.tmpl")
}

// GetConfigTemplate returns the elmos.yaml configuration template.
func GetConfigTemplate() ([]byte, error) {
	return Templates.ReadFile("templates/configs/elmos.yaml.tmpl")
}
