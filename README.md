# Permissio.io Go SDK

The official Go SDK for [Permissio.io](https://permissio.io) — a powerful authorization service for managing roles, permissions, and access control.

## Requirements

- Go 1.23+

## Installation

```bash
go get github.com/permissio/permissio-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/enforcement"
	"github.com/permissio/permissio-go/pkg/permissio"
)

func main() {
	cfg := config.NewConfigBuilder("permis_key_your_api_key_here").
		WithProjectID("your-project-id").
		WithEnvironmentID("your-environment-id").
		Build()

	client := permissio.New(cfg)

	user := enforcement.UserBuilder("user@example.com").Build()
	resource := enforcement.ResourceBuilder("document").Build()

	allowed, err := client.Check(user, enforcement.Action("read"), resource)
	if err != nil {
		log.Fatal(err)
	}

	if allowed {
		fmt.Println("Access granted!")
	} else {
		fmt.Println("Access denied!")
	}
}
```

## Auto-Scope Detection

If you omit `WithProjectID` and `WithEnvironmentID`, the SDK fetches them automatically from your API key on the first API call:

```go
cfg := config.NewConfigBuilder("permis_key_your_api_key_here").Build()
client := permissio.New(cfg)
// Project and environment IDs are fetched from /v1/api-key/scope on first use
```

To trigger scope initialization eagerly (e.g. at startup to catch config errors early):

```go
ctx := context.Background()
if err := client.Init(ctx); err != nil {
	log.Fatalf("Failed to initialize client: %v", err)
}
```

## ABAC (Attribute-Based Access Control)

```go
import "github.com/permissio/permissio-go/pkg/enforcement"

// User with attributes
user := enforcement.UserBuilder("user@example.com").
	WithAttribute("department", "engineering").
	WithAttribute("level", 5).
	Build()

// Resource with attributes and tenant
resource := enforcement.ResourceBuilder("document").
	WithKey("doc-123").
	WithTenant("acme-corp").
	WithAttribute("classification", "confidential").
	Build()

// Optional context
ctx := enforcement.ContextBuilder().
	With("ip_address", "192.168.1.1").
	With("request_time", "2026-03-15T12:00:00Z").
	Build()

allowed, err := client.Check(user, enforcement.Action("read"), resource)
```

### Enforcement builders

| Function | Description |
|----------|-------------|
| `enforcement.UserBuilder(key)` | Fluent builder for `User`; supports `.WithAttribute()`, `.WithAttributes()` |
| `enforcement.ResourceBuilder(type)` | Fluent builder for `Resource`; supports `.WithKey()`, `.WithTenant()`, `.WithAttribute()`, `.WithAttributes()` |
| `enforcement.ContextBuilder()` | Fluent builder for `Context`; supports `.With()`, `.WithData()` |

## API Management

All API operations require a `context.Context` as the first argument.

### Users

```go
ctx := context.Background()

// List users
users, err := client.Api.Users.List(ctx, nil)

// Get a user
user, err := client.Api.Users.Get(ctx, "user@example.com")

// Create a user
user, err := client.Api.Users.Create(ctx, &models.UserCreate{
	Key:       "user@example.com",
	Email:     "user@example.com",
	FirstName: "Jane",
	LastName:  "Doe",
})

// Sync user (upsert)
user, err := client.Api.Users.SyncUser(ctx, models.NewUserCreate("user@example.com").
	SetEmail("user@example.com").
	SetFirstName("Jane"))

// Top-level convenience: sync user and assign roles in one call
user, err := client.SyncUser(ctx, models.UserCreate{Key: "user@example.com"},
	[]models.RoleAssignmentCreate{
		{Role: "editor", Tenant: "acme-corp"},
	})

// Assign / unassign a role
_, err = client.Api.Users.AssignRole(ctx, "user@example.com", "editor", "acme-corp")
err = client.Api.Users.UnassignRole(ctx, "user@example.com", "editor", "acme-corp")

// Delete a user
err = client.Api.Users.Delete(ctx, "user@example.com")
```

### Tenants

```go
// List / get
tenants, err := client.Api.Tenants.List(ctx, nil)
tenant, err  := client.Api.Tenants.Get(ctx, "acme-corp")

// Create / update / delete
tenant, err = client.Api.Tenants.Create(ctx, &models.TenantCreate{
	Key:  "acme-corp",
	Name: "Acme Corporation",
})
tenant, err = client.Api.Tenants.Update(ctx, "acme-corp", &models.TenantUpdate{Name: "ACME Corp"})
err         = client.Api.Tenants.Delete(ctx, "acme-corp")

// Sync (upsert)
tenant, err = client.Api.Tenants.Sync(ctx, &models.TenantCreate{Key: "acme-corp", Name: "ACME"})
```

### Roles

```go
// List / get
roles, err := client.Api.Roles.List(ctx, nil)
role, err  := client.Api.Roles.Get(ctx, "editor")

// Create / update / delete
role, err = client.Api.Roles.Create(ctx, &models.RoleCreate{
	Key:         "editor",
	Name:        "Editor",
	Permissions: []string{"document:read", "document:write"},
})
role, err = client.Api.Roles.Update(ctx, "editor", &models.RoleUpdate{Name: "Content Editor"})
err       = client.Api.Roles.Delete(ctx, "editor")

// Sync (upsert)
role, err = client.Api.Roles.Sync(ctx, &models.RoleCreate{Key: "editor", Name: "Editor"})

// Permission management
err  = client.Api.Roles.AddPermission(ctx, "editor", "document:delete")
err  = client.Api.Roles.RemovePermission(ctx, "editor", "document:delete")
perms, err := client.Api.Roles.GetPermissions(ctx, "editor")

// Role inheritance (extends)
err  = client.Api.Roles.AddExtends(ctx, "editor", "viewer")
err  = client.Api.Roles.RemoveExtends(ctx, "editor", "viewer")
exts, err := client.Api.Roles.GetExtends(ctx, "editor")
```

### Resources

```go
// List / get
resources, err := client.Api.Resources.List(ctx, nil)
resource, err  := client.Api.Resources.Get(ctx, "document")

// Create / update / delete
resource, err = client.Api.Resources.Create(ctx, &models.ResourceCreate{
	Key:  "document",
	Name: "Document",
	Actions: map[string]models.ResourceAction{
		"read":   {Name: "Read"},
		"write":  {Name: "Write"},
		"delete": {Name: "Delete"},
	},
})
resource, err = client.Api.Resources.Update(ctx, "document", &models.ResourceUpdate{Name: "Doc"})
err           = client.Api.Resources.Delete(ctx, "document")

// Sync (upsert)
resource, err = client.Api.Resources.Sync(ctx, &models.ResourceCreate{Key: "document"})
```

### Role Assignments

```go
// Assign / unassign
assignment, err := client.Api.RoleAssignments.Assign(ctx, &models.RoleAssignmentCreate{
	User:   "user@example.com",
	Role:   "editor",
	Tenant: "acme-corp",
})

err = client.Api.RoleAssignments.Unassign(ctx, "user@example.com", "editor", "acme-corp")

// Unassign scoped to a resource instance
err = client.Api.RoleAssignments.UnassignWithResource(ctx,
	"user@example.com", "editor", "acme-corp", "document", "doc-123")

// List (with filters)
assignments, err := client.Api.RoleAssignments.List(ctx, &models.RoleAssignmentListParams{
	User:   "user@example.com",
	Tenant: "acme-corp",
})

// Convenience listing methods
assignments, err = client.Api.RoleAssignments.ListByUser(ctx, "user@example.com", nil)
assignments, err = client.Api.RoleAssignments.ListByTenant(ctx, "acme-corp", nil)
assignments, err = client.Api.RoleAssignments.ListByResource(ctx, "document", "doc-123", nil)

// Detailed listing
detailed, err := client.Api.RoleAssignments.ListDetailed(ctx, nil)

// Get by ID
single, err := client.Api.RoleAssignments.GetByID(ctx, "assignment-id")

// Utility helpers
roles, err  := client.Api.RoleAssignments.GetUserRoles(ctx, "user@example.com", nil)
users, err  := client.Api.RoleAssignments.GetRoleUsers(ctx, "editor", nil)
hasRole, err := client.Api.RoleAssignments.HasRole(ctx, "user@example.com", "editor", nil)

// Bulk operations
_, err = client.Api.RoleAssignments.BulkAssign(ctx, []models.RoleAssignmentCreate{
	{User: "alice@example.com", Role: "editor", Tenant: "acme-corp"},
	{User: "bob@example.com",   Role: "viewer", Tenant: "acme-corp"},
})
_, err = client.Api.RoleAssignments.BulkUnassign(ctx, []models.RoleAssignmentCreate{
	{User: "alice@example.com", Role: "editor", Tenant: "acme-corp"},
})
```

## Permission Checking

| Method | Signature | Description |
|--------|-----------|-------------|
| `Check` | `(user, action, resource) (bool, error)` | Simple permission check |
| `CheckWithContext` | `(ctx, user, action, resource) (bool, error)` | Check with explicit context |
| `CheckWithDetails` | `(ctx, user, action, resource) (*CheckResponse, error)` | Check with full response (reason, matched roles) |
| `CheckAndThrow` | `(ctx, user, action, resource) error` | Returns an error if access is denied |
| `BulkCheck` | `(ctx, []CheckRequest) (*BulkCheckResponse, error)` | Bulk permission checks |
| `GetPermissions` | `(ctx, GetPermissionsRequest) (*GetPermissionsResponse, error)` | All permissions for a user |

```go
// CheckAndThrow — useful in middleware
if err := client.CheckAndThrow(ctx, user, enforcement.Action("write"), resource); err != nil {
	// err is set when access is denied
	http.Error(w, "Forbidden", http.StatusForbidden)
	return
}
```

## Gin Middleware Example

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/enforcement"
	"github.com/permissio/permissio-go/pkg/permissio"
)

var permissioClient *permissio.Client

func AuthMiddleware(action, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userKey := c.GetHeader("X-User")
		if userKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-User header"})
			c.Abort()
			return
		}

		user := enforcement.UserBuilder(userKey).Build()
		resource := enforcement.ResourceBuilder(resourceType).Build()

		allowed, err := permissioClient.Check(user, enforcement.Action(action), resource)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			c.Abort()
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	cfg := config.NewConfigBuilder("permis_key_your_api_key_here").Build()
	permissioClient = permissio.New(cfg)

	router := gin.Default()
	router.POST("/posts", AuthMiddleware("create", "Post"), func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
	})
	router.Run(":8000")
}
```

## Configuration Options

| Builder method | Description | Default |
|----------------|-------------|---------|
| `WithApiUrl(url)` | Base API URL | `https://api.permissio.io` |
| `WithProjectID(id)` | Project ID | Auto-fetched |
| `WithEnvironmentID(id)` | Environment ID | Auto-fetched |
| `WithTimeout(duration)` | Request timeout | 30s |
| `WithDebug(enabled)` | Enable debug logging | `false` |
| `WithRetryAttempts(n)` | Retry attempts | 3 |
| `WithLogger(logger)` | Custom `*zap.Logger` | `nil` |
| `WithHTTPClient(client)` | Custom `*http.Client` | `http.DefaultClient` |

Use `Build()` for a plain config or `BuildWithValidation()` to return an error if required fields are missing.

## License

MIT License - see [LICENSE](LICENSE) for details.
