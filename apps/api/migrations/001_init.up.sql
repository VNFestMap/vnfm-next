CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT '',
    avatar_url VARCHAR(500) NOT NULL DEFAULT '',
    language_preference VARCHAR(8) NOT NULL DEFAULT 'zh',
    theme_preference VARCHAR(16) NOT NULL DEFAULT 'system',
    display_membership_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS clubs (
    id BIGSERIAL PRIMARY KEY,
    country VARCHAR(16) NOT NULL DEFAULT 'china',
    name VARCHAR(255) NOT NULL,
    school VARCHAR(255) NOT NULL DEFAULT '',
    province VARCHAR(255) NOT NULL DEFAULT '',
    prefecture VARCHAR(255) NOT NULL DEFAULT '',
    city VARCHAR(255) NOT NULL DEFAULT '',
    type VARCHAR(50) NOT NULL DEFAULT 'school',
    info TEXT NOT NULL DEFAULT '',
    remark TEXT NOT NULL DEFAULT '',
    logo_key VARCHAR(500) NOT NULL DEFAULT '',
    external_links TEXT NOT NULL DEFAULT '',
    contact_hidden BOOLEAN NOT NULL DEFAULT TRUE,
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_clubs_country_province ON clubs (country, province);
CREATE INDEX IF NOT EXISTS idx_clubs_name ON clubs (name);

CREATE TABLE IF NOT EXISTS club_memberships (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    club_id BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    country VARCHAR(16) NOT NULL DEFAULT 'china',
    role VARCHAR(32) NOT NULL DEFAULT 'member',
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    contact_account VARCHAR(255) NOT NULL DEFAULT '',
    apply_reason TEXT NOT NULL DEFAULT '',
    apply_role VARCHAR(32) NOT NULL DEFAULT 'member',
    reviewed_by BIGINT,
    reviewed_at TIMESTAMPTZ,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at TIMESTAMPTZ,
    UNIQUE (user_id, club_id)
);

CREATE INDEX IF NOT EXISTS idx_memberships_club_status ON club_memberships (club_id, status);
CREATE INDEX IF NOT EXISTS idx_memberships_user ON club_memberships (user_id);

CREATE TABLE IF NOT EXISTS club_verification_codes (
    id BIGSERIAL PRIMARY KEY,
    club_id BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    code VARCHAR(40) NOT NULL UNIQUE,
    max_uses INT NOT NULL DEFAULT 1,
    use_count INT NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    club_id BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    location VARCHAR(255) NOT NULL DEFAULT '',
    cover_key VARCHAR(500) NOT NULL DEFAULT '',
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ,
    register_until TIMESTAMPTZ,
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_events_starts ON events (starts_at DESC);
CREATE INDEX IF NOT EXISTS idx_events_club ON events (club_id);

CREATE TABLE IF NOT EXISTS event_registrations (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (event_id, user_id)
);

CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL DEFAULT '',
    link VARCHAR(500) NOT NULL DEFAULT '',
    related_type VARCHAR(50) NOT NULL DEFAULT '',
    related_id BIGINT NOT NULL DEFAULT 0,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_notif_user ON notifications (user_id, is_read, created_at DESC);

ALTER TABLE users
    ADD CONSTRAINT users_display_membership_fk
    FOREIGN KEY (display_membership_id) REFERENCES club_memberships(id) ON DELETE SET NULL;
