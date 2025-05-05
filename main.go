package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Only load .env file in development environment
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v\n", err)
		}
	}
}

func main() {
	dsn, exists := os.LookupEnv("DSN")
	if !exists {
		log.Fatalf("No DSN Provided\n")
	}
	log.Printf("DSN %s loaded successfully!!", dsn)
	redisHost, exists := os.LookupEnv("REDIS_HOST")
	if !exists {
		log.Fatalf("No Redis Host Provided\n")
	}
	redisPassword, exists := os.LookupEnv("REDIS_PASSWORD")
	if !exists {
		log.Fatalf("No Redis Password Provided\n")
	}
	rsClient, err := getRSClient(redisHost, redisPassword)
	if err != nil {
		log.Fatalf("Error getting Redis Client: %v\n", err)
	}

	var idx Index = newDBIndex(dsn, false, rsClient)
	router := gin.Default()
	if os.Getenv("ENV") == "development" {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://127.0.0.1:3000", "http://localhost:3000"},
			AllowMethods:     []string{"POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Add any other required headers here
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	router.GET("/documents/top10/*any", func(c *gin.Context) {
		corpusHandler(c.Writer, c.Request)
	})

	router.GET("/stats", func(c *gin.Context) {
		type StatRequestBody struct {
			Amount uint
		}

		var statRequestBody StatRequestBody
		if err := c.BindJSON(&statRequestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		stats := idx.getStatResults(statRequestBody.Amount)

		c.IndentedJSON(200, gin.H{"success": "true", "results": stats})
	})

	err = router.Run(":8080")
	if err != nil {
		return
	}
}
