//go:build tools

package tools

// This file exists to ensure that 'go mod tidy'
// includes specific Go-based tools in the go.mod file,
// even if they are not directly imported by the application code.
// These tools are typically run via 'go generate' or scripts.

import (
    _ "github.com/a-h/templ/cmd/templ" // Used by templ compiler in go generate
    _ "github.com/pubg/protoc-gen-jsonschema" // The JSON schema generator
)