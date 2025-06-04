package router

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"testing"
)

// TestSplitPathMemoryUsage tests memory usage patterns of the splitPath caching system
func TestSplitPathMemoryUsage(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	// Clear cache and reset stats
	router.ClearSplitPathCache()
	router.ResetSplitPathStats()

	// Get initial memory stats
	var m1 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Perform many splitPath operations to fill cache
	paths := []string{
		"/api/v1/users",
		"/api/v1/users/123",
		"/api/v1/posts",
		"/api/v1/posts/456",
		"/admin/dashboard",
		"/admin/users",
		"/public/assets/css",
		"/public/assets/js",
		"/health/check",
		"/metrics/prometheus",
	}

	// Fill cache with various paths
	for i := 0; i < 100; i++ {
		for _, path := range paths {
			router.splitPath(path)
		}
	}

	// Get memory stats after cache operations
	var m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m2)

	// Check cache stats
	cacheSize, hitRatio, totalHits, totalMisses, _ := router.GetSplitPathCacheStats()
	t.Logf("Cache size: %d", cacheSize)
	t.Logf("Hit ratio: %d%%", hitRatio)
	t.Logf("Total hits: %d", totalHits)
	t.Logf("Total misses: %d", totalMisses)

	// Memory usage analysis
	memUsed := m2.Alloc - m1.Alloc
	t.Logf("Memory used by cache operations: %d bytes", memUsed)
	t.Logf("Allocs during test: %d", m2.Mallocs-m1.Mallocs)
	t.Logf("Frees during test: %d", m2.Frees-m1.Frees)

	// Ensure cache size is reasonable
	if cacheSize > 1000 {
		t.Errorf("Cache size %d exceeds maximum expected size", cacheSize)
	}

	// Ensure hit ratio is high for repeated operations
	if hitRatio < 90 {
		t.Errorf("Hit ratio %d%% is too low for repeated operations", hitRatio)
	}
}

// TestSplitPathCacheEviction tests cache eviction behavior
func TestSplitPathCacheEviction(t *testing.T) {
	router := NewRouter().(*DefaultRouter)

	// Set small cache size for testing eviction
	router.SetSplitPathCacheConfig(50, 50) // max 50 entries, evict 50%
	router.ClearSplitPathCache()
	router.ResetSplitPathStats()

	// Fill cache beyond capacity
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("/path/to/resource/%d", i)
		router.splitPath(path)
	}

	cacheSize, _, _, _, totalRequests := router.GetSplitPathCacheStats()
	t.Logf("Cache size after filling: %d", cacheSize)
	t.Logf("Total requests: %d", totalRequests)

	// Cache should have been evicted and size should be reasonable
	if cacheSize > 50 {
		t.Errorf("Cache size %d should not exceed configured maximum of 50", cacheSize)
	}

	// Reset to default config
	router.SetSplitPathCacheConfig(1000, 33)
}

// BenchmarkSplitPathMemoryAllocation benchmarks memory allocations
func BenchmarkSplitPathMemoryAllocation(b *testing.B) {
	router := NewRouter().(*DefaultRouter)
	router.ClearSplitPathCache()

	paths := []string{
		"/api/v1/users/123/posts/456",
		"/admin/dashboard/settings",
		"/public/assets/css/style.css",
		"/health/check/detailed",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		path := paths[i%len(paths)]
		router.splitPath(path)
	}
}

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
