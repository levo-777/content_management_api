-- This migration creates the posts table and post_media junction table

CREATE TABLE posts (
    -- id is the primary key for the table
    id SERIAL PRIMARY KEY,
    -- title is the title of the post
    title VARCHAR(255) NOT NULL,
    -- content is the content of the post
    content TEXT NOT NULL,
    -- author is the author of the post
    author VARCHAR(100),
    -- created_at is the timestamp when the post was created
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- updated_at is the timestamp when the post was last updated
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create post_media junction table for many-to-many relationship
CREATE TABLE post_media (
    -- post_id references the posts table
    post_id INTEGER NOT NULL,
    -- media_id references the media table
    media_id INTEGER NOT NULL,
    -- created_at is the timestamp when the relationship was created
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- Primary key is composite of post_id and media_id
    PRIMARY KEY (post_id, media_id),
    -- Foreign key constraints with cascade delete
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE CASCADE
);

-- Add indexes for performance
CREATE INDEX idx_posts_author ON posts(author);
CREATE INDEX idx_posts_created_at ON posts(created_at);
CREATE INDEX idx_posts_title ON posts(title);
CREATE INDEX idx_post_media_post_id ON post_media(post_id);
CREATE INDEX idx_post_media_media_id ON post_media(media_id);
