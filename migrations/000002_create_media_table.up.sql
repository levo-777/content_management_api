-- This migration creates the media table

CREATE TABLE media (
    -- id is the primary key for the table
    id SERIAL PRIMARY KEY,
    -- url is the file location/path
    url VARCHAR(255) NOT NULL,
    -- type identifies the media type (image, video, etc.)
    type VARCHAR(50) NOT NULL,
    -- created_at is the timestamp when the media was created
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- updated_at is the timestamp when the media was last updated
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for performance
CREATE INDEX idx_media_type ON media(type);
CREATE INDEX idx_media_created_at ON media(created_at);