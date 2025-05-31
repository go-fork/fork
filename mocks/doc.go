// Package mocks provides mock implementations for the http context interfaces.
//
// This package includes mock implementations for:
// - httpContext.Context - via MockContext
// - httpContext.Request - via MockRequest
// - httpContext.Response - via MockResponse
//
// These mocks are designed to facilitate testing of middleware and components
// that depend on the HTTP context interfaces. They track method calls and manipulations
// so that tests can assert on expected behavior.
//
// Usage Example:
//
//	func TestMiddleware(t *testing.T) {
//	    // Create a new request
//	    req, _ := http.NewRequest("GET", "/test", nil)
//
//	    // Create a mock ResponseWriter
//	    rw := httptest.NewRecorder()
//
//	    // Create a mock context
//	    ctx := mocks.NewMockContext(rw, req)
//
//	    // Run middleware with the mock context
//	    middleware(ctx)
//
//	    // Assert on expected behavior
//	    if !ctx.NextCalled {
//	        t.Error("Expected middleware to call Next()")
//	    }
//
//	    if ctx.StatusCode != http.StatusOK {
//	        t.Errorf("Expected status 200, got %d", ctx.StatusCode)
//	    }
//	}
//
// The mocks in this package implement all methods required by the corresponding interfaces
// and track state changes for testing purposes. They are designed to be simple to use
// while providing enough information for comprehensive testing.
package mocks
