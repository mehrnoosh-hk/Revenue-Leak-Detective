//go:build ignore

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type QueryStructInfo struct {
	Name     string
	Fields   []FieldInfo
	FileName string
	Entity   string // e.g., "user", "event" extracted from filename
}

type FieldInfo struct {
	Name     string
	Type     string
	JSONTag  string
	Comments string
}

const queryStructTemplate = `
// {{.Name}} represents query-specific parameters/results for {{.Entity}} operations
// This struct is auto-generated from sqlc queries and should be reviewed before using
type {{.Name}} struct {
{{range .Fields}}	{{.Name}} {{.Type}} {{.JSONTag}}{{if .Comments}} // {{.Comments}}{{end}}
{{end}}}
`

func main() {
	sqlcDir := "internal/db/sqlc"
	domainDir := "internal/domain/models"

	// Find all *.sql.go files in sqlc directory
	sqlFiles, err := findSqlGoFiles(sqlcDir)
	if err != nil {
		fmt.Printf("Error finding sql.go files: %v\n", err)
		os.Exit(1)
	}

	// Parse all query structs from *.sql.go files
	allQueryStructs := make(map[string][]QueryStructInfo) // key: entity, value: structs
	for _, file := range sqlFiles {
		queryStructs, err := parseQueryStructs(file)
		if err != nil {
			fmt.Printf("Error parsing query structs from %s: %v\n", file, err)
			continue
		}

		// Group structs by entity
		for _, structInfo := range queryStructs {
			allQueryStructs[structInfo.Entity] = append(allQueryStructs[structInfo.Entity], structInfo)
		}
	}

	// Add query structs to existing domain model files
	totalProcessed := 0
	for entity, structs := range allQueryStructs {
		domainFile := filepath.Join(domainDir, entity+".go")
		err = addQueryStructsToFile(domainFile, structs)
		if err != nil {
			fmt.Printf("Error adding query structs to %s: %v\n", domainFile, err)
			continue
		}
		totalProcessed += len(structs)
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Successfully added %d query structs across %d entities\n", totalProcessed, len(allQueryStructs))
}

// findSqlGoFiles finds all *.sql.go files in the given directory
func findSqlGoFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql.go") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// parseQueryStructs extracts query-specific structs from a *.sql.go file
func parseQueryStructs(filePath string) ([]QueryStructInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var queryStructs []QueryStructInfo
	fileName := filepath.Base(filePath)
	entity := extractEntityFromFileName(fileName)

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if structType, ok := x.Type.(*ast.StructType); ok {
				structName := x.Name.Name

				// Filter for query-specific structs (not enum wrappers or base models)
				if isQueryStruct(structName) {
					structInfo := QueryStructInfo{
						Name:     structName,
						Fields:   parseFields(structType),
						FileName: fileName,
						Entity:   entity,
					}
					queryStructs = append(queryStructs, structInfo)
				}
			}
		}
		return true
	})

	return queryStructs, nil
}

// isQueryStruct determines if a struct is a query-specific struct
func isQueryStruct(structName string) bool {
	// Skip enum wrapper types and base model types
	if strings.HasPrefix(structName, "Null") {
		return false
	}

	// Include structs that end with common query patterns
	queryPatterns := []string{
		"Params",   // CreateUserParams, UpdateEventParams
		"Row",      // GetUserByIdRow, CreateUserRow
		"Result",   // Custom result types
		"Response", // Custom response types
	}

	for _, pattern := range queryPatterns {
		if strings.HasSuffix(structName, pattern) {
			return true
		}
	}

	return false
}

// extractEntityFromFileName extracts entity name from filename like "users.sql.go" -> "user"
func extractEntityFromFileName(fileName string) string {
	// Remove .sql.go extension
	name := strings.TrimSuffix(fileName, ".sql.go")

	// Convert plural to singular (basic rules)
	if strings.HasSuffix(name, "s") && len(name) > 1 {
		name = strings.TrimSuffix(name, "s")
	}

	return name
}

// parseFields parses struct fields and converts types appropriately
func parseFields(structType *ast.StructType) []FieldInfo {
	var fields []FieldInfo

	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			fieldType := getFieldType(field.Type)
			jsonTag := getJSONTag(field.Tag)
			convertedType := convertType(fieldType)

			fields = append(fields, FieldInfo{
				Name:    name.Name,
				Type:    convertedType,
				JSONTag: jsonTag,
			})
		}
	}

	return fields
}

// getFieldType extracts the type from an AST expression
func getFieldType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", getFieldType(t.X), t.Sel.Name)
	case *ast.ArrayType:
		elemType := getFieldType(t.Elt)
		return "[]" + elemType
	case *ast.StarExpr:
		elemType := getFieldType(t.X)
		return "*" + elemType
	default:
		return "interface{}"
	}
}

// convertType converts sqlc types to domain types
func convertType(sqlcType string) string {
	conversions := map[string]string{
		"pgtype.UUID":        "uuid.UUID",
		"pgtype.Timestamptz": "time.Time",
		"pgtype.Numeric":     "float32",
		"string":             "string",
		"int32":              "int32",
		"int64":              "int64",
		"bool":               "bool",
		"[]byte":             "interface{}", // For JSON data fields
		"*string":            "*string",
		"*int32":             "*int32",
		"*int64":             "*int64",
		"*bool":              "*bool",
	}

	if converted, ok := conversions[sqlcType]; ok {
		return converted
	}

	// Handle nullable enum types
	if strings.HasPrefix(sqlcType, "Null") && strings.HasSuffix(sqlcType, "Enum") {
		baseType := strings.TrimPrefix(sqlcType, "Null")
		return "*" + baseType
	}

	return sqlcType
}

// getJSONTag extracts JSON tag from struct field tag
func getJSONTag(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}

	tagStr := tag.Value
	tagStr = strings.Trim(tagStr, "`")

	// Extract json tag
	parts := strings.Split(tagStr, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "json:") {
			return "`" + part + "`"
		}
	}

	return ""
}

// addQueryStructsToFile adds query structs to an existing domain model file
func addQueryStructsToFile(filePath string, queryStructs []QueryStructInfo) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("domain model file %s does not exist", filePath)
	}

	// Generate new struct definitions
	var newStructsContent strings.Builder
	tmpl, err := template.New("queryStruct").Parse(queryStructTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	for _, structInfo := range queryStructs {
		err = tmpl.Execute(&newStructsContent, structInfo)
		if err != nil {
			return fmt.Errorf("failed to execute template for %s: %w", structInfo.Name, err)
		}
	}

	// Append new structs to existing file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(newStructsContent.String())
	if err != nil {
		return fmt.Errorf("failed to write new structs: %w", err)
	}

	return nil
}
