package webber

import (
	"net/http"
	"time"

	"github.com/DeshErBojhaa/tradeshift/webber/core"
	"github.com/gorilla/mux"
)

// Server wraps router from gorilla and native http server from go.
type Server struct {
	router     *mux.Router
	httpServer *http.Server
	mediaType  string
}

// NewServer returns new instance of gorilla mux.
func NewServer(listenAddress, mediaType string) *Server {
	r := mux.NewRouter()
	r.NotFoundHandler = notFoundHandler(mediaType)
	r.MethodNotAllowedHandler = methodNotAllowedHandler(mediaType)

	return &Server{
		router: r,
		httpServer: &http.Server{
			Addr:         listenAddress,
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
		},
		mediaType: mediaType,
	}
}

// GET attaches router to corresponding handler.
func (s *Server) GET(path string, h core.Handler) {
	s.register(path, h, core.MethodGet)
}

// POST attaches router to corresponding handler.
func (s *Server) POST(path string, h core.Handler) {
	s.register(path, h, core.MethodPost)
}

// PUT changes the parent of a given node
func (s *Server) PUT(path string, h core.Handler) {
	s.register(path, h, core.MethodUpdate)
}

// Serve starts the service
func (s *Server) Serve() error {
	s.httpServer.Handler = s.router
	return s.httpServer.ListenAndServe()
}

func (s *Server) register(path string, h core.Handler, method string) {
	s.router.HandleFunc(path, wrap(h)).Methods(method)
}

func wrap(handler core.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(r)
		resp := handler(req)
		resp(w)
	}
}

func notFoundHandler(mediaType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(core.HeaderContentType, mediaType)
		w.Header().Set(core.HeaderXContentTypeOptions, core.NoSniff)
		w.WriteHeader(http.StatusNotFound)
	})
}

func methodNotAllowedHandler(mediaType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(core.HeaderContentType, mediaType)
		w.Header().Set(core.HeaderXContentTypeOptions, core.NoSniff)
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}
