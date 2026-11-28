// Package developer holds no code: this test keeps the component diagram in
// diagrams.md true to the code. It lists every Go package of the module,
// sorts each into the diagram's group by path, and compares the packages
// named in each group's label and the imports between groups with the
// diagram's nodes and arrows.
package developer

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

const module = "github.com/NguyenTrongPhuc552003/elmos/"

// groups maps a package path prefix to the diagram node that holds it; the
// longest matching prefix wins.
var groups = map[string]string{
	"cmd/":         "entry",
	"sdk/":         "sdk",
	"core/app":     "app",
	"core/ui":      "ui",
	"core/domain/": "domain",
	"core/context": "workspace",
	"core/config":  "workspace",
	"core/infra/":  "infra",
	"assets":       "assets",
	"third_party/": "thirdparty",
}

func groupOf(pkg string) string {
	best := ""
	for prefix := range groups {
		if strings.HasPrefix(pkg, prefix) && len(prefix) > len(best) {
			best = prefix
		}
	}
	return groups[best]
}

// packages returns each package directory of the module (relative to the
// root) with the module packages it imports, test files excluded.
func packages(t *testing.T, root string) map[string]map[string]bool {
	t.Helper()
	out := map[string]map[string]bool{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "testdata" || d.Name() == "node_modules") {
			return filepath.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, filepath.Dir(path))
		if err != nil {
			return err
		}
		pkg := filepath.ToSlash(rel)
		if out[pkg] == nil {
			out[pkg] = map[string]bool{}
		}
		for _, imp := range f.Imports {
			if p := strings.Trim(imp.Path.Value, `"`); strings.HasPrefix(p, module) {
				out[pkg][strings.TrimPrefix(p, module)] = true
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

var (
	nodeRe = regexp.MustCompile(`^\s*(\w+)\["([^"]*)"\]\s*$`)
	edgeRe = regexp.MustCompile(`^\s*(\w+)\s*-->\s*(\w+)\s*$`)
)

// diagram returns the component diagram's nodes (with the package names in
// their labels) and arrows.
func diagram(t *testing.T) (map[string][]string, map[string]bool) {
	t.Helper()
	text, err := os.ReadFile("diagrams.md")
	if err != nil {
		t.Fatal(err)
	}
	// a Windows checkout may have CRLF line ends
	s := strings.ReplaceAll(string(text), "\r\n", "\n")
	start := strings.Index(s, "## Components")
	if start < 0 {
		t.Fatal("diagrams.md has no \"## Components\" section")
	}
	s = s[start:]
	begin := strings.Index(s, "```mermaid\n")
	if begin < 0 {
		t.Fatal("no mermaid block under \"## Components\"")
	}
	end := strings.Index(s[begin+1:], "```")
	if end < 0 {
		t.Fatal("the mermaid block under \"## Components\" is not closed")
	}
	nodes := map[string][]string{}
	edges := map[string]bool{}
	for _, line := range strings.Split(s[begin:begin+1+end], "\n") {
		if m := nodeRe.FindStringSubmatch(line); m != nil {
			// "Title<br/>prefix/ name, name,<br/>name": the names follow
			// the first <br/> and an optional path prefix ending in "/".
			parts := strings.SplitN(m[2], "<br/>", 2)
			if len(parts) < 2 {
				t.Fatalf("node %s: label has no package list: %q", m[1], m[2])
			}
			list := strings.ReplaceAll(parts[1], "<br/>", " ")
			if f := strings.Fields(list); len(f) > 0 && strings.HasSuffix(f[0], "/") {
				list = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(list), f[0]))
			}
			for _, name := range strings.Split(list, ",") {
				nodes[m[1]] = append(nodes[m[1]], strings.TrimSpace(name))
			}
		} else if m := edgeRe.FindStringSubmatch(line); m != nil {
			edges[m[1]+" --> "+m[2]] = true
		}
	}
	return nodes, edges
}

func sorted(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func TestComponentDiagramMatchesCode(t *testing.T) {
	pkgs := packages(t, filepath.Join("..", ".."))
	nodes, edges := diagram(t)

	wantNames := map[string]map[string]bool{}
	wantEdges := map[string]bool{}
	for pkg, imports := range pkgs {
		g := groupOf(pkg)
		if g == "" {
			t.Errorf("package %s belongs to no group: add its prefix to groups", pkg)
			continue
		}
		if wantNames[g] == nil {
			wantNames[g] = map[string]bool{}
		}
		wantNames[g][filepath.Base(pkg)] = true
		for imp := range imports {
			if h := groupOf(imp); h != "" && h != g {
				wantEdges[g+" --> "+h] = true
			}
		}
	}

	for g, want := range wantNames {
		got := map[string]bool{}
		for _, n := range nodes[g] {
			got[n] = true
		}
		for _, n := range sorted(want) {
			if !got[n] {
				t.Errorf("node %s: package %s is missing from its label", g, n)
			}
		}
		for _, n := range sorted(got) {
			if !want[n] {
				t.Errorf("node %s: its label names %s, which is no package in this group", g, n)
			}
		}
	}
	for g := range nodes {
		if wantNames[g] == nil {
			t.Errorf("node %s holds no package", g)
		}
	}
	for _, e := range sorted(wantEdges) {
		if !edges[e] {
			t.Errorf("arrow missing: %s (the code has this import)", e)
		}
	}
	for _, e := range sorted(edges) {
		if !wantEdges[e] {
			t.Errorf("arrow %s: no import between these groups", e)
		}
	}
}
