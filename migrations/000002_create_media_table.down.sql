-- This migration drops the media table

-- Drop indexes first
DROP INDEX IF EXISTS idx_media_created_at;
DROP INDEX IF EXISTS idx_media_type;

-- Drop the media table
DROP TABLE IF EXISTS media;
