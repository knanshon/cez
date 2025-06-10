package handlers

import (
	"net/http"
)

// Registerer is an interface that defines a method for registering HTTP handlers
type Registerer interface {
	RegisterHandlers(mux *http.ServeMux)
}