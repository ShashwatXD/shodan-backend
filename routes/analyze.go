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

func AnalyzeTextHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Text string `json:"text" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}

		client := &http.Client{Timeout: 10 * time.Second}
		body, _ := json.Marshal(map[string]string{"text": payload.Text})
		req, _ := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8000/predict", bytes.NewBuffer(body))
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

func AnalyzeBatchHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var payload struct {
			Texts []string `json:"texts" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
			return
		}

		client := &http.Client{Timeout: 20 * time.Second}
		body, _ := json.Marshal(map[string][]string{"text": payload.Texts})
		req, _ := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8000/predict", bytes.NewBuffer(body))
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

		var items []interface{}
		switch v := out.(type) {
		case []interface{}:
			items = v
		case map[string]interface{}:
			if r, ok := v["results"].([]interface{}); ok {
				items = r
			} else {
				items = []interface{}{v}
			}
		default:
			items = []interface{}{v}
		}

		for i, it := range items {
			var text string
			if i < len(payload.Texts) {
				text = payload.Texts[i]
			}
			outputBytes, _ := json.Marshal(it)
			confidence, model := extractConfidenceAndModel(it)
			r := models.Result{
				Text:       text,
				Output:     string(outputBytes),
				Confidence: confidence,
				Model:      model,
				CreatedAt:  time.Now(),
			}
			_ = db.Create(&r)
		}

		c.JSON(http.StatusOK, gin.H{"results": out})
	}
}

func extractConfidenceAndModel(obj interface{}) (float64, string) {
	var confidence float64
	var model string
	if m, ok := obj.(map[string]interface{}); ok {
		if v, ok := m["confidence"]; ok {
			switch t := v.(type) {
			case float64:
				confidence = t
			case float32:
				confidence = float64(t)
			case int:
				confidence = float64(t)
			}
		}
		if v, ok := m["model"].(string); ok {
			model = v
		}
	}
	return confidence, model
}
