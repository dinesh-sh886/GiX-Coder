package architecture

import (
	"path/filepath"
	"strings"
	"testing"
)

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
			for otherName := range layers {
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
