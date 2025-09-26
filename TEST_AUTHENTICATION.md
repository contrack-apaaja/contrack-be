# 🧪 Complete Authentication Testing Guide

## 📋 **Prerequisites**

1. **Create your .env file**:
```bash
# Copy the example
cp env.example .env
```

2. **Update your .env file** with these values:
```env
# Server Configuration
PORT=8080

# Database Configuration (REPLACE WITH YOUR ACTUAL PASSWORD)
DATABASE_URL=postgresql://postgres.qzuqwvvurdouhcvaybju:YOUR_ACTUAL_SUPABASE_PASSWORD@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres

# JWT Configuration (Use this secure secret)
JWT_SECRET=c0ntr4ck-h4ck4th0n-2024-sup3r-s3cr3t-jwt-k3y-f0r-pr0duct10n-us3-r4nd0m-g3n3r4t3d-k3y

# Supabase Configuration (Optional)
SUPABASE_URL=
SUPABASE_KEY=
```

## 🚀 **Start the Server**

```bash
go run ./cmd/server
```

You should see output like:
```
Successfully connected to PostgreSQL database
Users table created successfully
Database migrations completed successfully
🚀 Server starting on port 8080
📚 API Documentation:
   POST /api/auth/register - Register new user
   POST /api/auth/login    - Login user
   POST /api/auth/refresh  - Refresh token
   GET  /api/profile       - Get user profile (protected)
   GET  /api/hello         - Health check
```

## 🧪 **Testing Steps**

### **Step 1: Health Check**
```bash
curl http://localhost:8080/api/hello
```
**Expected Response:**
```json
{"message":"Hello from Gin!"}
```

### **Step 2: Register a New User**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "daffa@contrack.com",
    "password": "hackathon2024"
  }'
```

**Expected Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "daffa@contrack.com",
    "created_at": "2024-09-26T10:30:00Z",
    "updated_at": "2024-09-26T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### **Step 3: Try to Register Same User Again (Should Fail)**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "daffa@contrack.com",
    "password": "hackathon2024"
  }'
```

**Expected Response (409 Conflict):**
```json
{
  "error": "User already exists"
}
```

### **Step 4: Login with the User**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "daffa@contrack.com",
    "password": "hackathon2024"
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "Login successful",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "daffa@contrack.com",
    "created_at": "2024-09-26T10:30:00Z",
    "updated_at": "2024-09-26T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**💡 SAVE THE TOKEN** from this response for the next steps!

### **Step 5: Try Login with Wrong Password (Should Fail)**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "daffa@contrack.com",
    "password": "wrongpassword"
  }'
```

**Expected Response (401 Unauthorized):**
```json
{
  "error": "Invalid email or password"
}
```

### **Step 6: Access Protected Profile Endpoint**
```bash
# Replace YOUR_JWT_TOKEN with the actual token from Step 4
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Expected Response (200 OK):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "daffa@contrack.com",
    "created_at": "2024-09-26T10:30:00Z",
    "updated_at": "2024-09-26T10:30:00Z"
  }
}
```

### **Step 7: Try Protected Endpoint Without Token (Should Fail)**
```bash
curl -X GET http://localhost:8080/api/profile
```

**Expected Response (401 Unauthorized):**
```json
{
  "error": "Authorization header is required"
}
```

### **Step 8: Try Protected Endpoint with Invalid Token (Should Fail)**
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer invalid.token.here"
```

**Expected Response (401 Unauthorized):**
```json
{
  "error": "Invalid or expired token"
}
```

### **Step 9: Refresh Token**
```bash
# Replace YOUR_JWT_TOKEN with the actual token
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_JWT_TOKEN"
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "Token refreshed successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### **Step 10: Test Protected Users Endpoint**
```bash
# Replace YOUR_JWT_TOKEN with a valid token
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🔧 **PowerShell Testing (Windows)**

If you're using PowerShell, here are the equivalent commands:

### Register User:
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/auth/register" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"email":"daffa@contrack.com","password":"hackathon2024"}'
```

### Login:
```powershell
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/login" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"email":"daffa@contrack.com","password":"hackathon2024"}'

$token = $response.token
Write-Host "Token: $token"
```

### Access Profile:
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/profile" `
  -Method GET `
  -Headers @{"Authorization"="Bearer $token"}
```

## 🐛 **Troubleshooting**

### **Database Connection Issues:**
- Make sure you've replaced `[YOUR-PASSWORD]` in the DATABASE_URL
- Check that your Supabase database is accessible
- Verify the connection string format

### **JWT Token Issues:**
- Tokens expire after 24 hours by default
- Make sure you're using the `Bearer ` prefix
- Check that JWT_SECRET is set correctly

### **Common HTTP Status Codes:**
- `200`: Success
- `201`: Created (successful registration)
- `400`: Bad Request (invalid JSON or validation errors)
- `401`: Unauthorized (invalid credentials or token)
- `409`: Conflict (user already exists)
- `500`: Internal Server Error (database or server issues)

## 📱 **Testing with Postman**

1. **Import this collection** into Postman:
   - Create a new collection called "Contrack Auth"
   - Add the following requests:

2. **Register Request:**
   - Method: POST
   - URL: `http://localhost:8080/api/auth/register`
   - Headers: `Content-Type: application/json`
   - Body: Raw JSON
   ```json
   {
     "email": "daffa@contrack.com",
     "password": "hackathon2024"
   }
   ```

3. **Login Request:**
   - Method: POST
   - URL: `http://localhost:8080/api/auth/login`
   - Headers: `Content-Type: application/json`
   - Body: Raw JSON (same as register)

4. **Profile Request:**
   - Method: GET
   - URL: `http://localhost:8080/api/profile`
   - Headers: `Authorization: Bearer {{token}}`
   - Use Postman variables to store the token

## ✅ **Success Indicators**

Your authentication system is working correctly if:

1. ✅ Server starts without errors
2. ✅ Database connection is established
3. ✅ User registration works and returns a token
4. ✅ Login works with correct credentials
5. ✅ Login fails with wrong credentials
6. ✅ Protected endpoints require valid tokens
7. ✅ Token refresh works
8. ✅ Invalid tokens are rejected

## 🚀 **Next Steps**

Once testing is complete, you can:
1. Integrate with your frontend application
2. Add more protected endpoints
3. Implement additional features like password reset
4. Add rate limiting for production
5. Set up proper logging and monitoring
