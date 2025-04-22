package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PromptRequest struct {
	Model  *string `json:"model"`
	Prompt string  `json:"prompt"`
	Stream *bool   `json:"stream"`
}

func (p *PromptRequest) SetDefaults() {
	if p.Model == nil {
		defaultModel := "llama2"
		p.Model = &defaultModel
	}
	if p.Stream == nil {
		defaultStream := false
		p.Stream = &defaultStream
	}
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func HandlePrompt(c *gin.Context) {
	// Parse request from GraphQL or REST API
	var p PromptRequest
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	p.SetDefaults()

	// Prepare payload for Ollama
	body, _ := json.Marshal(OllamaRequest{
		Model:  *p.Model,
		Prompt: p.Prompt,
		Stream: *p.Stream,
	})

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Failed to call Ollama:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LLM not reachable"})
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var fullResponse string

	for decoder.More() {
		var chunk map[string]any
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
