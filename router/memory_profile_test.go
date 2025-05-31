package router

import (
	"fmt"
	"runtime"
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
