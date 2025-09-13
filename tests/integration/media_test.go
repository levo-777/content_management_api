package integration

import (
	"cms-backend/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)



func TestMediaIntegration(t *testing.T) {
	
	clearTables()

	t.Run("Create Media", func(t *testing.T) {
		
		body := `{
			"url": "http://example.com/test.jpg",
			"type": "image"
		}`

		
		req := httptest.NewRequest("POST", "/api/v1/media", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status 201, got %d: %s", w.Code, w.Body.String())
		}

		
		var response models.Media
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.URL != "http://example.com/test.jpg" {
			t.Errorf("Expected URL 'http://example.com/test.jpg', got %s", response.URL)
		}
	})

    t.Run("Get All Media", func(t *testing.T) {
		
		clearTables()
		
		
		media1 := models.Media{
			URL:  "http://example.com/test1.jpg",
			Type: "image",
		}
		media2 := models.Media{
			URL:  "http://example.com/test2.mp4",
			Type: "video",
		}
		
		testDB.Create(&media1)
		testDB.Create(&media2)
		
		
		req := httptest.NewRequest("GET", "/api/v1/media", nil)
		w := httptest.NewRecorder()
		
		
		router.ServeHTTP(w, req)
		
		
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}
		
		var response []models.Media
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		
		if len(response) != 2 {
			t.Errorf("Expected 2 media items, got %d", len(response))
		}
		
		
		foundImage := false
		foundVideo := false
		for _, media := range response {
			if media.Type == "image" && media.URL == "http://example.com/test1.jpg" {
				foundImage = true
			}
			if media.Type == "video" && media.URL == "http://example.com/test2.mp4" {
				foundVideo = true
			}
		}
		
		if !foundImage {
			t.Error("Expected to find image media item")
		}
		if !foundVideo {
			t.Error("Expected to find video media item")
		}
	})

  
}

