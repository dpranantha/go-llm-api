package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Prompt struct {
	Prompt string `json:"prompt"`
}

func HandlePrompt(c *fiber.Ctx) error {
	var p Prompt
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
	}

	// Prepare payload for Ollama
	reqBody := map[string]any{
		"model":  "llama2",
		"prompt": p.Prompt,
		"stream": true,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error calling Ollama:", err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to call LLM"})
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var fullResponse string

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

	return c.JSON(fiber.Map{"response": fullResponse})
}
