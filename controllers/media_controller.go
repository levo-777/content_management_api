package controllers

import (
    "cms-backend/models"
    "cms-backend/utils"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func GetMedia(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var media []models.Media

    if err := db.Find(&media).Error; err != nil {
        c.JSON(http.StatusInternalServerError, utils.HTTPError{
            Code:    http.StatusInternalServerError,
            Message: err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, media)
}

func GetMediaByID(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "Invalid media ID",
        })
        return
    }

    var media models.Media
    if err := db.First(&media, uint(id)).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, utils.HTTPError{
                Code:    http.StatusNotFound,
                Message: "Media not found",
            })
        } else {
            c.JSON(http.StatusInternalServerError, utils.HTTPError{
                Code:    http.StatusInternalServerError,
                Message: err.Error(),
            })
        }
        return
    }
    c.JSON(http.StatusOK, media)
}

func CreateMedia(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    
    var media models.Media
    if err := c.ShouldBindJSON(&media); err != nil {
        c.JSON(http.StatusBadRequest, utils.HTTPError{
            Code:    http.StatusBadRequest,
            Message: err.Error(),
        })
        return
    }

    if media.URL == "" || media.Type == "" {
        c.JSON(http.StatusBadRequest, utils.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "URL and type are required",
        })
        return
    }

    tx := db.Begin()
    if err := tx.Create(&media).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, utils.HTTPError{
            Code:    http.StatusInternalServerError,
            Message: err.Error(),
        })
        return
    }
    tx.Commit()
    c.JSON(http.StatusCreated, media)
}

func DeleteMedia(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.HTTPError{
            Code:    http.StatusBadRequest,
            Message: "Invalid media ID",
        })
        return
    }

    var media models.Media
    if err := db.First(&media, uint(id)).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, utils.HTTPError{
                Code:    http.StatusNotFound,
                Message: "Media not found",
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
    if err := tx.Delete(&media).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, utils.HTTPError{
            Code:    http.StatusInternalServerError,
            Message: err.Error(),
        })
        return
    }
    tx.Commit()
    c.JSON(http.StatusOK, utils.MessageResponse{
        Message: "Media deleted successfully",
    })
}

