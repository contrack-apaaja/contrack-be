# 🎭 Role System for Frontend Rendering

## Overview

The backend now provides user roles for frontend rendering purposes without any authorization restrictions. All authenticated users have full access to all endpoints, but the role information is still available for frontend UI customization.

## User Roles

### 1. REGULAR (Default Role)
- **Description**: Standard users
- **Frontend Usage**: Show basic UI elements
- **Backend Access**: Full access to all endpoints

### 2. LEGAL
- **Description**: Legal professionals
- **Frontend Usage**: Show legal-specific UI elements
- **Backend Access**: Full access to all endpoints

### 3. MANAGEMENT
- **Description**: Management personnel
- **Frontend Usage**: Show management-specific UI elements
- **Backend Access**: Full access to all endpoints

## Backend Changes

### ✅ What's Still Available:
- **User Registration**: New users get `REGULAR` role by default
- **Login Response**: Includes user role information
- **Database**: Role column and constraints maintained
- **Authentication**: JWT tokens still work
- **User Context**: Role information available in requests

### ❌ What's Removed:
- **Authorization Middleware**: No more permission checking
- **403 Forbidden Responses**: All endpoints accessible to all authenticated users
- **Role-Based Restrictions**: No backend restrictions based on roles

## Frontend Integration

### Login Response Structure
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "role": "REGULAR",
      "created_at": "2024-09-26T10:30:00Z",
      "updated_at": "2024-09-26T10:30:00Z"
    },
    "token": "jwt.token.here"
  }
}
```

### Frontend Role-Based Rendering

#### React/Next.js Example
```javascript
// Get user role from login response
const { user } = loginResponse.data;
const userRole = user.role;

// Conditional rendering based on role
const Dashboard = () => {
  return (
    <div>
      <h1>Dashboard</h1>
      
      {/* Show for all users */}
      <section>
        <h2>Contracts</h2>
        <ContractList />
      </section>
      
      {/* Show only for LEGAL and MANAGEMENT */}
      {(userRole === 'LEGAL' || userRole === 'MANAGEMENT') && (
        <section>
          <h2>Legal Tools</h2>
          <LegalTools />
        </section>
      )}
      
      {/* Show only for MANAGEMENT */}
      {userRole === 'MANAGEMENT' && (
        <section>
          <h2>Admin Panel</h2>
          <AdminPanel />
        </section>
      )}
    </div>
  );
};
```

#### Vue.js Example
```vue
<template>
  <div>
    <h1>Dashboard</h1>
    
    <!-- Show for all users -->
    <section>
      <h2>Contracts</h2>
      <ContractList />
    </section>
    
    <!-- Show only for LEGAL and MANAGEMENT -->
    <section v-if="userRole === 'LEGAL' || userRole === 'MANAGEMENT'">
      <h2>Legal Tools</h2>
      <LegalTools />
    </section>
    
    <!-- Show only for MANAGEMENT -->
    <section v-if="userRole === 'MANAGEMENT'">
      <h2>Admin Panel</h2>
      <AdminPanel />
    </section>
  </div>
</template>

<script>
export default {
  computed: {
    userRole() {
      return this.$store.state.user.role;
    }
  }
}
</script>
```

#### Angular Example
```typescript
// Component
export class DashboardComponent {
  userRole: string = 'REGULAR';
  
  constructor(private authService: AuthService) {
    this.userRole = this.authService.getUserRole();
  }
  
  isLegalOrManagement(): boolean {
    return this.userRole === 'LEGAL' || this.userRole === 'MANAGEMENT';
  }
  
  isManagement(): boolean {
    return this.userRole === 'MANAGEMENT';
  }
}
```

```html
<!-- Template -->
<div>
  <h1>Dashboard</h1>
  
  <!-- Show for all users -->
  <section>
    <h2>Contracts</h2>
    <app-contract-list></app-contract-list>
  </section>
  
  <!-- Show only for LEGAL and MANAGEMENT -->
  <section *ngIf="isLegalOrManagement()">
    <h2>Legal Tools</h2>
    <app-legal-tools></app-legal-tools>
  </section>
  
  <!-- Show only for MANAGEMENT -->
  <section *ngIf="isManagement()">
    <h2>Admin Panel</h2>
    <app-admin-panel></app-admin-panel>
  </section>
</div>
```

## API Endpoints (All Accessible)

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/refresh` - Refresh token
- `GET /api/profile` - Get user profile

### Contracts
- `GET /api/contracts/` - List contracts
- `POST /api/contracts/` - Create contract
- `GET /api/contracts/:id` - Get contract
- `PUT /api/contracts/:id` - Update contract
- `DELETE /api/contracts/:id` - Delete contract
- `POST /api/contracts/:id/status` - Change contract status
- `GET /api/contracts/stats` - Get contract statistics

### Clause Templates
- `GET /api/clauses/` - List clause templates
- `POST /api/clauses/` - Create clause template
- `GET /api/clauses/:id` - Get clause template
- `PUT /api/clauses/:id` - Update clause template
- `DELETE /api/clauses/:id` - Delete clause template
- `GET /api/clauses/search` - Search clause templates
- `GET /api/clauses/types` - Get clause types

### And many more...

## Role Assignment

### For Development/Testing
```sql
-- Update user role in database
UPDATE users SET role = 'LEGAL' WHERE email = 'user@example.com';
UPDATE users SET role = 'MANAGEMENT' WHERE email = 'admin@example.com';
```

### For Production
You can create an admin interface to assign roles, or use database management tools.

## Frontend State Management

### Redux (React)
```javascript
// Store
const initialState = {
  user: {
    id: null,
    email: null,
    role: 'REGULAR',
    token: null
  }
};

// Actions
export const loginSuccess = (userData) => ({
  type: 'LOGIN_SUCCESS',
  payload: userData
});

// Reducer
const authReducer = (state = initialState, action) => {
  switch (action.type) {
    case 'LOGIN_SUCCESS':
      return {
        ...state,
        user: action.payload.user,
        token: action.payload.token
      };
    default:
      return state;
  }
};
```

### Vuex (Vue.js)
```javascript
// Store
const store = new Vuex.Store({
  state: {
    user: {
      id: null,
      email: null,
      role: 'REGULAR',
      token: null
    }
  },
  mutations: {
    LOGIN_SUCCESS(state, userData) {
      state.user = userData.user;
      state.token = userData.token;
    }
  },
  actions: {
    login({ commit }, credentials) {
      // Login logic
      commit('LOGIN_SUCCESS', response.data);
    }
  }
});
```

## Benefits of This Approach

1. **Frontend Flexibility**: Complete control over UI rendering
2. **No Backend Restrictions**: All endpoints accessible
3. **Role Information**: Still available for UI customization
4. **Simple Implementation**: No complex authorization logic
5. **Easy Testing**: All endpoints work for all users
6. **Frontend-Driven**: UI logic controlled by frontend

## Migration Notes

- **Existing Users**: All existing users keep their roles
- **New Users**: Still get `REGULAR` role by default
- **Database**: Role system maintained for frontend use
- **API**: All endpoints now accessible to all authenticated users
- **Frontend**: Can use roles for UI customization

---

**Summary**: The backend now provides role information for frontend rendering without any authorization restrictions. All authenticated users have full access to all endpoints, but the role system is still available for frontend UI customization.
