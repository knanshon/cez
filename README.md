# cez - Connect RPC API Explorer

This project implements a simple Connect RPC service in Go and explores how to generate an API explorer UI for it. The service is built around a gRPC-like API defined using Protocol Buffers and exposed using [Connect](https://connectrpc.com/). It includes both a backend and a simple web front-end.

## Project Structure

- **api/greeter/v1/service.proto** – Protocol Buffers definition for the GreeterService. See [service.proto](api/greeter/v1/service.proto).
- **cmd/server/main.go** – Entry point for the RPC server. See [main.go](cmd/server/main.go).
- **internal/greeter/service.go** – Service implementation for the GreeterService. See [service.go](internal/greeter/service.go).
- **gen/** – Contains generated Go code for the Protocol Buffers files and Connect RPC handler.
  - [service.pb.go](gen/api/greeter/v1/service.pb.go)
  - [service.connect.go](gen/api/greeter/v1/greeterv1connect/service.connect.go)
- **generate.go** – A helper file that orchestrates code generation via `go generate ./...` (invoking both `buf generate` and the JSON schema generation script).
- **tools.go** – Declares build dependencies so that Go modules include necessary code-generation tools.
- **buf.gen.yaml & buf.yaml** – Configuration files for [buf](https://docs.buf.build/) to generate code.
- **go.mod** – Go module file.
- **bootstrap.sh** – A script to set up non-Go dependencies (like protoc and Buf CLI) and to generate all code, building the server.
- **web/static/** – Static assets for the web front-end (HTML, CSS). See [index.html](web/static/index.html) and [style.css](web/static/style.css).
- **LICENSE** – Project license.

## Prerequisites

- **Go 1.24.2** (or later)
- [Buf](https://docs.buf.build/installation) (for generating code)
- Dependencies will be managed via Go modules.

## Bootstrap & Code Generation

A new bootstrap script sets up the project by installing non-Go dependencies and ensuring that Go-based tools are installed. It also runs the code generation and builds the server.

To run the bootstrap process, use:

```sh
./bootstrap.sh
```

This script does the following:
- Installs non-Go dependencies like `protoc` and the Buf CLI (using Homebrew on macOS).
- Runs `go mod tidy` to sync build dependencies (including those declared in `tools.go`).
- Executes `go generate ./...` (as configured in [generate.go](generate.go)) which:
  - Runs `buf generate` to produce the Go and Connect RPC client code in the `gen/` directory.
  - Executes the script at `./scripts/generate_schemas.sh` to clear and generate fresh JSON schema files from `api/greeter/v1/service.proto` into the `web/schemas` directory.

**Note:** If you encounter a "Permission denied" error with the JSON schemas script, make it executable by running:

```sh
chmod +x scripts/generate_schemas.sh
```

## Running the Server

After bootstrapping, you can launch the server with:

```sh
./bin/server
```

The server will start on port `8080` and provide:
- A standard HTTP endpoint at `/api/hello`
- The Connect RPC endpoint for the GreeterService at `/greeter.v1.GreeterService/`

## Usage

### Connect RPC Client

Clients can connect using the Connect protocol at `http://localhost:8080`. The generated Connect client code is available in [gen/api/greeter/v1/greeterv1connect/service.connect.go](gen/api/greeter/v1/greeterv1connect/service.connect.go).

### Web Front-End

A basic web front-end is available under the `web` directory. It features a simple HTML page ([index.html](web/static/index.html)) styled with [style.css](web/static/style.css) and uses [HTMX](https://htmx.org/) for dynamic interactions.

## Development

- Service implementation is found in [internal/greeter/service.go](internal/greeter/service.go).
- The Protocol Buffers definition is in [api/greeter/v1/service.proto](api/greeter/v1/service.proto).
- Adjust or add endpoints via [cmd/server/main.go](cmd/server/main.go).

## License

This project is licensed under the GNU General Public License v3. See [LICENSE](LICENSE) for details.