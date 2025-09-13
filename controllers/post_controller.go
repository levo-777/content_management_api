package controllers

import (
	"cms-backend/models"
	"cms-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPosts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var posts []models.Post

	title := c.Query("title")
	author := c.Query("author")

	query := db
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if author != "" {
		query = query.Where("author = ?", author)
	}

	if err := query.Preload("Media").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "Invalid post ID",
        })
        return
    }

    var post models.Post
    if err := db.Preload("Media").First(&post, uint(id)).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, utils.HTTPError{
                Code:    http.StatusNotFound,
                Message: "Post not found",
            })
        } else {
            c.JSON(http.StatusInternalServerError, utils.HTTPError{
                Code:    http.StatusInternalServerError,
                Message: err.Error(),
            })
        }
        return
    }
    c.JSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    
    var post models.Post
    if err := c.ShouldBindJSON(&post); err != nil {
        c.JSON(http.StatusBadRequest, utils.HTTPError{
            Code:    http.StatusBadRequest,
            Message: err.Error(),
        })
        return
    }

    if post.Title == "" || post.Content == "" {
        c.JSON(http.StatusBadRequest, utils.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "Title and content are required",
        })
        return
    }

    tx := db.Begin()
    if err := tx.Create(&post).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, utils.HTTPError{
            Code:    http.StatusInternalServerError,
            Message: err.Error(),
        })
        return
    }
    tx.Commit()
    c.JSON(http.StatusCreated, post)
}

func UpdatePost(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "Invalid post ID",
        })
        return
    }

    var post models.Post
    if err := db.First(&post, uint(id)).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, utils.HTTPError{
                Code:    http.StatusNotFound,
                Message: "Post not found",
            })
        } else {
            c.JSON(http.StatusInternalServerError, utils.HTTPError{
                Code:    http.StatusInternalServerError,
                Message: err.Error(),
            })
        }
        return
    }

    var updateData models.Post
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, utils.HTTPError{
            Code:    http.StatusBadRequest,
            Message: err.Error(),
        })
        return
    }

    if updateData.Title != "" {
        post.Title = updateData.Title
    }
    if updateData.Content != "" {
        post.Content = updateData.Content
    }
    if updateData.Author != "" {
        post.Author = updateData.Author
    }

    tx := db.Begin()
    if err := tx.Save(&post).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, utils.HTTPError{
            Code:    http.StatusInternalServerError,
            Message: err.Error(),
        })
        return
    }
    tx.Commit()
    c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "Invalid post ID",
        })
        return
    }

    var post models.Post
    if err := db.First(&post, uint(id)).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, utils.HTTPError{
                Code:    http.StatusNotFound,
                Message: "Post not found",
            })
        } else {
            c.JSON(http.StatusInternalServerError, utils.HTTPError{
                Code:    http.StatusInternalServerError,
                Message: err.Error(),
            })
        }
        return
    }

    tx := db.Begin()
    if err := tx.Delete(&post).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, utils.HTTPError{
            Code:    http.StatusInternalServerError,
            Message: err.Error(),
        })
        return
    }
    tx.Commit()
    c.JSON(http.StatusOK, utils.MessageResponse{
        Message: "Post deleted successfully",
    })
}