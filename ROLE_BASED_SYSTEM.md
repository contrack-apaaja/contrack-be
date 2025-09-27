# 🔐 Role-Based Access Control (RBAC) System

## Overview

This document describes the comprehensive role-based access control system implemented in the Contrack backend. The system provides three distinct user roles with different permission levels for managing contracts and clauses.

## User Roles

### 1. REGULAR (Default Role)
- **Description**: Standard users with basic access
- **Permissions**: 
  - ✅ Read access to all contracts and clauses
  - ✅ View dashboard and statistics
  - ✅ Access AI analysis features
  - ❌ Cannot create, update, or delete contracts/clauses

### 2. LEGAL
- **Description**: Legal professionals with contract management privileges
- **Permissions**:
  - ✅ All REGULAR permissions
  - ✅ Create, update, and delete contracts
  - ✅ Create, update, and delete clause templates
  - ✅ Change contract status
  - ✅ Create contract versions

### 3. MANAGEMENT
- **Description**: Management personnel with full administrative privileges
- **Permissions**:
  - ✅ All LEGAL permissions
  - ✅ Full system access
  - ✅ Administrative functions

## Implementation Details

### Database Schema

```sql
-- Users table with role column
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'REGULAR' 
        CHECK (role IN ('REGULAR', 'LEGAL', 'MANAGEMENT')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Role Enum

```go
type UserRole string

const (
    RoleRegular    UserRole = "REGULAR"
    RoleLegal      UserRole = "LEGAL"
    RoleManagement UserRole = "MANAGEMENT"
)

// Methods
func (r UserRole) IsValid() bool
func (r UserRole) CanUpdateContracts() bool
func (r UserRole) String() string
```

### Authentication Flow

1. **Registration**: New users automatically get `REGULAR` role
2. **Login Response**: Includes user role information
3. **Token Validation**: JWT tokens include user context with role
4. **Authorization**: Middleware checks role permissions for protected endpoints

## API Endpoints by Role

### Public Endpoints (No Authentication)
- `GET /api/hello` - Health check
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Token refresh

### Protected Endpoints (All Authenticated Users)
- `GET /api/profile` - User profile
- `GET /api/users` - List users
- `GET /api/clauses/` - List clause templates
- `GET /api/clauses/:id` - Get clause template
- `GET /api/clauses/search` - Search clause templates
- `GET /api/clauses/types` - Get clause types
- `GET /api/contracts/` - List contracts
- `GET /api/contracts/:id` - Get contract
- `GET /api/contracts/stats` - Contract statistics
- `GET /api/contracts/:id/status-history` - Contract status history
- `GET /api/contract-versions/:baseId` - Get contract versions
- `GET /api/ai/analysis/:id` - Get AI analysis
- `GET /api/ai/analyses` - List AI analyses
- `GET /api/ai/stats` - AI analysis statistics
- `GET /api/dashboard/status-counts` - Dashboard status counts
- `GET /api/dashboard/contracts` - Dashboard contract list

### LEGAL & MANAGEMENT Only Endpoints
- `POST /api/clauses/` - Create clause template
- `PUT /api/clauses/:id` - Update clause template
- `DELETE /api/clauses/:id` - Delete clause template
- `PATCH /api/clauses/:id/toggle-status` - Toggle clause status
- `POST /api/contracts/` - Create contract
- `PUT /api/contracts/:id` - Update contract
- `DELETE /api/contracts/:id` - Delete contract
- `POST /api/contracts/:id/status` - Change contract status
- `POST /api/contract-versions/:baseId` - Create contract version

## Middleware Implementation

### Authentication Middleware
```go
func AuthMiddleware(jwtService *jwt.Service) gin.HandlerFunc
```
- Validates JWT tokens
- Extracts user information
- Fetches user role from database
- Stores user context for downstream handlers

### Role-Based Authorization Middleware
```go
// Require specific roles
func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc

// Require contract update permissions (LEGAL or MANAGEMENT)
func RequireContractUpdatePermission() gin.HandlerFunc
```

## Database Migration

The system includes automatic migration for existing users:

```go
func MigrateExistingUsersToRegularRole() error
```

This function:
1. Adds the `role` column if it doesn't exist
2. Sets default value to `REGULAR`
3. Updates existing users to have `REGULAR` role
4. Adds appropriate constraints

## Usage Examples

### Frontend Integration

```javascript
// Login response now includes role
const loginResponse = await api.post('/auth/login', {
  email: 'user@example.com',
  password: 'password123'
});

console.log(loginResponse.data.user.role); // "REGULAR", "LEGAL", or "MANAGEMENT"

// Check permissions before making requests
if (user.role === 'LEGAL' || user.role === 'MANAGEMENT') {
  // User can create/update contracts and clauses
  await api.post('/clauses', clauseData);
} else {
  // Show read-only interface
  console.log('User can only view contracts and clauses');
}
```

### Backend Permission Checking

```go
// In a controller
func (c *Controller) SomeAction(ctx *gin.Context) {
    userRole, exists := middleware.GetUserRole(ctx)
    if !exists {
        utils.UnauthorizedResponse(ctx, "User role not found")
        return
    }
    
    if !userRole.CanUpdateContracts() {
        utils.ForbiddenResponse(ctx, "Insufficient permissions")
        return
    }
    
    // Proceed with action
}
```

## Error Responses

### 401 Unauthorized
```json
{
  "status": "error",
  "message": "Invalid or expired token",
  "error": {
    "code": "UNAUTHORIZED"
  }
}
```

### 403 Forbidden
```json
{
  "status": "error",
  "message": "Only LEGAL and MANAGEMENT users can update contracts and clauses",
  "error": {
    "code": "FORBIDDEN"
  }
}
```

## Security Considerations

1. **Role Validation**: All roles are validated against the enum
2. **Database Constraints**: Role column has CHECK constraints
3. **JWT Security**: Tokens include user context but role is fetched from database
4. **Middleware Chain**: Authentication → Authorization → Business Logic
5. **Default Permissions**: New users start with minimal permissions (REGULAR)

## Testing the System

### Test User Creation
```bash
# Register a new user (gets REGULAR role by default)
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'

# Login to get token and role
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### Test Permission Enforcement
```bash
# Try to create a clause as REGULAR user (should fail)
curl -X POST http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"clause_code": "TEST_001", "title": "Test Clause", "type": "Test", "content": "Test content"}'

# Response: 403 Forbidden
```

## Future Enhancements

1. **Role Assignment**: Admin interface to assign roles
2. **Permission Granularity**: More specific permissions per role
3. **Role Hierarchy**: Inheritance-based permissions
4. **Audit Logging**: Track role-based actions
5. **Dynamic Permissions**: Runtime permission configuration

## Migration Guide

### For Existing Users
1. Run the server - existing users will automatically get `REGULAR` role
2. No manual intervention required
3. Database migration is automatic and safe

### For New Deployments
1. Role system is enabled by default
2. All new users get `REGULAR` role
3. Update user roles as needed through database or admin interface

---

**Note**: This role-based system provides a solid foundation for access control while maintaining backward compatibility with existing functionality.
