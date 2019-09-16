// Package webber provides the router and attached handlers to request
// Not actually necessary for this simple usecase. But added anyway
// keeping scalability in mind.
package webber

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// BasicRequest wraps the http request with any route param associated with the request
type BasicRequest struct {
	httpRequest *http.Request
	pathParams  map[string]string
}

// NewRequest returns a 'BasicRequest' where path parameters are stored separately
func NewRequest(r *http.Request) *BasicRequest {
	return &BasicRequest{
		httpRequest: r,
		pathParams:  mux.Vars(r),
	}
}

// Header gets the header from the request
func (r *BasicRequest) Header(key string) string {
	return r.httpRequest.Header.Get(key)
}

// PathParam gets the optional path parameters
func (r *BasicRequest) PathParam(key string) (string, bool) {
	v, ok := r.pathParams[key]
	return v, ok
}

// JSON marshals the body into json format
func (r *BasicRequest) JSON(target interface{}) error {
	return json.NewDecoder(r.httpRequest.Body).Decode(target)
}
