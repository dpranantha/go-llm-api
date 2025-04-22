package see

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
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
		defaultStream := true
		p.Stream = &defaultStream
	}
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func HandlePromptSSE(c *fiber.Ctx) error {
	// Parse request from GraphQL or REST API
	var p PromptRequest
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
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
		log.Println("Error calling Ollama:", err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to call LLM"})
	}
	defer resp.Body.Close()

	// Set headers for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	writeSSEData := func(data string) {
		c.Response().BodyWriter().Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
	}

	decoder := json.NewDecoder(resp.Body)

	// Write data to client and flush
	for decoder.More() {
		var chunk map[string]any
		if err := decoder.Decode(&chunk); err != nil {
			break
		}

		if part, ok := chunk["response"].(string); ok {
			// Write the chunk to the response
			writeSSEData(part) // Write the data to the SSE stream
		}
	}

	// Send final message indicating stream completion
	writeSSEData("[done]")
	return nil
}
