# AI Analysis API Documentation

## Overview
Fitur AI rekomendasi untuk analisis risiko klausul kontrak menggunakan OpenAI GPT-3.5-turbo. Sistem ini menganalisis metadata klausul dan memberikan penilaian risiko beserta rekomendasi.

## Endpoints

### 1. Analisis Risiko Klausul
**POST** `/api/ai/analyze`

Menganalisis klausul untuk potensi risiko hukum menggunakan AI.

#### Request Body
```json
{
  "clause_id": 1
}
```

#### Response
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
      "analysis_summary": "Klausul ini memiliki risiko sedang terkait dengan ketidakjelasan dalam definisi kewajiban...",
      "identified_risks": [
        "Ketidakjelasan dalam definisi kewajiban",
        "Potensi konflik dengan regulasi yang berlaku",
        "Ketidakseimbangan dalam alokasi risiko"
      ],
      "recommendations": [
        "Tambahkan definisi yang lebih spesifik untuk kewajiban",
        "Sertakan klausul force majeure yang komprehensif",
        "Tinjau ulang alokasi risiko untuk keseimbangan yang lebih baik"
      ],
      "legal_implications": "Klausul ini dapat menimbulkan ketidakpastian hukum...",
      "compliance_notes": "Perlu memastikan kepatuhan terhadap UU Perlindungan Konsumen...",
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
      "content": "Pembayaran harus dilakukan dalam waktu 30 hari...",
      "is_active": true,
      "created_at": "2024-01-10T09:00:00Z",
      "updated_at": "2024-01-10T09:00:00Z"
    }
  }
}
```

### 2. Get Analisis by ID
**GET** `/api/ai/analysis/{id}`

Mengambil analisis spesifik berdasarkan ID.

#### Response
```json
{
  "success": true,
  "message": "Analysis retrieved successfully",
  "data": {
    "analysis": { /* ... */ },
    "clause": { /* ... */ }
  }
}
```

### 3. Get Analisis by Clause ID
**GET** `/api/ai/analysis/clause/{clause_id}`

Mengambil analisis terbaru untuk klausul tertentu.

#### Response
```json
{
  "success": true,
  "message": "Analysis retrieved successfully",
  "data": {
    "analysis": { /* ... */ },
    "clause": { /* ... */ }
  }
}
```

### 4. Get Daftar Analisis
**GET** `/api/ai/analyses`

Mengambil daftar analisis dengan pagination dan filtering.

#### Query Parameters
- `clause_id` (optional): Filter berdasarkan ID klausul
- `risk_level` (optional): Filter berdasarkan level risiko (low, medium, high, critical)
- `min_risk_score` (optional): Skor risiko minimum (0-100)
- `max_risk_score` (optional): Skor risiko maksimum (0-100)
- `min_confidence` (optional): Skor kepercayaan minimum (0-100)
- `page` (optional): Nomor halaman (default: 1)
- `limit` (optional): Item per halaman (default: 10, max: 100)
- `sort_by` (optional): Field untuk sorting (id, clause_id, risk_level, risk_score, created_at, updated_at)
- `sort_dir` (optional): Arah sorting (asc, desc)

#### Example Request
```
GET /api/ai/analyses?risk_level=high&min_risk_score=70&page=1&limit=20&sort_by=created_at&sort_dir=desc
```

#### Response
```json
{
  "success": true,
  "message": "Analyses retrieved successfully",
  "data": {
    "analyses": [
      {
        "id": 1,
        "clause_id": 1,
        "risk_level": "high",
        "risk_score": 75.5,
        "analysis_summary": "...",
        "identified_risks": ["..."],
        "recommendations": ["..."],
        "legal_implications": "...",
        "compliance_notes": "...",
        "confidence_score": 85.0,
        "model_version": "gpt-3.5-turbo-v1.0",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 50,
      "total_pages": 3,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

### 5. Delete Analisis
**DELETE** `/api/ai/analysis/{id}`

Menghapus analisis berdasarkan ID.

#### Response
```json
{
  "success": true,
  "message": "Analysis deleted successfully",
  "data": null
}
```

### 6. Get Statistik Analisis
**GET** `/api/ai/stats`

Mengambil statistik tentang analisis AI.

#### Response
```json
{
  "success": true,
  "message": "Statistics retrieved successfully",
  "data": {
    "total_analyses": 150,
    "risk_distribution": {
      "low": 45,
      "medium": 60,
      "high": 35,
      "critical": 10
    },
    "average_risk_score": 52.3,
    "average_confidence": 78.5
  }
}
```

## Risk Levels

### Low (0-24)
- Risiko minimal
- Klausul aman untuk digunakan
- Tidak ada rekomendasi khusus

### Medium (25-49)
- Risiko sedang
- Perlu perhatian namun masih dapat diterima
- Rekomendasi perbaikan opsional

### High (50-74)
- Risiko tinggi
- Perlu revisi sebelum digunakan
- Rekomendasi perbaikan wajib

### Critical (75-100)
- Risiko kritis
- Tidak disarankan untuk digunakan
- Perlu revisi menyeluruh

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request format",
  "error": "clause_id is required"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Clause not found",
  "error": "Clause with ID 999 does not exist"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "AI analysis failed",
  "error": "OpenAI API error: 401 - Invalid API key"
}
```

## Authentication

Semua endpoint AI memerlukan authentication. Sertakan token JWT dalam header:

```
Authorization: Bearer <your-jwt-token>
```

## Rate Limiting

- Maksimal 10 analisis per menit per user
- Maksimal 100 analisis per jam per user
- Maksimal 1000 analisis per hari per user

## Environment Variables

Pastikan environment variables berikut diset:

```bash
OPENAI_API_KEY=sk-proj-your-openai-api-key
```

## Database Schema

### clause_risk_analyses Table

```sql
CREATE TABLE clause_risk_analyses (
    id SERIAL PRIMARY KEY,
    clause_id INTEGER NOT NULL REFERENCES clause_templates(id) ON DELETE CASCADE,
    risk_level VARCHAR(20) NOT NULL CHECK (risk_level IN ('low', 'medium', 'high', 'critical')),
    risk_score DECIMAL(5,2) NOT NULL CHECK (risk_score >= 0 AND risk_score <= 100),
    analysis_summary TEXT NOT NULL,
    identified_risks JSONB NOT NULL DEFAULT '[]',
    recommendations JSONB NOT NULL DEFAULT '[]',
    legal_implications TEXT NOT NULL,
    compliance_notes TEXT NOT NULL,
    confidence_score DECIMAL(5,2) NOT NULL CHECK (confidence_score >= 0 AND confidence_score <= 100),
    model_version VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Usage Examples

### 1. Analisis Klausul Baru
```bash
curl -X POST http://localhost:8080/api/ai/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"clause_id": 1}'
```

### 2. Get Analisis Terbaru untuk Klausul
```bash
curl -X GET http://localhost:8080/api/ai/analysis/clause/1 \
  -H "Authorization: Bearer <token>"
```

### 3. Filter Analisis Berisiko Tinggi
```bash
curl -X GET "http://localhost:8080/api/ai/analyses?risk_level=high&min_risk_score=70" \
  -H "Authorization: Bearer <token>"
```

## Best Practices

1. **Analisis Berkala**: Lakukan analisis ulang untuk klausul yang sudah ada
2. **Review Manual**: Selalu review hasil AI dengan ahli hukum
3. **Update Model**: Pantau versi model AI untuk update terbaru
4. **Backup Data**: Backup hasil analisis secara berkala
5. **Monitoring**: Pantau penggunaan API untuk optimasi biaya

## Troubleshooting

### OpenAI API Error
- Pastikan API key valid dan memiliki kredit cukup
- Check rate limit OpenAI
- Verify network connectivity

### Database Error
- Pastikan migration sudah dijalankan
- Check database connection
- Verify table schema

### Authentication Error
- Pastikan token JWT valid
- Check token expiration
- Verify user permissions
