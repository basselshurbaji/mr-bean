-- +goose Up
CREATE TABLE gear (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id),
    type_id    TEXT        NOT NULL,
    name       TEXT        NOT NULL,
    brand      TEXT,
    model      TEXT,
    year       TEXT        CHECK (year ~ '^\d{4}$'),
    notes      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE stations (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id),
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE station_gear (
    station_id UUID    NOT NULL REFERENCES stations(id) ON DELETE CASCADE,
    gear_id    UUID    NOT NULL REFERENCES gear(id),
    position   INTEGER NOT NULL,
    PRIMARY KEY (station_id, gear_id),
    CONSTRAINT uq_station_gear_position UNIQUE (station_id, position)
);

-- +goose Down
DROP TABLE station_gear;
DROP TABLE stations;
DROP TABLE gear;
