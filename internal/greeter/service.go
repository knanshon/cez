package greeter

import (
    "context"
    "fmt"

    // Import the generated protobuf messages
    greeterv1 "github.com/knanshon/cez/gen/api/greeter/v1"
    // Import the connectrpc package for connect.Request/Response
    "connectrpc.com/connect"
)

// Service implements the GreeterServiceHandler interface.
type Service struct{}

// NewService creates a new Greeter service instance.
func NewService() *Service {
    return &Service{}
}

// Greet implements the Greet RPC method defined in service.proto.
func (s *Service) Greet(ctx context.Context, req *connect.Request[greeterv1.GreetRequest]) (*connect.Response[greeterv1.GreetResponse], error) {
    fmt.Printf("Greeter Service received request: %s\n", req.Msg.Name)
    greeting := fmt.Sprintf("Hello, %s!", req.Msg.Name)
    return connect.NewResponse(&greeterv1.GreetResponse{Greeting: greeting}), nil
}