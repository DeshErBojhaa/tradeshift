package core

import "net/http"

// Package level constants
const (
	HeaderContentType         = "Content-Type"
	HeaderXContentTypeOptions = "X-Content-Type-Options"
	NoSniff                   = "nosniff"
	MediaTypeJSON             = "application/json"
	MethodGet                 = "GET"
	MethodPost                = "POST"
	MethodUpdate              = "PUT"
)

// ResponseWriter ...
type ResponseWriter func(w http.ResponseWriter)

// Handler ...
type Handler func(r Request) ResponseWriter

// Request ...
type Request interface {
	PathParam(key string) (string, bool)
	JSON(target interface{}) error
	Header(key string) string
}
