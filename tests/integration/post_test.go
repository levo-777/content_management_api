package integration

import (
	"cms-backend/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)



func TestPostIntegration(t *testing.T) {
	
	clearTables()

	
	_ = createTestMedia(t)

	t.Run("Create Post with Media", func(t *testing.T) {
		
		postBody := `{
			"title": "Test Post with Media",
			"content": "This is a test post with media attachment",
			"author": "Test Author"
		}`

		req := httptest.NewRequest("POST", "/api/v1/posts", strings.NewReader(postBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status 201, got %d: %s", w.Code, w.Body.String())
		}

		var response models.Post
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Title != "Test Post with Media" {
			t.Errorf("Expected title 'Test Post with Media', got %s", response.Title)
		}

		if response.Author != "Test Author" {
			t.Errorf("Expected author 'Test Author', got %s", response.Author)
		}
	})

	t.Run("Get Posts with Filter", func(t *testing.T) {
		
		clearTables()

		
		post1 := models.Post{
			Title:   "Filtered Post 1",
			Content: "Content 1",
			Author:  "Author A",
		}
		post2 := models.Post{
			Title:   "Filtered Post 2",
			Content: "Content 2",
			Author:  "Author B",
		}

		testDB.Create(&post1)
		testDB.Create(&post2)

		
		req := httptest.NewRequest("GET", "/api/v1/posts?author=Author%20A", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		var response []models.Post
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response) != 1 {
			t.Errorf("Expected 1 post, got %d", len(response))
		}

		if response[0].Author != "Author A" {
			t.Errorf("Expected author 'Author A', got %s", response[0].Author)
		}

		// Test filtering by title
		req = httptest.NewRequest("GET", "/api/v1/posts?title=Filtered", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response) != 2 {
			t.Errorf("Expected 2 posts, got %d", len(response))
		}
	})
}

// Helper function to create test media
func createTestMedia(t *testing.T) uint {
	body := `{
		"url": "http://example.com/test.jpg",
		"type": "image"
	}`

	req := httptest.NewRequest("POST", "/api/v1/media", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create test media, status: %d, body: %s", w.Code, w.Body.String())
	}

	var response models.Media
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to create test media: %v", err)
	}

	return response.ID
}


