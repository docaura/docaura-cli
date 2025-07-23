package analyzer

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

// Analyzer analyzes Go source code and extracts documentation information.
type Analyzer struct {
	fset *token.FileSet
}

// New creates a new code analyzer instance.
func New() *Analyzer {
	return &Analyzer{
		fset: token.NewFileSet(),
	}
}

// AnalyzePackage analyzes a Go package in the specified directory and returns
// comprehensive package information including functions, types, constants,
// variables, and documentation.
func (a *Analyzer) AnalyzePackage(dir string) (*PackageInfo, error) {
	pkgs, err := parser.ParseDir(a.fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse directory %q: %w", dir, err)
	}

	pkg := findMainPackage(pkgs)
	if pkg == nil {
		return nil, fmt.Errorf("no Go package found in %q", dir)
	}

	// Create documentation from parsed package
	docPkg := doc.New(pkg, "./", 0)

	info := &PackageInfo{
		Name:        pkg.Name,
		Path:        dir,
		Description: cleanDoc(docPkg.Doc),
		Imports:     extractImports(pkg),
	}

	a.populatePackageInfo(info, docPkg)

	return info, nil
}

// populatePackageInfo populates the PackageInfo with data from the doc.Package.
func (a *Analyzer) populatePackageInfo(info *PackageInfo, docPkg *doc.Package) {
	// Analyze functions
	for _, fn := range docPkg.Funcs {
		fnInfo := a.analyzeFunctionDecl(fn)
		info.Functions = append(info.Functions, fnInfo)
	}

	// Analyze types and their methods
	for _, typ := range docPkg.Types {
		typeInfo := a.analyzeTypeDecl(typ)
		info.Types = append(info.Types, typeInfo)

		// Add methods to functions list
		for _, method := range typ.Methods {
			methodInfo := a.analyzeFunctionDecl(method)
			methodInfo.IsMethod = true
			methodInfo.Receiver = typ.Name
			info.Functions = append(info.Functions, methodInfo)
		}
	}

	// Analyze constants
	for _, c := range docPkg.Consts {
		constInfo := a.analyzeConstantDecl(c)
		info.Constants = append(info.Constants, constInfo...)
	}

	// Analyze variables
	for _, v := range docPkg.Vars {
		varInfo := a.analyzeVariableDecl(v)
		info.Variables = append(info.Variables, varInfo...)
	}
}

// analyzeFunctionDecl analyzes a function declaration and returns function information.
func (a *Analyzer) analyzeFunctionDecl(fn *doc.Func) FunctionInfo {
	info := FunctionInfo{
		Name:        fn.Name,
		Description: cleanDoc(fn.Doc),
		IsExported:  ast.IsExported(fn.Name),
		Examples:    extractExamplesFromDoc(fn.Doc),
	}

	if fn.Decl != nil && fn.Decl.Type != nil {
		info.Signature = a.getFunctionSignature(fn.Decl)
		info.Parameters = extractParameters(fn.Decl.Type.Params)
		info.Returns = extractReturns(fn.Decl.Type.Results)
	}

	return info
}

// analyzeTypeDecl analyzes a type declaration and returns type information.
func (a *Analyzer) analyzeTypeDecl(typ *doc.Type) TypeInfo {
	info := TypeInfo{
		Name:        typ.Name,
		Description: cleanDoc(typ.Doc),
		IsExported:  ast.IsExported(typ.Name),
	}

	if typ.Decl != nil {
		for _, spec := range typ.Decl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			info.Kind = getTypeKind(ts.Type)
			if structType, ok := ts.Type.(*ast.StructType); ok {
				info.Fields = extractStructFields(structType)
			}
		}
	}

	// Extract method names
	for _, method := range typ.Methods {
		info.Methods = append(info.Methods, method.Name)
	}

	return info
}

// analyzeConstantDecl analyzes a constant declaration and returns constant information.
func (a *Analyzer) analyzeConstantDecl(c *doc.Value) []ConstantInfo {
	var constants []ConstantInfo

	for _, spec := range c.Decl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for i, name := range vs.Names {
			constInfo := ConstantInfo{
				Name:        name.Name,
				Description: cleanDoc(c.Doc),
				IsExported:  ast.IsExported(name.Name),
			}

			if vs.Type != nil {
				constInfo.Type = typeToString(vs.Type)
			}

			if i < len(vs.Values) && vs.Values[i] != nil {
				constInfo.Value = exprToString(vs.Values[i])
			}

			constants = append(constants, constInfo)
		}
	}

	return constants
}

// analyzeVariableDecl analyzes a variable declaration and returns variable information.
func (a *Analyzer) analyzeVariableDecl(v *doc.Value) []VariableInfo {
	var variables []VariableInfo

	for _, spec := range v.Decl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for _, name := range vs.Names {
			varInfo := VariableInfo{
				Name:        name.Name,
				Description: cleanDoc(v.Doc),
				IsExported:  ast.IsExported(name.Name),
			}

			if vs.Type != nil {
				varInfo.Type = typeToString(vs.Type)
			}

			variables = append(variables, varInfo)
		}
	}

	return variables
}

// getFunctionSignature builds a function signature string from an AST function declaration.
func (a *Analyzer) getFunctionSignature(decl *ast.FuncDecl) string {
	var parts []string

	parts = append(parts, "func")

	// Add receiver if it's a method
	if decl.Recv != nil {
		recv := fieldListToString(decl.Recv)
		parts = append(parts, fmt.Sprintf("(%s)", recv))
	}

	// Add function name
	parts = append(parts, decl.Name.Name)

	// Add parameters
	if decl.Type.Params != nil {
		params := fieldListToString(decl.Type.Params)
		parts = append(parts, fmt.Sprintf("(%s)", params))
	} else {
		parts = append(parts, "()")
	}

	// Add return types
	if decl.Type.Results != nil {
		results := fieldListToString(decl.Type.Results)
		if len(decl.Type.Results.List) == 1 && len(decl.Type.Results.List[0].Names) == 0 {
			parts = append(parts, results)
		} else {
			parts = append(parts, fmt.Sprintf("(%s)", results))
		}
	}

	return strings.Join(parts, " ")
}

// findMainPackage finds the main (non-test) package from a map of packages.
func findMainPackage(pkgs map[string]*ast.Package) *ast.Package {
	for name, pkg := range pkgs {
		if !strings.HasSuffix(name, "_test") {
			return pkg
		}
	}
	return nil
}
