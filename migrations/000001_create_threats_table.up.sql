CREATE TABLE IF NOT EXISTS threats (
    id BIGSERIAL PRIMARY KEY,
    md5 VARCHAR(32) UNIQUE NOT NULL CHECK (md5 ~ '^[a-fA-F0-9]{32}$'),
    sha256 VARCHAR(64) CHECK (sha256 IS NULL OR sha256 ~ '^[a-fA-F0-9]{64}$'),
    name VARCHAR(255) NOT NULL, 
    family VARCHAR(100), 
    type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    confidence SMALLINT CHECK(confidence >= 0 AND confidence <= 100), 
    description TEXT,
    tags TEXT[] NOT NULL DEFAULT '{}',
    source VARCHAR(255), 
    first_seen TIMESTAMPTZ,
    last_seen TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
