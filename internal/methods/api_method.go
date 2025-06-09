package methods

import (
	"fmt"
)

type ApiMethod struct {
    Service        string 
    Method         string 
    RequestSchema  string
    ResponseSchema string 
}

func (m ApiMethod) ToMap() map[string]string {
    derivedEndpoint := fmt.Sprintf("/%s/%s", m.Service, m.Method)
    return map[string]string{
        "service":        m.Service,
        "method":         m.Method,
        "requestSchema":  m.RequestSchema,
        "responseSchema": m.ResponseSchema,
        "endpoint":       derivedEndpoint, // Add the derived endpoint
    }
}