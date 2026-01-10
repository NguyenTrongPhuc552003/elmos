// Package context provides build context management for elmos.
package context

import (
	"reflect"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *Error
		want string
	}{
		// TODO: Add test cases.
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
	tests := []struct {
		name    string
		e       *Error
		wantErr bool
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.codeString(); got != tt.want {
				t.Errorf("Error.codeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigError(t *testing.T) {
	type args struct {
		msg   string
		cause error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConfigError(tt.args.msg, tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageError(t *testing.T) {
	type args struct {
		msg   string
		cause error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ImageError(tt.args.msg, tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImageError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoError(t *testing.T) {
	type args struct {
		msg   string
		cause error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RepoError(tt.args.msg, tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RepoError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildError(t *testing.T) {
	type args struct {
		msg   string
		cause error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildError(tt.args.msg, tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQEMUError(t *testing.T) {
	type args struct {
		msg   string
		cause error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QEMUError(tt.args.msg, tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QEMUError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleError(t *testing.T) {
	type args struct {
		msg   string
		cause error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ModuleError(tt.args.msg, tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModuleError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDependencyError(t *testing.T) {
	type args struct {
		msg   string
		cause error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DependencyError(tt.args.msg, tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DependencyError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionError(t *testing.T) {
	type args struct {
		msg   string
		cause error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PermissionError(tt.args.msg, tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PermissionError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppError(t *testing.T) {
	type args struct {
		msg   string
		cause error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppError(tt.args.msg, tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRootfsError(t *testing.T) {
	type args struct {
		msg   string
		cause error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RootfsError(tt.args.msg, tt.args.cause); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RootfsError() = %v, want %v", got, tt.want)
			}
		})
	}
}
