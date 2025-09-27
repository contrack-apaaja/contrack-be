# Contract Recommendations API

## 📋 **Overview**

Endpoint untuk mengambil semua rekomendasi AI dari contract berdasarkan ID contract.

## 🚀 **Endpoint**

```
GET /api/ai/contract/{contract_id}/recommendations
```

## 📝 **Request**

### **Headers**
```
Authorization: Bearer <token>
Content-Type: application/json
```

### **URL Parameters**
- `contract_id` (int, required): ID contract yang ingin diambil rekomendasinya

### **Example Request**
```bash
curl -X GET "http://localhost:8080/api/ai/contract/1/recommendations" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 📤 **Response**

### **Success Response (200)**
```json
{
  "status": "success",
  "message": "Contract recommendations retrieved successfully",
  "data": {
    "contract_id": 1,
    "overall_risk_level": "high",
    "overall_risk_score": 68.0,
    "total_clauses": 10,
    "clause_recommendations": [
      {
        "clause_id": 4,
        "risk_level": "high",
        "risk_score": 85.0,
        "analysis_summary": "Klausul ini memiliki risiko tinggi karena melibatkan berbagai aspek seperti perhitungan harga kontrak, pembayaran, dan pembiayaan yang kompleks.",
        "recommendations": [
          "Menyediakan perincian yang lebih jelas mengenai perhitungan PPN dan PPh untuk menghindari perselisihan di masa depan",
          "Menetapkan ketentuan yang lebih spesifik mengenai penggunaan uang muka maksimum 10%",
          "Menyusun mekanisme penagihan biaya produksi material yang lebih terstruktur dan transparan",
          "Mengurangi jumlah rekening pembayaran agar meminimalkan risiko keamanan pembayaran"
        ],
        "identified_risks": [
          "Risiko perhitungan PPN dan PPh yang tidak jelas dan dapat menimbulkan perselisihan",
          "Risiko pembayaran uang muka yang tinggi tanpa ketentuan yang jelas mengenai penggunaannya",
          "Risiko penagihan biaya produksi material yang ambigu dan dapat disalahgunakan",
          "Risiko keamanan pembayaran karena terdapat beberapa rekening yang harus dipilih tanpa penjelasan lebih lanjut"
        ],
        "legal_implications": "Klausul ini dapat menimbulkan sengketa antara pihak-pihak yang terlibat jika tidak dijelaskan dengan lebih rinci. Selain itu, ketidakjelasan dalam pembayaran dan penagihan biaya dapat merugikan salah satu pihak dan melanggar prinsip keadilan kontrak.",
        "compliance_notes": "Penting untuk memastikan bahwa semua ketentuan dalam klausul ini sesuai dengan hukum yang berlaku dan tidak melanggar prinsip-prinsip kontrak yang adil.",
        "confidence_score": 90.0,
        "created_at": "2025-09-26T17:58:30.185004Z"
      }
    ],
    "overall_recommendations": [
      "Menyediakan perincian yang lebih jelas mengenai perhitungan PPN dan PPh untuk menghindari perselisihan di masa depan",
      "Menetapkan ketentuan yang lebih spesifik mengenai penggunaan uang muka maksimum 10%",
      "Menyusun mekanisme penagihan biaya produksi material yang lebih terstruktur dan transparan",
      "Mengurangi jumlah rekening pembayaran agar meminimalkan risiko keamanan pembayaran"
    ],
    "key_risks": [
      "Risiko perhitungan PPN dan PPh yang tidak jelas dan dapat menimbulkan perselisihan",
      "Risiko pembayaran uang muka yang tinggi tanpa ketentuan yang jelas mengenai penggunaannya",
      "Risiko penagihan biaya produksi material yang ambigu dan dapat disalahgunakan"
    ],
    "created_at": "2025-09-26T17:58:30.185004Z"
  }
}
```

### **Error Responses**

#### **400 Bad Request - Invalid Contract ID**
```json
{
  "status": "error",
  "message": "Invalid contract ID",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": "Contract ID must be a valid integer"
  }
}
```

#### **401 Unauthorized**
```json
{
  "status": "error",
  "message": "User not authenticated",
  "error": {
    "code": "AUTH_ERROR",
    "details": null
  }
}
```

#### **404 Not Found - No Recommendations**
```json
{
  "status": "error",
  "message": "No recommendations found for this contract",
  "error": {
    "code": "NOT_FOUND",
    "details": "No AI analysis found for this contract"
  }
}
```

#### **500 Internal Server Error**
```json
{
  "status": "error",
  "message": "Failed to get contract recommendations",
  "error": {
    "code": "DATABASE_ERROR",
    "details": "Database connection error"
  }
}
```

## 📊 **Response Fields**

### **Main Response**
- `contract_id` (int): ID contract yang dianalisis
- `overall_risk_level` (string): Level risiko keseluruhan contract (low/medium/high/critical)
- `overall_risk_score` (float): Skor risiko rata-rata (0-100)
- `total_clauses` (int): Jumlah klausul yang dianalisis
- `clause_recommendations` (array): Array rekomendasi per klausul
- `overall_recommendations` (array): Array rekomendasi keseluruhan (deduplicated)
- `key_risks` (array): Array risiko utama yang teridentifikasi (deduplicated)
- `created_at` (string): Timestamp analisis terbaru

### **Clause Recommendation Object**
- `clause_id` (int): ID klausul
- `risk_level` (string): Level risiko klausul (low/medium/high/critical)
- `risk_score` (float): Skor risiko klausul (0-100)
- `analysis_summary` (string): Ringkasan analisis klausul
- `recommendations` (array): Array rekomendasi spesifik untuk klausul
- `identified_risks` (array): Array risiko yang teridentifikasi
- `legal_implications` (string): Implikasi hukum
- `compliance_notes` (string): Catatan kepatuhan
- `confidence_score` (float): Skor kepercayaan analisis (0-100)
- `created_at` (string): Timestamp analisis

## 🔧 **Usage Examples**

### **JavaScript/Fetch**
```javascript
const response = await fetch('http://localhost:8080/api/ai/contract/1/recommendations', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer YOUR_TOKEN'
  }
});

const data = await response.json();
console.log(data);
```

### **Python/Requests**
```python
import requests

url = "http://localhost:8080/api/ai/contract/1/recommendations"
headers = {
    "Authorization": "Bearer YOUR_TOKEN"
}

response = requests.get(url, headers=headers)
data = response.json()
print(data)
```

### **cURL**
```bash
curl -X GET "http://localhost:8080/api/ai/contract/1/recommendations" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🎯 **Use Cases**

1. **Dashboard Contract Overview**: Menampilkan ringkasan rekomendasi untuk contract tertentu
2. **Contract Review**: Melihat semua rekomendasi AI untuk contract sebelum approval
3. **Risk Assessment**: Menganalisis tingkat risiko keseluruhan contract
4. **Compliance Check**: Memverifikasi kepatuhan contract terhadap regulasi

## ⚠️ **Notes**

- Endpoint ini memerlukan authentication (Bearer token)
- Data yang dikembalikan berdasarkan analisis AI yang sudah ada di database
- Jika tidak ada analisis untuk contract, akan mengembalikan 404
- Rekomendasi dan risiko di-deduplicate untuk menghindari duplikasi
- Data diurutkan berdasarkan timestamp terbaru (DESC)

## 🔗 **Related Endpoints**

- `POST /api/ai/analyze-contract` - Membuat analisis contract baru
- `GET /api/ai/analyses` - Melihat semua analisis dengan pagination
- `GET /api/ai/stats` - Statistik analisis AI
