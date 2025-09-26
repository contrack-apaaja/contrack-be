# 🔧 CORS Troubleshooting Guide

## Problem Solved ✅

The CORS error you encountered:
```
Access to XMLHttpRequest at 'http://localhost:8080/api/clauses' from origin 'http://localhost:3000' has been blocked by CORS policy: No 'Access-Control-Allow-Origin' header is present on the requested resource.
```

Has been fixed by implementing a comprehensive CORS middleware.

## What Was Fixed

1. **Enhanced CORS Headers**: Added all necessary headers for cross-origin requests
2. **OPTIONS Handling**: Proper preflight request handling
3. **Origin Flexibility**: Dynamic origin handling for development
4. **Credentials Support**: Enabled for authenticated requests

## Steps to Apply the Fix

1. **Restart your server**:
```bash
go run ./cmd/server
```

2. **Test from your frontend** (localhost:3000):
```javascript
// This should now work without CORS errors
fetch('http://localhost:8080/api/clauses', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer YOUR_JWT_TOKEN',
    'Content-Type': 'application/json'
  }
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));
```

## CORS Headers Now Included

- `Access-Control-Allow-Origin`: Dynamic origin or wildcard
- `Access-Control-Allow-Credentials`: true
- `Access-Control-Allow-Methods`: GET, POST, PUT, DELETE, PATCH, OPTIONS
- `Access-Control-Allow-Headers`: All common headers including Authorization
- `Access-Control-Expose-Headers`: Response headers accessible to frontend
- `Access-Control-Max-Age`: Caches preflight for 24 hours

## Frontend Integration Examples

### JavaScript/Fetch
```javascript
// GET request
const getClauses = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/clauses', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    console.log('Clauses:', data);
    return data;
  } catch (error) {
    console.error('Error fetching clauses:', error);
  }
};

// POST request
const createClause = async (clauseData) => {
  try {
    const response = await fetch('http://localhost:8080/api/clauses', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(clauseData)
    });
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error creating clause:', error);
  }
};
```

### Axios
```javascript
import axios from 'axios';

// Configure axios defaults
axios.defaults.baseURL = 'http://localhost:8080/api';
axios.defaults.headers.common['Authorization'] = `Bearer ${localStorage.getItem('token')}`;

// GET request
const getClauses = async () => {
  try {
    const response = await axios.get('/clauses');
    return response.data;
  } catch (error) {
    console.error('Error:', error.response?.data || error.message);
  }
};

// POST request
const createClause = async (clauseData) => {
  try {
    const response = await axios.post('/clauses', clauseData);
    return response.data;
  } catch (error) {
    console.error('Error:', error.response?.data || error.message);
  }
};
```

### React Hook Example
```javascript
import { useState, useEffect } from 'react';

const useClauses = () => {
  const [clauses, setClauses] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchClauses = async () => {
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/api/clauses', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        }
      });
      
      const data = await response.json();
      
      if (data.status === 'success') {
        setClauses(data.data.clause_templates);
      } else {
        setError(data.message);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchClauses();
  }, []);

  return { clauses, loading, error, refetch: fetchClauses };
};
```

## Production Considerations

For production, consider using the more restrictive CORS middleware:

```go
// In cmd/server/main.go, replace the CORS middleware with:
allowedOrigins := []string{
    "https://yourdomain.com",
    "https://app.yourdomain.com",
}
r.Use(middleware.CORSMiddlewareWithOrigins(allowedOrigins))
```

## Common Issues & Solutions

### 1. Still Getting CORS Errors
- Make sure you restarted the server after the changes
- Clear browser cache
- Check browser developer tools for the actual error

### 2. Authentication Issues
- Ensure JWT token is valid and not expired
- Check that the Authorization header is properly formatted: `Bearer <token>`

### 3. Preflight Requests Failing
- The middleware now handles OPTIONS requests properly
- Make sure your frontend is sending the correct headers

### 4. Credentials Not Working
- The middleware now includes `Access-Control-Allow-Credentials: true`
- Make sure your frontend includes `credentials: 'include'` if needed

## Testing CORS

You can test CORS directly from browser console:

```javascript
// Test from browser console on localhost:3000
fetch('http://localhost:8080/api/hello')
  .then(response => response.json())
  .then(data => console.log('CORS working:', data))
  .catch(error => console.error('CORS error:', error));
```

Your CORS issues should now be resolved! 🎉
