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

The client supports optional caching for schema requests to improve performance in scenarios with frequent schema lookups.

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
// Fetch schema with caching support (automatically handles cache hits/misses)
schema, err := msoClient.GetSchemaWithCache("schema-id")

// Invalidate a specific schema from cache (e.g., after updates)
msoClient.InvalidateSchemaCache("schema-id")

// Clear all cached schemas (useful for bulk operations or cleanup)
msoClient.InvalidateAllSchemaCache()

// Get cache statistics for monitoring
hits, misses, invalidations, hitRatio := msoClient.GetCacheStats()

// Log current cache statistics
msoClient.LogCacheStats()
```