package main

import (
    "fmt"
    "net/http"

	"github.com/a-h/templ"

    // Import the generated Connect handler creator
    greeterv1connect "github.com/knanshon/cez/gen/api/greeter/v1/greeterv1connect"
    // Import the internal greeter service implementation
    "github.com/knanshon/cez/internal/greeter"

	"github.com/knanshon/cez/web/templates"
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
	
	// Serve the API Explorer page
	mux.Handle("/explorer", templ.Handler(templates.ExplorerPage()))
	// Serve static files from web/static (e.g., CSS, htmx.min.js)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	// Serve generated JSON schemas (e.g., GreetRequest.json)
	mux.Handle("/schemas/", http.StripPrefix("/schemas/", http.FileServer(http.Dir("web/schemas"))))

    fmt.Println("Server listening on :8080")
    err := http.ListenAndServe(":8080", mux)
    if err != nil {
        fmt.Printf("Server failed: %v\n", err)
    }
}