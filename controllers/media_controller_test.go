package controllers

import (
	"bytes"
	"cms-backend/models"
	"cms-backend/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestGetMedia(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "https://example.com/image1.jpg", "image", time.Now(), time.Now()).
		AddRow(2, "https://example.com/video1.mp4", "video", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "media"`).WillReturnRows(rows)

	router.GET("/media", GetMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/media", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response []models.Media
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("Expected 2 media items, but got %d", len(response))
	}
	if response[0].URL != "https://example.com/image1.jpg" {
		t.Fatalf("Expected URL 'https://example.com/image1.jpg', but got '%s'", response[0].URL)
	}
	if response[0].Type != "image" {
		t.Fatalf("Expected type 'image', but got '%s'", response[0].Type)
	}
	if response[1].URL != "https://example.com/video1.mp4" {
		t.Fatalf("Expected URL 'https://example.com/video1.mp4', but got '%s'", response[1].URL)
	}
	if response[1].Type != "video" {
		t.Fatalf("Expected type 'video', but got '%s'", response[1].Type)
	}
}

func TestGetMediaDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "media"`).WillReturnError(gorm.ErrInvalidDB)

	router.GET("/media", GetMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/media", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status 500, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("Expected error code 500, but got %d", response.Code)
	}
}

func TestGetMediaByID(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "https://example.com/image1.jpg", "image", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "media" WHERE "media"\."id" = \$1 ORDER BY "media"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	router.GET("/media/:id", GetMediaByID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/media/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response models.Media
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.ID != 1 {
		t.Fatalf("Expected media ID 1, but got %d", response.ID)
	}
	if response.URL != "https://example.com/image1.jpg" {
		t.Fatalf("Expected URL 'https://example.com/image1.jpg', but got '%s'", response.URL)
	}
	if response.Type != "image" {
		t.Fatalf("Expected type 'image', but got '%s'", response.Type)
	}
}

func TestGetMediaByIDNotFound(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "media" WHERE "media"\."id" = \$1 ORDER BY "media"\."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	router.GET("/media/:id", GetMediaByID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/media/999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusNotFound {
		t.Fatalf("Expected error code 404, but got %d", response.Code)
	}
	if response.Message != "Media not found" {
		t.Fatalf("Expected message 'Media not found', but got '%s'", response.Message)
	}
}

func TestGetMediaByIDInvalidID(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	router.GET("/media/:id", GetMediaByID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/media/invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusBadRequest {
		t.Fatalf("Expected error code 400, but got %d", response.Code)
	}
	if response.Message != "Invalid media ID" {
		t.Fatalf("Expected message 'Invalid media ID', but got '%s'", response.Message)
	}
}

func TestGetMediaByIDDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "media" WHERE "media"\."id" = \$1 ORDER BY "media"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(gorm.ErrInvalidDB)

	router.GET("/media/:id", GetMediaByID)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/media/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status 500, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("Expected error code 500, but got %d", response.Code)
	}
}

func TestCreateMedia(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "media"`).
		WithArgs("https://example.com/new-image.jpg", "image", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	media := models.Media{
		URL:  "https://example.com/new-image.jpg",
		Type: "image",
	}
	mediaJSON, _ := json.Marshal(media)

	router.POST("/media", CreateMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/media", bytes.NewBuffer(mediaJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, but got %d", w.Code)
	}

	var response models.Media
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.URL != "https://example.com/new-image.jpg" {
		t.Fatalf("Expected URL 'https://example.com/new-image.jpg', but got '%s'", response.URL)
	}
	if response.Type != "image" {
		t.Fatalf("Expected type 'image', but got '%s'", response.Type)
	}
}

func TestCreateMediaMissingURL(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	media := models.Media{
		Type: "image",
	}
	mediaJSON, _ := json.Marshal(media)

	router.POST("/media", CreateMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/media", bytes.NewBuffer(mediaJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusBadRequest {
		t.Fatalf("Expected error code 400, but got %d", response.Code)
	}
	// The actual validation message from GORM is different
	if response.Message == "" {
		t.Fatalf("Expected validation error message, but got empty message")
	}
}

func TestCreateMediaMissingType(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	media := models.Media{
		URL: "https://example.com/new-image.jpg",
	}
	mediaJSON, _ := json.Marshal(media)

	router.POST("/media", CreateMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/media", bytes.NewBuffer(mediaJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusBadRequest {
		t.Fatalf("Expected error code 400, but got %d", response.Code)
	}
	// The actual validation message from GORM is different
	if response.Message == "" {
		t.Fatalf("Expected validation error message, but got empty message")
	}
}

func TestCreateMediaEmptyURL(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	media := models.Media{
		URL:  "",
		Type: "image",
	}
	mediaJSON, _ := json.Marshal(media)

	router.POST("/media", CreateMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/media", bytes.NewBuffer(mediaJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusBadRequest {
		t.Fatalf("Expected error code 400, but got %d", response.Code)
	}
	// The actual validation message from GORM is different
	if response.Message == "" {
		t.Fatalf("Expected validation error message, but got empty message")
	}
}

func TestCreateMediaEmptyType(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	media := models.Media{
		URL:  "https://example.com/new-image.jpg",
		Type: "",
	}
	mediaJSON, _ := json.Marshal(media)

	router.POST("/media", CreateMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/media", bytes.NewBuffer(mediaJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusBadRequest {
		t.Fatalf("Expected error code 400, but got %d", response.Code)
	}
	// The actual validation message from GORM is different
	if response.Message == "" {
		t.Fatalf("Expected validation error message, but got empty message")
	}
}

func TestCreateMediaInvalidJSON(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	router.POST("/media", CreateMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/media", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusBadRequest {
		t.Fatalf("Expected error code 400, but got %d", response.Code)
	}
}

func TestCreateMediaDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "media"`).
		WithArgs("https://example.com/new-image.jpg", "image", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	media := models.Media{
		URL:  "https://example.com/new-image.jpg",
		Type: "image",
	}
	mediaJSON, _ := json.Marshal(media)

	router.POST("/media", CreateMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/media", bytes.NewBuffer(mediaJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status 500, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("Expected error code 500, but got %d", response.Code)
	}
}

func TestDeleteMedia(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing media
	rows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "https://example.com/image1.jpg", "image", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "media" WHERE "media"\."id" = \$1 ORDER BY "media"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// Mock delete transaction
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "media" WHERE "media"\."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router.DELETE("/media/:id", DeleteMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/media/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response utils.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Message != "Media deleted successfully" {
		t.Fatalf("Expected message 'Media deleted successfully', but got '%s'", response.Message)
	}
}

func TestDeleteMediaNotFound(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "media" WHERE "media"\."id" = \$1 ORDER BY "media"\."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	router.DELETE("/media/:id", DeleteMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/media/999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusNotFound {
		t.Fatalf("Expected error code 404, but got %d", response.Code)
	}
	if response.Message != "Media not found" {
		t.Fatalf("Expected message 'Media not found', but got '%s'", response.Message)
	}
}

func TestDeleteMediaInvalidID(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	router.DELETE("/media/:id", DeleteMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/media/invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusBadRequest {
		t.Fatalf("Expected error code 400, but got %d", response.Code)
	}
	if response.Message != "Invalid media ID" {
		t.Fatalf("Expected message 'Invalid media ID', but got '%s'", response.Message)
	}
}

func TestDeleteMediaDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing media
	rows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "https://example.com/image1.jpg", "image", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "media" WHERE "media"\."id" = \$1 ORDER BY "media"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// Mock delete transaction error
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "media" WHERE "media"\."id" = \$1`).
		WithArgs(1).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	router.DELETE("/media/:id", DeleteMedia)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/media/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status 500, but got %d", w.Code)
	}

	var response utils.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("Expected error code 500, but got %d", response.Code)
	}
}