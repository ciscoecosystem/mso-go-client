package client

import (
	"net/url"
	"testing"
)

var TestBaseUrls = [...]string{
	"https://ndo.host.cisco",
	"https://ndo.host.cisco/",
	"https://ndo.host.cisco//",
	"https://ndo.host.cisco///",
}

func AssertFullUrl(t *testing.T, baseUrl string, platform string, method string, path string, expected string) {
	url, err := url.Parse(baseUrl)
	if err != nil {
		t.Fatal(err)
	}
	ndclient := &Client{
		BaseURL:  url,
		platform: platform,
	}

	actual, err := ndclient.MakeFullUrl(method, path)
	if actual != expected || err != nil {
		t.Errorf(`MakeFullUrl("%s", "%s") %s = %q, %v, expected %#q`, method, path, platform, actual, err, expected)
	}
}

func TestMakeFullUrl_Login(t *testing.T) {
	expected := "https://ndo.host.cisco/login"
	for _, baseUrl := range TestBaseUrls {
		AssertFullUrl(t, baseUrl, "nd", "POST", "/login", expected)
		AssertFullUrl(t, baseUrl, "mso", "POST", "/login", expected)
	}
}

func TestMakeFullUrl_Get(t *testing.T) {
	expected := "https://ndo.host.cisco/templates/123"
	expected_nd := "https://ndo.host.cisco/mso/templates/123"
	paths := [...]string{
		"templates/123",
		"/templates/123",
		"///templates/123",
	}
	for _, baseUrl := range TestBaseUrls {
		for _, path := range paths {
			AssertFullUrl(t, baseUrl, "nd", "GET", path, expected_nd)
			AssertFullUrl(t, baseUrl, "mso", "GET", path, expected)
		}
	}
}

func TestMakeFullUrl_Patch(t *testing.T) {
	expected := "https://ndo.host.cisco/templates/123?validate=false"
	expected_nd := "https://ndo.host.cisco/mso/templates/123?validate=false"
	path := "/templates/123"
	for _, baseUrl := range TestBaseUrls {
		AssertFullUrl(t, baseUrl, "nd", "PATCH", path, expected_nd)
		AssertFullUrl(t, baseUrl, "mso", "PATCH", path, expected)
	}
}

func TestMakeFullUrl_PatchExtraQuery(t *testing.T) {
	expected := "https://ndo.host.cisco/templates/123?extra=query&validate=false"
	expected_nd := "https://ndo.host.cisco/mso/templates/123?extra=query&validate=false"
	path := "templates/123?extra=query"
	for _, baseUrl := range TestBaseUrls {
		AssertFullUrl(t, baseUrl, "nd", "PATCH", path, expected_nd)
		AssertFullUrl(t, baseUrl, "mso", "PATCH", path, expected)
	}
}

// Test 1: ThreadSafeCache Basic Operations Tests
func TestThreadSafeCache_BasicOperations(t *testing.T) {
	cache := NewThreadSafeCache()

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

func TestThreadSafeCache_Delete(t *testing.T) {
	cache := NewThreadSafeCache()

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

// Test 2: ThreadSafeCache Pattern Deletion Tests
func TestThreadSafeCache_DeletePattern(t *testing.T) {
	cache := NewThreadSafeCache()

	// Set multiple items with schema pattern and other items
	cache.Set("schema_123", "data1")
	cache.Set("schema_456", "data2")
	cache.Set("other_789", "data3")
	cache.Set("schema_abc", "data4")

	cloneFunc := func(item interface{}) (interface{}, error) {
		return item, nil
	}

	// Verify all items exist before deletion
	_, found1, _ := cache.Get("schema_123", cloneFunc)
	_, found2, _ := cache.Get("schema_456", cloneFunc)
	_, found3, _ := cache.Get("other_789", cloneFunc)
	_, found4, _ := cache.Get("schema_abc", cloneFunc)

	if !found1 || !found2 || !found3 || !found4 {
		t.Error("All items should exist before pattern deletion")
	}

	// Delete schema pattern
	cache.DeletePattern("schema_")

	// Verify schema items are deleted but other items remain
	_, found1, _ = cache.Get("schema_123", cloneFunc)
	_, found2, _ = cache.Get("schema_456", cloneFunc)
	_, found3, _ = cache.Get("other_789", cloneFunc)
	_, found4, _ = cache.Get("schema_abc", cloneFunc)

	if found1 || found2 || found4 {
		t.Error("Schema items should be deleted")
	}
	if !found3 {
		t.Error("Other items should remain")
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
		Cache:        NewThreadSafeCache(),
		cacheEnabled: false,
	}

	// Should not panic when cache is disabled
	client.InvalidateSchemaCache("test-schema")
	client.InvalidateAllSchemaCache()

	// Should return zero stats when cache is disabled
	hits, misses, invalidations, hitRatio := client.GetCacheStats()
	if hits != 0 || misses != 0 || invalidations != 0 || hitRatio != 0 {
		t.Errorf("Cache stats should be zero when disabled, got hits=%d, misses=%d, invalidations=%d, hitRatio=%.1f",
			hits, misses, invalidations, hitRatio)
	}
}

func TestClient_CacheInvalidationMethods(t *testing.T) {
	client := &Client{
		Cache:        NewThreadSafeCache(),
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

	// Test bulk invalidation (should only invalidate schema items)
	client.InvalidateAllSchemaCache()
	_, found2, _ = client.Cache.Get("schema_456", cloneFunc)
	_, found3, _ := client.Cache.Get("other_789", cloneFunc)

	if found2 {
		t.Error("Schema 456 should be invalidated by bulk operation")
	}
	if !found3 {
		t.Error("Other items should remain after schema bulk operation")
	}
}
