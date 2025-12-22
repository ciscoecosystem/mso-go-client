package client

import (
	"testing"
)

// Test 1: Cache Basic Operations Tests
func TestCache_BasicOperations(t *testing.T) {
	cache := NewCache()

	// Test Set and Get
	testValue := "test_data"
	cache.Set("key1", testValue)

	cloneFunc := func(item interface{}) (interface{}, error) {
		return item.(string), nil
	}

	result, found, err := cache.Get("key1", cloneFunc)
	if !found || err != nil || result.(string) != testValue {
		t.Errorf("Expected to find cached value, got found=%v, err=%v, result=%v", found, err, result)
	}

	// Test cache miss
	_, found, _ = cache.Get("nonexistent", cloneFunc)
	if found {
		t.Error("Should not find nonexistent key")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := NewCache()

	// Set and verify item exists
	cache.Set("key1", "value1")
	cloneFunc := func(item interface{}) (interface{}, error) {
		return item, nil
	}

	_, found, _ := cache.Get("key1", cloneFunc)
	if !found {
		t.Error("Item should exist before deletion")
	}

	// Delete and verify item is gone
	cache.Delete("key1")
	_, found, _ = cache.Get("key1", cloneFunc)
	if found {
		t.Error("Item should not exist after deletion")
	}

	// Verify invalidation was recorded
	_, _, invalidations, _ := cache.GetStats()
	if invalidations != 1 {
		t.Errorf("Expected 1 invalidation, got %d", invalidations)
	}
}

// Test 2: Cache Clear Functionality Tests
func TestCache_Clear(t *testing.T) {
	cache := NewCache()

	// Set multiple items
	cache.Set("schema_123", "data1")
	cache.Set("schema_456", "data2")
	cache.Set("other_789", "data3")

	cloneFunc := func(item interface{}) (interface{}, error) {
		return item, nil
	}

	// Verify all items exist before clearing
	_, found1, _ := cache.Get("schema_123", cloneFunc)
	_, found2, _ := cache.Get("schema_456", cloneFunc)
	_, found3, _ := cache.Get("other_789", cloneFunc)

	if !found1 || !found2 || !found3 {
		t.Error("All items should exist before clearing")
	}

	// Clear cache
	cache.Clear()

	// Verify all items are deleted
	_, found1, _ = cache.Get("schema_123", cloneFunc)
	_, found2, _ = cache.Get("schema_456", cloneFunc)
	_, found3, _ = cache.Get("other_789", cloneFunc)

	if found1 || found2 || found3 {
		t.Error("All items should be deleted after clearing")
	}

	// Verify cache size is zero
	if cache.Size() != 0 {
		t.Errorf("Cache size should be 0 after clear, got %d", cache.Size())
	}
}

// Test 3: Cache Enable/Disable Option Tests
func TestCacheEnabled_Option(t *testing.T) {
	client := &Client{}

	// Test enabling cache
	CacheEnabled(true)(client)
	if !client.cacheEnabled {
		t.Error("Cache should be enabled after CacheEnabled(true)")
	}

	// Test disabling cache
	CacheEnabled(false)(client)
	if client.cacheEnabled {
		t.Error("Cache should be disabled after CacheEnabled(false)")
	}
}

// Test 4: Client Cache Integration Tests
func TestClient_InvalidateSchemaCache_Disabled(t *testing.T) {
	client := &Client{
		Cache:        NewCache(),
		cacheEnabled: false,
	}

	// Should not panic when cache is disabled
	client.InvalidateSchemaCache("test-schema")
	client.ClearCache()

	// Should return zero stats when cache is disabled - test through the cache directly
	hits, misses, invalidations, hitRatio := client.Cache.GetStats()
	if hits != 0 || misses != 0 || invalidations != 0 || hitRatio != 0 {
		t.Errorf("Cache stats should be zero when disabled, got hits=%d, misses=%d, invalidations=%d, hitRatio=%.1f",
			hits, misses, invalidations, hitRatio)
	}
}

func TestClient_CacheInvalidationMethods(t *testing.T) {
	client := &Client{
		Cache:        NewCache(),
		cacheEnabled: true,
	}

	// Manually add items to cache to test invalidation
	client.Cache.Set("schema_123", "test_data_123")
	client.Cache.Set("schema_456", "test_data_456")
	client.Cache.Set("other_789", "test_data_789")

	cloneFunc := func(item interface{}) (interface{}, error) {
		return item, nil
	}

	// Verify items exist
	_, found1, _ := client.Cache.Get("schema_123", cloneFunc)
	_, found2, _ := client.Cache.Get("schema_456", cloneFunc)
	if !found1 || !found2 {
		t.Error("Schema items should exist before invalidation")
	}

	// Test individual invalidation
	client.InvalidateSchemaCache("123")
	_, found1, _ = client.Cache.Get("schema_123", cloneFunc)
	_, found2, _ = client.Cache.Get("schema_456", cloneFunc)

	if found1 {
		t.Error("Schema 123 should be invalidated")
	}
	if !found2 {
		t.Error("Schema 456 should still exist")
	}

	// Test bulk invalidation using ClearCache (clears all items)
	client.ClearCache()
	_, found2, _ = client.Cache.Get("schema_456", cloneFunc)
	_, found3, _ := client.Cache.Get("other_789", cloneFunc)

	if found2 {
		t.Error("Schema 456 should be invalidated by clear operation")
	}
	if found3 {
		t.Error("Other items should also be cleared by clear operation")
	}
}

// Test 5: Cache Statistics Tests
func TestCache_Statistics(t *testing.T) {
	cache := NewCache()

	cloneFunc := func(item interface{}) (interface{}, error) {
		return item, nil
	}

	// Initial stats should be zero
	hits, misses, invalidations, hitRatio := cache.GetStats()
	if hits != 0 || misses != 0 || invalidations != 0 || hitRatio != 0 {
		t.Errorf("Initial stats should be zero, got hits=%d, misses=%d, invalidations=%d, hitRatio=%.1f",
			hits, misses, invalidations, hitRatio)
	}

	// Add item and test cache hit
	cache.Set("key1", "value1")
	_, found, _ := cache.Get("key1", cloneFunc)
	if !found {
		t.Error("Should find cached item")
	}

	// Test cache miss
	_, found, _ = cache.Get("nonexistent", cloneFunc)
	if found {
		t.Error("Should not find nonexistent item")
	}

	// Check updated stats
	hits, misses, invalidations, hitRatio = cache.GetStats()
	if hits != 1 || misses != 1 {
		t.Errorf("Expected 1 hit and 1 miss, got hits=%d, misses=%d", hits, misses)
	}
	if hitRatio != 50.0 {
		t.Errorf("Expected hit ratio of 50.0%%, got %.1f%%", hitRatio)
	}
}

// Test 6: Cache Memory Statistics Tests
func TestCache_MemoryStatistics(t *testing.T) {
	cache := NewCache()

	// Test with empty cache
	totalBytes, totalMB, avgBytesPerItem, _ := cache.GetMemoryStats()
	if totalBytes != 0 || totalMB != 0 || avgBytesPerItem != 0 {
		t.Errorf("Empty cache should have zero memory stats, got totalBytes=%d, totalMB=%.2f, avgBytesPerItem=%.2f",
			totalBytes, totalMB, avgBytesPerItem)
	}

	// Add some data
	testData := []byte("test data for memory calculation")
	cache.Set("key1", testData)

	totalBytes, totalMB, avgBytesPerItem, _ = cache.GetMemoryStats()
	expectedBytes := int64(len(testData))
	if totalBytes != expectedBytes {
		t.Errorf("Expected totalBytes=%d, got %d", expectedBytes, totalBytes)
	}
	if avgBytesPerItem != float64(expectedBytes) {
		t.Errorf("Expected avgBytesPerItem=%.2f, got %.2f", float64(expectedBytes), avgBytesPerItem)
	}
}