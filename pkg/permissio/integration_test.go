//go:build integration

package permissio_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/enforcement"
	"github.com/permissio/permissio-go/pkg/models"
	permissio "github.com/permissio/permissio-go/pkg/permissio"
)

const (
	envScopedKey = "permis_key_d39064912cd9d1f0052a98430e3eb7d689a350d84f2d0a018843541b5da3e5ef"
	apiURL       = "http://localhost:3001"
)

// newClient creates a new test client with the env-scoped API key.
func newClient(t *testing.T) *permissio.Client {
	t.Helper()
	cfg := config.NewConfigBuilder(envScopedKey).
		WithApiUrl(apiURL).
		WithRetryAttempts(0).
		Build()
	return permissio.New(cfg)
}

// ts returns a timestamp-based unique suffix.
func ts() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

// TestIntegration runs all 22 integration scenarios as subtests, using a single
// shared timestamp so all created resources share a consistent prefix that is
// cleaned up in TestMain / t.Cleanup.
func TestIntegration(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	// -----------------------------------------------------------------------
	// 1. API key scope auto-fetch
	// -----------------------------------------------------------------------
	t.Run("01_scope_auto_fetch", func(t *testing.T) {
		err := client.Init(ctx)
		if err != nil {
			t.Fatalf("Init() failed: %v", err)
		}
		projID, envID, err := client.GetScope(ctx)
		if err != nil {
			t.Fatalf("GetScope() failed: %v", err)
		}
		if projID == "" {
			t.Error("project_id is empty")
		}
		if envID == "" {
			t.Error("environment_id is empty")
		}
		t.Logf("scope: project_id=%s environment_id=%s", projID, envID)
	})

	suffix := ts()
	userKey := "test-user-" + suffix
	tenantKey := "test-tenant-" + suffix
	resourceKey := "test-resource-" + suffix
	roleKey := "test-role-" + suffix

	// cleanup runs after the test function returns (i.e. after all subtests finish).
	t.Cleanup(func() {
		cleanCtx := context.Background()
		// ignore errors during cleanup
		_ = client.Api.RoleAssignments.Unassign(cleanCtx, userKey, roleKey, tenantKey)
		_ = client.Api.Roles.Delete(cleanCtx, roleKey)
		_ = client.Api.Resources.Delete(cleanCtx, resourceKey)
		_ = client.Api.Tenants.Delete(cleanCtx, tenantKey)
		_ = client.Api.Users.Delete(cleanCtx, userKey)
	})

	// -----------------------------------------------------------------------
	// 2. Users CRUD
	// -----------------------------------------------------------------------
	t.Run("02_users_crud", func(t *testing.T) {
		// Create
		created, err := client.Api.Users.Create(ctx, &models.UserCreate{
			Key:       userKey,
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		})
		if err != nil {
			t.Fatalf("Users.Create() failed: %v", err)
		}
		if created.Key != userKey {
			t.Errorf("expected key %q, got %q", userKey, created.Key)
		}

		// List
		list, err := client.Api.Users.List(ctx, nil)
		if err != nil {
			t.Fatalf("Users.List() failed: %v", err)
		}
		found := false
		for _, u := range list.Data {
			if u.Key == userKey {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created user %q not found in list", userKey)
		}

		// Get
		got, err := client.Api.Users.Get(ctx, userKey)
		if err != nil {
			t.Fatalf("Users.Get() failed: %v", err)
		}
		if got.Key != userKey {
			t.Errorf("expected key %q, got %q", userKey, got.Key)
		}

		// Sync (upsert via SyncUser)
		synced, err := client.Api.Users.SyncUser(ctx, models.UserCreate{
			Key:       userKey,
			Email:     "updated@example.com",
			FirstName: "Updated",
			LastName:  "User",
		})
		if err != nil {
			t.Fatalf("Users.SyncUser() failed: %v", err)
		}
		if synced.Key != userKey {
			t.Errorf("expected key %q after sync, got %q", userKey, synced.Key)
		}
	})

	// -----------------------------------------------------------------------
	// 3. Tenants CRUD
	// -----------------------------------------------------------------------
	t.Run("03_tenants_crud", func(t *testing.T) {
		// Create
		created, err := client.Api.Tenants.Create(ctx, &models.TenantCreate{
			Key:  tenantKey,
			Name: "Test Tenant " + suffix,
		})
		if err != nil {
			t.Fatalf("Tenants.Create() failed: %v", err)
		}
		if created.Key != tenantKey {
			t.Errorf("expected key %q, got %q", tenantKey, created.Key)
		}

		// List
		list, err := client.Api.Tenants.List(ctx, nil)
		if err != nil {
			t.Fatalf("Tenants.List() failed: %v", err)
		}
		found := false
		for _, ten := range list.Data {
			if ten.Key == tenantKey {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created tenant %q not found in list", tenantKey)
		}

		// Get
		got, err := client.Api.Tenants.Get(ctx, tenantKey)
		if err != nil {
			t.Fatalf("Tenants.Get() failed: %v", err)
		}
		if got.Key != tenantKey {
			t.Errorf("expected key %q, got %q", tenantKey, got.Key)
		}
	})

	// -----------------------------------------------------------------------
	// 4. Resources CRUD
	// -----------------------------------------------------------------------
	t.Run("04_resources_crud", func(t *testing.T) {
		// Create
		created, err := client.Api.Resources.Create(ctx, &models.ResourceCreate{
			Key:     resourceKey,
			Name:    "Test Resource " + suffix,
			Actions: []string{"read", "write", "delete"},
		})
		if err != nil {
			t.Fatalf("Resources.Create() failed: %v", err)
		}
		if created.Key != resourceKey {
			t.Errorf("expected key %q, got %q", resourceKey, created.Key)
		}

		// List
		list, err := client.Api.Resources.List(ctx, nil)
		if err != nil {
			t.Fatalf("Resources.List() failed: %v", err)
		}
		found := false
		for _, res := range list.Data {
			if res.Key == resourceKey {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created resource %q not found in list", resourceKey)
		}
	})

	// -----------------------------------------------------------------------
	// 5. Roles CRUD
	// -----------------------------------------------------------------------
	permission := resourceKey + ":read"

	t.Run("05_roles_crud", func(t *testing.T) {
		// Create role with permissions
		created, err := client.Api.Roles.Create(ctx, &models.RoleCreate{
			Key:         roleKey,
			Name:        "Test Role " + suffix,
			Permissions: []string{permission},
		})
		if err != nil {
			t.Fatalf("Roles.Create() failed: %v", err)
		}
		if created.Key != roleKey {
			t.Errorf("expected key %q, got %q", roleKey, created.Key)
		}

		// List
		list, err := client.Api.Roles.List(ctx, nil)
		if err != nil {
			t.Fatalf("Roles.List() failed: %v", err)
		}
		found := false
		for _, r := range list.Data {
			if r.Key == roleKey {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created role %q not found in list", roleKey)
		}

		// Get
		got, err := client.Api.Roles.Get(ctx, roleKey)
		if err != nil {
			t.Fatalf("Roles.Get() failed: %v", err)
		}
		if got.Key != roleKey {
			t.Errorf("expected key %q, got %q", roleKey, got.Key)
		}
		hasPermission := false
		for _, p := range got.Permissions {
			if p == permission {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			t.Errorf("role %q does not have permission %q; got %v", roleKey, permission, got.Permissions)
		}
	})

	// -----------------------------------------------------------------------
	// 6. Role Assignments — assign and list
	// -----------------------------------------------------------------------
	t.Run("06_role_assignments_assign_list", func(t *testing.T) {
		// Assign
		assigned, err := client.Api.RoleAssignments.Assign(ctx, &models.RoleAssignmentCreate{
			User:   userKey,
			Role:   roleKey,
			Tenant: tenantKey,
		})
		if err != nil {
			t.Fatalf("RoleAssignments.Assign() failed: %v", err)
		}
		if assigned.User != userKey {
			t.Errorf("expected user %q, got %q", userKey, assigned.User)
		}
		if assigned.Role != roleKey {
			t.Errorf("expected role %q, got %q", roleKey, assigned.Role)
		}

		// List
		list, err := client.Api.RoleAssignments.List(ctx, &models.RoleAssignmentListParams{
			User: userKey,
		})
		if err != nil {
			t.Fatalf("RoleAssignments.List() failed: %v", err)
		}
		found := false
		for _, ra := range list {
			if ra.User == userKey && ra.Role == roleKey {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("role assignment user=%q role=%q not found in list", userKey, roleKey)
		}
	})

	// -----------------------------------------------------------------------
	// 7. check() — allowed case
	// -----------------------------------------------------------------------
	t.Run("07_check_allowed", func(t *testing.T) {
		allowed, err := client.Check(
			enforcement.User{Key: userKey},
			enforcement.Action("read"),
			enforcement.Resource{Type: resourceKey, Tenant: tenantKey},
		)
		if err != nil {
			t.Fatalf("Check() failed: %v", err)
		}
		if !allowed {
			t.Errorf("expected allowed=true, got false")
		}
	})

	// -----------------------------------------------------------------------
	// 8. check() — denied case (missing action)
	// -----------------------------------------------------------------------
	t.Run("08_check_denied", func(t *testing.T) {
		allowed, err := client.Check(
			enforcement.User{Key: userKey},
			enforcement.Action("delete"),
			enforcement.Resource{Type: resourceKey, Tenant: tenantKey},
		)
		if err != nil {
			t.Fatalf("Check() failed: %v", err)
		}
		if allowed {
			t.Errorf("expected allowed=false, got true")
		}
	})

	// -----------------------------------------------------------------------
	// 9. BulkCheck — 3 checks at once
	// -----------------------------------------------------------------------
	t.Run("09_bulk_check", func(t *testing.T) {
		checks := []models.CheckRequest{
			{User: userKey, Action: "read", Resource: resourceKey, Tenant: tenantKey},
			{User: userKey, Action: "write", Resource: resourceKey, Tenant: tenantKey},
			{User: userKey, Action: "delete", Resource: resourceKey, Tenant: tenantKey},
		}
		resp, err := client.BulkCheck(ctx, checks)
		if err != nil {
			t.Fatalf("BulkCheck() failed: %v", err)
		}
		if len(resp.Results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(resp.Results))
		}
		// read → allowed (role has read permission)
		if !resp.Results[0].Response.Allowed {
			t.Errorf("bulk check[0] read: expected allowed=true, got false")
		}
		// write → denied (role only has read)
		if resp.Results[1].Response.Allowed {
			t.Errorf("bulk check[1] write: expected allowed=false, got true")
		}
		// delete → denied
		if resp.Results[2].Response.Allowed {
			t.Errorf("bulk check[2] delete: expected allowed=false, got true")
		}
	})

	// -----------------------------------------------------------------------
	// 10. GetPermissions — returns roles + permissions
	// -----------------------------------------------------------------------
	t.Run("10_get_permissions", func(t *testing.T) {
		resp, err := client.GetPermissions(ctx, models.GetPermissionsRequest{
			User:   userKey,
			Tenant: tenantKey,
		})
		if err != nil {
			t.Fatalf("GetPermissions() failed: %v", err)
		}
		if len(resp.Roles) == 0 {
			t.Error("expected at least one role in GetPermissions response")
		}
		if len(resp.Permissions) == 0 {
			t.Error("expected at least one permission in GetPermissions response")
		}
		foundRole := false
		for _, r := range resp.Roles {
			if r == roleKey {
				foundRole = true
				break
			}
		}
		if !foundRole {
			t.Errorf("role %q not found in GetPermissions response: %v", roleKey, resp.Roles)
		}
		foundPerm := false
		for _, p := range resp.Permissions {
			if p == permission {
				foundPerm = true
				break
			}
		}
		if !foundPerm {
			t.Errorf("permission %q not found in GetPermissions response: %v", permission, resp.Permissions)
		}
	})

	// -----------------------------------------------------------------------
	// 11. SyncUser — convenience method on Client
	// -----------------------------------------------------------------------
	syncUserKey := "test-sync-user-" + suffix
	t.Cleanup(func() {
		_ = client.Api.RoleAssignments.Unassign(context.Background(), syncUserKey, roleKey, tenantKey)
		_ = client.Api.Users.Delete(context.Background(), syncUserKey)
	})

	t.Run("11_sync_user", func(t *testing.T) {
		synced, err := client.SyncUser(ctx, models.UserCreate{
			Key:       syncUserKey,
			Email:     "sync@example.com",
			FirstName: "Sync",
			LastName:  "User",
		}, []models.RoleAssignmentCreate{
			{Role: roleKey, Tenant: tenantKey},
		})
		if err != nil {
			t.Fatalf("SyncUser() failed: %v", err)
		}
		if synced.Key != syncUserKey {
			t.Errorf("expected key %q, got %q", syncUserKey, synced.Key)
		}

		// Verify role was assigned
		list, err := client.Api.RoleAssignments.List(ctx, &models.RoleAssignmentListParams{
			User: syncUserKey,
		})
		if err != nil {
			t.Fatalf("RoleAssignments.List() for sync user failed: %v", err)
		}
		found := false
		for _, ra := range list {
			if ra.User == syncUserKey && ra.Role == roleKey {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("sync user %q does not have role %q assigned", syncUserKey, roleKey)
		}
	})

	// -----------------------------------------------------------------------
	// 12. Role Assignment unassign — verify removed
	// -----------------------------------------------------------------------
	t.Run("12_role_assignment_unassign", func(t *testing.T) {
		err := client.Api.RoleAssignments.Unassign(ctx, userKey, roleKey, tenantKey)
		if err != nil {
			t.Fatalf("RoleAssignments.Unassign() failed: %v", err)
		}

		// Verify removed
		list, err := client.Api.RoleAssignments.List(ctx, &models.RoleAssignmentListParams{
			User: userKey,
			Role: roleKey,
		})
		if err != nil {
			t.Fatalf("RoleAssignments.List() after unassign failed: %v", err)
		}
		for _, ra := range list {
			if ra.User == userKey && ra.Role == roleKey {
				t.Errorf("role assignment user=%q role=%q still exists after unassign", userKey, roleKey)
			}
		}
	})
}
