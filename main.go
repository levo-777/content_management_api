package main

import (
	"cms-backend/models"
	"cms-backend/routes"
	"cms-backend/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	if env == "development" {
		log.Println("Running AutoMigrate...")
		if err := db.AutoMigrate(&models.Page{}, &models.Post{}, &models.Media{}); err != nil {
			log.Fatalf("Failed to automigrate database: %v", err)
		}
	}

	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	routes.InitializeRoutes(router, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
