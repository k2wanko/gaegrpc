package gaegrpc

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// HeaderKey is http.Request ID
const HeaderKey = "x-gae-grpc-id"

type requestKey struct{}

var (
	reqs = make(map[string]*http.Request)
	mu   sync.RWMutex
)

func newContextWithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestKey{}, r)
}

// RequestFromContext returns *http.Request
func RequestFromContext(ctx context.Context) *http.Request {
	if r, ok := ctx.Value(requestKey{}).(*http.Request); ok {
		return r
	}
	return nil
}

func requestIDFromContext(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		v := md[HeaderKey]
		if len(v) > 0 {
			return v[0]
		}

	}
	return ""
}

func newAppContext(ctx context.Context) context.Context {
	id := requestIDFromContext(ctx)
	if id != "" {
		mu.RLock()
		r := reqs[id]
		mu.RUnlock()
		ctx = newContextWithRequest(ctx, r)
		ctx = appengine.WithContext(ctx, r)
	}
	return ctx
}

func injectAppContext() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			resp, err = handler(newAppContext(ctx), req)
			return
		}),
		grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
			err = handler(srv, &wrapServerStream{
				ctx:          newAppContext(ss.Context()),
				ServerStream: ss,
			})
			return
		}),
	}
}

type wrapServerStream struct {
	ctx context.Context
	grpc.ServerStream
}

func (wss *wrapServerStream) Context() context.Context {
	return wss.ctx
}

// requestID returns ID, ID is pointer address.
func requestID(r *http.Request) string {
	return fmt.Sprintf("%x", &r)
}

// NewServer returns grpc.Server for App Engine
func NewServer(opt ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(append(injectAppContext(), opt...)...)
}

// NewRequest returns http.Request for GRPC, set the http.Request on memory
func NewRequest(r *http.Request) *http.Request {
	id := requestID(r)
	mu.Lock()
	reqs[id] = r
	mu.Unlock()
	r.Header.Add(HeaderKey, id)
	return r
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

// DeleteRequest deletes the http.Request on memory
func DeleteRequest(r *http.Request) {
	mu.Lock()
	delete(reqs, requestID(r))
	mu.Unlock()
}

type wrapHandler struct {
	h http.Handler
}

func (s *wrapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.h.ServeHTTP(NewWrapResponseWriter(w), NewRequest(r))
	DeleteRequest(r)
}

// NewWrapHandler returns http.Handler for App Engine
func NewWrapHandler(h http.Handler) http.Handler {
	return &wrapHandler{
		h: h,
	}
}
