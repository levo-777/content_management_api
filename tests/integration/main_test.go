package integration

import (
	"cms-backend/models"
	"cms-backend/routes"
	"cms-backend/utils"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)


var (
	testDB *gorm.DB
	router *gin.Engine
)



func TestMain(m *testing.M) {
	
	setup()

	
	code := m.Run()

	
	cleanup()

	os.Exit(code)
}

func setup() {
	
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	
	gin.SetMode(gin.TestMode)

	
	var err error
	testDB, err = utils.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	
	if err := testDB.AutoMigrate(&models.Media{}, &models.Page{}, &models.Post{}); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	
	router = gin.New()
	routes.InitializeRoutes(router, testDB)
}

func cleanup() {
	if testDB != nil {
		
		sqlDB, err := testDB.DB()
		if err != nil {
			log.Printf("Error getting SQL database: %v", err)
			return
		}

		
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
}

func clearTables() {
	if testDB != nil {
		
		testDB.Exec("DELETE FROM post_media")
		testDB.Exec("DELETE FROM posts")
		testDB.Exec("DELETE FROM media")
		testDB.Exec("DELETE FROM pages")
	}
}


