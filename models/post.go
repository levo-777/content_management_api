package models

import "time"

type Post struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `gorm:"size:255;not null" json:"title" binding:"required"`
    Content   string    `gorm:"type:text;not null" json:"content" binding:"required"`
    Author    string    `gorm:"size:100" json:"author"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Media     []Media   `gorm:"many2many:post_media" json:"media"`
}
