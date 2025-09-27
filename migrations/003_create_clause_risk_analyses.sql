-- Create clause_risk_analyses table for AI analysis results
CREATE TABLE IF NOT EXISTS clause_risk_analyses (
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

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_clause_risk_analyses_clause_id ON clause_risk_analyses(clause_id);
CREATE INDEX IF NOT EXISTS idx_clause_risk_analyses_risk_level ON clause_risk_analyses(risk_level);
CREATE INDEX IF NOT EXISTS idx_clause_risk_analyses_risk_score ON clause_risk_analyses(risk_score);
CREATE INDEX IF NOT EXISTS idx_clause_risk_analyses_created_at ON clause_risk_analyses(created_at);
CREATE INDEX IF NOT EXISTS idx_clause_risk_analyses_confidence_score ON clause_risk_analyses(confidence_score);

-- Create composite index for common queries
CREATE INDEX IF NOT EXISTS idx_clause_risk_analyses_clause_created ON clause_risk_analyses(clause_id, created_at DESC);

-- Add comments for documentation
COMMENT ON TABLE clause_risk_analyses IS 'Stores AI-powered risk analysis results for clause templates';
COMMENT ON COLUMN clause_risk_analyses.clause_id IS 'Reference to the clause template being analyzed';
COMMENT ON COLUMN clause_risk_analyses.risk_level IS 'Categorical risk level: low, medium, high, critical';
COMMENT ON COLUMN clause_risk_analyses.risk_score IS 'Numerical risk score from 0-100';
COMMENT ON COLUMN clause_risk_analyses.analysis_summary IS 'AI-generated summary of the risk analysis';
COMMENT ON COLUMN clause_risk_analyses.identified_risks IS 'JSON array of specific risks identified';
COMMENT ON COLUMN clause_risk_analyses.recommendations IS 'JSON array of recommendations provided by AI';
COMMENT ON COLUMN clause_risk_analyses.legal_implications IS 'AI-generated explanation of legal implications';
COMMENT ON COLUMN clause_risk_analyses.compliance_notes IS 'AI-generated compliance-related notes';
COMMENT ON COLUMN clause_risk_analyses.confidence_score IS 'AI confidence in the analysis (0-100)';
COMMENT ON COLUMN clause_risk_analyses.model_version IS 'Version of the AI model used for analysis';
