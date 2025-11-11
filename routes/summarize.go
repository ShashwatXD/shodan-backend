package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shodan-backend/models"
)

func SummarizeTextHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Text string `json:"text" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}

		client := &http.Client{Timeout: 20 * time.Second}
		body, _ := json.Marshal(map[string]string{"text": payload.Text})
		req, _ := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8000/summarize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact model service", "details": err.Error()})
			return
		}
		defer resp.Body.Close()
		respBytes, _ := io.ReadAll(resp.Body)

		var out interface{}
		if err := json.Unmarshal(respBytes, &out); err != nil {
			c.JSON(http.StatusOK, gin.H{"raw": string(respBytes)})
			return
		}

		outputBytes, _ := json.Marshal(out)
		confidence, model := extractConfidenceAndModel(out)
		r := models.Result{
			Text:       payload.Text,
			Output:     string(outputBytes),
			Confidence: confidence,
			Model:      model,
			CreatedAt:  time.Now(),
		}
		_ = db.Create(&r)

		c.JSON(http.StatusOK, gin.H{"result": out})
	}
}
