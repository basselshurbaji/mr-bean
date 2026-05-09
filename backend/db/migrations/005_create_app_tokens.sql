-- +goose Up
CREATE TABLE app_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    app_name   TEXT NOT NULL,
    revoked    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX app_tokens_user_id_idx ON app_tokens(user_id);

-- +goose Down
DROP TABLE IF EXISTS app_tokens;