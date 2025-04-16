package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func HandlePrompt(c *gin.Context) {
	var req PromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	body, _ := json.Marshal(OllamaRequest{
		Model:  "llama2",
		Prompt: req.Prompt,
	})

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Failed to call Ollama:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LLM not reachable"})
		return
	}
	defer resp.Body.Close()

	var fullResponse string
	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err != nil {
			log.Println("Decode error:", err)
			break
		}
		if part, ok := chunk["response"].(string); ok {
			fullResponse += part
		}
	}

	c.JSON(http.StatusOK, gin.H{"response": fullResponse})
}
