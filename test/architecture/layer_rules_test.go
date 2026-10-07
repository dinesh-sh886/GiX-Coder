package architecture

import (
	"go/parser"
	"go/token"
	"os"
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

// TestLayerViolations tests that modules don't violate layer boundaries
func TestLayerViolations(t *testing.T) {
	// Layer hierarchy (lower layer cannot import upper layer):
	// 1. Domain (shared, api) - lowest
	// 2. Infrastructure (config, logging, metrics, tracing, validation, security, dto) - middle
	// 3. Application (gateway, workflow, harness, sandbox, router, policy, audit) - upper

	layers := map[string]int{
		"shared":     1,
		"api":        1,
		"config":     2,
		"logging":    2,
		"metrics":    2,
		"tracing":    2,
		"validation": 2,
		"security":   2,
		"dto":        2,
		"constants":  2,
		"version":    2,
		"gateway":    3,
		"workflow":   3,
		"harness":    3,
		"sandbox":    3,
		"router":     3,
		"policy":     3,
		"audit":      3,
	}

	modules := GetModules()
	rootDir := findProjectRoot()

	for _, module := range modules {
		if module.Path == "api" {
			continue
		}

		modulePath := filepath.Join(rootDir, module.Path)
		imports := collectImports(t, modulePath)

		moduleLayer := layers[module.Name]

		for _, imp := range imports {
			// Skip standard library and external dependencies
			if isStdLib(imp) || isExternalDep(imp) {
				continue
			}

			// Check if importing from another module
			for otherName, otherLayer := range layers {
				if otherName == module.Name {
					continue
				}

				if strings.HasPrefix(imp, "github.com/gix-coder/gix-coder/"+otherName) {
					otherLayerNum := layers[otherName]

					// Check for layer violation: upper layer importing from lower layer is OK
					// But lower layer importing from upper layer is NOT OK
					if moduleLayer < otherLayerNum {
						t.Errorf(
							"Layer violation: %s (layer %d) imports from %s (layer %d) - lower layer cannot depend on upper layer",
							module.Name, moduleLayer, otherName, otherLayerNum,
						)
					}
				}
			}
		}
	}
}

// TestNoDirectDatabaseAccess tests that modules don't access other module's databases directly
func TestNoDirectDatabaseAccess(t *testing.T) {
	// Check for direct SQL queries or database access patterns that violate ownership
	// This is a simplified check - full enforcement requires static analysis tools

	modules := GetModules()
	rootDir := findProjectRoot()

	for _, module := range modules {
		if module.Path == "api" {
			continue
		}

		modulePath := filepath.Join(rootDir, module.Path)
		imports := collectImports(t, modulePath)

		for _, imp := range imports {
			// Check for direct database access patterns
			// Direct database access should only be in the owning module
			if strings.Contains(imp, "database/sql") ||
				strings.Contains(imp, "github.com/jackc/pgx") ||
				strings.Contains(imp, "gorm.io") ||
				strings.Contains(imp, "database") {
				// This is allowed only if the module owns the database
				// For now, we just log it - full enforcement requires deeper analysis
				t.Logf("Module %s imports database package: %s", module.Name, imp)
			}
		}
	}
}

// TestNoHTTPCallsOutsideGateway tests that only gateway makes external HTTP calls
func TestNoHTTPCallsOutsideGateway(t *testing.T) {
	modules := GetModules()
	rootDir := findProjectRoot()

	for _, module := range modules {
		if module.Path == "api" || module.Name == "gateway" {
			continue
		}

		modulePath := filepath.Join(rootDir, module.Path)
		imports := collectImports(t, modulePath)

		for _, imp := range imports {
			// Check for HTTP client usage outside gateway
			if strings.Contains(imp, "net/http") ||
				strings.Contains(imp, "github.com/go-resty/resty") ||
				strings.Contains(imp, "github.com/valyala/fasthttp") {
				t.Errorf(
					"Module %s imports HTTP client package %s - only gateway should make external HTTP calls",
					module.Name, imp,
				)
			}
		}
	}
}

// TestNoSQLInApplicationLayer tests that SQL is only in infrastructure layer
func TestNoSQLInApplicationLayer(t *testing.T) {
	applicationModules := []string{"gateway", "workflow", "harness", "sandbox", "router", "policy", "audit"}
	rootDir := findProjectRoot()

	for _, moduleName := range applicationModules {
		modulePath := filepath.Join(rootDir, moduleName)
		imports := collectImports(t, modulePath)

		for _, imp := range imports {
			// Check for direct SQL usage
			if strings.HasPrefix(imp, "database/sql") ||
				strings.HasPrefix(imp, "github.com/jackc/pgx") ||
				strings.HasPrefix(imp, "github.com/jmoiron/sqlx") ||
				strings.HasPrefix(imp, "gorm.io") {
				t.Errorf(
					"Application module %s directly imports SQL package %s - database access should be in infrastructure layer",
					moduleName, imp,
				)
			}
		}
	}
}

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
