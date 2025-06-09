package main

import (
    "fmt"
	"log"
    "net/http"

	"github.com/a-h/templ"

    greeterv1connect "github.com/knanshon/cez/gen/api/greeter/v1/greeterv1connect"
    "github.com/knanshon/cez/internal/greeter"
    "github.com/knanshon/cez/internal/schemas"
    "github.com/knanshon/cez/internal/methods"
	"github.com/knanshon/cez/web/templates"
)

var schemaManager *schemas.Manager

func main() {    
	// Initialize the schema manager at startup
    var err error
    schemaManager, err = schemas.NewManager("web/schemas")
    if err != nil {
		log.Fatalf("Failed to initialize schema manager: %v", err)
    }

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

	mux.HandleFunc("/api-explorer/methods", func(w http.ResponseWriter, r *http.Request) {
		// This is a hardcoded list for the POC.
		// In a real application, you would dynamically discover these from your loaded Protobuf definitions
		// or from a configuration. For simplicity and to move quickly, hardcoding is fine for now.
		apiMethods := []methods.ApiMethod{
			{
				Service:        "greeter.v1.GreeterService",
				Method:         "Greet",
				RequestSchema:  "greeter.v1.GreetRequest",   // The fully qualified name of the request schema
				ResponseSchema: "greeter.v1.GreetResponse", // The fully qualified name of the response schema
			},
		}

		var methodsForTemplate []map[string]string
    	
		for _, m := range apiMethods {
			methodsForTemplate = append(methodsForTemplate, m.ToMap())
    	}

		// Render the APIMethodSelector Templ component with the prepared data
		component := templates.APIMethodSelector(methodsForTemplate)

		// Render the component and write the HTML to the response writer
		err := component.Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering APIMethodSelector: %v", err)
			http.Error(w, "Failed to render API method selector", http.StatusInternalServerError)
		}
	})

    fmt.Println("Server listening on :8080")
    err = http.ListenAndServe(":8080", mux)
    if err != nil {
        fmt.Printf("Server failed: %v\n", err)
    }
}