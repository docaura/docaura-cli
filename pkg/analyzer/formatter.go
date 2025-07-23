package analyzer

import (
	"fmt"
	"go/ast"
	"strings"
)

// fieldListToString converts a field list to a string representation.
func fieldListToString(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}

	parts := make([]string, 0, len(fields.List))
	for _, field := range fields.List {
		fieldType := typeToString(field.Type)
		if len(field.Names) == 0 {
			parts = append(parts, fieldType)
		} else {
			for _, name := range field.Names {
				parts = append(parts, fmt.Sprintf("%s %s", name.Name, fieldType))
			}
		}
	}

	return strings.Join(parts, ", ")
}

// typeToString converts an AST expression representing a type to its string representation.
func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", typeToString(t.Key), typeToString(t.Value))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", typeToString(t.X), t.Sel.Name)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType:
		switch t.Dir {
		case ast.SEND:
			return "chan<- " + typeToString(t.Value)
		case ast.RECV:
			return "<-chan " + typeToString(t.Value)
		default:
			return "chan " + typeToString(t.Value)
		}
	case *ast.FuncType:
		return buildFuncTypeString(t)
	case *ast.Ellipsis:
		return "..." + typeToString(t.Elt)
	default:
		return "unknown"
	}
}

// buildFuncTypeString builds a string representation of a function type.
func buildFuncTypeString(ft *ast.FuncType) string {
	var parts []string
	parts = append(parts, "func")

	if ft.Params != nil {
		params := fieldListToString(ft.Params)
		parts = append(parts, fmt.Sprintf("(%s)", params))
	} else {
		parts = append(parts, "()")
	}

	if ft.Results != nil {
		results := fieldListToString(ft.Results)
		if len(ft.Results.List) == 1 && len(ft.Results.List[0].Names) == 0 {
			parts = append(parts, results)
		} else {
			parts = append(parts, fmt.Sprintf("(%s)", results))
		}
	}

	return strings.Join(parts, " ")
}

// exprToString converts an AST expression to its string representation.
func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return e.Value
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", exprToString(e.X), e.Sel.Name)
	case *ast.CallExpr:
		return exprToString(e.Fun) + "(...)"
	default:
		return "..."
	}
}

// getTypeKind determines the kind of a type expression.
func getTypeKind(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.StructType:
		return "struct"
	case *ast.InterfaceType:
		return "interface"
	case *ast.ArrayType:
		return "array"
	case *ast.MapType:
		return "map"
	case *ast.ChanType:
		return "channel"
	case *ast.FuncType:
		return "function"
	default:
		return "alias"
	}
}

// cleanDoc cleans and normalizes documentation strings.
func cleanDoc(doc string) string {
	if doc == "" {
		return ""
	}

	// Normalize line endings and trim whitespace
	doc = strings.TrimSpace(doc)
	doc = strings.ReplaceAll(doc, "\r\n", "\n")

	// Split into lines and remove empty lines
	lines := strings.Split(doc, "\n")
	cleaned := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}

	return strings.Join(cleaned, "\n")
}
