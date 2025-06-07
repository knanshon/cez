# cez - Connect RPC Greeter Service

This project implements a simple Connect RPC GreeterService using Go. The service is hosted under the `cmd/server` directory.

## Project Structure

- **cmd/server/main.go** – Entry point for the RPC server.
- **internal/api/greeter/** – Contains the service implementation and generated handler code from the proto definition.
- **go.mod** – Module file declaring the module as `cez`.

## Prerequisites

- Go 1.24.2 or later installed.
- [Optional] bufbuild/connect-go package if not vendored.

## Launch Instructions

To run the server, execute the following commands from the project root:

```sh

```

The server will start on port `8080` and will be accessible via the Connect RPC endpoint for the GreeterService.

## Usage

Clients can connect to the service at `http://localhost:8080` using the Connect protocol.

Enjoy experimenting with your simple Connect RPC service!