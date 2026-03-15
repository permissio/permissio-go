# Changelog

All notable changes to the Permissio.io Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

No unreleased changes at this time.

---

## [0.1.0-alpha.1] - 2025-03-15

### Added
- **`permissio.New(cfg)` client**: Main entry point; returns a `*Client` with full API access
- **`config.NewConfigBuilder(token)` (fluent)**: Full configuration support — `WithApiUrl()`, `WithProjectID()`, `WithEnvironmentID()`, `WithTimeout()`, `WithDebug()`, `WithRetryAttempts()`, `WithThrowOnError()`, `WithCustomHeader()`, `WithCustomHeaders()`, `WithLogger()`, `WithHTTPClient()`, `Build()`, `BuildWithValidation()`
- **Permission checking**:
  - `Check(user, action, resource)` — simple boolean permission check
  - `CheckWithContext(ctx, user, action, resource)` — context-aware boolean check
  - `CheckWithDetails(ctx, user, action, resource)` — full `CheckResponse` with client-side RBAC evaluation (role inheritance, wildcards)
  - `BulkCheck(ctx, checks)` — batch permission checks
  - `CheckAndThrow(ctx, user, action, resource)` — returns error on denial
  - `GetPermissions(ctx, request)` — retrieve all permissions for a user
- **Auto-scope detection**: Automatically fetches `ProjectID` and `EnvironmentID` from `/v1/api-key/scope` on first API call using double-checked locking
- **`Init(ctx)` method**: Explicitly trigger scope initialization before first check
- **`SyncUser(ctx, user, roles)` method**: Create or update a user with role assignments in one call
- **Enforcement builders** (`pkg/enforcement`):
  - `UserBuilder(key)` → `.WithAttribute()`, `.WithAttributes()`, `.Build()` → `User`
  - `ResourceBuilder(resourceType)` → `.WithKey()`, `.WithTenant()`, `.WithAttribute()`, `.WithAttributes()`, `.Build()` → `Resource`
  - `ContextBuilder()` → `.With()`, `.WithData()`, `.Build()` → `Context`
  - `Action` — typed string alias for action names
- **Users API** (`Api.Users`): `List()`, `Get()`, `Create()`, `Update()`, `Delete()`, `SyncUser()`, `AssignRole()`, `UnassignRole()`, `GetRoles()`, `AddTenant()`, `RemoveTenant()`, `GetTenants()`
- **Tenants API** (`Api.Tenants`): `List()`, `Get()`, `Create()`, `Update()`, `Delete()`, `Sync()`, `AddUser()`, `RemoveUser()`, `GetUsers()`
- **Roles API** (`Api.Roles`): `List()`, `Get()`, `Create()`, `Update()`, `Delete()`, `Sync()`, `GetPermissions()`, `AddPermission()`, `RemovePermission()`, `GetExtends()`, `AddExtends()`, `RemoveExtends()`
- **Resources API** (`Api.Resources`): `List()`, `Get()`, `Create()`, `Update()`, `Delete()`, `Sync()`, `GetActions()`, `AddAction()`, `RemoveAction()`, `CreateInstance()`, `GetInstance()`, `DeleteInstance()`
- **Role Assignments API** (`Api.RoleAssignments`): `List()`, `ListByUser()`, `ListByTenant()`, `ListByResource()`, `ListDetailed()`, `GetByID()`, `Assign()`, `Unassign()`, `UnassignWithResource()`, `BulkAssign()`, `BulkUnassign()`, `HasRole()`, `GetUserRoles()`, `GetRoleUsers()`
- **Structured logging**: `go.uber.org/zap` logger integration; custom logger injectable via `WithLogger()`
- **Gin middleware example**: Built-in middleware pattern for Gin-based applications
- **Go 1.23+ module**: Full Go 1.23 compatibility with `go.uber.org/zap v1.27.0` and `github.com/gin-gonic/gin v1.11.0`
- **Examples**: Basic usage and Gin integration examples
