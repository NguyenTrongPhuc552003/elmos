package homebrew

import (
	"testing"
)

func TestClearCache(t *testing.T) {
	r := &Resolver{
		cache: map[string]string{
			"llvm": "/opt/homebrew/opt/llvm",
			"qemu": "/opt/homebrew/opt/qemu",
		},
	}

	if len(r.cache) != 2 {
		t.Errorf("cache len = %d, want 2", len(r.cache))
	}

	r.ClearCache()

	if len(r.cache) != 0 {
		t.Errorf("after ClearCache, cache len = %d, want 0", len(r.cache))
	}
}

func TestGetBinWithCachedPrefix(t *testing.T) {
	r := &Resolver{
		cache: map[string]string{
			"llvm": "/opt/homebrew/opt/llvm",
		},
	}

	got := r.GetBin("llvm")
	want := "/opt/homebrew/opt/llvm/bin"

	if got != want {
		t.Errorf("GetBin() = %q, want %q", got, want)
	}
}

func TestGetSbinWithCachedPrefix(t *testing.T) {
	r := &Resolver{
		cache: map[string]string{
			"e2fsprogs": "/opt/homebrew/opt/e2fsprogs",
		},
	}

	got := r.GetSbin("e2fsprogs")
	want := "/opt/homebrew/opt/e2fsprogs/sbin"

	if got != want {
		t.Errorf("GetSbin() = %q, want %q", got, want)
	}
}

func TestGetIncludeWithCachedPrefix(t *testing.T) {
	r := &Resolver{
		cache: map[string]string{
			"libelf": "/opt/homebrew/opt/libelf",
		},
	}

	got := r.GetInclude("libelf")
	want := "/opt/homebrew/opt/libelf/include"

	if got != want {
		t.Errorf("GetInclude() = %q, want %q", got, want)
	}
}

func TestGetLibWithCachedPrefix(t *testing.T) {
	r := &Resolver{
		cache: map[string]string{
			"libelf": "/opt/homebrew/opt/libelf",
		},
	}

	got := r.GetLib("libelf")
	want := "/opt/homebrew/opt/libelf/lib"

	if got != want {
		t.Errorf("GetLib() = %q, want %q", got, want)
	}
}

func TestGetLibexecBinWithCachedPrefix(t *testing.T) {
	r := &Resolver{
		cache: map[string]string{
			"gnu-sed": "/opt/homebrew/opt/gnu-sed",
		},
	}

	got := r.GetLibexecBin("gnu-sed")
	want := "/opt/homebrew/opt/gnu-sed/libexec/gnubin"

	if got != want {
		t.Errorf("GetLibexecBin() = %q, want %q", got, want)
	}
}

func TestGetBinEmptyPrefix(t *testing.T) {
	// When prefix is already cached as empty, should return ""
	r := &Resolver{
		cache: map[string]string{
			"nonexistent": "",
		},
	}

	got := r.GetBin("nonexistent")
	if got != "" {
		t.Errorf("GetBin() with empty prefix = %q, want empty", got)
	}
}
