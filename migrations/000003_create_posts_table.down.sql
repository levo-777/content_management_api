-- This migration drops the posts table and post_media junction table

-- Drop indexes first
DROP INDEX IF EXISTS idx_post_media_media_id;
DROP INDEX IF EXISTS idx_post_media_post_id;
DROP INDEX IF EXISTS idx_posts_title;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_author;

-- Drop post_media junction table first (due to foreign key constraints)
DROP TABLE IF EXISTS post_media;

-- Drop posts table after post_media is removed
DROP TABLE IF EXISTS posts;