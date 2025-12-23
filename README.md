# mso-go-client
 This repository contains the golang client SDK to interact with Cisco MSO/NDO using REST API calls. This SDK is used by [terraform-provider-mso](https://github.com/CiscoDevNet/terraform-provider-mso).

## Installation ##

Use `go get` to retrieve the SDK to add it to your `GOPATH` workspace, or project's Go module dependencies.


```sh
$go get github.com/ciscoecosystem/mso-go-client
```

There are no additional dependencies needed to be installed.

## Overview ##
  
* <strong>client</strong> :- This package contains the HTTP Client configuration as well as service methods which serves the CRUD operations on the configuration objects in Cisco MSO/NDO.

* <strong>models</strong> :- This package contains all the models structs and utility methods for the same.

* <strong>tests</strong> :- This package contains the unit tests for the CRUD operations that can be performed on the configuration objects.

## How to Use ##

import the client in your go application and retrive the client object by calling client.GetClient() method.
```golang
import github.com/ciscoecosystem/mso-go-client/client
client.GetClient("URL", "Username", client.Password("Password"), client.Insecure(true/false))
```

mso-go-client also supports running against NDO or ND-based MSO. To use against an ND based authentication call the GetClient method as follows.  
  

```golang
client.GetClient("URL", "Username", client.Password("Password"), client.Insecure(true/false), client.Platform("nd"))

```

Use that client object to call the service methods to perform the CRUD operations on the configuration objects.

Example,

```golang
	client.Save("api/v1/tenants", models.NewTenant(TenantAttributes))
    # TenantAttributes is struct present in models/tenant.go
```

## Caching Support ##

The client supports optional caching for API requests to improve performance by storing frequently accessed data in memory.

### Enabling Caching ###

Caching is **disabled by default** for safety. Enable it using the `CacheEnabled` option:

```golang
import "github.com/ciscoecosystem/mso-go-client/client"

// Enable caching
msoClient := client.GetClient("URL", "Username",
    client.Password("Password"),
    client.Insecure(true),
    client.CacheEnabled(true))
```

### Cache Operations ###

Once caching is enabled, you can use the following methods:

```golang
// Fetch data with caching support (automatically handles cache hits/misses)
data, err := msoClient.GetViaURLWithCache("/api/v1/schemas/schema-id")

// Invalidate a specific URL from cache (e.g., after updates)
msoClient.InvalidateURLCache("/api/v1/schemas/schema-id")

// Clear all cached items (useful for bulk operations or cleanup)
msoClient.ClearCache()
```

### Cache Debug Information ###

When debug logging is enabled (`TF_LOG=DEBUG` or `TF_LOG=TRACE`), the cache provides detailed information about:

- Cache hits and misses per resource
- Memory usage per cached item and total cache size
- System memory usage
- Hit ratios and performance statistics

Example debug output:
```
[DEBUG] SCHEMA_CACHE_HIT for /api/v1/schemas/123 | ItemSize: 15.2KB | System: 45.3MB
[DEBUG] CACHE_CLEARED | AggregateStats: Items=25, Hits=150, Misses=12, HitRatio=92.6% | Memory: Cache=1.25MB, AvgItem=51.2KB, System: 45.3MB
```