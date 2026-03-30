-- 0001_create_articles.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS articles (
    id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    title        TEXT        NOT NULL,
    slug         TEXT        NOT NULL UNIQUE,
    summary      TEXT        NOT NULL,
    content      TEXT        NOT NULL,
    tags         TEXT[]      NOT NULL DEFAULT '{}',
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_articles_slug         ON articles (slug);
CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles (published_at DESC NULLS LAST);
CREATE INDEX IF NOT EXISTS idx_articles_tags         ON articles USING gin (tags);

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER articles_updated_at
    BEFORE UPDATE ON articles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
