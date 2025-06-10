package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"

	"github.com/knanshon/cez/internal/methods"
	"github.com/knanshon/cez/internal/schemas"

	explorer "github.com/knanshon/cez/web/templates/explorer"
	explorer_components "github.com/knanshon/cez/web/templates/explorer/components"
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
	serviceName := r.Header.Get("HX-Data-Service")
	methodName := r.Header.Get("HX-Data-Method")
	apiEndpoint := r.Header.Get("HX-Data-Endpoint")

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

// Register API Explorer handlers
func (h *ExplorerHandlers) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/explorer", h.handleExplorerPage)
	mux.HandleFunc("/explorer/components/api-method-selector", h.handleApiMethodSelectorComponent)
	mux.HandleFunc("/explorer/components/api-form-builder", h.handleApiFormBuilderComponent)
}
