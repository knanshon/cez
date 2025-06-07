package main

import (
    "fmt"
    "net/http"

    // Import the generated Connect handler creator
    greeterv1connect "github.com/knanshon/cez/gen/api/greeter/v1/greeterv1connect"
    // Import the internal greeter service implementation
    "github.com/knanshon/cez/internal/greeter"
)

func main() {
    mux := http.NewServeMux()

    // 1. Standard Mux Test Endpoint
    mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from the Go server!")
    })

    // 2. Connect RPC Handler for GreeterService
    greeterSvc := greeter.NewService() // Instantiate your service
    // NewGreeterServiceHandler takes the service implementation
    path, handler := greeterv1connect.NewGreeterServiceHandler(greeterSvc)
    mux.Handle(path, handler)

    fmt.Println("Server listening on :8080")
    err := http.ListenAndServe(":8080", mux)
    if err != nil {
        fmt.Printf("Server failed: %v\n", err)
    }
}