-- +goose Up
CREATE TABLE beans (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL REFERENCES users(id),
    name          TEXT        NOT NULL,
    roaster       TEXT,
    origin        TEXT,
    process       TEXT,
    roast_level   TEXT,
    tasting_notes TEXT,
    notes         TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE beans;
