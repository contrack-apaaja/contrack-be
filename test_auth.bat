@echo off
echo 🧪 Testing Contrack Authentication Service
echo.

echo 📡 Step 1: Health Check
curl -s http://localhost:8080/api/hello
echo.
echo.

echo 👤 Step 2: Register User
curl -s -X POST http://localhost:8080/api/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@contrack.com\",\"password\":\"hackathon2024\"}"
echo.
echo.

echo 🔐 Step 3: Login User  
curl -s -X POST http://localhost:8080/api/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@contrack.com\",\"password\":\"hackathon2024\"}"
echo.
echo.

echo ⚠️  Copy the token from above and run:
echo curl -X GET http://localhost:8080/api/profile -H "Authorization: Bearer YOUR_TOKEN"
echo.

pause
