// Package mocks provides reusable mock implementations for http interfaces.
package mocks

import (
	"net/http"
	"sync"

	httpCtx "go.fork.vn/fork/context"
)

// MockAdapter is a reusable mock implementation of adapter.Adapter for testing.
// It implements all methods required by the adapter.Adapter interface and tracks
// information about method calls for testing assertions.
type MockAdapter struct {
	NameVal         string
	HandlerFuncs    map[string]func(ctx httpCtx.Context)
	Middlewares     []func(ctx httpCtx.Context)
	Handler         http.Handler
	RunCalled       bool
	RunTLSCalled    bool
	ShutdownCalled  bool
	ServeHTTPCalled bool
	Mutex           sync.RWMutex
}

// NewMockAdapter creates a new MockAdapter instance with default values.
func NewMockAdapter(name string) *MockAdapter {
	return &MockAdapter{
		NameVal:      name,
		HandlerFuncs: make(map[string]func(ctx httpCtx.Context)),
		Middlewares:  make([]func(ctx httpCtx.Context), 0),
	}
}

// Name returns the name of the adapter
func (a *MockAdapter) Name() string {
	return a.NameVal
}

// Run mocks starting the HTTP server
func (a *MockAdapter) Run(addr string) error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	a.RunCalled = true
	return nil
}

// RunTLS mocks starting the HTTPS server
func (a *MockAdapter) RunTLS(addr, certFile, keyFile string, alpnProtocols ...string) error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	a.RunTLSCalled = true
	return nil
}

// ServeHTTP mocks handling HTTP requests
func (a *MockAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	a.ServeHTTPCalled = true

	if a.Handler != nil {
		a.Handler.ServeHTTP(w, r)
	}
}

// HandleFunc mocks registering a handler function
func (a *MockAdapter) HandleFunc(method, path string, handler func(ctx httpCtx.Context)) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	key := method + ":" + path
	a.HandlerFuncs[key] = handler
}

// Use mocks adding middleware
func (a *MockAdapter) Use(middleware func(ctx httpCtx.Context)) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	a.Middlewares = append(a.Middlewares, middleware)
}

// SetHandler mocks setting the main handler
func (a *MockAdapter) SetHandler(handler http.Handler) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	a.Handler = handler
}

// Shutdown mocks gracefully shutting down the HTTP server
func (a *MockAdapter) Shutdown() error {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()
	a.ShutdownCalled = true
	return nil
}
