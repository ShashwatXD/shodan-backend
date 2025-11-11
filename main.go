package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"shodan-backend/database"
	"shodan-backend/routes"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}

	r := gin.Default()

	// Register routes
	r.POST("/analyze/text", routes.AnalyzeTextHandler(db))
	r.POST("/analyze/batch", routes.AnalyzeBatchHandler(db))
	r.POST("/summarize", routes.SummarizeTextHandler(db))
	r.GET("/history", func(c *gin.Context) {
		var results []interface{}
		var rs []map[string]interface{}
		if err := db.Raw("SELECT id, text, output, confidence, model, created_at FROM results").Scan(&rs).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch history", "details": err.Error()})
			return
		}
		for _, r := range rs {
			results = append(results, r)
		}
		c.JSON(200, gin.H{"results": results})
	})

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
