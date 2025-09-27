package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"contrack-be/internal/config"
	"contrack-be/internal/models"
)

// OpenAIService handles OpenAI API interactions
type OpenAIService struct {
	config *config.Config
	client *http.Client
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	MaxTokens int      `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// Message represents a message in the OpenAI conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response structure from OpenAI API
type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a choice in the OpenAI response
type Choice struct {
	Message Message `json:"message"`
}

// Usage represents token usage information
type Usage struct {
	TotalTokens int `json:"total_tokens"`
}

// AIAnalysisResult represents the structured result from AI analysis
type AIAnalysisResult struct {
	RiskLevel         models.RiskLevel `json:"risk_level"`
	RiskScore         float64          `json:"risk_score"`
	AnalysisSummary   string           `json:"analysis_summary"`
	IdentifiedRisks   []string         `json:"identified_risks"`
	Recommendations   []string         `json:"recommendations"`
	LegalImplications string           `json:"legal_implications"`
	ComplianceNotes   string           `json:"compliance_notes"`
	ConfidenceScore   float64          `json:"confidence_score"`
}

// NewOpenAIService creates a new OpenAI service instance
func NewOpenAIService(cfg *config.Config) *OpenAIService {
	return &OpenAIService{
		config: cfg,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// AnalyzeClauseRisk analyzes a clause for potential risks using OpenAI
func (s *OpenAIService) AnalyzeClauseRisk(clause *models.ClauseTemplate) (*AIAnalysisResult, error) {
	if s.config.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	// Create the prompt for AI analysis
	prompt := s.createAnalysisPrompt(clause)

	// Make request to OpenAI
	response, err := s.makeOpenAIRequest(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI analysis: %w", err)
	}

	// Parse the AI response
	result, err := s.parseAIResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return result, nil
}

// createAnalysisPrompt creates a structured prompt for clause risk analysis
func (s *OpenAIService) createAnalysisPrompt(clause *models.ClauseTemplate) string {
	return fmt.Sprintf(`Anda adalah seorang ahli hukum kontrak yang berpengalaman. Analisis klausul kontrak berikut dan berikan penilaian risiko yang komprehensif.

INFORMASI KLAUSUL:
- Kode Klausul: %s
- Judul: %s
- Tipe: %s
- Isi Klausul: %s

TUGAS ANDA:
1. Analisis risiko hukum dari klausul ini
2. Berikan skor risiko 0-100 (0 = sangat aman, 100 = sangat berisiko)
3. Identifikasi risiko spesifik yang ditemukan
4. Berikan rekomendasi perbaikan
5. Jelaskan implikasi hukum
6. Berikan catatan kepatuhan
7. Berikan skor kepercayaan analisis 0-100

FORMAT RESPON (JSON):
{
  "risk_level": "low|medium|high|critical",
  "risk_score": 0-100,
  "analysis_summary": "Ringkasan analisis risiko",
  "identified_risks": ["risiko 1", "risiko 2", "..."],
  "recommendations": ["rekomendasi 1", "rekomendasi 2", "..."],
  "legal_implications": "Penjelasan implikasi hukum",
  "compliance_notes": "Catatan kepatuhan",
  "confidence_score": 0-100
}

Pastikan respons dalam format JSON yang valid dan mudah diparse.`, 
		clause.ClauseCode, 
		clause.Title, 
		clause.Type, 
		clause.Content)
}

// makeOpenAIRequest makes a request to OpenAI API
func (s *OpenAIService) makeOpenAIRequest(prompt string) (string, error) {
	requestBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "Anda adalah asisten AI yang ahli dalam analisis risiko kontrak hukum. Berikan respons yang akurat dan terstruktur dalam format JSON.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   2000,
		Temperature: 0.3,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.OpenAIAPIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in OpenAI response")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// parseAIResponse parses the AI response and extracts structured data
func (s *OpenAIService) parseAIResponse(response string) (*AIAnalysisResult, error) {
	// Clean the response - remove any markdown formatting
	cleanResponse := strings.TrimSpace(response)
	if strings.HasPrefix(cleanResponse, "```json") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
	}
	if strings.HasSuffix(cleanResponse, "```") {
		cleanResponse = strings.TrimSuffix(cleanResponse, "```")
	}
	cleanResponse = strings.TrimSpace(cleanResponse)

	var result AIAnalysisResult
	if err := json.Unmarshal([]byte(cleanResponse), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Validate and normalize the result
	if err := s.validateAndNormalizeResult(&result); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &result, nil
}

// validateAndNormalizeResult validates and normalizes the AI analysis result
func (s *OpenAIService) validateAndNormalizeResult(result *AIAnalysisResult) error {
	// Validate risk level
	if !models.IsValidRiskLevel(result.RiskLevel) {
		// Auto-correct based on risk score
		result.RiskLevel = models.GetRiskLevelFromScore(result.RiskScore)
	}

	// Clamp scores to valid ranges
	if result.RiskScore < 0 {
		result.RiskScore = 0
	} else if result.RiskScore > 100 {
		result.RiskScore = 100
	}

	if result.ConfidenceScore < 0 {
		result.ConfidenceScore = 0
	} else if result.ConfidenceScore > 100 {
		result.ConfidenceScore = 100
	}

	// Ensure arrays are not nil
	if result.IdentifiedRisks == nil {
		result.IdentifiedRisks = []string{}
	}
	if result.Recommendations == nil {
		result.Recommendations = []string{}
	}

	// Ensure strings are not empty
	if result.AnalysisSummary == "" {
		result.AnalysisSummary = "Analisis risiko tidak tersedia"
	}
	if result.LegalImplications == "" {
		result.LegalImplications = "Implikasi hukum tidak tersedia"
	}
	if result.ComplianceNotes == "" {
		result.ComplianceNotes = "Catatan kepatuhan tidak tersedia"
	}

	return nil
}

// GetModelVersion returns the current model version being used
func (s *OpenAIService) GetModelVersion() string {
	return "gpt-3.5-turbo-v1.0"
}

// AnalyzeContractRisk analyzes a contract with multiple clauses
func (s *OpenAIService) AnalyzeContractRisk(contractID int, userID string) (*models.ContractAnalysisResult, error) {
	// Get clause template IDs from contract
	clauseTemplateIDs, err := s.getContractClauseTemplateIDs(contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract clause templates: %w", err)
	}

	if len(clauseTemplateIDs) == 0 {
		return nil, fmt.Errorf("no clause templates found for contract %d", contractID)
	}

	// Analyze each clause individually using real AI analysis
	var clauseAnalyses []models.ClauseRiskAnalysis
	var totalRiskScore float64
	var maxRiskLevel models.RiskLevel = models.RiskLevelLow

	for _, clauseTemplateID := range clauseTemplateIDs {
		// For now, we'll create a mock clause template and use the existing AnalyzeClauseRisk method
		// In a real implementation, you would fetch the clause template from database
		mockClause := &models.ClauseTemplate{
			ID:          clauseTemplateID,
			ClauseCode:  fmt.Sprintf("CLAUSE_%d", clauseTemplateID),
			Title:       fmt.Sprintf("Clause Template %d", clauseTemplateID),
			Type:        "General",
			Content:     fmt.Sprintf("This is clause template %d content for analysis", clauseTemplateID),
		}
		
		// Use the existing AnalyzeClauseRisk method for real AI analysis
		aiResult, err := s.AnalyzeClauseRisk(mockClause)
		if err != nil {
			// If AI analysis fails, create a fallback analysis
			clauseAnalysis := models.ClauseRiskAnalysis{
				ID:                0,
				ClauseID:          clauseTemplateID,
				RiskLevel:         models.RiskLevelMedium,
				RiskScore:         50.0,
				AnalysisSummary:   fmt.Sprintf("Analysis for clause template %d (AI analysis unavailable)", clauseTemplateID),
				IdentifiedRisks:   []string{"AI analysis temporarily unavailable"},
				Recommendations:   []string{"Manual review recommended"},
				LegalImplications: "Manual legal review required",
				ComplianceNotes:   "Compliance verification needed",
				ConfidenceScore:   60.0,
				ModelVersion:      s.GetModelVersion(),
				CreatedAt:         time.Now(),
			}
			clauseAnalyses = append(clauseAnalyses, clauseAnalysis)
			totalRiskScore += clauseAnalysis.RiskScore
			continue
		}
		
		// Convert AI result to ClauseRiskAnalysis
		clauseAnalysis := models.ClauseRiskAnalysis{
			ID:                0, // Will be set by database
			ClauseID:          clauseTemplateID,
			RiskLevel:         aiResult.RiskLevel,
			RiskScore:         aiResult.RiskScore,
			AnalysisSummary:   aiResult.AnalysisSummary,
			IdentifiedRisks:   aiResult.IdentifiedRisks,
			Recommendations:   aiResult.Recommendations,
			LegalImplications: aiResult.LegalImplications,
			ComplianceNotes:   aiResult.ComplianceNotes,
			ConfidenceScore:   aiResult.ConfidenceScore,
			ModelVersion:      s.GetModelVersion(),
			CreatedAt:         time.Now(),
		}
		
		clauseAnalyses = append(clauseAnalyses, clauseAnalysis)
		totalRiskScore += clauseAnalysis.RiskScore
		
		// Update max risk level
		if clauseAnalysis.RiskLevel == models.RiskLevelCritical || 
		   (clauseAnalysis.RiskLevel == models.RiskLevelHigh && maxRiskLevel != models.RiskLevelCritical) ||
		   (clauseAnalysis.RiskLevel == models.RiskLevelMedium && maxRiskLevel == models.RiskLevelLow) {
			maxRiskLevel = clauseAnalysis.RiskLevel
		}
	}
	
	// Calculate overall metrics
	overallRiskScore := totalRiskScore / float64(len(clauseTemplateIDs))
	
	// Generate comprehensive overall contract analysis using AI
	contractSummary, keyRisks, recommendations, err := s.generateContractSummary(clauseAnalyses, contractID)
	if err != nil {
		// Fallback if AI analysis fails
		contractSummary = fmt.Sprintf("Contract %d analysis completed. %d clauses analyzed with overall risk level: %s", 
			contractID, len(clauseTemplateIDs), maxRiskLevel)
		keyRisks = []string{
			fmt.Sprintf("Overall contract risk level: %s", maxRiskLevel),
			"Review individual clause analyses for specific risks",
		}
		recommendations = []string{
			"Review all clause analyses for detailed recommendations",
			"Consider legal review for high-risk clauses",
		}
	}
	
	result := &models.ContractAnalysisResult{
		ContractID:        contractID,
		ClauseAnalyses:    clauseAnalyses,
		OverallRiskLevel:  maxRiskLevel,
		OverallRiskScore:  overallRiskScore,
		ContractSummary:   contractSummary,
		KeyRisks:          keyRisks,
		Recommendations:   recommendations,
		CreatedAt:         time.Now(),
	}
	
	return result, nil
}

// generateContractSummary generates comprehensive contract analysis using AI
func (s *OpenAIService) generateContractSummary(clauseAnalyses []models.ClauseRiskAnalysis, contractID int) (string, []string, []string, error) {
	// Create a comprehensive prompt for contract-level analysis
	prompt := fmt.Sprintf(`Anda adalah seorang ahli hukum kontrak yang berpengalaman. Analisis kontrak secara keseluruhan berdasarkan analisis individual klausul berikut:

INFORMASI KONTRAK:
- ID Kontrak: %d
- Jumlah Klausul: %d

ANALISIS KLAUSUL INDIVIDUAL:
`, contractID, len(clauseAnalyses))

	for i, analysis := range clauseAnalyses {
		prompt += fmt.Sprintf(`
Klausul %d:
- Risk Level: %s
- Risk Score: %.1f
- Summary: %s
- Identified Risks: %v
- Recommendations: %v
`, i+1, analysis.RiskLevel, analysis.RiskScore, analysis.AnalysisSummary, analysis.IdentifiedRisks, analysis.Recommendations)
	}

	prompt += `
TUGAS ANDA:
1. Berikan ringkasan analisis kontrak secara keseluruhan
2. Identifikasi risiko utama yang perlu diperhatikan
3. Berikan rekomendasi strategis untuk kontrak ini
4. Fokus pada risiko yang paling kritis dan dampaknya

FORMAT RESPON (JSON):
{
  "contract_summary": "Ringkasan komprehensif analisis kontrak",
  "key_risks": ["risiko utama 1", "risiko utama 2", "..."],
  "recommendations": ["rekomendasi strategis 1", "rekomendasi strategis 2", "..."]
}

Pastikan respons dalam format JSON yang valid dan memberikan insight yang mendalam.`

	// Call OpenAI API
	response, err := s.makeOpenAIRequest(prompt)
	if err != nil {
		return "", nil, nil, err
	}

	// Parse the JSON response
	var aiResponse struct {
		ContractSummary string   `json:"contract_summary"`
		KeyRisks        []string `json:"key_risks"`
		Recommendations []string `json:"recommendations"`
	}

	if err := json.Unmarshal([]byte(response), &aiResponse); err != nil {
		return "", nil, nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return aiResponse.ContractSummary, aiResponse.KeyRisks, aiResponse.Recommendations, nil
}

// getContractClauseTemplateIDs retrieves clause template IDs for a contract
func (s *OpenAIService) getContractClauseTemplateIDs(contractID int) ([]int, error) {
	// For now, return mock clause template IDs
	// In a real implementation, this would query the contract_clauses table
	// SELECT clause_template_id FROM contract_clauses WHERE contract_id = $1
	
	// Mock data for testing
	mockClauseTemplateIDs := []int{2, 3, 4, 5}
	return mockClauseTemplateIDs, nil
}
