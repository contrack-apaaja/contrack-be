# 🔐 Complete Authorization Matrix

## Overview
This document lists all the authorizations implemented in the Contrack backend role-based access control system.

## User Roles

### 1. REGULAR (Default Role)
- **Description**: Standard users with read-only access
- **Default Assignment**: All new users get this role automatically

### 2. LEGAL
- **Description**: Legal professionals with contract management privileges
- **Manual Assignment**: Must be assigned through database or admin interface

### 3. MANAGEMENT
- **Description**: Management personnel with full administrative privileges
- **Manual Assignment**: Must be assigned through database or admin interface

## Authorization Matrix

### 🔓 Public Endpoints (No Authentication Required)

| Endpoint | Method | Description | All Roles |
|----------|--------|-------------|-----------|
| `/api/hello` | GET | Health check | ✅ |
| `/api/auth/register` | POST | User registration | ✅ |
| `/api/auth/login` | POST | User login | ✅ |
| `/api/auth/refresh` | POST | Token refresh | ✅ |

### 🔒 Protected Endpoints (Authentication Required)

#### User Management
| Endpoint | Method | Description | REGULAR | LEGAL | MANAGEMENT |
|----------|--------|-------------|---------|-------|------------|
| `/api/profile` | GET | Get user profile | ✅ | ✅ | ✅ |
| `/api/users` | GET | List all users | ✅ | ✅ | ✅ |

#### Clause Templates - Read Operations
| Endpoint | Method | Description | REGULAR | LEGAL | MANAGEMENT |
|----------|--------|-------------|---------|-------|------------|
| `/api/clauses/` | GET | List clause templates | ✅ | ✅ | ✅ |
| `/api/clauses/:id` | GET | Get clause template by ID | ✅ | ✅ | ✅ |
| `/api/clauses/by-code/:code` | GET | Get clause template by code | ✅ | ✅ | ✅ |
| `/api/clauses/search` | GET | Search clause templates | ✅ | ✅ | ✅ |
| `/api/clauses/types` | GET | Get clause types | ✅ | ✅ | ✅ |

#### Clause Templates - Write Operations
| Endpoint | Method | Description | REGULAR | LEGAL | MANAGEMENT |
|----------|--------|-------------|---------|-------|------------|
| `/api/clauses/` | POST | Create clause template | ❌ | ✅ | ✅ |
| `/api/clauses/:id` | PUT | Update clause template | ❌ | ✅ | ✅ |
| `/api/clauses/:id` | DELETE | Delete clause template | ❌ | ✅ | ✅ |
| `/api/clauses/:id/toggle-status` | PATCH | Toggle clause status | ❌ | ✅ | ✅ |

#### Contracts - Read Operations
| Endpoint | Method | Description | REGULAR | LEGAL | MANAGEMENT |
|----------|--------|-------------|---------|-------|------------|
| `/api/contracts/stats` | GET | Get contract statistics | ✅ | ✅ | ✅ |
| `/api/contracts/` | GET | List contracts | ✅ | ✅ | ✅ |
| `/api/contracts/:id` | GET | Get contract by ID | ✅ | ✅ | ✅ |
| `/api/contracts/:id/status-history` | GET | Get contract status history | ✅ | ✅ | ✅ |

#### Contracts - Write Operations
| Endpoint | Method | Description | REGULAR | LEGAL | MANAGEMENT |
|----------|--------|-------------|---------|-------|------------|
| `/api/contracts/` | POST | Create contract | ❌ | ✅ | ✅ |
| `/api/contracts/:id` | PUT | Update contract | ❌ | ✅ | ✅ |
| `/api/contracts/:id` | DELETE | Delete contract | ❌ | ✅ | ✅ |
| `/api/contracts/:id/status` | POST | Change contract status | ❌ | ✅ | ✅ |

#### Contract Versioning
| Endpoint | Method | Description | REGULAR | LEGAL | MANAGEMENT |
|----------|--------|-------------|---------|-------|------------|
| `/api/contract-versions/:baseId` | GET | Get contract versions | ✅ | ✅ | ✅ |
| `/api/contract-versions/:baseId` | POST | Create contract version | ❌ | ✅ | ✅ |

#### Stakeholders
| Endpoint | Method | Description | REGULAR | LEGAL | MANAGEMENT |
|----------|--------|-------------|---------|-------|------------|
| `/api/stakeholders/` | POST | Create stakeholder | ✅ | ✅ | ✅ |
| `/api/stakeholders/` | GET | List stakeholders | ✅ | ✅ | ✅ |
| `/api/stakeholders/:id` | GET | Get stakeholder | ✅ | ✅ | ✅ |
| `/api/stakeholders/:id` | PUT | Update stakeholder | ✅ | ✅ | ✅ |
| `/api/stakeholders/:id` | DELETE | Delete stakeholder | ✅ | ✅ | ✅ |
| `/api/stakeholders/types` | GET | Get stakeholder types | ✅ | ✅ | ✅ |

#### AI Analysis
| Endpoint | Method | Description | REGULAR | LEGAL | MANAGEMENT |
|----------|--------|-------------|---------|-------|------------|
| `/api/ai/analyze` | POST | Analyze clause risk | ✅ | ✅ | ✅ |
| `/api/ai/analysis/:id` | GET | Get analysis by ID | ✅ | ✅ | ✅ |
| `/api/ai/analysis/clause/:clause_id` | GET | Get analysis by clause ID | ✅ | ✅ | ✅ |
| `/api/ai/analyses` | GET | List analyses | ✅ | ✅ | ✅ |
| `/api/ai/analysis/:id` | DELETE | Delete analysis | ✅ | ✅ | ✅ |
| `/api/ai/stats` | GET | Get analysis statistics | ✅ | ✅ | ✅ |

#### Dashboard
| Endpoint | Method | Description | REGULAR | LEGAL | MANAGEMENT |
|----------|--------|-------------|---------|-------|------------|
| `/api/dashboard/status-counts` | GET | Get status counts | ✅ | ✅ | ✅ |
| `/api/dashboard/contracts` | GET | Get contract list | ✅ | ✅ | ✅ |

## Permission Summary

### REGULAR Users Can:
- ✅ **Read** all contracts and clauses
- ✅ **View** dashboard and statistics
- ✅ **Access** AI analysis features
- ✅ **Manage** stakeholders
- ❌ **Cannot** create, update, or delete contracts
- ❌ **Cannot** create, update, or delete clause templates
- ❌ **Cannot** change contract status
- ❌ **Cannot** create contract versions

### LEGAL Users Can:
- ✅ **All REGULAR permissions**
- ✅ **Create, update, and delete** contracts
- ✅ **Create, update, and delete** clause templates
- ✅ **Change** contract status
- ✅ **Create** contract versions

### MANAGEMENT Users Can:
- ✅ **All LEGAL permissions**
- ✅ **Full system access**
- ✅ **Administrative functions**

## Implementation Details

### Middleware Chain
1. **CORS Middleware** - Handles cross-origin requests
2. **Authentication Middleware** - Validates JWT tokens and fetches user role
3. **Authorization Middleware** - Checks role-based permissions
4. **Business Logic** - Controller handles the actual request

### Error Responses

#### 401 Unauthorized
```json
{
  "status": "error",
  "message": "Invalid or expired token",
  "error": {
    "code": "UNAUTHORIZED"
  }
}
```

#### 403 Forbidden
```json
{
  "status": "error",
  "message": "Only LEGAL and MANAGEMENT users can update contracts and clauses",
  "error": {
    "code": "FORBIDDEN"
  }
}
```

### Database Constraints
- Role column has CHECK constraint: `role IN ('REGULAR', 'LEGAL', 'MANAGEMENT')`
- Default role for new users: `REGULAR`
- Existing users automatically migrated to `REGULAR` role

## Testing Authorization

### Test REGULAR User Permissions
```bash
# Register and login as REGULAR user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "regular@example.com", "password": "password123"}'

# Try to create clause (should fail with 403)
curl -X POST http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"clause_code": "TEST_001", "title": "Test", "type": "Test", "content": "Test"}'

# Try to read clauses (should work)
curl -X GET http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Test LEGAL/MANAGEMENT User Permissions
```bash
# Update user role in database to LEGAL or MANAGEMENT
UPDATE users SET role = 'LEGAL' WHERE email = 'user@example.com';

# Now user can create/update contracts and clauses
curl -X POST http://localhost:8080/api/clauses \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"clause_code": "TEST_001", "title": "Test", "type": "Test", "content": "Test"}'
```

## Security Features

1. **JWT Token Validation** - All protected endpoints require valid tokens
2. **Role-Based Access Control** - Permissions checked at middleware level
3. **Database Constraints** - Role validation at database level
4. **Default Permissions** - New users start with minimal permissions
5. **Audit Trail** - All actions are logged with user context

## Future Enhancements

1. **Role Assignment API** - Admin interface to assign roles
2. **Permission Granularity** - More specific permissions per role
3. **Role Hierarchy** - Inheritance-based permissions
4. **Audit Logging** - Track role-based actions
5. **Dynamic Permissions** - Runtime permission configuration

---

**Total Endpoints Protected**: 25+ endpoints with role-based authorization
**Security Level**: Enterprise-grade with proper error handling
**Backward Compatibility**: ✅ Maintained for existing users
