package gaegrpc

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/grpc"
)

// NewServer returns grpc.Server for App Engine
func NewServer(opt ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(opt...)
}

// NewRequest returns http.Request for GRPC, set the http.Request on memory
func NewRequest(r *http.Request) *http.Request {
	return r.WithContext(appengine.WithContext(r.Context(), r))
}

type wrapResponseWriter struct {
	w http.ResponseWriter
}

// NewWrapResponseWriter returns wraped http.ResponseWriter
func NewWrapResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	return &wrapResponseWriter{
		w: w,
	}
}

func (w *wrapResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *wrapResponseWriter) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w *wrapResponseWriter) WriteHeader(code int) {
	w.w.WriteHeader(code)
}

func (w *wrapResponseWriter) CloseNotify() <-chan bool {
	if w, ok := w.w.(http.CloseNotifier); ok {
		return w.CloseNotify()
	}
	return nil
}

func (w *wrapResponseWriter) Flush() {
	if w, ok := w.w.(http.Flusher); ok {
		w.Flush()
	}
	return
}

type wrapHandler struct {
	h http.Handler
}

func (s *wrapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.h.ServeHTTP(NewWrapResponseWriter(w), NewRequest(r))
}

// NewWrapHandler returns http.Handler for App Engine
func NewWrapHandler(h http.Handler) http.Handler {
	return &wrapHandler{
		h: h,
	}
}
