#!/bin/bash
# This script converts your .proto file into JSON schemas.

set -euo pipefail # Exit immediately if a command fails, or if an unset variable is used.

# Define the specific output directory where you want the schemas.
# The protoc-gen-jsonschema tool does NOT automatically create subdirectories based on Protobuf package.
SCHEMA_OUTPUT_DIR="web/schemas/greeter/v1"

echo "Generating JSON schemas to: $SCHEMA_OUTPUT_DIR"

# Create the directory if it doesn't exist
mkdir -p "$SCHEMA_OUTPUT_DIR" || { echo "Error: Failed to create schema directory $SCHEMA_OUTPUT_DIR"; exit 1; }

# Run protoc with the jsonschema generator
# --proto_path specifies the base directory where your .proto files are found.
# --jsonschema_out points directly to the desired output directory for the schemas.
protoc \
  --proto_path=api \
  --jsonschema_out="$SCHEMA_OUTPUT_DIR" \
  api/greeter/v1/service.proto

echo "JSON schema generation complete."