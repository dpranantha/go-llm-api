package see

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		defaultStream := true
		p.Stream = &defaultStream
	}
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Non standard SSE with POST due to long context is possible
func HandleLLMStreamSSE(c *gin.Context) {
	// Read the prompt from query or body
	var p PromptRequest
	if err := c.BindJSON(&p); err != nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	p.SetDefaults()

	body, _ := json.Marshal(OllamaRequest{
		Model:  *p.Model,
		Prompt: p.Prompt,
		Stream: *p.Stream,
	})

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintf(c.Writer, "data: %s\n\n", "[Error contacting LLM]")
		c.Writer.Flush()
		return
	}
	defer resp.Body.Close()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()
	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		var chunk map[string]any
		if err := decoder.Decode(&chunk); err != nil {
			break
		}
		if part, ok := chunk["response"].(string); ok {
			fmt.Fprintf(c.Writer, "data: %s\n\n", part)
			c.Writer.Flush()
		}
	}

	fmt.Fprintf(c.Writer, "data: [done]\n\n")
	c.Writer.Flush()
}
