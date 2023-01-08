package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pArtour/go-book/database"
	"github.com/pArtour/go-book/routes"
	"log"
	"os"
)

const CollectionName = "books"

// Main is the entry point for the application.
func main() {
	// Load .env file
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get mongo url from env file
	mongoUrl := os.Getenv("MONGO_URL")
	if mongoUrl == "" {
		log.Fatal("MONGO_URL is not set")
	}

	// Get port from env file
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	// Init database connection
	client := database.InitDBConnection(mongoUrl)

	// Open collection
	collection := database.OpenCollection(client, CollectionName)

	// Init router
	router := gin.Default()

	routes.BookRoutes(router, collection)

	err = router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}

}
