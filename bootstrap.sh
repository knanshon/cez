#!/bin/bash

# bootstrap.sh
# This script installs non-Go specific dependencies (protoc and Buf CLI) for macOS.

set -euo pipefail # Exit immediately if a command fails, or if an unset variable is used.

echo "Starting macOS bootstrap script..."

# --- Ensure Homebrew is installed ---
if ! command -v brew &> /dev/null; then
    echo "Homebrew not found. Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    echo "Homebrew installed. Please run 'brew doctor' if you encounter any issues."
else
    echo "Homebrew is already installed."
fi
#brew update # Ensure Homebrew is up-to-date

# --- Install Protobuf Compiler (protoc) ---
echo "--- Installing Protobuf Compiler (protoc) ---"
if brew list protobuf &> /dev/null; then
  echo "protoc is already installed via Homebrew. Skipping."
else
  echo "Installing protoc via Homebrew..."
  brew install protobuf
fi
protoc --version || { echo "Error: protoc installation failed or not found in PATH."; exit 1; }

# --- Install Buf CLI ---
echo "--- Installing Buf CLI ---"
if brew list buf &> /dev/null; then
  echo "Buf CLI is already installed via Homebrew. Skipping."
else
  echo "Installing Buf CLI via Homebrew..."
  brew install bufbuild/buf/buf
fi
buf --version || { echo "Error: Buf CLI installation failed or not found in PATH."; exit 1; }

echo "Bootstrap complete for non-Go dependencies."
echo "Ensuring your Go-based tools (templ, protoc-gen-jsonschema) are installed and your project code is generated"

go mod tidy || { echo "Error: 'go mod tidy' failed."; exit 1; }
echo "go.mod updated and Go module dependencies synced."

go generate ./... || { echo "Error: 'go generate ./...' failed."; exit 1; }
echo "Protobuf/Connect code and JSON schemas generated."

go build -o bin/server ./cmd/server || { echo "Error: Go server build failed."; exit 1; }
echo "Go server built to bin/server."

echo "Bootstrap script completed successfully."
echo "You can now run your Go server with './bin/server'."