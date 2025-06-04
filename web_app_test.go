package fork_test

import (
	"fmt"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go.fork.vn/fork"
	forkContext "go.fork.vn/fork/context"
	forkErrors "go.fork.vn/fork/errors"
	fork_mocks "go.fork.vn/fork/mocks"
	forkRouter "go.fork.vn/fork/router"
)

// TestWebApp_NewWebApp tests the creation of a new WebApp instance
func TestWebApp_NewWebApp(t *testing.T) {
	t.Run("creates new app successfully", func(t *testing.T) {
		app := fork.NewWebApp()
		assert.NotNil(t, app)

		// Verify default config is set
		config := app.GetConfig()
		assert.NotNil(t, config)
		assert.True(t, config.GracefulShutdown.Enabled)
		assert.Equal(t, 30, config.GracefulShutdown.Timeout)
	})
}

// TestWebApp_SetAdapter tests setting an adapter
func TestWebApp_SetAdapter(t *testing.T) {
	app := fork.NewWebApp()
	mockAdapter := fork_mocks.NewMockAdapter(t)

	mockAdapter.EXPECT().SetHandler(mock.AnythingOfType("*router.DefaultRouter")).Once()

	app.SetAdapter(mockAdapter)

	// Verify adapter is set
	assert.Equal(t, mockAdapter, app.GetAdapter())
	mockAdapter.AssertExpectations(t)
}

// TestWebApp_Use tests middleware registration
func TestWebApp_Use(t *testing.T) {
	app := fork.NewWebApp()

	middlewareCalled := false
	middleware := func(ctx forkContext.Context) {
		middlewareCalled = true
		ctx.Next()
	}

	app.Use(middleware)

	// Middleware should be registered but not called yet
	assert.False(t, middlewareCalled)
}

// TestWebApp_HTTPMethods tests all HTTP method handlers
func TestWebApp_HTTPMethods(t *testing.T) {
	methods := []struct {
		name       string
		method     string
		setupRoute func(*fork.WebApp, string, forkRouter.HandlerFunc)
	}{
		{"GET", fork.MethodGet, func(app *fork.WebApp, path string, handler forkRouter.HandlerFunc) { app.GET(path, handler) }},
		{"POST", fork.MethodPost, func(app *fork.WebApp, path string, handler forkRouter.HandlerFunc) { app.POST(path, handler) }},
		{"PUT", fork.MethodPut, func(app *fork.WebApp, path string, handler forkRouter.HandlerFunc) { app.PUT(path, handler) }},
		{"DELETE", fork.MethodDelete, func(app *fork.WebApp, path string, handler forkRouter.HandlerFunc) { app.DELETE(path, handler) }},
		{"PATCH", fork.MethodPatch, func(app *fork.WebApp, path string, handler forkRouter.HandlerFunc) { app.PATCH(path, handler) }},
		{"HEAD", fork.MethodHead, func(app *fork.WebApp, path string, handler forkRouter.HandlerFunc) { app.HEAD(path, handler) }},
		{"OPTIONS", fork.MethodOptions, func(app *fork.WebApp, path string, handler forkRouter.HandlerFunc) { app.OPTIONS(path, handler) }},
	}

	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			app := fork.NewWebApp()

			handlerFunc := func(ctx forkContext.Context) {
				ctx.JSON(200, map[string]string{"message": "success"})
			}

			// This should not panic
			assert.NotPanics(t, func() {
				method.setupRoute(app, "/test", handlerFunc)
			})
		})
	}
}

// TestWebApp_Any tests the Any method that registers all HTTP methods
func TestWebApp_Any(t *testing.T) {
	app := fork.NewWebApp()

	handlerFunc := func(ctx forkContext.Context) {
		ctx.JSON(200, map[string]string{"message": "any method"})
	}

	// Should not panic when registering Any route
	assert.NotPanics(t, func() {
		app.Any("/any", handlerFunc)
	})
}

// TestWebApp_Group tests router grouping functionality
func TestWebApp_Group(t *testing.T) {
	app := fork.NewWebApp()

	group := app.Group("/api/v1")
	assert.NotNil(t, group)

	// Group should be able to register routes
	assert.NotPanics(t, func() {
		group.Handle(fork.MethodGet, "/users", func(ctx forkContext.Context) {
			ctx.JSON(200, map[string]string{"message": "users"})
		})
	})
}

// TestWebApp_Static tests static file serving
func TestWebApp_Static(t *testing.T) {
	app := fork.NewWebApp()

	// Should not panic when setting up static routes
	assert.NotPanics(t, func() {
		app.Static("/static", "./public")
	})
}

// TestWebApp_ServeHTTP tests the HTTP handler functionality
func TestWebApp_ServeHTTP(t *testing.T) {
	app := fork.NewWebApp()

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Add a simple route
	app.GET("/test", func(ctx forkContext.Context) {
		ctx.JSON(200, map[string]string{"message": "test"})
	})

	// This should not panic
	assert.NotPanics(t, func() {
		app.ServeHTTP(w, req)
	})
}

// TestWebApp_ConcurrentAccess tests thread safety
func TestWebApp_ConcurrentAccess(t *testing.T) {
	app := fork.NewWebApp()

	var wg sync.WaitGroup
	numGoroutines := 100

	// Test concurrent middleware registration
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			app.Use(func(ctx forkContext.Context) {
				ctx.Set("index", index)
				ctx.Next()
			})
		}(i)
	}

	// Test concurrent route registration
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			path := fmt.Sprintf("/test-%d", index)
			app.GET(path, func(ctx forkContext.Context) {
				ctx.JSON(200, map[string]int{"index": index})
			})
		}(i)
	}

	wg.Wait()

	// No assertions needed - test passes if no race conditions occur
}

// TestWebApp_MiddlewareChain tests middleware execution order
func TestWebApp_MiddlewareChain(t *testing.T) {
	app := fork.NewWebApp()
	var executionOrder []string

	// Add middlewares in order
	app.Use(func(ctx forkContext.Context) {
		executionOrder = append(executionOrder, "middleware1")
		ctx.Next()
	})

	app.Use(func(ctx forkContext.Context) {
		executionOrder = append(executionOrder, "middleware2")
		ctx.Next()
	})

	app.GET("/test", func(ctx forkContext.Context) {
		executionOrder = append(executionOrder, "handler")
		ctx.JSON(200, map[string]string{"message": "ok"})
	})

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	// Verify execution order
	expected := []string{"middleware1", "middleware2", "handler"}
	assert.Equal(t, expected, executionOrder)
}

// TestWebApp_ErrorHandling tests error handling scenarios
func TestWebApp_ErrorHandling(t *testing.T) {
	app := fork.NewWebApp()

	// Test route that returns an error
	app.GET("/error", func(ctx forkContext.Context) {
		err := forkErrors.NewHttpError(500, "Internal Server Error", nil, nil)
		ctx.Error(err)
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
}

// TestWebApp_ContextValues tests context value storage and retrieval
func TestWebApp_ContextValues(t *testing.T) {
	app := fork.NewWebApp()

	app.Use(func(ctx forkContext.Context) {
		ctx.Set("user_id", 12345)
		ctx.Set("user_name", "john_doe")
		ctx.Next()
	})

	app.GET("/profile", func(ctx forkContext.Context) {
		userID := ctx.GetInt("user_id")
		userName := ctx.GetString("user_name")

		ctx.JSON(200, map[string]interface{}{
			"user_id":   userID,
			"user_name": userName,
		})
	})

	req := httptest.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// TestWebApp_AbortMiddleware tests middleware abort functionality
func TestWebApp_AbortMiddleware(t *testing.T) {
	app := fork.NewWebApp()
	var handlerCalled bool

	app.Use(func(ctx forkContext.Context) {
		// Simulate authentication failure
		ctx.JSON(401, map[string]string{"error": "unauthorized"})
		ctx.Abort()
	})

	app.GET("/protected", func(ctx forkContext.Context) {
		handlerCalled = true
		ctx.JSON(200, map[string]string{"message": "protected resource"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.False(t, handlerCalled, "Handler should not be called when middleware aborts")
}

// TestWebApp_RouterGroup tests router group functionality
func TestWebApp_RouterGroup(t *testing.T) {
	app := fork.NewWebApp()

	// Create API v1 group
	v1 := app.Group("/api/v1")
	v1.Use(func(ctx forkContext.Context) {
		ctx.Set("api_version", "v1")
		ctx.Next()
	})

	v1.Handle(fork.MethodGet, "/users", func(ctx forkContext.Context) {
		version := ctx.GetString("api_version")
		ctx.JSON(200, map[string]string{"version": version, "endpoint": "users"})
	})

	// Create API v2 group
	v2 := app.Group("/api/v2")
	v2.Use(func(ctx forkContext.Context) {
		ctx.Set("api_version", "v2")
		ctx.Next()
	})

	v2.Handle(fork.MethodGet, "/users", func(ctx forkContext.Context) {
		version := ctx.GetString("api_version")
		ctx.JSON(200, map[string]string{"version": version, "endpoint": "users"})
	})

	// Test v1 endpoint
	req1 := httptest.NewRequest("GET", "/api/v1/users", nil)
	w1 := httptest.NewRecorder()
	app.ServeHTTP(w1, req1)
	assert.Equal(t, 200, w1.Code)

	// Test v2 endpoint
	req2 := httptest.NewRequest("GET", "/api/v2/users", nil)
	w2 := httptest.NewRecorder()
	app.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)
}

// TestWebApp_ParameterHandling tests URL parameter handling
func TestWebApp_ParameterHandling(t *testing.T) {
	app := fork.NewWebApp()

	app.GET("/users/:id", func(ctx forkContext.Context) {
		userID := ctx.Param("id")
		ctx.JSON(200, map[string]string{"user_id": userID})
	})

	app.GET("/posts/:postId/comments/:commentId", func(ctx forkContext.Context) {
		postID := ctx.Param("postId")
		commentID := ctx.Param("commentId")
		ctx.JSON(200, map[string]string{
			"post_id":    postID,
			"comment_id": commentID,
		})
	})

	// Test single parameter
	req1 := httptest.NewRequest("GET", "/users/123", nil)
	w1 := httptest.NewRecorder()
	app.ServeHTTP(w1, req1)
	assert.Equal(t, 200, w1.Code)

	// Test multiple parameters
	req2 := httptest.NewRequest("GET", "/posts/456/comments/789", nil)
	w2 := httptest.NewRecorder()
	app.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)
}

// TestWebApp_QueryParameters tests query parameter handling
func TestWebApp_QueryParameters(t *testing.T) {
	app := fork.NewWebApp()

	app.GET("/search", func(ctx forkContext.Context) {
		query := ctx.Query("q")
		page := ctx.DefaultQuery("page", "1")
		limit := ctx.DefaultQuery("limit", "10")

		ctx.JSON(200, map[string]string{
			"query": query,
			"page":  page,
			"limit": limit,
		})
	})

	req := httptest.NewRequest("GET", "/search?q=golang&page=2", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// BenchmarkWebApp_SimpleRoute benchmarks a simple route handler
func BenchmarkWebApp_SimpleRoute(b *testing.B) {
	app := fork.NewWebApp()
	app.GET("/test", func(ctx forkContext.Context) {
		ctx.JSON(200, map[string]string{"message": "hello"})
	})

	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkWebApp_WithMiddleware benchmarks route with middleware
func BenchmarkWebApp_WithMiddleware(b *testing.B) {
	app := fork.NewWebApp()

	// Add multiple middlewares
	app.Use(func(ctx forkContext.Context) {
		ctx.Set("timestamp", time.Now())
		ctx.Next()
	})

	app.Use(func(ctx forkContext.Context) {
		ctx.Set("request_id", "req-123")
		ctx.Next()
	})

	app.GET("/test", func(ctx forkContext.Context) {
		ctx.JSON(200, map[string]string{"message": "hello"})
	})

	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkWebApp_ParameterRoute benchmarks parameterized routes
func BenchmarkWebApp_ParameterRoute(b *testing.B) {
	app := fork.NewWebApp()
	app.GET("/users/:id/posts/:postId", func(ctx forkContext.Context) {
		userID := ctx.Param("id")
		postID := ctx.Param("postId")
		ctx.JSON(200, map[string]string{
			"user_id": userID,
			"post_id": postID,
		})
	})

	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// TestWebApp_GracefulShutdown tests graceful shutdown functionality
func TestWebApp_GracefulShutdown(t *testing.T) {
	app := fork.NewWebApp()
	mockAdapter := fork_mocks.NewMockAdapter(t)

	mockAdapter.EXPECT().SetHandler(mock.AnythingOfType("*router.DefaultRouter")).Maybe()
	mockAdapter.EXPECT().Shutdown().Return(nil).Once()

	app.SetAdapter(mockAdapter)

	// Test shutdown
	err := app.Shutdown()
	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

// TestWebApp_ConfigDefaults tests default configuration values
func TestWebApp_ConfigDefaults(t *testing.T) {
	app := fork.NewWebApp()

	// Test that app is created with default config
	assert.NotNil(t, app)

	config := app.GetConfig()
	assert.NotNil(t, config)
	assert.True(t, config.GracefulShutdown.Enabled)
	assert.Equal(t, 30, config.GracefulShutdown.Timeout)
	assert.True(t, config.GracefulShutdown.WaitForConnections)
}

// TestWebApp_SetConfig tests configuration management
func TestWebApp_SetConfig(t *testing.T) {
	app := fork.NewWebApp()

	customConfig := &fork.WebAppConfig{
		GracefulShutdown: fork.GracefulShutdownConfig{
			Enabled:            false,
			Timeout:            60,
			WaitForConnections: false,
		},
	}

	app.SetConfig(customConfig)

	config := app.GetConfig()
	assert.False(t, config.GracefulShutdown.Enabled)
	assert.Equal(t, 60, config.GracefulShutdown.Timeout)
	assert.False(t, config.GracefulShutdown.WaitForConnections)
}

// TestWebApp_ConnectionTracking tests connection tracking functionality
func TestWebApp_ConnectionTracking(t *testing.T) {
	app := fork.NewWebApp()

	// Initially should have 0 connections
	assert.Equal(t, int32(0), app.GetActiveConnections())

	// Track a connection
	app.TrackConnection()
	assert.Equal(t, int32(1), app.GetActiveConnections())

	// Track another connection
	app.TrackConnection()
	assert.Equal(t, int32(2), app.GetActiveConnections())

	// Untrack a connection
	app.UntrackConnection()
	assert.Equal(t, int32(1), app.GetActiveConnections())

	// Untrack the last connection
	app.UntrackConnection()
	assert.Equal(t, int32(0), app.GetActiveConnections())
}

// TestWebApp_ShutdownTimeout tests shutdown timeout configuration
func TestWebApp_ShutdownTimeout(t *testing.T) {
	app := fork.NewWebApp()

	// Set custom timeout
	app.SetShutdownTimeout(45 * time.Second)

	config := app.GetConfig()
	assert.Equal(t, 45, config.GracefulShutdown.Timeout)
}

// TestWebApp_IsShuttingDown tests shutdown state tracking
func TestWebApp_IsShuttingDown(t *testing.T) {
	app := fork.NewWebApp()

	// Initially should not be shutting down
	assert.False(t, app.IsShuttingDown())

	// After calling GracefulShutdown, it should be marked as shutting down
	// Note: This test may need to be adjusted based on actual implementation
	mockAdapter := fork_mocks.NewMockAdapter(t)
	mockAdapter.EXPECT().SetHandler(mock.AnythingOfType("*router.DefaultRouter")).Maybe()
	mockAdapter.EXPECT().Shutdown().Return(nil).Once()

	app.SetAdapter(mockAdapter)
	app.GracefulShutdown()

	assert.True(t, app.IsShuttingDown())
}

// TestWebApp_EnableSecurityMiddleware tests security middleware activation
func TestWebApp_EnableSecurityMiddleware(t *testing.T) {
	app := fork.NewWebApp()

	// Should not panic when enabling security middleware
	assert.NotPanics(t, func() {
		app.EnableSecurityMiddleware()
	})
}

// TestWebApp_Router tests router access
func TestWebApp_Router(t *testing.T) {
	app := fork.NewWebApp()

	router := app.Router()
	assert.NotNil(t, router)

	// Router should be able to handle routes
	assert.NotPanics(t, func() {
		router.Handle(fork.MethodGet, "/test", func(ctx forkContext.Context) {
			ctx.JSON(200, map[string]string{"message": "test"})
		})
	})
}

// TestWebApp_NewContext tests context creation
func TestWebApp_NewContext(t *testing.T) {
	app := fork.NewWebApp()

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	ctx := app.NewContext(w, req)
	assert.NotNil(t, ctx)
}

// TestWebApp_CleanupResources tests resource cleanup
func TestWebApp_CleanupResources(t *testing.T) {
	app := fork.NewWebApp()

	// Should not panic during cleanup
	assert.NotPanics(t, func() {
		app.CleanupResources()
	})
}

// Helper function to create a test context
func createTestContext(method, path string) (forkContext.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	return forkContext.NewContext(w, req), w
}

// TestWebApp_Integration tests integration scenarios
func TestWebApp_Integration(t *testing.T) {
	app := fork.NewWebApp()

	// Set custom configuration
	config := &fork.WebAppConfig{
		GracefulShutdown: fork.GracefulShutdownConfig{
			Enabled:            true,
			Timeout:            30,
			WaitForConnections: true,
		},
	}
	app.SetConfig(config)

	// Add global middleware
	app.Use(func(ctx forkContext.Context) {
		ctx.Header("X-App-Name", "Fork Framework")
		ctx.Next()
	})

	// Add API routes
	api := app.Group("/api")
	api.Use(func(ctx forkContext.Context) {
		ctx.Header("X-API-Version", "1.0")
		ctx.Next()
	})

	api.Handle(fork.MethodGet, "/health", func(ctx forkContext.Context) {
		ctx.JSON(200, map[string]string{"status": "healthy"})
	})

	api.Handle(fork.MethodPost, "/users", func(ctx forkContext.Context) {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := ctx.Bind(&user); err != nil {
			ctx.JSON(400, map[string]string{"error": "invalid request"})
			return
		}
		ctx.JSON(201, map[string]string{"message": "user created"})
	})

	// Test health check
	req1 := httptest.NewRequest("GET", "/api/health", nil)
	w1 := httptest.NewRecorder()
	app.ServeHTTP(w1, req1)
	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, "Fork Framework", w1.Header().Get("X-App-Name"))
	assert.Equal(t, "1.0", w1.Header().Get("X-API-Version"))

	// Test user creation (would need actual request body)
	req2 := httptest.NewRequest("POST", "/api/users", nil)
	w2 := httptest.NewRecorder()
	app.ServeHTTP(w2, req2)
	// Status depends on implementation - could be 400 due to empty body
}

// TestWebApp_ErrorScenarios tests various error scenarios
func TestWebApp_ErrorScenarios(t *testing.T) {
	t.Run("serve without adapter", func(t *testing.T) {
		app := fork.NewWebApp()

		err := app.Serve()
		assert.Error(t, err)
		assert.Equal(t, fork.ErrAdapterNotSet, err)
	})

	t.Run("runTLS without adapter", func(t *testing.T) {
		app := fork.NewWebApp()

		err := app.RunTLS("cert.pem", "key.pem")
		assert.Error(t, err)
		assert.Equal(t, fork.ErrAdapterNotSet, err)
	})

	t.Run("runTLS with invalid cert files", func(t *testing.T) {
		app := fork.NewWebApp()
		mockAdapter := fork_mocks.NewMockAdapter(t)

		mockAdapter.EXPECT().SetHandler(mock.AnythingOfType("*router.DefaultRouter")).Once()

		app.SetAdapter(mockAdapter)

		err := app.RunTLS("", "")
		assert.Error(t, err)
		assert.Equal(t, fork.ErrInvalidCertificate, err)
	})
}

// TestWebApp_ConfigValidation tests configuration validation
func TestWebApp_ConfigValidation(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := &fork.WebAppConfig{
			GracefulShutdown: fork.GracefulShutdownConfig{
				Enabled:            true,
				Timeout:            30,
				WaitForConnections: true,
				SignalBufferSize:   1,
			},
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid timeout", func(t *testing.T) {
		config := &fork.WebAppConfig{
			GracefulShutdown: fork.GracefulShutdownConfig{
				Enabled:            true,
				Timeout:            -1, // Invalid negative timeout
				WaitForConnections: true,
				SignalBufferSize:   1,
			},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Equal(t, fork.ErrInvalidConfiguration, err)
	})

	t.Run("invalid signal buffer size", func(t *testing.T) {
		config := &fork.WebAppConfig{
			GracefulShutdown: fork.GracefulShutdownConfig{
				Enabled:            true,
				Timeout:            30,
				WaitForConnections: true,
				SignalBufferSize:   0, // Invalid buffer size
			},
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Equal(t, fork.ErrInvalidConfiguration, err)
	})
}
