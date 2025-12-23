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

	// Cache is cleared - no way to directly verify size with current API,
	// but the Get operations above confirm all items were removed
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
func TestClient_InvalidateURLCache_Disabled(t *testing.T) {
	client := &Client{
		Cache:        NewCache(),
		cacheEnabled: false,
	}

	// Should not panic when cache is disabled
	client.InvalidateURLCache("/api/v1/schemas/test-schema")
	client.ClearCache()

	// Cache operations should work even when disabled (just won't affect anything)
	// Test passes if no panics occur
}

func TestClient_CacheInvalidationMethods(t *testing.T) {
	client := &Client{
		Cache:        NewCache(),
		cacheEnabled: true,
	}

	// Manually add items to cache to test invalidation
	client.Cache.Set("/api/v1/schemas/123", "test_data_123")
	client.Cache.Set("/api/v1/schemas/456", "test_data_456")
	client.Cache.Set("/api/v1/templates/789", "test_data_789")

	cloneFunc := func(item interface{}) (interface{}, error) {
		return item, nil
	}

	// Verify items exist
	_, found1, _ := client.Cache.Get("/api/v1/schemas/123", cloneFunc)
	_, found2, _ := client.Cache.Get("/api/v1/schemas/456", cloneFunc)
	if !found1 || !found2 {
		t.Error("Schema items should exist before invalidation")
	}

	// Test individual URL invalidation
	client.InvalidateURLCache("/api/v1/schemas/123")
	_, found1, _ = client.Cache.Get("/api/v1/schemas/123", cloneFunc)
	_, found2, _ = client.Cache.Get("/api/v1/schemas/456", cloneFunc)

	if found1 {
		t.Error("Schema 123 should be invalidated")
	}
	if !found2 {
		t.Error("Schema 456 should still exist")
	}

	// Test bulk invalidation using ClearCache (clears all items)
	client.ClearCache()
	_, found2, _ = client.Cache.Get("/api/v1/schemas/456", cloneFunc)
	_, found3, _ := client.Cache.Get("/api/v1/templates/789", cloneFunc)

	if found2 {
		t.Error("Schema 456 should be invalidated by clear operation")
	}
	if found3 {
		t.Error("Other items should also be cleared by clear operation")
	}
}

// Test 5: Cache Hit/Miss Behavior Tests
func TestCache_HitMissBehavior(t *testing.T) {
	cache := NewCache()

	cloneFunc := func(item interface{}) (interface{}, error) {
		return item, nil
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

	// Cache behavior is correct - hits return items, misses return false
}

// Test 6: Cache Data Storage Tests
func TestCache_DataStorage(t *testing.T) {
	cache := NewCache()

	cloneFunc := func(item interface{}) (interface{}, error) {
		return item, nil
	}

	// Test storing byte data (typical use case for JSON caching)
	testData := []byte(`{"test": "data for memory calculation"}`)
	cache.Set("key1", testData)

	result, found, err := cache.Get("key1", cloneFunc)
	if !found || err != nil {
		t.Error("Should retrieve cached byte data")
	}

	retrievedData := result.([]byte)
	if string(retrievedData) != string(testData) {
		t.Errorf("Retrieved data doesn't match original. Expected %s, got %s", string(testData), string(retrievedData))
	}
}

// Test 7: GetViaURLWithCache Method Tests
func TestClient_GetViaURLWithCache_Disabled(t *testing.T) {
	// Create a properly initialized client similar to other tests
	client := GetClient("https://test.example.com", "testuser", CacheEnabled(false))

	// Test should use regular GetViaURL when caching disabled
	// This test would require mocking HTTP requests, so we'll just verify the method exists
	// and doesn't panic when called with caching disabled
	_, err := client.GetViaURLWithCache("/test/url")
	// We expect an error since we don't have a real HTTP client setup,
	// but the important thing is that it doesn't panic
	if err == nil {
		t.Error("Expected error for non-working HTTP client")
	}
}

func TestClient_DetectURLResourceType(t *testing.T) {
	client := &Client{}

	testCases := []struct {
		url          string
		expectedType string
	}{
		{"/api/v1/schemas/123", "SCHEMA"},
		{"/api/v1/templates/456", "TEMPLATE"},
		{"/api/v1/sites/789", "SITE"},
		{"/api/v1/users/abc", "USER"},
		{"/api/v1/tenants/def", "TENANT"},
		{"/api/v1/labels/ghi", "LABEL"},
		{"/api/v1/remote-locations/jkl", "REMOTE_LOCATION"},
		{"/api/v1/unknown/mno", "RESOURCE"},
	}

	for _, tc := range testCases {
		result := client.detectURLResourceType(tc.url)
		if result != tc.expectedType {
			t.Errorf("detectURLResourceType(%s) = %s, expected %s", tc.url, result, tc.expectedType)
		}
	}
}