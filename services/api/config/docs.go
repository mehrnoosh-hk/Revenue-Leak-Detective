package config

import (
	"fmt"
	"reflect"
	"strings"
)

// FieldInfo contains information about a configuration field
type FieldInfo struct {
	Name        string
	Type        string
	Description string
	Default     string
	Required    bool
	Environment string
	Example     string
	Validation  string
	Options     []string
}

// ConfigDocs generates documentation for the configuration system
func ConfigDocs() string {
	var docs strings.Builder

	docs.WriteString("# Configuration Documentation\n\n")
	docs.WriteString("This documentation is auto-generated from the configuration structs.\n\n")

	// Generate docs for each config struct
	docs.WriteString(generateStructDocs("HTTPConfig", reflect.TypeOf(HTTPConfig{})))
	docs.WriteString(generateStructDocs("DatabaseConfig", reflect.TypeOf(DatabaseConfig{})))
	docs.WriteString(generateStructDocs("EnvironmentConfig", reflect.TypeOf(EnvironmentConfig{})))
	docs.WriteString(generateStructDocs("BuildInfoConfig", reflect.TypeOf(BuildInfoConfig{})))

	return docs.String()
}

// generateStructDocs generates documentation for a specific struct
func generateStructDocs(structName string, t reflect.Type) string {
	var docs strings.Builder

	docs.WriteString(fmt.Sprintf("## %s\n\n", structName))

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		info := extractFieldInfo(field)
		docs.WriteString(generateFieldDoc(info))
	}

	docs.WriteString("\n")
	return docs.String()
}

// extractFieldInfo extracts information from a struct field
func extractFieldInfo(field reflect.StructField) FieldInfo {
	info := FieldInfo{
		Name: field.Name,
		Type: field.Type.String(),
	}

	// Extract description from comments
	comments := strings.Split(field.Tag.Get("comment"), "\n")
	if len(comments) > 0 {
		info.Description = strings.TrimSpace(comments[0])
	}

	// Extract other information from comments
	for _, comment := range comments {
		comment = strings.TrimSpace(comment)
		if strings.HasPrefix(comment, "Default:") {
			info.Default = strings.TrimSpace(strings.TrimPrefix(comment, "Default:"))
		} else if strings.HasPrefix(comment, "Environment variable:") {
			info.Environment = strings.TrimSpace(strings.TrimPrefix(comment, "Environment variable:"))
		} else if strings.HasPrefix(comment, "Options:") {
			options := strings.TrimSpace(strings.TrimPrefix(comment, "Options:"))
			info.Options = strings.Split(options, ", ")
		} else if strings.HasPrefix(comment, "Must be between") {
			info.Validation = comment
		}
	}

	// Extract validation tags
	if validate := field.Tag.Get("validate"); validate != "" {
		info.Validation = validate
		if strings.Contains(validate, "required") {
			info.Required = true
		}
	}

	// Extract example
	if example := field.Tag.Get("example"); example != "" {
		info.Example = example
	}

	return info
}

// generateFieldDoc generates documentation for a single field
func generateFieldDoc(info FieldInfo) string {
	var doc strings.Builder

	doc.WriteString(fmt.Sprintf("### %s\n\n", info.Name))

	if info.Description != "" {
		doc.WriteString(fmt.Sprintf("**Description:** %s\n\n", info.Description))
	}

	doc.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", info.Type))

	if info.Required {
		doc.WriteString("**Required:** Yes\n\n")
	} else {
		doc.WriteString("**Required:** No\n\n")
	}

	if info.Default != "" {
		doc.WriteString(fmt.Sprintf("**Default:** `%s`\n\n", info.Default))
	}

	if info.Environment != "" {
		doc.WriteString(fmt.Sprintf("**Environment Variable:** `%s`\n\n", info.Environment))
	}

	if info.Example != "" {
		doc.WriteString(fmt.Sprintf("**Example:** `%s`\n\n", info.Example))
	}

	if info.Validation != "" {
		doc.WriteString(fmt.Sprintf("**Validation:** %s\n\n", info.Validation))
	}

	if len(info.Options) > 0 {
		doc.WriteString(fmt.Sprintf("**Options:** %s\n\n", strings.Join(info.Options, ", ")))
	}

	doc.WriteString("---\n\n")
	return doc.String()
}

// GenerateExampleConfig generates an example configuration
func GenerateExampleConfig() string {
	example := `# Example Configuration

## HTTP Configuration
API_HOST=0.0.0.0
API_PORT=3030

## Database Configuration
# Option 1: Using individual parameters
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=revenue_leak_detective_dev
POSTGRES_SSL=disable

# Option 2: Using connection URL (takes precedence)
# POSTGRES_URL=postgresql://user:pass@host:5432/dbname?sslmode=disable

## Environment Configuration
ENVIRONMENT=development
LOG_LEVEL=INFO
DEBUG=false
CONFIG_VERSION=1.0.0

## Build Information (auto-populated)
GIT_COMMIT_HASH=a1b2c3d
GIT_COMMIT_FULL=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0
GIT_COMMIT_DATE=2024-01-15T10:30:00Z
GIT_COMMIT_DATE_SHORT=2024-01-15
GIT_COMMIT_MESSAGE=feat: add new configuration system
GIT_BRANCH=main
GIT_TAG=v1.0.0
GIT_DIRTY=false
BUILD_TIMESTAMP=2024-01-15T10:30:00Z`

	return example
}
