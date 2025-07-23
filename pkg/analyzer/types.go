package analyzer

type PackageInfo struct {
	Name        string         `json:"name"`
	Path        string         `json:"path"`
	Description string         `json:"description"`
	Functions   []FunctionInfo `json:"functions"`
	Types       []TypeInfo     `json:"types"`
	Constants   []ConstantInfo `json:"constants"`
	Variables   []VariableInfo `json:"variables"`
	Examples    []ExampleInfo  `json:"examples"`
	Imports     []string       `json:"imports"`
}

// FunctionInfo represents information about a function or method.
type FunctionInfo struct {
	Name        string          `json:"name"`
	Signature   string          `json:"signature"`
	Description string          `json:"description"`
	Parameters  []ParameterInfo `json:"parameters"`
	Returns     []ReturnInfo    `json:"returns"`
	Examples    []string        `json:"examples"`
	IsExported  bool            `json:"is_exported"`
	IsMethod    bool            `json:"is_method"`
	Receiver    string          `json:"receiver,omitempty"`
}

// TypeInfo represents information about a type declaration.
type TypeInfo struct {
	Name        string      `json:"name"`
	Kind        string      `json:"kind"` // struct, interface, alias, etc.
	Description string      `json:"description"`
	Fields      []FieldInfo `json:"fields,omitempty"`
	Methods     []string    `json:"methods,omitempty"`
	IsExported  bool        `json:"is_exported"`
}

// FieldInfo represents information about a struct field.
type FieldInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Tag         string `json:"tag,omitempty"`
	Description string `json:"description"`
}

// ParameterInfo represents information about a function parameter.
type ParameterInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ReturnInfo represents information about a function return value.
type ReturnInfo struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ConstantInfo represents information about a constant declaration.
type ConstantInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
	IsExported  bool   `json:"is_exported"`
}

// VariableInfo represents information about a variable declaration.
type VariableInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	IsExported  bool   `json:"is_exported"`
}

// ExampleInfo represents information about a code example.
type ExampleInfo struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Doc  string `json:"doc"`
}
