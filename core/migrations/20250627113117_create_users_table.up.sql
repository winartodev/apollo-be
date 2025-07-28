CREATE TABLE IF NOT EXISTS users
(
    id                SERIAL PRIMARY KEY,
    username          VARCHAR(225) NOT NULL UNIQUE,
    email             VARCHAR(255) NOT NULL UNIQUE,
    phone_number      VARCHAR(16)  NOT NULL UNIQUE,
    first_name        VARCHAR(100) NOT NULL DEFAULT '',
    last_name         VARCHAR(100) NOT NULL DEFAULT '',
    is_active         BOOLEAN      NOT NULL DEFAULT TRUE,
    is_email_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    is_phone_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    password          VARCHAR(255) NOT NULL,
    refresh_token     VARCHAR(255),
    last_login        BIGINT       NULL,
    created_at        BIGINT       NOT NULL,
    updated_at        BIGINT       NULL,
    deleted_at        BIGINT       NULL
);

CREATE INDEX idx_users_is_active ON users (is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_email_verified ON users (is_email_verified);
CREATE INDEX idx_users_phone_verified ON users (is_phone_verified);