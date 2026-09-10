-- Legacy import bookkeeping: original club identity is (id, country),
-- and old user ids are not NextMoe OIDC ids.

ALTER TABLE clubs
    ADD COLUMN IF NOT EXISTS legacy_id BIGINT,
    ADD COLUMN IF NOT EXISTS legacy_country VARCHAR(16) NOT NULL DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_clubs_legacy
    ON clubs (legacy_country, legacy_id)
    WHERE legacy_id IS NOT NULL;

ALTER TABLE events
    ADD COLUMN IF NOT EXISTS legacy_id BIGINT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_events_legacy
    ON events (legacy_id)
    WHERE legacy_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS legacy_user_map (
    legacy_id BIGINT PRIMARY KEY,
    username VARCHAR(255) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL DEFAULT '',
    oidc_id BIGINT REFERENCES users(id),
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_legacy_user_map_email ON legacy_user_map (email)
    WHERE email <> '';
CREATE INDEX IF NOT EXISTS idx_legacy_user_map_oidc ON legacy_user_map (oidc_id)
    WHERE oidc_id IS NOT NULL;
