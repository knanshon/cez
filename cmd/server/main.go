package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	greeterv1connect "github.com/knanshon/cez/gen/api/greeter/v1/greeterv1connect"
	"github.com/knanshon/cez/internal/greeter"
	"github.com/knanshon/cez/internal/handlers"
	"github.com/knanshon/cez/internal/schemas"
)

func main() {
	// Setup logging
	logger := log.New(os.Stdout, "API_SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize schema manager
	schemaManager, err := schemas.NewManager("web/schemas")
	if err != nil {
		logger.Fatalf("Failed to initialize schema manager: %v", err)
	}

	// Instantiate your handlers struct
	explorerHandlers := handlers.NewExplorerHandlers(schemaManager, logger)

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

	// 3. Register static file servers
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("/schemas/", http.StripPrefix("/schemas/", http.FileServer(http.Dir("web/schemas"))))

	// 4. Register Explorer handlers
	explorerHandlers.RegisterHandlers(mux)

	// 5. Start the server
	logger.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
