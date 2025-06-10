package main

import (
    "fmt"
	"log"
	"os"
    "net/http"

	"github.com/a-h/templ"

    greeterv1connect "github.com/knanshon/cez/gen/api/greeter/v1/greeterv1connect"
    "github.com/knanshon/cez/internal/greeter"
    "github.com/knanshon/cez/internal/schemas"
    "github.com/knanshon/cez/internal/methods"
	"github.com/knanshon/cez/web/templates"
)

// ExplorerHandlers holds dependencies required by the explorer's HTTP handlers.
type ExplorerHandlers struct {
	SchemaManager *schemas.Manager
	Logger        *log.Logger
}

// NewExplorerHandlers creates a new instance of ExplorerHandlers.
func NewExplorerHandlers(sm *schemas.Manager, logger *log.Logger) *ExplorerHandlers {
	return &ExplorerHandlers{
		SchemaManager: sm,
		Logger:        logger,
	}
}

// HandleExplorerPage serves the base API explorer HTML page.
func (h *ExplorerHandlers) HandleExplorerPage(w http.ResponseWriter, r *http.Request) {
	templ.Handler(templates.ExplorerPage()).ServeHTTP(w, r)
}

// HandleAPIMethods serves the list of API methods (as HTML) for HTMX.
func (h *ExplorerHandlers) HandleAPIMethods(w http.ResponseWriter, r *http.Request) {
	// Define your internal representation of API methods using the ApiMethod struct
	internalMethods := []methods.ApiMethod{
		{
			Service:        "greeter.v1.GreeterService",
			Method:         "Greet",
			RequestSchema:  "greeter.v1.GreetRequest",
			ResponseSchema: "greeter.v1.GreetResponse",
		},
		// Add other API methods here if you expand your API
	}

	// Convert []ApiMethod to []map[string]string using the ApiMethod's ToMap() method
	// This prepares the data for the Templ component.
	var methodsForTemplate []map[string]string
	for _, m := range internalMethods {
		methodsForTemplate = append(methodsForTemplate, m.ToMap())
	}

	// Render the APIMethodSelector Templ component with the prepared data
	component := templates.APIMethodSelector(methodsForTemplate)

	if err := component.Render(r.Context(), w); err != nil {
		h.Logger.Printf("Error rendering APIMethodSelector: %v", err)
		http.Error(w, "Failed to render API method selector", http.StatusInternalServerError)
	}
}

// HandleFormLoader serves the dynamic form based on the selected schema.
func (h *ExplorerHandlers) HandleFormLoader(w http.ResponseWriter, r *http.Request) {
	schemaName := r.URL.Query().Get("method")
	serviceName := r.Header.Get("HX-Data-Service")
	methodName := r.Header.Get("HX-Data-Method")
	apiEndpoint := r.Header.Get("HX-Data-Endpoint")

	if schemaName == "" {
		templ.Handler(templates.FormBuilder("No schema selected", nil, "", "", "")).ServeHTTP(w, r)
		return
	}

	schemaData, ok := h.SchemaManager.GetSchema(schemaName)
	if !ok {
		h.Logger.Printf("Error: Schema not found for %s", schemaName)
		http.Error(w, fmt.Sprintf("Schema not found: %s", schemaName), http.StatusBadRequest)
		return
	}

	component := templates.FormBuilder(schemaName, schemaData, serviceName, methodName, apiEndpoint)

	if err := component.Render(r.Context(), w); err != nil {
		h.Logger.Printf("Error rendering FormBuilder: %v", err)
		http.Error(w, "Failed to render API form", http.StatusInternalServerError)
	}
}

func main() {    
	// Setup logging
	logger := log.New(os.Stdout, "API_SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize schema manager
	schemaManager, err := schemas.NewManager("web/schemas")
	if err != nil {
		logger.Fatalf("Failed to initialize schema manager: %v", err)
	}

	// Instantiate your handlers struct
	explorerHandlers := NewExplorerHandlers(schemaManager, logger)

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

	// 4. Register API Explorer handlers using methods on the explorerHandlers struct
	mux.HandleFunc("/explorer", explorerHandlers.HandleExplorerPage)
	mux.HandleFunc("/api-explorer/methods", explorerHandlers.HandleAPIMethods) // Returns HTML now
	mux.HandleFunc("/api-explorer/form", explorerHandlers.HandleFormLoader)

	logger.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}