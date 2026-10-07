package architecture

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func findProjectRoot() string {
	dir, _ := filepath.Abs(".")
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

func collectImports(t *testing.T, dir string) []string {
	var imports []string

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		if strings.Contains(path, "migrations") {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return nil
		}

		for _, imp := range f.Imports {
			impPath := strings.Trim(imp.Path.Value, `"`)
			imports = append(imports, impPath)
		}

		return nil
	})

	return imports
}

func isStdLib(imp string) bool {
	stdLibs := []string{
		"archive", "bufio", "bytes", "compress", "container", "context",
		"crypto", "database", "debug", "encoding", "errors", "expvar",
		"flag", "fmt", "hash", "html", "image", "index", "io", "log",
		"math", "mime", "net", "os", "path", "plugin", "reflect", "regexp",
		"runtime", "sort", "strconv", "strings", "sync", "syscall", "testing",
		"text", "time", "unicode", "unsafe",
	}

	for _, std := range stdLibs {
		if strings.HasPrefix(imp, std) {
			return true
		}
	}
	return false
}

func isExternalDep(imp string) bool {
	return strings.Contains(imp, ".")
}
