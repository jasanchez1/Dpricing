BEGIN;

CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    auth_id VARCHAR,
    email VARCHAR,
    phone VARCHAR,
    language VARCHAR,
    country VARCHAR
);
 COMMIT;