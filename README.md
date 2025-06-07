# cez - Connect RPC Greeter Service

This project implements a simple Connect RPC service in Go. The service is built around a gRPC-like API defined using Protocol Buffers and exposed using [Connect](https://connectrpc.com/). It includes both a backend and a simple web front-end.

## Project Structure

- **api/greeter/v1/service.proto** – Protocol Buffers definition for the GreeterService. See [service.proto](api/greeter/v1/service.proto).
- **cmd/server/main.go** – Entry point for the RPC server. See [main.go](cmd/server/main.go).
- **internal/greeter/service.go** – Service implementation for the GreeterService. See [service.go](internal/greeter/service.go).
- **gen/** – Contains generated Go code for the protocol buffers and Connect RPC handler.
  - [service.pb.go](gen/api/greeter/v1/service.pb.go)
  - [service.connect.go](gen/api/greeter/v1/greeterv1connect/service.connect.go)
- **buf.gen.yaml & buf.yaml** – Configuration files for [buf](https://docs.buf.build/) to generate code.
- **go.mod** – Go module file.
- **web/static/** – Static assets for the web front-end (HTML, CSS). See [index.html](web/static/index.html) and [style.css](web/static/style.css).
- **LICENSE** – Project license.

## Prerequisites

- **Go 1.24.2** (or later)
- [Buf](https://docs.buf.build/installation) (for generating code)
- Dependencies will be managed via Go modules.

## Code Generation

This project uses buf to generate Go and Connect code from the `.proto` definitions.

To generate the code, run:

```sh
buf generate
```

This will generate files in the `gen/` directory according to the configuration in [buf.gen.yaml](buf.gen.yaml).

## Running the Server

You can launch the server with:

```sh
go run cmd/server/main.go
```

The server will start on port `8080` and host:
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