-- 0002_create_contacts.up.sql
CREATE TABLE IF NOT EXISTS contacts (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL,
    subject    TEXT        NOT NULL DEFAULT '',
    message    TEXT        NOT NULL,
    status     TEXT        NOT NULL DEFAULT 'pending'
                           CHECK (status IN ('pending', 'processed', 'failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contacts_status     ON contacts (status);
CREATE INDEX IF NOT EXISTS idx_contacts_created_at ON contacts (created_at DESC);
