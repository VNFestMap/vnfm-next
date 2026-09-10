DROP INDEX IF EXISTS idx_legacy_user_map_oidc;
DROP INDEX IF EXISTS idx_legacy_user_map_email;
DROP TABLE IF EXISTS legacy_user_map;

DROP INDEX IF EXISTS idx_events_legacy;
ALTER TABLE events DROP COLUMN IF EXISTS legacy_id;

DROP INDEX IF EXISTS idx_clubs_legacy;
ALTER TABLE clubs DROP COLUMN IF EXISTS legacy_id;
ALTER TABLE clubs DROP COLUMN IF EXISTS legacy_country;
