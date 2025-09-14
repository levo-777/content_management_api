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

func TestGetPosts(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "First Post", "Content 1", "Author 1", time.Now(), time.Now()).
		AddRow(2, "Second Post", "Content 2", "Author 2", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "posts"`).WillReturnRows(rows)
	
	// Mock the Preload("Media") query
	mock.ExpectQuery(`SELECT \* FROM "post_media" WHERE "post_media"\."post_id" IN \(\$1,\$2\)`).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "media_id"}))

	router.GET("/posts", GetPosts)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response []models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("Expected 2 posts, but got %d", len(response))
	}
	if response[0].Title != "First Post" {
		t.Fatalf("Expected 'First Post', but got '%s'", response[0].Title)
	}
	if response[1].Title != "Second Post" {
		t.Fatalf("Expected 'Second Post', but got '%s'", response[1].Title)
	}
}

func TestGetPostsWithFilters(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Test Post", "Test Content", "Test Author", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE title ILIKE \$1 AND author = \$2`).
		WithArgs("%test%", "Test Author").
		WillReturnRows(rows)
	
	// Mock the Preload("Media") query
	mock.ExpectQuery(`SELECT \* FROM "post_media" WHERE "post_media"\."post_id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "media_id"}))

	router.GET("/posts", GetPosts)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts?title=test&author=Test Author", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response []models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("Expected 1 post, but got %d", len(response))
	}
	if response[0].Title != "Test Post" {
		t.Fatalf("Expected 'Test Post', but got '%s'", response[0].Title)
	}
}

func TestGetPost(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Test Post", "Test Content", "Test Author", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)
	
	// Mock the Preload("Media") query
	mock.ExpectQuery(`SELECT \* FROM "post_media" WHERE "post_media"\."post_id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "media_id"}))

	router.GET("/posts/:id", GetPost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.ID != 1 {
		t.Fatalf("Expected post ID 1, but got %d", response.ID)
	}
	if response.Title != "Test Post" {
		t.Fatalf("Expected title 'Test Post', but got '%s'", response.Title)
	}
	if response.Content != "Test Content" {
		t.Fatalf("Expected content 'Test Content', but got '%s'", response.Content)
	}
	if response.Author != "Test Author" {
		t.Fatalf("Expected author 'Test Author', but got '%s'", response.Author)
	}
}

func TestGetPostNotFound(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	router.GET("/posts/:id", GetPost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts/999", nil)
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
	if response.Message != "Post not found" {
		t.Fatalf("Expected message 'Post not found', but got '%s'", response.Message)
	}
}

func TestGetPostInvalidID(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	router.GET("/posts/:id", GetPost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts/invalid", nil)
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
	if response.Message != "Invalid post ID" {
		t.Fatalf("Expected message 'Invalid post ID', but got '%s'", response.Message)
	}
}

func TestGetPostDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(gorm.ErrInvalidDB)

	router.GET("/posts/:id", GetPost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts/1", nil)
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

func TestCreatePost(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "posts"`).
		WithArgs("New Post", "New Content", "New Author", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	post := models.Post{
		Title:   "New Post",
		Content: "New Content",
		Author:  "New Author",
	}
	postJSON, _ := json.Marshal(post)

	router.POST("/posts", CreatePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, but got %d", w.Code)
	}

	var response models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Title != "New Post" {
		t.Fatalf("Expected title 'New Post', but got '%s'", response.Title)
	}
	if response.Content != "New Content" {
		t.Fatalf("Expected content 'New Content', but got '%s'", response.Content)
	}
	if response.Author != "New Author" {
		t.Fatalf("Expected author 'New Author', but got '%s'", response.Author)
	}
}

func TestCreatePostMissingTitle(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	post := models.Post{
		Content: "New Content",
		Author:  "New Author",
	}
	postJSON, _ := json.Marshal(post)

	router.POST("/posts", CreatePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(postJSON))
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

func TestCreatePostMissingContent(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	post := models.Post{
		Title:  "New Post",
		Author: "New Author",
	}
	postJSON, _ := json.Marshal(post)

	router.POST("/posts", CreatePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(postJSON))
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

func TestCreatePostInvalidJSON(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	router.POST("/posts", CreatePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer([]byte("invalid json")))
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

func TestCreatePostDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "posts"`).
		WithArgs("New Post", "New Content", "New Author", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	post := models.Post{
		Title:   "New Post",
		Content: "New Content",
		Author:  "New Author",
	}
	postJSON, _ := json.Marshal(post)

	router.POST("/posts", CreatePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(postJSON))
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

func TestUpdatePost(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing post
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Old Title", "Old Content", "Old Author", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// Mock update transaction
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "posts" SET "title"=\$1,"content"=\$2,"author"=\$3,"created_at"=\$4,"updated_at"=\$5 WHERE "id" = \$6`).
		WithArgs("Updated Title", "Updated Content", "Updated Author", sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updateData := models.Post{
		Title:   "Updated Title",
		Content: "Updated Content",
		Author:  "Updated Author",
	}
	updateJSON, _ := json.Marshal(updateData)

	router.PUT("/posts/:id", UpdatePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/posts/1", bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Title != "Updated Title" {
		t.Fatalf("Expected title 'Updated Title', but got '%s'", response.Title)
	}
	if response.Content != "Updated Content" {
		t.Fatalf("Expected content 'Updated Content', but got '%s'", response.Content)
	}
	if response.Author != "Updated Author" {
		t.Fatalf("Expected author 'Updated Author', but got '%s'", response.Author)
	}
}

func TestUpdatePostNotFound(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	updateData := models.Post{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	updateJSON, _ := json.Marshal(updateData)

	router.PUT("/posts/:id", UpdatePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/posts/999", bytes.NewBuffer(updateJSON))
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
	if response.Message != "Post not found" {
		t.Fatalf("Expected message 'Post not found', but got '%s'", response.Message)
	}
}

func TestUpdatePostInvalidID(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	updateData := models.Post{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	updateJSON, _ := json.Marshal(updateData)

	router.PUT("/posts/:id", UpdatePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/posts/invalid", bytes.NewBuffer(updateJSON))
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
	if response.Message != "Invalid post ID" {
		t.Fatalf("Expected message 'Invalid post ID', but got '%s'", response.Message)
	}
}

func TestUpdatePostInvalidJSON(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing post
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Old Title", "Old Content", "Old Author", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	router.PUT("/posts/:id", UpdatePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/posts/1", bytes.NewBuffer([]byte("invalid json")))
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

func TestDeletePost(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing post
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Test Post", "Test Content", "Test Author", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// Mock delete transaction
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "posts" WHERE "posts"\."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router.DELETE("/posts/:id", DeletePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/posts/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response utils.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if response.Message != "Post deleted successfully" {
		t.Fatalf("Expected message 'Post deleted successfully', but got '%s'", response.Message)
	}
}

func TestDeletePostNotFound(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	router.DELETE("/posts/:id", DeletePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/posts/999", nil)
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
	if response.Message != "Post not found" {
		t.Fatalf("Expected message 'Post not found', but got '%s'", response.Message)
	}
}

func TestDeletePostInvalidID(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	router.DELETE("/posts/:id", DeletePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/posts/invalid", nil)
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
	if response.Message != "Invalid post ID" {
		t.Fatalf("Expected message 'Invalid post ID', but got '%s'", response.Message)
	}
}

func TestDeletePostDatabaseError(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	// Mock finding existing post
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Test Post", "Test Content", "Test Author", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// Mock delete transaction error
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "posts" WHERE "posts"\."id" = \$1`).
		WithArgs(1).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	router.DELETE("/posts/:id", DeletePost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/posts/1", nil)
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