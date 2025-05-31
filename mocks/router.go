// Package mocks provides reusable mock implementations for http interfaces.
package mocks

import (
	"net/http"
	"strings"
	"sync"

	forkCtx "go.fork.vn/fork/context"
	"go.fork.vn/fork/router"
)

// MockRouter is a reusable mock implementation of router.Router for testing.
// It implements all methods required by the router.Router interface and tracks
// information about method calls for testing assertions.
type MockRouter struct {
	Handlers        map[string][]router.HandlerFunc
	Middlewares     []router.HandlerFunc
	StaticDirs      map[string]string
	GroupRouters    map[string]*MockRouter
	ServeHTTPCalled bool
	Mutex           sync.RWMutex
}

// NewMockRouter creates a new MockRouter instance with default values.
func NewMockRouter() *MockRouter {
	return &MockRouter{
		Handlers:     make(map[string][]router.HandlerFunc),
		Middlewares:  make([]router.HandlerFunc, 0),
		StaticDirs:   make(map[string]string),
		GroupRouters: make(map[string]*MockRouter),
	}
}

// Handle registers a handler for a method and path
func (r *MockRouter) Handle(method string, path string, handlers ...router.HandlerFunc) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	key := method + ":" + path
	r.Handlers[key] = handlers
}

// Group creates a new router group with a prefix
func (r *MockRouter) Group(prefix string) router.Router {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	groupRouter := NewMockRouter()
	r.GroupRouters[prefix] = groupRouter
	return groupRouter
}

// Use adds middleware to the router
func (r *MockRouter) Use(middleware ...router.HandlerFunc) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Middlewares = append(r.Middlewares, middleware...)
}

// Static serves static files from a root directory
func (r *MockRouter) Static(prefix string, root string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.StaticDirs[prefix] = root
}

// Route represents a registered route
type MockRoute struct {
	Method      string
	Path        string
	HandlerFunc []router.HandlerFunc
}

// Routes returns all registered routes
func (r *MockRouter) Routes() []router.Route {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	routes := make([]router.Route, 0, len(r.Handlers))
	for key, handlers := range r.Handlers {
		parts := strings.Split(key, ":")
		if len(parts) != 2 {
			continue
		}

		routes = append(routes, router.Route{
			Method:  parts[0],
			Path:    parts[1],
			Handler: handlers[0], // Lấy handler đầu tiên từ danh sách
		})
	}

	// Add routes from groups
	for prefix, groupRouter := range r.GroupRouters {
		groupRoutes := groupRouter.Routes()
		for _, route := range groupRoutes {
			route.Path = prefix + route.Path
			routes = append(routes, route)
		}
	}

	return routes
}

// ServeHTTP implements http.Handler
func (r *MockRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Mutex.Lock()
	r.ServeHTTPCalled = true
	r.Mutex.Unlock()

	// Create context
	ctx := NewMockContext(w, req)

	// Find handler
	key := req.Method + ":" + req.URL.Path
	r.Mutex.RLock()
	handlers, exists := r.Handlers[key]
	r.Mutex.RUnlock()

	if exists && len(handlers) > 0 {
		// Convert router.HandlerFunc to []func(httpContext.Context)
		contextHandlers := make([]func(forkCtx.Context), len(handlers))
		for i, h := range handlers {
			handlerCopy := h // Capture the handler in a new variable to avoid closure issues
			contextHandlers[i] = func(c forkCtx.Context) {
				handlerCopy(c)
			}
		}

		// Set handlers on context
		ctx.SetHandlers(contextHandlers)
		// Start execution
		ctx.Next()
	} else {
		// No handler found, return 404
		ctx.Status(http.StatusNotFound)
	}
}

// Find tìm route theo method và path
func (r *MockRouter) Find(method, path string) router.HandlerFunc {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	key := method + ":" + path
	handlers, exists := r.Handlers[key]

	if exists && len(handlers) > 0 {
		return handlers[0]
	}

	// Kiểm tra trong các group router
	for prefix, groupRouter := range r.GroupRouters {
		if strings.HasPrefix(path, prefix) {
			// Kiểm tra nếu path có matching với prefix của group
			relativePath := path[len(prefix):]
			if handler := groupRouter.Find(method, relativePath); handler != nil {
				return handler
			}
		}
	}

	return nil
}
