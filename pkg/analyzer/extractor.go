package analyzer

import (
	"go/ast"
	"strings"
)

// extractImports extracts all import paths from a package.
func extractImports(pkg *ast.Package) []string {
	importSet := make(map[string]bool)

	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			importSet[path] = true
		}
	}

	imports := make([]string, 0, len(importSet))
	for imp := range importSet {
		imports = append(imports, imp)
	}

	return imports
}

// extractParameters extracts parameter information from a field list.
func extractParameters(fields *ast.FieldList) []ParameterInfo {
	if fields == nil {
		return nil
	}

	var params []ParameterInfo
	for _, field := range fields.List {
		paramType := typeToString(field.Type)

		if len(field.Names) == 0 {
			// Anonymous parameter
			params = append(params, ParameterInfo{
				Type: paramType,
			})
		} else {
			for _, name := range field.Names {
				params = append(params, ParameterInfo{
					Name: name.Name,
					Type: paramType,
				})
			}
		}
	}

	return params
}

// extractReturns extracts return type information from a field list.
func extractReturns(fields *ast.FieldList) []ReturnInfo {
	if fields == nil {
		return nil
	}

	returns := make([]ReturnInfo, 0, len(fields.List))
	for _, field := range fields.List {
		returns = append(returns, ReturnInfo{
			Type: typeToString(field.Type),
		})
	}

	return returns
}

// extractStructFields extracts field information from a struct type.
func extractStructFields(structType *ast.StructType) []FieldInfo {
	if structType.Fields == nil {
		return nil
	}

	var fields []FieldInfo

	for _, field := range structType.Fields.List {
		fieldType := typeToString(field.Type)
		var tag string
		if field.Tag != nil {
			tag = field.Tag.Value
		}

		if len(field.Names) == 0 {
			// Embedded field
			fields = append(fields, FieldInfo{
				Type: fieldType,
				Tag:  tag,
			})
		} else {
			for _, name := range field.Names {
				fields = append(fields, FieldInfo{
					Name: name.Name,
					Type: fieldType,
					Tag:  tag,
				})
			}
		}
	}

	return fields
}

// extractExamplesFromDoc extracts code examples from documentation comments.
func extractExamplesFromDoc(doc string) []string {
	if doc == "" {
		return nil
	}

	var examples []string
	lines := strings.Split(doc, "\n")

	var inExample bool
	var currentExample strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for example start markers
		if strings.HasPrefix(trimmed, "Example:") ||
			strings.HasPrefix(trimmed, "Usage:") ||
			strings.Contains(trimmed, "```go") {
			inExample = true
			currentExample.Reset()
			continue
		}

		if inExample {
			// Check for example end markers
			if strings.Contains(trimmed, "```") ||
				(trimmed == "" && currentExample.Len() > 0) {
				if currentExample.Len() > 0 {
					examples = append(examples, currentExample.String())
					currentExample.Reset()
				}
				inExample = false
				continue
			}

			// Extract indented code lines
			if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
				cleaned := strings.TrimPrefix(strings.TrimPrefix(line, "    "), "\t")
				currentExample.WriteString(cleaned)
				currentExample.WriteString("\n")
			}
		}
	}

	return examples
}
