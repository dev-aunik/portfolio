-- 0001_create_articles.down.sql
DROP TRIGGER IF EXISTS articles_updated_at ON articles;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS articles;
