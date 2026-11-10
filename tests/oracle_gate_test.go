package test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOracleQualityGate is a meta-test over the test suite itself. It encodes
// the repository's "accuracy and oracle" rules and fails the build when the
// suite drifts into weaker assertions:
//
//  1. Every Test* function must make at least one value-level assertion
//     (assert.Equal, require.Contains, assert.Regexp, ...) — existence-only
//     oracles such as NotNil/NotEmpty/NoError cannot detect a wrong value and
//     are therefore not counted.
//  2. Plan-only tests (those built on SharedPlan) must never call
//     terraform.Destroy: a shared hermetic plan has nothing to destroy, and a
//     stray destroy against the real backend could tear down live
//     infrastructure.
//  3. Tests that apply real resources (terraform.InitAndApply) must call
//     integrationSkip so `go test -short` (unit mode, pre-commit hook, CI) can
//     never create GCP resources.
func TestOracleQualityGate(t *testing.T) {
	t.Parallel()

	strongAssertions := map[string]bool{
		"Equal": true, "NotEqual": true,
		"Contains": true, "NotContains": true,
		"Regexp": true, "MatchRegexp": true, "NotRegexp": true,
		"Len":    true,
		"Subset": true, "ElementsMatch": true,
		"Greater": true, "GreaterOrEqual": true,
		"Less": true, "LessOrEqual": true,
		"InDelta": true, "JSONEq": true,
		"True": true, "False": true,
		// Deliberately NOT counted as value-level oracles: they only prove
		// that something exists / a call succeeded, never that it is correct.
		"Nil": false, "NotNil": false,
		"Error": false, "NoError": false,
		"Empty": false, "NotEmpty": false,
	}

	files, err := filepath.Glob("*_test.go")
	require.NoError(t, err, "must be able to list test files")
	require.NotEmpty(t, files, "expected *_test.go files in tests/")

	var violations []string

	for _, file := range files {
		fset := token.NewFileSet()
		parsed, err := parser.ParseFile(fset, file, nil, parser.SkipObjectResolution)
		require.NoError(t, err, "test file %s must parse", file)

		for _, decl := range parsed.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "Test") || fn.Name.Name == "TestMain" {
				continue
			}

			var (
				hasStrong       bool
				usesSharedPlan  bool
				usesDestroy     bool
				usesApply       bool
				usesIntegration bool
			)

			ast.Inspect(fn, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				// Resolve callee as pkg.Sel or ident.
				var pkg, sel string
				switch callee := call.Fun.(type) {
				case *ast.SelectorExpr:
					if ident, ok := callee.X.(*ast.Ident); ok {
						pkg, sel = ident.Name, callee.Sel.Name
					}
				case *ast.Ident:
					sel = callee.Name
				}

				if (pkg == "assert" || pkg == "require") && strongAssertions[sel] {
					hasStrong = true
				}
				if pkg == "terraform" {
					switch sel {
					case "Destroy":
						usesDestroy = true
					case "InitAndApply", "InitAndApplyE":
						usesApply = true
					}
				}
				if sel == "SharedPlan" {
					usesSharedPlan = true
				}
				if sel == "integrationSkip" {
					usesIntegration = true
				}
				return true
			})

			name := fn.Name.Name
			if !hasStrong {
				violations = append(violations,
					file+":"+name+": no value-level assertion found; add an exact-value oracle (assert.Equal/Contains/Regexp/...) instead of presence-only checks")
			}
			if usesSharedPlan && usesDestroy {
				violations = append(violations,
					file+":"+name+": plan-only test calls terraform.Destroy; shared plan tests must never destroy")
			}
			if usesApply && !usesIntegration {
				violations = append(violations,
					file+":"+name+": applies real resources but never calls integrationSkip; -short runs must not touch GCP")
			}
		}
	}

	assert.Len(t, violations, 0, "oracle quality gate failed:\n  %s", strings.Join(violations, "\n  "))
}
