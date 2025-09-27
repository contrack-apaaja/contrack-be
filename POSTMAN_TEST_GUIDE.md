# Panduan Test AI Analysis dengan Postman

## 1. Setup Environment

### a. Set API Key di .env
Pastikan file `.env` Anda berisi:
```bash
OPENAI_API_KEY=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
DATABASE_URL=your-database-url
JWT_SECRET=your-jwt-secret
PORT=8080
```

### b. Jalankan Migration Database
```bash
./run_migration.sh
```
**Fungsi script ini**: Membuat tabel `clause_risk_analyses` di database untuk menyimpan hasil analisis AI.

### c. Start Server
```bash
go run cmd/server/main.go
```

## 2. Test dengan Postman

### Step 1: Register User (jika belum ada)
**POST** `http://localhost:8080/api/auth/register`

Headers:
```
Content-Type: application/json
```

Body (raw JSON):
```json
{
  "email": "test@example.com",
  "password": "password123",
  "name": "Test User"
}
```

### Step 2: Login untuk mendapatkan JWT Token
**POST** `http://localhost:8080/api/auth/login`

Headers:
```
Content-Type: application/json
```

Body (raw JSON):
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

Response akan memberikan JWT token:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "test@example.com",
      "name": "Test User"
    }
  }
}
```

**Copy token ini untuk digunakan di request berikutnya!**

### Step 3: Buat Clause Template
**POST** `http://localhost:8080/api/clauses`

Headers:
```
Content-Type: application/json
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

Body (raw JSON):
```json
{
  "clause_code": "PAYMENT_001",
  "title": "Ketentuan Pembayaran",
  "type": "Payment",
  "content": "Pembayaran harus dilakukan dalam waktu 30 hari setelah invoice diterima. Jika pembayaran terlambat, akan dikenakan bunga 2% per bulan. Pihak yang membayar bertanggung jawab penuh atas semua biaya yang timbul akibat keterlambatan pembayaran.",
  "is_active": true
}
```

Response:
```json
{
  "success": true,
  "message": "Clause template created successfully",
  "data": {
    "id": 1,
    "clause_code": "PAYMENT_001",
    "title": "Ketentuan Pembayaran",
    "type": "Payment",
    "content": "Pembayaran harus dilakukan dalam waktu 30 hari setelah invoice diterima...",
    "is_active": true,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

**Catat ID clause yang dibuat (misalnya: 1)**

### Step 4: Analisis Clause dengan AI
**POST** `http://localhost:8080/api/ai/analyze`

Headers:
```
Content-Type: application/json
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

Body (raw JSON):
```json
{
  "clause_id": 1
}
```

Response:
```json
{
  "success": true,
  "message": "Analysis completed successfully",
  "data": {
    "analysis": {
      "id": 1,
      "clause_id": 1,
      "risk_level": "medium",
      "risk_score": 65.5,
      "analysis_summary": "Klausul pembayaran ini memiliki risiko sedang terkait dengan ketidakjelasan dalam definisi kewajiban dan potensi konflik dengan regulasi yang berlaku...",
      "identified_risks": [
        "Bunga 2% per bulan mungkin terlalu tinggi dan dapat melanggar prinsip keseimbangan kontrak",
        "Tanggung jawab penuh dapat menimbulkan ketidakadilan dan tidak proporsional",
        "Tidak ada klausul force majeure untuk situasi di luar kendali"
      ],
      "recommendations": [
        "Pertimbangkan untuk mengurangi tingkat bunga menjadi maksimal 1% per bulan",
        "Tambahkan klausul force majeure yang komprehensif",
        "Sertakan mekanisme penyelesaian sengketa yang adil",
        "Tinjau ulang alokasi risiko untuk keseimbangan yang lebih baik"
      ],
      "legal_implications": "Klausul ini dapat melanggar prinsip keseimbangan kontrak dan berpotensi tidak dapat dilaksanakan di pengadilan...",
      "compliance_notes": "Perlu memastikan kepatuhan terhadap UU Perlindungan Konsumen dan regulasi perbankan terkait tingkat bunga...",
      "confidence_score": 85.0,
      "model_version": "gpt-3.5-turbo-v1.0",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "clause": {
      "id": 1,
      "clause_code": "PAYMENT_001",
      "title": "Ketentuan Pembayaran",
      "type": "Payment",
      "content": "Pembayaran harus dilakukan dalam waktu 30 hari setelah invoice diterima...",
      "is_active": true,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  }
}
```

### Step 5: Get Analisis by ID
**GET** `http://localhost:8080/api/ai/analysis/1`

Headers:
```
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

### Step 6: Get Analisis by Clause ID
**GET** `http://localhost:8080/api/ai/analysis/clause/1`

Headers:
```
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

### Step 7: Get Daftar Analisis
**GET** `http://localhost:8080/api/ai/analyses`

Headers:
```
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

Query Parameters (optional):
- `risk_level=high`
- `min_risk_score=70`
- `page=1`
- `limit=10`

### Step 8: Get Statistik Analisis
**GET** `http://localhost:8080/api/ai/stats`

Headers:
```
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

## 3. Test Cases untuk Berbagai Jenis Klausul

### Test Case 1: Klausul Risiko Rendah
```json
{
  "clause_code": "LOW_RISK_001",
  "title": "Klausul Risiko Rendah",
  "type": "General",
  "content": "Para pihak sepakat untuk menyelesaikan sengketa melalui mediasi terlebih dahulu sebelum mengajukan gugatan ke pengadilan."
}
```

### Test Case 2: Klausul Risiko Tinggi
```json
{
  "clause_code": "HIGH_RISK_001",
  "title": "Klausul Risiko Tinggi",
  "type": "Liability",
  "content": "Pihak pertama tidak bertanggung jawab atas kerugian apapun yang timbul dari kontrak ini, termasuk kerugian langsung, tidak langsung, konsekuensial, atau punitif, tanpa batas waktu dan tanpa pengecualian apapun."
}
```

### Test Case 3: Klausul Risiko Sedang
```json
{
  "clause_code": "MEDIUM_RISK_001",
  "title": "Klausul Risiko Sedang",
  "type": "Payment",
  "content": "Pembayaran dilakukan dalam 14 hari setelah barang diterima. Keterlambatan pembayaran akan dikenakan denda 1% per hari."
}
```

## 4. Troubleshooting

### Error 401 Unauthorized
- Pastikan JWT token valid
- Check token expiration
- Pastikan header Authorization: Bearer <token>

### Error 404 Not Found
- Pastikan clause_id yang digunakan sudah ada
- Check apakah clause sudah dibuat sebelumnya

### Error 500 Internal Server Error
- Check apakah OpenAI API key valid
- Pastikan database connection berjalan
- Check server logs untuk detail error

### Error 400 Bad Request
- Pastikan format JSON request benar
- Check apakah semua field required sudah diisi

## 5. Tips Testing

1. **Simpan JWT Token**: Copy token dari response login dan gunakan di semua request berikutnya
2. **Test Berurutan**: Lakukan test sesuai urutan (register → login → create clause → analyze)
3. **Variasi Klausul**: Test dengan berbagai jenis klausul untuk melihat perbedaan analisis
4. **Monitor Response**: Perhatikan risk_level, risk_score, dan confidence_score
5. **Check Database**: Verifikasi data tersimpan dengan benar di database

## 6. Expected Results

### Risk Level Distribution
- **Low (0-24)**: Klausul aman, minimal risiko
- **Medium (25-49)**: Risiko sedang, perlu perhatian
- **High (50-74)**: Risiko tinggi, perlu revisi
- **Critical (75-100)**: Risiko kritis, tidak disarankan

### Confidence Score
- **80-100**: Analisis sangat dapat dipercaya
- **60-79**: Analisis cukup dapat dipercaya
- **40-59**: Analisis kurang dapat dipercaya
- **0-39**: Analisis tidak dapat dipercaya

## 7. Postman Collection

Anda dapat membuat Postman collection dengan struktur berikut:

```
AI Analysis API
├── Auth
│   ├── Register User
│   └── Login User
├── Clause Management
│   ├── Create Clause
│   ├── Get Clause
│   └── List Clauses
└── AI Analysis
    ├── Analyze Clause
    ├── Get Analysis by ID
    ├── Get Analysis by Clause ID
    ├── List Analyses
    ├── Delete Analysis
    └── Get Statistics
```

Setiap request dalam collection dapat menggunakan environment variables:
- `{{base_url}}` = http://localhost:8080
- `{{jwt_token}}` = JWT token dari login response
- `{{clause_id}}` = ID clause yang dibuat
- `{{analysis_id}}` = ID analisis yang dibuat
