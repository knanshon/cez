package schemas

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "strings"
)

// Manager holds the loaded JSON schemas.
type Manager struct {
    LoadedSchemas map[string]json.RawMessage
}

// NewManager loads schemas from the specified base directory and returns a new Manager.
// This function should be called once at application startup.
func NewManager(schemaBaseDir string) (*Manager, error) {
    m := &Manager{
        LoadedSchemas: make(map[string]json.RawMessage),
    }

    // Iterate through known schema paths. For POC, this is hardcoded.
    // In a real app, you might scan for all subdirectories.
    schemaDirsToLoad := []string{
        filepath.Join(schemaBaseDir, "greeter", "v1"),
        // Add other schema directories here if you expand your APIs
    }

    for _, currentSchemaDir := range schemaDirsToLoad {
        files, err := ioutil.ReadDir(currentSchemaDir)
        if err != nil {
            // Log a warning or error, depending on whether missing schemas are fatal
            fmt.Printf("Warning: Failed to read schema directory %s: %v\n", currentSchemaDir, err)
            continue // Continue to next directory
        }

        for _, file := range files {
            if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
                schemaName := strings.TrimSuffix(file.Name(), ".json") // e.g., "GreetRequest"

                // Dynamically construct the fully qualified schema name based on the directory structure
                // e.g., "web/schemas/greeter/v1/GreetRequest.json" -> "greeter.v1.GreetRequest"
                relativePath, err := filepath.Rel(schemaBaseDir, currentSchemaDir)
                if err != nil {
                    return nil, fmt.Errorf("failed to get relative path for %s: %w", currentSchemaDir, err)
                }
                // Replace path separators with dots for the package part
                protoPackageName := strings.ReplaceAll(relativePath, string(filepath.Separator), ".")

                fullyQualifiedName := fmt.Sprintf("%s.%s", protoPackageName, schemaName)

                filePath := filepath.Join(currentSchemaDir, file.Name())

                data, err := ioutil.ReadFile(filePath)
                if err != nil {
                    return nil, fmt.Errorf("failed to read schema file %s: %w", filePath, err)
                }
                m.LoadedSchemas[fullyQualifiedName] = data
                fmt.Printf("Loaded schema: %s from %s\n", fullyQualifiedName, filePath)
            }
        }
    }
    return m, nil
}

// GetSchema returns the raw JSON schema data for a given fully qualified name.
func (m *Manager) GetSchema(fullyQualifiedName string) (json.RawMessage, bool) {
    data, ok := m.LoadedSchemas[fullyQualifiedName]
    return data, ok
}