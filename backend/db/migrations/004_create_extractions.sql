-- +goose Up
CREATE TABLE extractions (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users(id),
    bean_id      UUID        NOT NULL REFERENCES beans(id),
    dose_in      FLOAT8      NOT NULL CHECK (dose_in > 0),
    yield_out    FLOAT8      NOT NULL CHECK (yield_out > 0),
    time         FLOAT8      NOT NULL CHECK (time > 0),
    target_time  FLOAT8      NOT NULL CHECK (target_time > 0),
    grind_size   FLOAT8      NOT NULL CHECK (grind_size > 0),
    pre_infusion BOOLEAN     NOT NULL DEFAULT FALSE,
    tasting_note TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE extraction_gear (
    extraction_id UUID NOT NULL REFERENCES extractions(id) ON DELETE CASCADE,
    gear_id       UUID NOT NULL REFERENCES gear(id),
    PRIMARY KEY (extraction_id, gear_id)
);

-- +goose Down
DROP TABLE extraction_gear;
DROP TABLE extractions;
