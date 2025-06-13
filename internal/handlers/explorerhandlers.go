package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/a-h/templ"

	"github.com/knanshon/cez/internal/methods"
	"github.com/knanshon/cez/internal/schemas"

	explorer "github.com/knanshon/cez/web/templates/explorer"
	explorer_components "github.com/knanshon/cez/web/templates/explorer/components"

	greeterv1 "github.com/knanshon/cez/gen/api/greeter/v1"
	greeterv1connect "github.com/knanshon/cez/gen/api/greeter/v1/greeterv1connect"
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

func (h *ExplorerHandlers) renderTemplComponent(w http.ResponseWriter, r *http.Request, component templ.Component, errMsg string) {
	if err := component.Render(r.Context(), w); err != nil {
		h.Logger.Printf("Error rendering %s: %v", errMsg, err)
		http.Error(w, fmt.Sprintf("Failed to render %s", errMsg), http.StatusInternalServerError)
	}
}

// Serves the base API explorer HTML page.
func (h *ExplorerHandlers) handleExplorerPage(w http.ResponseWriter, r *http.Request) {
	templ.Handler(explorer.Page()).ServeHTTP(w, r)
}

// Serves the list of API methods (as HTML) for HTMX.
func (h *ExplorerHandlers) handleApiMethodSelectorComponent(w http.ResponseWriter, r *http.Request) {
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

	var methodsForTemplate []map[string]string
	for _, m := range internalMethods {
		methodsForTemplate = append(methodsForTemplate, m.ToMap())
	}

	component := explorer_components.ApiMethodSelector(methodsForTemplate)
	h.renderTemplComponent(w, r, component, "API method selector")
}

// Serves the dynamic form based on the selected schema.
func (h *ExplorerHandlers) handleApiFormBuilderComponent(w http.ResponseWriter, r *http.Request) {
	schemaName := r.URL.Query().Get("method")
	serviceName := r.Header.Get("Hx-Data-Service")
	methodName := r.Header.Get("Hx-Data-Method")
	apiEndpoint := r.Header.Get("Hx-Data-Endpoint")

	if schemaName == "" {
		templ.Handler(explorer_components.ApiFormBuilder("No schema selected", nil, "", "", "")).ServeHTTP(w, r)
		return
	}

	schemaData, ok := h.SchemaManager.GetSchema(schemaName)
	if !ok {
		h.Logger.Printf("Error: Schema not found for %s", schemaName)
		http.Error(w, fmt.Sprintf("Schema not found: %s", schemaName), http.StatusBadRequest)
		return
	}

	component := explorer_components.ApiFormBuilder(schemaName, schemaData, serviceName, methodName, apiEndpoint)

	h.renderTemplComponent(w, r, component, "API form")
}

func (h *ExplorerHandlers) handleApiCallComponent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		h.Logger.Printf("Method not allowed for /explorer/components/api-call: %s", r.Method)
		return
	}

	// Get the target service and method from HTMX headers sent by the form
	serviceName := r.Header.Get("Hx-Data-Service")
	methodName := r.Header.Get("Hx-Data-Method")

	if serviceName == "" || methodName == "" {
		http.Error(w, "Missing Hx-Data-Service or Hx-Data-Method headers", http.StatusBadRequest)
		h.Logger.Println("Missing service or method headers in API call.")
		return
	}

	// Parse the form data submitted by HTMX
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		h.Logger.Printf("Failed to parse form data: %v", err)
		return
	}

	// --- Dynamic RPC Call Dispatch (Simplified for POC) ---
	// In a full solution, you'd have a dispatch table or use reflection
	// to map serviceName/methodName to the correct Protobuf message type
	// and Connect RPC client method. For this POC, we'll hardcode for GreeterService/Greet.

	var jsonResponse string
	var rpcErr error // To capture RPC error if any

	if serviceName == "greeter.v1.GreeterService" && methodName == "Greet" {
		// Extract data for GreetRequest
		name := r.FormValue("name") // 'name' is the field name from api_form_builder.templ

		// Construct the Protobuf request message
		greetRequest := &greeterv1.GreetRequest{
			Name: name,
		}

		// Initialize the Connect RPC client for GreeterService
		// Use http.DefaultClient for simplicity. The base URL is the current host.
		client := greeterv1connect.NewGreeterServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s", r.Host), // RPC endpoint on the same server
		)

		// Make the RPC call
		req := connect.NewRequest(greetRequest)
		res, err := client.Greet(r.Context(), req) // Pass the request context
		if err != nil {
			rpcErr = err // Store the RPC error
			if connectErr, ok := err.(*connect.Error); ok {
				jsonResponse = fmt.Sprintf(`{"error": "%s", "code": "%s", "details": "%s"}`,
					connectErr.Message(), connectErr.Code(), connectErr.Details())
			} else {
				jsonResponse = fmt.Sprintf(`{"error": "RPC call failed: %v"}`, err)
			}
		} else {
			// Marshal the successful RPC response message into pretty JSON
			jsonBytes, err := json.MarshalIndent(res.Msg, "", "  ")
			if err != nil {
				h.Logger.Printf("Failed to marshal Greet response: %v", err)
				http.Error(w, "Failed to marshal RPC response", http.StatusInternalServerError)
				return
			}
			jsonResponse = string(jsonBytes)
		}
	} else {
		jsonResponse = fmt.Sprintf(`{"error": "Unsupported API: %s/%s"}`, serviceName, methodName)
		http.Error(w, jsonResponse, http.StatusBadRequest)
		h.Logger.Printf("Unsupported API call: %s/%s", serviceName, methodName)
		return
	}

	// Render the ResponseViewer Templ component with the JSON response
	// The ResponseViewer will update the #api-response-viewer div in the frontend.
	component := explorer_components.ApiResponseViewer(jsonResponse) // Call from explorer_components package
	h.renderTemplComponent(w, r, component, "API response viewer")

	if rpcErr != nil {
		// If there was an RPC error, return an appropriate HTTP status code
		// This is important for HTMX to handle errors if needed (e.g., hx-target="this" hx-swap="outerHTML" hx-on::error)
		if connectErr, ok := rpcErr.(*connect.Error); ok {
			// Map common Connect errors to HTTP status codes
			switch connectErr.Code() {
			case connect.CodeInvalidArgument, connect.CodeNotFound, connect.CodeAlreadyExists:
				w.WriteHeader(http.StatusBadRequest)
			case connect.CodePermissionDenied, connect.CodeUnauthenticated:
				w.WriteHeader(http.StatusUnauthorized)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// Register API Explorer handlers
func (h *ExplorerHandlers) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/explorer", h.handleExplorerPage)
	mux.HandleFunc("/explorer/components/api-method-selector", h.handleApiMethodSelectorComponent)
	mux.HandleFunc("/explorer/components/api-form-builder", h.handleApiFormBuilderComponent)
	mux.HandleFunc("/explorer/components/api-call", h.handleApiCallComponent)
}
