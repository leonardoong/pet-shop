CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL,
    user_type  VARCHAR(10) NOT NULL CHECK (user_type IN ('customer', 'admin')),
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked    BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_refresh_tokens_hash   ON refresh_tokens (token_hash);
CREATE        INDEX idx_refresh_tokens_user   ON refresh_tokens (user_id, user_type);
