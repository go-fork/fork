package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.fork.vn/fork/context"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	if router == nil {
		t.Fatal("Expected router to be created, got nil")
	}

	// Kiểm tra kết quả
	r, ok := router.(*DefaultRouter)
	if !ok {
		t.Fatal("Expected DefaultRouter type")
	}

	if r.basePath != "" {
		t.Errorf("Expected empty basePath, got '%s'", r.basePath)
	}

	if len(r.routes) != 0 {
		t.Errorf("Expected empty routes, got %d routes", len(r.routes))
	}

	if len(r.middlewares) != 0 {
		t.Errorf("Expected empty middlewares, got %d middlewares", len(r.middlewares))
	}

	if len(r.groups) != 0 {
		t.Errorf("Expected empty groups, got %d groups", len(r.groups))
	}
}

func TestDefaultRouter_Handle(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	// Đăng ký một handler
	handlerCalled := false
	handler := func(ctx context.Context) {
		handlerCalled = true
		ctx.String(http.StatusOK, "OK")
	}

	router.Handle("GET", "/test", handler)

	// Kiểm tra route đã được đăng ký
	if len(router.routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(router.routes))
	}

	route := router.routes[0]
	if route.Method != "GET" {
		t.Errorf("Expected method GET, got %s", route.Method)
	}

	if route.Path != "/test" {
		t.Errorf("Expected path /test, got %s", route.Path)
	}

	// Kiểm tra handler đã được thiết lập
	if route.Handler == nil {
		t.Fatal("Expected handler to be set")
	}

	// Kiểm tra handler hoạt động
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	ctx := context.NewContext(w, r)

	route.Handler(ctx)

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", w.Body.String())
	}
}

func TestDefaultRouter_Group(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	// Tạo một group
	group := router.Group("/api")

	// Kiểm tra group
	g, ok := group.(*DefaultRouter)
	if !ok {
		t.Fatal("Expected DefaultRouter type for group")
	}

	if g.basePath != "/api" {
		t.Errorf("Expected basePath /api, got '%s'", g.basePath)
	}

	// Kiểm tra group được thêm vào router
	if len(router.groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(router.groups))
	}

	// Kiểm tra group hoạt động
	handlerCalled := false
	group.Handle("GET", "/users", func(ctx context.Context) {
		handlerCalled = true
		ctx.String(http.StatusOK, "Users")
	})

	// Kiểm tra route được thêm vào group
	if len(g.routes) != 1 {
		t.Fatalf("Expected 1 route in group, got %d", len(g.routes))
	}

	groupRoute := g.routes[0]
	if groupRoute.Method != "GET" {
		t.Errorf("Expected method GET, got %s", groupRoute.Method)
	}

	if groupRoute.Path != "/api/users" {
		t.Errorf("Expected path /api/users, got %s", groupRoute.Path)
	}

	// Kiểm tra handler đã được thiết lập
	if groupRoute.Handler == nil {
		t.Fatal("Expected handler to be set")
	}

	// Kiểm tra handler hoạt động
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/users", nil)
	ctx := context.NewContext(w, r)

	groupRoute.Handler(ctx)

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}

	// Kiểm tra nested group
	subGroup := group.Group("/v1")
	sg, ok := subGroup.(*DefaultRouter)
	if !ok {
		t.Fatal("Expected DefaultRouter type for subgroup")
	}

	if sg.basePath != "/api/v1" {
		t.Errorf("Expected basePath /api/v1, got '%s'", sg.basePath)
	}

	// Kiểm tra subgroup được thêm vào group
	if len(g.groups) != 1 {
		t.Errorf("Expected 1 subgroup, got %d", len(g.groups))
	}
}

func TestDefaultRouter_Use(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	// Thêm middleware
	middlewareCalled := false
	router.Use(func(ctx context.Context) {
		middlewareCalled = true
		ctx.Next()
	})

	// Kiểm tra middleware được thêm vào
	if len(router.middlewares) != 1 {
		t.Fatalf("Expected 1 middleware, got %d", len(router.middlewares))
	}

	// Đăng ký handler
	handlerCalled := false
	router.Handle("GET", "/test", func(ctx context.Context) {
		handlerCalled = true
		ctx.String(http.StatusOK, "OK")
	})

	// Gọi handler qua ServeHTTP để xác nhận middleware
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, r)

	// Kiểm tra middleware và handler đã được gọi
	if !middlewareCalled {
		t.Error("Expected middleware to be called")
	}

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
}

func TestDefaultRouter_ServeHTTP(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	// Đăng ký một handler
	handlerCalled := false
	router.Handle("GET", "/test", func(ctx context.Context) {
		handlerCalled = true
		ctx.String(http.StatusOK, "OK")
	})

	// Gọi ServeHTTP
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, r)

	// Kiểm tra handler đã được gọi
	if !handlerCalled {
		t.Error("Expected handler to be called")
	}

	// Kiểm tra response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", w.Body.String())
	}

	// Kiểm tra với path không tồn tại
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/not-found", nil)

	router.ServeHTTP(w, r)

	// Kiểm tra response
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDefaultRouter_Static(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	// Đăng ký static route
	router.Static("/static", "./testdata")

	// Kiểm tra route đã được đăng ký
	if len(router.routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(router.routes))
	}

	route := router.routes[0]
	if route.Method != "GET" {
		t.Errorf("Expected method GET, got %s", route.Method)
	}

	if route.Path != "/static/*filepath" {
		t.Errorf("Expected path /static/*filepath, got %s", route.Path)
	}
}

func TestDefaultRouter_Routes(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	// Đăng ký một số route
	router.Handle("GET", "/users", func(ctx context.Context) {})
	router.Handle("POST", "/users", func(ctx context.Context) {})

	// Tạo group và đăng ký route
	group := router.Group("/api")
	group.Handle("GET", "/products", func(ctx context.Context) {})

	// Lấy tất cả route
	routes := router.Routes()

	// Kiểm tra số lượng route
	if len(routes) != 3 {
		t.Fatalf("Expected 3 routes, got %d", len(routes))
	}

	// Kiểm tra từng route
	expectedRoutes := map[string]string{
		"GET /users":        "",
		"POST /users":       "",
		"GET /api/products": "",
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, ok := expectedRoutes[key]; !ok {
			t.Errorf("Unexpected route: %s", key)
		} else {
			delete(expectedRoutes, key)
		}
	}

	if len(expectedRoutes) > 0 {
		t.Errorf("Missing routes: %v", expectedRoutes)
	}
}

func TestDefaultRouter_pathMatch(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	testCases := []struct {
		pattern  string
		path     string
		expected bool
	}{
		{"/users", "/users", true},
		{"/users", "/user", false},
		{"/static/*filepath", "/static/css/style.css", true},
		{"/static/*filepath", "/public/css/style.css", false},
		{"/api/v1", "/api/v1", true},
		{"/api/v1", "/api/v2", false},
	}

	for _, tc := range testCases {
		result := router.pathMatch(tc.pattern, tc.path)
		if result != tc.expected {
			t.Errorf("pathMatch(%q, %q) = %v, expected %v", tc.pattern, tc.path, result, tc.expected)
		}
	}
}

func TestDefaultRouter_calculateAbsolutePath(t *testing.T) {
	testCases := []struct {
		basePath     string
		relativePath string
		expected     string
	}{
		{"", "/users", "/users"},
		{"/api", "/users", "/api/users"},
		{"/api/", "users", "/api/users"},
		{"/api", "", "/api"},
		{"", "", ""},
	}

	for _, tc := range testCases {
		router := &DefaultRouter{basePath: tc.basePath}
		result := router.calculateAbsolutePath(tc.relativePath)
		if result != tc.expected {
			t.Errorf("calculateAbsolutePath(%q) with basePath %q = %q, expected %q",
				tc.relativePath, tc.basePath, result, tc.expected)
		}
	}
}

func TestDefaultRouter_combineHandlers(t *testing.T) {
	router := &DefaultRouter{
		middlewares: []HandlerFunc{
			func(ctx context.Context) { /* middleware 1 */ },
			func(ctx context.Context) { /* middleware 2 */ },
		},
	}

	handlers := []HandlerFunc{
		func(ctx context.Context) { /* handler 1 */ },
		func(ctx context.Context) { /* handler 2 */ },
	}

	combined := router.combineHandlers(handlers)

	if len(combined) != 4 {
		t.Fatalf("Expected 4 handlers, got %d", len(combined))
	}
}

// TestEnhancedPathMatching tests the advanced path matching functionality
func TestEnhancedPathMatching(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	testCases := []struct {
		pattern  string
		path     string
		expected bool
	}{
		// Static routes
		{"/users", "/users", true},
		{"/products", "/products", true},
		{"/products", "/produc", false},

		// Named parameters
		{"/users/:id", "/users/123", true},
		{"/users/:id", "/users/abc", true},
		{"/users/:id", "/users/", false},
		{"/users/:id", "/users", false},

		// Regex constraints
		{"/users/:id<\\d+>", "/users/123", true},
		{"/users/:id<\\d+>", "/users/abc", false},
		{"/users/:id<[a-z]+>", "/users/abc", true},
		{"/users/:id<[a-z]+>", "/users/123", false},

		// Optional parameters
		{"/users/:id?", "/users/123", true},
		{"/users/:id?", "/users", true},
		{"/optional/:param?/test", "/optional/test", true},
		{"/optional/:param?/test", "/optional/value/test", true},

		// Wildcard parameters
		{"/files/*filepath", "/files/images/logo.png", true},
		{"/files/*filepath", "/files", true},
		{"/static/*path", "/static/css/style.css", true},
		{"/static/*path", "/public/css/style.css", false},

		// Complex combinations
		{"/api/:version?/users", "/api/users", true},
		{"/api/:version?/users", "/api/v1/users", true},
		{"/posts/:year<\\d{4}>/:month<\\d{2}>/:day<\\d{2}>", "/posts/2023/05/15", true},
		{"/posts/:year<\\d{4}>/:month<\\d{2}>/:day<\\d{2}>", "/posts/2023/5/15", false},
		{"/users/:id<\\d+>/profile/:section?", "/users/123/profile", true},
		{"/users/:id<\\d+>/profile/:section?", "/users/123/profile/about", true},
		{"/users/:id<\\d+>/profile/:section?", "/users/abc/profile", false},
	}

	for _, tc := range testCases {
		result := router.pathMatch(tc.pattern, tc.path)
		if result != tc.expected {
			t.Errorf("pathMatch(%q, %q) = %v, expected %v", tc.pattern, tc.path, result, tc.expected)
		}
	}
}

// TestExtractParams tests the parameter extraction functionality
func TestExtractParams(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	testCases := []struct {
		pattern string
		path    string
		params  map[string]string
	}{
		{
			"/users/:id",
			"/users/123",
			map[string]string{"id": "123"},
		},
		{
			"/users/:id<\\d+>",
			"/users/456",
			map[string]string{"id": "456"},
		},
		{
			"/posts/:year/:month/:day",
			"/posts/2023/05/15",
			map[string]string{"year": "2023", "month": "05", "day": "15"},
		},
		{
			"/api/:version?/users",
			"/api/users",
			map[string]string{"version": ""},
		},
		{
			"/api/:version?/users",
			"/api/v1/users",
			map[string]string{"version": "v1"},
		},
		{
			"/files/*filepath",
			"/files/images/logo.png",
			map[string]string{"filepath": "images/logo.png"},
		},
		{
			"/files/*filepath",
			"/files",
			map[string]string{"filepath": ""},
		},
		{
			"/users/:id/posts/:slug/:status?",
			"/users/123/posts/hello-world",
			map[string]string{"id": "123", "slug": "hello-world", "status": ""},
		},
		{
			"/users/:id<\\d+>/profile/:section?",
			"/users/123/profile",
			map[string]string{"id": "123", "section": ""},
		},
	}

	for _, tc := range testCases {
		params := router.extractParams(tc.pattern, tc.path)

		// Compare the maps
		if len(params) != len(tc.params) {
			t.Errorf("extractParams(%q, %q) returned %d params, expected %d",
				tc.pattern, tc.path, len(params), len(tc.params))
			continue
		}

		for k, expected := range tc.params {
			if actual, ok := params[k]; !ok || actual != expected {
				t.Errorf("extractParams(%q, %q): param[%q] = %q, expected %q",
					tc.pattern, tc.path, k, actual, expected)
			}
		}
	}
}

// TestContext_Param tests the Param and ParamMap methods in the context
func TestRouteParamsInContext(t *testing.T) {
	// Setup router
	r := NewRouter()

	// Define test route
	r.Handle("GET", "/users/:id/posts/:slug/*rest", func(ctx context.Context) {
		// Test individual param access
		id := ctx.Param("id")
		slug := ctx.Param("slug")
		rest := ctx.Param("rest")

		// Test ParamMap
		allParams := ctx.ParamMap()

		ctx.String(http.StatusOK, "id:%s,slug:%s,rest:%s,count:%d",
			id, slug, rest, len(allParams))
	})

	// Make request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/123/posts/hello-world/comments/1", nil)
	r.ServeHTTP(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expected := "id:123,slug:hello-world,rest:comments/1,count:3"
	if w.Body.String() != expected {
		t.Errorf("Expected body %q, got %q", expected, w.Body.String())
	}
}
