package router

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

// BenchmarkSplitPath tests the performance of the optimized splitPath method
func BenchmarkSplitPath(b *testing.B) {
	router := NewRouter().(*DefaultRouter)

	// Test cases with various path patterns
	testPaths := []string{
		"/",
		"/api",
		"/api/v1",
		"/api/v1/users",
		"/api/v1/users/123",
		"/api/v1/users/123/posts",
		"/api/v1/users/123/posts/456",
		"/static/css/main.css",
		"/static/js/app.js",
		"/admin/dashboard",
		"/admin/users/management",
		"/health",
		"/metrics",
		"/ping",
		"/favicon.ico",
		"/robots.txt",
		"/sitemap.xml",
		"/very/long/path/with/many/segments/that/tests/performance",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		path := testPaths[i%len(testPaths)]
		router.splitPath(path)
	}
}

// BenchmarkSplitPathCacheMiss tests performance when cache misses occur
func BenchmarkSplitPathCacheMiss(b *testing.B) {
	router := NewRouter().(*DefaultRouter)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Generate unique paths to force cache misses
		path := fmt.Sprintf("/unique/path/%d/test/%d", i, rand.Int())
		router.splitPath(path)
	}
}

// BenchmarkSplitPathCacheHit tests performance when cache hits occur
func BenchmarkSplitPathCacheHit(b *testing.B) {
	router := NewRouter().(*DefaultRouter)

	// Pre-populate cache
	testPaths := []string{
		"/api/v1/users",
		"/static/css/main.css",
		"/admin/dashboard",
	}

	for _, path := range testPaths {
		router.splitPath(path)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		path := testPaths[i%len(testPaths)]
		router.splitPath(path)
	}
}

// BenchmarkSplitPathCommonPaths tests performance for pre-computed common paths
func BenchmarkSplitPathCommonPaths(b *testing.B) {
	router := NewRouter().(*DefaultRouter)

	commonTestPaths := []string{
		"/",
		"/api",
		"/users",
		"/health",
		"/favicon.ico",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		path := commonTestPaths[i%len(commonTestPaths)]
		router.splitPath(path)
	}
}

// BenchmarkSplitPathConcurrent tests concurrent access performance
func BenchmarkSplitPathConcurrent(b *testing.B) {
	router := NewRouter().(*DefaultRouter)

	testPaths := []string{
		"/api/v1/users",
		"/api/v1/posts",
		"/static/assets",
		"/admin/settings",
		"/health/check",
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			path := testPaths[rand.Intn(len(testPaths))]
			router.splitPath(path)
		}
	})
}

// TestSplitPathCacheStats tests the cache statistics functionality
func TestSplitPathCacheStats(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	// Reset stats
	router.ResetSplitPathStats()

	// Test various paths
	testPaths := []string{
		"/api/v1/users", // Cache miss
		"/api/v1/users", // Cache hit
		"/static/css",   // Cache miss
		"/api/v1/users", // Cache hit
		"/health",       // Common path hit
	}

	for _, path := range testPaths {
		router.splitPath(path)
	}

	cacheSize, hitRatio, totalHits, totalMisses, totalRequests := router.GetSplitPathCacheStats()

	t.Logf("Cache size: %d", cacheSize)
	t.Logf("Hit ratio: %d%%", hitRatio)
	t.Logf("Total hits: %d", totalHits)
	t.Logf("Total misses: %d", totalMisses)
	t.Logf("Total requests: %d", totalRequests)

	if totalRequests != int64(len(testPaths)) {
		t.Errorf("Expected %d total requests, got %d", len(testPaths), totalRequests)
	}

	if hitRatio < 0 || hitRatio > 100 {
		t.Errorf("Hit ratio should be between 0-100, got %d", hitRatio)
	}
}

// TestSplitPathCacheConfig tests the cache configuration functionality
func TestSplitPathCacheConfig(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	// Test default config
	maxSize, evictPercent := router.GetSplitPathCacheConfig()
	if maxSize != 1000 {
		t.Errorf("Expected default max size 1000, got %d", maxSize)
	}
	if evictPercent != 33 {
		t.Errorf("Expected default evict percent 33, got %d", evictPercent)
	}

	// Test setting config
	router.SetSplitPathCacheConfig(500, 50)
	maxSize, evictPercent = router.GetSplitPathCacheConfig()
	if maxSize != 500 {
		t.Errorf("Expected max size 500, got %d", maxSize)
	}
	if evictPercent != 50 {
		t.Errorf("Expected evict percent 50, got %d", evictPercent)
	}
}

// TestSplitPathEdgeCases tests various edge cases and correctness
func TestSplitPathEdgeCases(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	testCases := []struct {
		path     string
		expected []string
	}{
		{"", []string{}},
		{"/", []string{}},
		{"//", []string{}},
		{"///", []string{}},
		{"/api", []string{"api"}},
		{"/api/", []string{"api"}},
		{"//api//", []string{"api"}},
		{"/api/v1", []string{"api", "v1"}},
		{"/api/v1/users", []string{"api", "v1", "users"}},
		{"api", []string{"api"}},
		{"api/v1", []string{"api", "v1"}},
		{"/api//v1///users/", []string{"api", "v1", "users"}},
	}

	for _, tc := range testCases {
		result := router.splitPath(tc.path)
		if len(result) != len(tc.expected) {
			t.Errorf("Path %q: expected %d segments, got %d", tc.path, len(tc.expected), len(result))
			continue
		}

		for i, expected := range tc.expected {
			if result[i] != expected {
				t.Errorf("Path %q: segment %d expected %q, got %q", tc.path, i, expected, result[i])
			}
		}
	}
}

// Benchmark comparing with naive string splitting
func BenchmarkSplitPathVsStringsSplit(b *testing.B) {
	router := NewRouter().(*DefaultRouter)
	testPath := "/api/v1/users/123/posts/456/comments"

	b.Run("OptimizedSplitPath", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			router.splitPath(testPath)
		}
	})

	b.Run("StringsSplit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			naiveSplit(testPath)
		}
	})
}

// naiveSplit simulates the old approach using strings.Split
func naiveSplit(path string) []string {
	if path == "" || path == "/" {
		return []string{}
	}

	// Remove leading and trailing slashes
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}
