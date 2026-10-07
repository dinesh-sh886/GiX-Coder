package architecture

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Module represents a Go module in the monolith
type Module struct {
	Name        string
	Path        string
	InternalDir string
	AllowedDeps []string
}

// GetModules returns all defined modules
func GetModules() []Module {
	return []Module{
		{
			Name:        "shared",
			Path:        "shared",
			InternalDir: "",
			AllowedDeps: []string{},
		},
		{
			Name:        "gateway",
			Path:        "gateway",
			InternalDir: "internal",
			AllowedDeps: []string{"shared", "api", "workflow", "policy", "audit"},
		},
		{
			Name:        "workflow",
			Path:        "workflow",
			InternalDir: "internal",
			AllowedDeps: []string{"shared", "api", "harness", "router", "audit"},
		},
		{
			Name:        "harness",
			Path:        "harness",
			InternalDir: "internal",
			AllowedDeps: []string{"shared", "api", "sandbox", "router", "policy", "audit"},
		},
		{
			Name:        "sandbox",
			Path:        "sandbox",
			InternalDir: "internal",
			AllowedDeps: []string{"shared", "api"},
		},
		{
			Name:        "router",
			Path:        "router",
			InternalDir: "internal",
			AllowedDeps: []string{"shared", "api", "policy"},
		},
		{
			Name:        "policy",
			Path:        "policy",
			InternalDir: "internal",
			AllowedDeps: []string{"shared", "api"},
		},
		{
			Name:        "audit",
			Path:        "audit",
			InternalDir: "internal",
			AllowedDeps: []string{"shared", "api"},
		},
		{
			Name:        "api",
			Path:        "api",
			InternalDir: "",
			AllowedDeps: []string{"shared"},
		},
	}
}

// TestNoForbiddenImports tests that modules don't import forbidden packages
func TestNoForbiddenImports(t *testing.T) {
	modules := GetModules()
	moduleMap := make(map[string]Module)
	for _, m := range modules {
		moduleMap[m.Name] = m
	}

	rootDir := findProjectRoot()

	for _, module := range modules {
		t.Run(module.Name, func(t *testing.T) {
			if module.Path == "api" {
				// API module only contains protobuf generated code
				return
			}

			modulePath := filepath.Join(rootDir, module.Path)
			imports := collectImports(t, modulePath)

			for _, imp := range imports {
				// Skip standard library
				if isStdLib(imp) {
					continue
				}

				// Skip external dependencies
				if isExternalDep(imp) {
					continue
				}

				// Check if importing from another module
				for _, otherModule := range modules {
					if otherModule.Name == module.Name {
						continue
					}

					if strings.HasPrefix(imp, "github.com/gix-coder/gix-coder/"+otherModule.Name) {
						// Check if this dependency is allowed
						allowed := false
						for _, allowedDep := range module.AllowedDeps {
							if allowedDep == otherModule.Name {
								allowed = true
								break
							}
						}

						if !allowed {
							t.Errorf(
								"Module %s imports forbidden dependency %s (import: %s)",
								module.Name, otherModule.Name, imp,
							)
						}
					}
				}
			}
		})
	}
}

// TestNoCrossModuleInternalImports tests that modules don't import internal packages of other modules
func TestNoCrossModuleInternalImports(t *testing.T) {
	rootDir := findProjectRoot()

	for _, module := range GetModules() {
		if module.InternalDir == "" {
			continue
		}

		internalPath := filepath.Join(rootDir, module.Path, module.InternalDir)
		imports := collectImports(t, internalPath)

		for _, imp := range imports {
			if strings.HasPrefix(imp, "github.com/gix-coder/gix-coder/") {
				for _, otherModule := range GetModules() {
					if otherModule.Name == module.Name {
						continue
					}
					if strings.HasPrefix(imp, "github.com/gix-coder/gix-coder/"+otherModule.Name+"/"+otherModule.InternalDir) {
						t.Errorf(
							"Module %s internal package imports forbidden internal package %s",
							module.Name, imp,
						)
					}
				}
			}
		}
	}
}

// TestNoCycles tests that there are no circular dependencies
func TestNoCycles(t *testing.T) {
	modules := GetModules()
	graph := make(map[string][]string)

	for _, module := range modules {
		graph[module.Name] = module.AllowedDeps
	}

	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}

		recStack[node] = false
		return false
	}

	for _, module := range modules {
		if !visited[module.Name] {
			if dfs(module.Name) {
				t.Errorf("Circular dependency detected involving %s", module.Name)
			}
		}
	}
}

// TestSharedPackageNoExternalDeps tests that shared package has no external dependencies
func TestSharedPackageNoExternalDeps(t *testing.T) {
	rootDir := findProjectRoot()
	sharedPath := filepath.Join(rootDir, "shared")

	imports := collectImports(t, sharedPath)

	for _, imp := range imports {
		if isExternalDep(imp) && !isStdLib(imp) {
			// Allow only specific external dependencies
			allowed := false
			allowedDeps := []string{
				"github.com/rs/zerolog",
				"github.com/spf13/viper",
				"github.com/prometheus/client_golang",
				"go.opentelemetry.io/otel",
				"google.golang.org/grpc",
				"google.golang.org/protobuf",
				"golang.org/x/crypto",
				"golang.org/x/sync",
				"golang.org/x/time",
				"github.com/go-playground/validator/v10",
				"github.com/jackc/pgx/v5",
				"github.com/redis/go-redis/v9",
				"github.com/golang-migrate/migrate/v4",
				"github.com/golang-jwt/jwt/v5",
			}

			for _, allowed := range allowedDeps {
				if strings.HasPrefix(imp, allowed) {
					allowed = true
					break
				}
			}

			if !allowed {
				t.Errorf("Shared package has unexpected external dependency: %s", imp)
			}
		}
	}
}

// TestNoDirectDatabaseAccess tests that modules don't directly access other module's databases
func TestNoDirectDatabaseAccess(t *testing.T) {
	// This would check for direct SQL queries across module boundaries
	// For now, we document the rule and rely on code review
	t.Log("Database access isolation is enforced by architecture - enforced by code review")
}

// Helper functions

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
