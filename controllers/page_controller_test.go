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

func TestGetPages(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "First Page", "Content 1", time.Now(), time.Now()).
		AddRow(2, "Second Page", "Content 2", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages"`).WillReturnRows(rows)

	router.GET("/pages", GetPages)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/pages", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response []models.Page
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("Expected 2 pages, but got %d", len(response))
	}
}

func TestGetPage(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Test Page", "Test Content", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	router.GET("/pages/:id", GetPage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/pages/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response models.Page
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.ID != 1 {
		t.Fatalf("Expected page ID 1, but got %d", response.ID)
	}
	if response.Title != "Test Page" {
		t.Fatalf("Expected title 'Test Page', but got '%s'", response.Title)
	}
}

func TestCreatePage(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "pages"`).
		WithArgs("New Page", "New Content", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	page := models.Page{
		Title:   "New Page",
		Content: "New Content",
	}
	pageJSON, _ := json.Marshal(page)

	router.POST("/pages", CreatePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer(pageJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, but got %d", w.Code)
	}

	var response models.Page
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Title != "New Page" {
		t.Fatalf("Expected title 'New Page', but got '%s'", response.Title)
	}
}

func TestUpdatePage(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing page
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Old Title", "Old Content", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// Mock update transaction
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "pages" SET "title"=\$1,"content"=\$2,"created_at"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
		WithArgs("Updated Title", "Updated Content", sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updateData := models.Page{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	updateJSON, _ := json.Marshal(updateData)

	router.PUT("/pages/:id", UpdatePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/pages/1", bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response models.Page
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Title != "Updated Title" {
		t.Fatalf("Expected title 'Updated Title', but got '%s'", response.Title)
	}
}

func TestDeletePage(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing page
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Test Page", "Test Content", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// Mock delete transaction
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "pages" WHERE "pages"\."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router.DELETE("/pages/:id", DeletePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/pages/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response utils.MessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Message != "Page deleted successfully" {
		t.Fatalf("Expected deletion message, but got '%s'", response.Message)
	}
}

// Error condition tests

func TestGetPagesDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "pages"`).WillReturnError(gorm.ErrInvalidDB)

	router.GET("/pages", GetPages)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/pages", nil)
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

func TestGetPageNotFound(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	router.GET("/pages/:id", GetPage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/pages/999", nil)
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
	if response.Message != "Page not found" {
		t.Fatalf("Expected message 'Page not found', but got '%s'", response.Message)
	}
}

func TestGetPageInvalidID(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	router.GET("/pages/:id", GetPage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/pages/invalid", nil)
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
	if response.Message != "Invalid page ID" {
		t.Fatalf("Expected message 'Invalid page ID', but got '%s'", response.Message)
	}
}

func TestGetPageDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(gorm.ErrInvalidDB)

	router.GET("/pages/:id", GetPage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/pages/1", nil)
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

func TestCreatePageInvalidJSON(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	router.POST("/pages", CreatePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer([]byte("invalid json")))
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

func TestCreatePageDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "pages"`).
		WithArgs("New Page", "New Content", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	page := models.Page{
		Title:   "New Page",
		Content: "New Content",
	}
	pageJSON, _ := json.Marshal(page)

	router.POST("/pages", CreatePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer(pageJSON))
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

func TestUpdatePageNotFound(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	updateData := models.Page{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	updateJSON, _ := json.Marshal(updateData)

	router.PUT("/pages/:id", UpdatePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/pages/999", bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
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
	if response.Message != "Page not found" {
		t.Fatalf("Expected message 'Page not found', but got '%s'", response.Message)
	}
}

func TestUpdatePageInvalidID(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	updateData := models.Page{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	updateJSON, _ := json.Marshal(updateData)

	router.PUT("/pages/:id", UpdatePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/pages/invalid", bytes.NewBuffer(updateJSON))
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
	if response.Message != "Invalid page ID" {
		t.Fatalf("Expected message 'Invalid page ID', but got '%s'", response.Message)
	}
}

func TestUpdatePageInvalidJSON(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing page
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Old Title", "Old Content", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	router.PUT("/pages/:id", UpdatePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/pages/1", bytes.NewBuffer([]byte("invalid json")))
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

func TestUpdatePageDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing page
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Old Title", "Old Content", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// Mock update transaction error
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "pages" SET "title"=\$1,"content"=\$2,"created_at"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
		WithArgs("Updated Title", "Updated Content", sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	updateData := models.Page{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	updateJSON, _ := json.Marshal(updateData)

	router.PUT("/pages/:id", UpdatePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/pages/1", bytes.NewBuffer(updateJSON))
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

func TestDeletePageNotFound(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	router.DELETE("/pages/:id", DeletePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/pages/999", nil)
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
	if response.Message != "Page not found" {
		t.Fatalf("Expected message 'Page not found', but got '%s'", response.Message)
	}
}

func TestDeletePageInvalidID(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	router.DELETE("/pages/:id", DeletePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/pages/invalid", nil)
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
	if response.Message != "Invalid page ID" {
		t.Fatalf("Expected message 'Invalid page ID', but got '%s'", response.Message)
	}
}

func TestDeletePageDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing page
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Test Page", "Test Content", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// Mock delete transaction error
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "pages" WHERE "pages"\."id" = \$1`).
		WithArgs(1).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	router.DELETE("/pages/:id", DeletePage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/pages/1", nil)
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
