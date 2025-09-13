package controllers

import (
	"cms-backend/models"
	"cms-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPages(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var pages []models.Page

	if err := db.Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, pages)
}

func GetPage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid page ID",
		})
		return
	}

	var page models.Page
	if err := db.First(&page, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, utils.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Page not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, utils.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, page)
}
func CreatePage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var page models.Page

	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	tx := db.Begin()
	if err := tx.Create(&page).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	tx.Commit()
	c.JSON(http.StatusCreated, page)
}

func UpdatePage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid page ID",
		})
		return
	}

	var page models.Page
	if err := db.First(&page, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, utils.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Page not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, utils.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}
		return
	}

	var updateData models.Page
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	page.Title = updateData.Title
	page.Content = updateData.Content

	tx := db.Begin()
	if err := tx.Save(&page).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, page)
}

func DeletePage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid page ID",
		})
		return
	}

	var page models.Page
	if err := db.First(&page, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, utils.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Page not found",
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
	if err := tx.Delete(&page).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, utils.MessageResponse{
		Message: "Page deleted successfully",
	})
}
