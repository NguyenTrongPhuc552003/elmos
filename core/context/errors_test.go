// Package context provides build context management for elmos.
package context

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	baseCause := errors.New("underlying cause")
	tests := []struct {
		name string
		e    *Error
		want string
	}{
		{
			name: "Error without cause",
			e:    &Error{Code: ErrCodeConfig, Message: "config invalid"},
			want: "[CONFIG] config invalid",
		},
		{
			name: "Error with cause",
			e:    &Error{Code: ErrCodeBuild, Message: "build failed", Cause: baseCause},
			want: "[BUILD] build failed: underlying cause",
		},
		{
			name: "Unknown error code",
			e:    &Error{Code: 999, Message: "unknown error"},
			want: "[ERROR] unknown error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	baseCause := errors.New("cause error")
	tests := []struct {
		name    string
		e       *Error
		wantErr bool
	}{
		{
			name:    "Error with cause",
			e:       &Error{Code: ErrCodeConfig, Message: "msg", Cause: baseCause},
			wantErr: true,
		},
		{
			name:    "Error without cause",
			e:       &Error{Code: ErrCodeConfig, Message: "msg", Cause: nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("Error.Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestError_codeString(t *testing.T) {
	tests := []struct {
		name string
		e    *Error
		want string
	}{
		{name: "CONFIG", e: &Error{Code: ErrCodeConfig}, want: "CONFIG"},
		{name: "IMAGE", e: &Error{Code: ErrCodeImage}, want: "IMAGE"},
		{name: "REPO", e: &Error{Code: ErrCodeRepo}, want: "REPO"},
		{name: "BUILD", e: &Error{Code: ErrCodeBuild}, want: "BUILD"},
		{name: "QEMU", e: &Error{Code: ErrCodeQEMU}, want: "QEMU"},
		{name: "MODULE", e: &Error{Code: ErrCodeModule}, want: "MODULE"},
		{name: "DEP", e: &Error{Code: ErrCodeDependency}, want: "DEP"},
		{name: "PERM", e: &Error{Code: ErrCodePermission}, want: "PERM"},
		{name: "APP", e: &Error{Code: ErrCodeApp}, want: "APP"},
		{name: "ROOTFS", e: &Error{Code: ErrCodeRootfs}, want: "ROOTFS"},
		{name: "Unknown", e: &Error{Code: 999}, want: "ERROR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.codeString(); got != tt.want {
				t.Errorf("Error.codeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// assertError checks error fields without using reflect.DeepEqual
func assertError(t *testing.T, got *Error, wantCode ErrCode, wantMsg string, wantCause error) {
	t.Helper()
	if got.Code != wantCode {
		t.Errorf("Code = %v, want %v", got.Code, wantCode)
	}
	if got.Message != wantMsg {
		t.Errorf("Message = %v, want %v", got.Message, wantMsg)
	}
	if got.Cause != wantCause {
		t.Errorf("Cause = %v, want %v", got.Cause, wantCause)
	}
}

func TestConfigError(t *testing.T) {
	cause := errors.New("io error")
	tests := []struct {
		name      string
		msg       string
		cause     error
		wantCode  ErrCode
		wantMsg   string
		wantCause error
	}{
		{
			name:      "With cause",
			msg:       "config load failed",
			cause:     cause,
			wantCode:  ErrCodeConfig,
			wantMsg:   "config load failed",
			wantCause: cause,
		},
		{
			name:      "Without cause",
			msg:       "missing config",
			cause:     nil,
			wantCode:  ErrCodeConfig,
			wantMsg:   "missing config",
			wantCause: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConfigError(tt.msg, tt.cause)
			assertError(t, got, tt.wantCode, tt.wantMsg, tt.wantCause)
		})
	}
}

func TestImageError(t *testing.T) {
	cause := errors.New("disk full")
	got := ImageError("create failed", cause)
	assertError(t, got, ErrCodeImage, "create failed", cause)

	got = ImageError("not found", nil)
	assertError(t, got, ErrCodeImage, "not found", nil)
}

func TestRepoError(t *testing.T) {
	cause := errors.New("git error")
	got := RepoError("clone failed", cause)
	assertError(t, got, ErrCodeRepo, "clone failed", cause)

	got = RepoError("not a git repo", nil)
	assertError(t, got, ErrCodeRepo, "not a git repo", nil)
}

func TestBuildError(t *testing.T) {
	cause := errors.New("compiler error")
	got := BuildError("kernel build failed", cause)
	assertError(t, got, ErrCodeBuild, "kernel build failed", cause)

	got = BuildError("no config", nil)
	assertError(t, got, ErrCodeBuild, "no config", nil)
}

func TestQEMUError(t *testing.T) {
	cause := errors.New("qemu exited")
	got := QEMUError("emulation failed", cause)
	assertError(t, got, ErrCodeQEMU, "emulation failed", cause)

	got = QEMUError("not installed", nil)
	assertError(t, got, ErrCodeQEMU, "not installed", nil)
}

func TestModuleError(t *testing.T) {
	cause := errors.New("link error")
	got := ModuleError("module build failed", cause)
	assertError(t, got, ErrCodeModule, "module build failed", cause)

	got = ModuleError("module not found", nil)
	assertError(t, got, ErrCodeModule, "module not found", nil)
}

func TestDependencyError(t *testing.T) {
	cause := errors.New("brew error")
	got := DependencyError("llvm missing", cause)
	assertError(t, got, ErrCodeDependency, "llvm missing", cause)

	got = DependencyError("qemu not installed", nil)
	assertError(t, got, ErrCodeDependency, "qemu not installed", nil)
}

func TestPermissionError(t *testing.T) {
	cause := errors.New("access denied")
	got := PermissionError("cannot write", cause)
	assertError(t, got, ErrCodePermission, "cannot write", cause)

	got = PermissionError("root required", nil)
	assertError(t, got, ErrCodePermission, "root required", nil)
}

func TestAppError(t *testing.T) {
	cause := errors.New("link error")
	got := AppError("app build failed", cause)
	assertError(t, got, ErrCodeApp, "app build failed", cause)

	got = AppError("app not found", nil)
	assertError(t, got, ErrCodeApp, "app not found", nil)
}

func TestRootfsError(t *testing.T) {
	cause := errors.New("debootstrap failed")
	got := RootfsError("rootfs creation failed", cause)
	assertError(t, got, ErrCodeRootfs, "rootfs creation failed", cause)

	got = RootfsError("rootfs not found", nil)
	assertError(t, got, ErrCodeRootfs, "rootfs not found", nil)
}
