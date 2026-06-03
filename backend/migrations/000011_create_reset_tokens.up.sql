CREATE TABLE IF NOT EXISTS reset_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) NOT NULL,
    token_hash  VARCHAR(64) UNIQUE NOT NULL,
    user_type   VARCHAR(10) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    used        BOOLEAN DEFAULT false,
    created_at  TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_reset_tokens_hash ON reset_tokens(token_hash);
CREATE INDEX idx_reset_tokens_email ON reset_tokens(email);
