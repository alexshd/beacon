package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type Candidate struct {
	FuncName     string
	TypeName     string
	IsComparable bool
	NeedsWrapper bool
	Receiver     string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: lawtest-gen <file.go>")
		fmt.Println()
		fmt.Println("Analyzes Go source file and generates lawtest skeleton")
		os.Exit(1)
	}

	filename := os.Args[1]
	candidates, pkgName, err := analyzeFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(candidates) == 0 {
		fmt.Println("No lawtest candidates found.")
		fmt.Println()
		fmt.Println("lawtest works with:")
		fmt.Println("  • func(T, T) T                - Binary functions")
		fmt.Println("  • func (T) Method(T) T        - Methods on types")
		fmt.Println("  • func ([]T) Method([]T) []T  - Methods on slices")
		os.Exit(0)
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  lawtest Candidates Found")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	for i, c := range candidates {
		fmt.Printf("%d. %s\n", i+1, c.FuncName)
		fmt.Printf("   Type: %s\n", c.TypeName)
		if c.NeedsWrapper {
			fmt.Printf("   ⚠️  Type is NOT comparable - needs wrapper (see example)\n")
		} else if c.IsComparable {
			fmt.Printf("   ✅ Type is comparable\n")
		} else {
			fmt.Printf("   ❓ Comparability unknown - may need wrapper\n")
		}
		fmt.Println()
	}

	// Generate test file
	testFilename := strings.TrimSuffix(filename, ".go") + "_law_test.go"
	content := generateTestFile(pkgName, candidates)

	err = os.WriteFile(testFilename, []byte(content), 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing test file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated: %s\n", testFilename)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Review generated tests")
	fmt.Println("  2. Complete TODO items")
	fmt.Println("  3. Verify operations should have tested properties")
	fmt.Println("  4. Run: go test -v")
	fmt.Println()
}

func analyzeFile(filename string) ([]Candidate, string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, "", err
	}

	var candidates []Candidate
	pkgName := node.Name.Name

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if c := analyzeFunc(x); c != nil {
				candidates = append(candidates, *c)
			}
		}
		return true
	})

	return candidates, pkgName, nil
}

func analyzeFunc(fn *ast.FuncDecl) *Candidate {
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return nil
	}

	// Check for func(T, T) T pattern
	if fn.Recv == nil && len(fn.Type.Params.List) == 2 {
		param1Type := exprToString(fn.Type.Params.List[0].Type)
		param2Type := exprToString(fn.Type.Params.List[1].Type)
		returnType := exprToString(fn.Type.Results.List[0].Type)

		if param1Type == param2Type && param1Type == returnType {
			return &Candidate{
				FuncName:     fn.Name.Name,
				TypeName:     param1Type,
				IsComparable: isLikelyComparable(param1Type),
				NeedsWrapper: isNonComparable(param1Type),
			}
		}
	}

	// Check for method(T) T pattern (with receiver)
	// Handles both: (T) Method(T) T and ([]T) Method([]T) []T
	if fn.Recv != nil && len(fn.Type.Params.List) == 1 {
		receiverType := exprToString(fn.Recv.List[0].Type)
		paramType := exprToString(fn.Type.Params.List[0].Type)
		returnType := exprToString(fn.Type.Results.List[0].Type)

		if receiverType == paramType && receiverType == returnType {
			return &Candidate{
				FuncName:     fn.Name.Name,
				TypeName:     receiverType,
				IsComparable: isLikelyComparable(receiverType),
				NeedsWrapper: isNonComparable(receiverType),
				Receiver:     receiverType,
			}
		}
	}

	return nil
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	default:
		return "unknown"
	}
}

func isLikelyComparable(typeName string) bool {
	// Basic comparable types
	comparable := []string{
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "string", "bool", "byte", "rune",
	}

	for _, c := range comparable {
		if typeName == c {
			return true
		}
	}

	// Pointers are comparable
	if strings.HasPrefix(typeName, "*") {
		return true
	}

	return false
}

func isNonComparable(typeName string) bool {
	// Known non-comparable types
	if strings.HasPrefix(typeName, "[]") {
		return true // slices
	}
	if strings.HasPrefix(typeName, "map[") {
		return true // maps
	}
	if strings.HasPrefix(typeName, "func(") {
		return true // functions
	}
	return false
}

func generateTestFile(pkgName string, candidates []Candidate) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "package %s\n\n", pkgName)
	sb.WriteString("import (\n")
	sb.WriteString("\t\"testing\"\n\n")
	sb.WriteString("\t\"github.com/alexshd/lawtest\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString("// This file was auto-generated by lawtest-gen\n")
	sb.WriteString("// Review and complete the TODO items before running tests\n\n")

	for _, c := range candidates {
		if c.NeedsWrapper {
			generateWrapperTests(&sb, c)
		} else {
			generateDirectTests(&sb, c)
		}
	}

	return sb.String()
}

func generateDirectTests(sb *strings.Builder, c Candidate) {
	funcToTest := c.FuncName
	if c.Receiver != "" {
		// For methods, we need to show pattern
		fmt.Fprintf(sb, "// TODO: Create wrapper function for method %s\n", c.FuncName)
		fmt.Fprintf(sb, "// func Wrap%s(a, b %s) %s {\n", c.FuncName, c.Receiver, c.Receiver)
		fmt.Fprintf(sb, "//     return a.%s(b)\n", c.FuncName)
		sb.WriteString("// }\n\n")
		funcToTest = "Wrap" + c.FuncName
	}

	// Immutability test
	fmt.Fprintf(sb, "func Test%sImmutability(t *testing.T) {\n", c.FuncName)
	fmt.Fprintf(sb, "\t// TODO: Create a generator for %s\n", c.TypeName)
	sb.WriteString("\tgen := func() " + c.TypeName + " {\n")
	sb.WriteString("\t\t// TODO: Return a valid instance of " + c.TypeName + "\n")
	sb.WriteString("\t\tpanic(\"implement generator\")\n")
	sb.WriteString("\t}\n\n")
	fmt.Fprintf(sb, "\tlawtest.ImmutableOp(t, %s, gen)\n", funcToTest)
	sb.WriteString("}\n\n")

	// Associativity test
	fmt.Fprintf(sb, "func Test%sAssociativity(t *testing.T) {\n", c.FuncName)
	fmt.Fprintf(sb, "\t// TODO: Verify that %s SHOULD be associative\n", c.FuncName)
	sb.WriteString("\t// Question: Does (a op b) op c = a op (b op c) for your operation?\n")
	sb.WriteString("\t// If NO, remove this test - not all operations are associative\n\n")
	sb.WriteString("\tgen := func() " + c.TypeName + " {\n")
	sb.WriteString("\t\t// TODO: Return a valid instance of " + c.TypeName + "\n")
	sb.WriteString("\t\tpanic(\"implement generator\")\n")
	sb.WriteString("\t}\n\n")
	fmt.Fprintf(sb, "\tlawtest.Associative(t, %s, gen)\n", funcToTest)
	sb.WriteString("}\n\n")
}

func generateWrapperTests(sb *strings.Builder, c Candidate) {
	wrapperName := c.TypeName + "Wrapper"
	fmt.Fprintf(sb, "// %s is NOT comparable - needs wrapper for lawtest\n", c.TypeName)
	sb.WriteString("// See config-merge-example for pattern\n\n")

	sb.WriteString("// TODO: Create wrapper type in your source file:\n")
	fmt.Fprintf(sb, "// type %s struct {\n", wrapperName)
	fmt.Fprintf(sb, "//     data %s\n", c.TypeName)
	sb.WriteString("// }\n\n")

	sb.WriteString("// TODO: Create wrapper function:\n")
	fmt.Fprintf(sb, "// func Wrap%s(a, b *%s) *%s {\n", c.FuncName, wrapperName, wrapperName)
	fmt.Fprintf(sb, "//     return &%s{data: %s(a.data, b.data)}\n", wrapperName, c.FuncName)
	sb.WriteString("// }\n\n")

	fmt.Fprintf(sb, "func Test%sImmutability(t *testing.T) {\n", c.FuncName)
	sb.WriteString("\t// TODO: Implement after creating wrapper type\n")
	sb.WriteString("\tt.Skip(\"TODO: Create wrapper type first - see comments above\")\n")
	sb.WriteString("}\n\n")

	fmt.Fprintf(sb, "func Test%sAssociativity(t *testing.T) {\n", c.FuncName)
	sb.WriteString("\t// TODO: Implement after creating wrapper type\n")
	sb.WriteString("\tt.Skip(\"TODO: Create wrapper type first - see comments above\")\n")
	sb.WriteString("}\n\n")
}
