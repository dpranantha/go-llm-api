package ws

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/websocket/v2"
)

func HandleLLMStreamFiber(c *websocket.Conn) {
	log.Println("WebSocket connection opened")

	for {
		// Read message from client
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Client disconnected or read error:", err)
			break
		}
		prompt := string(message)
		log.Println("Received prompt:", prompt)

		// Prepare request to Ollama
		body, _ := json.Marshal(map[string]any{
			"model":  "llama2",
			"prompt": prompt,
			"stream": true,
		})

		resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Println("Error contacting Ollama:", err)
			c.WriteMessage(websocket.TextMessage, []byte("[Error contacting LLM]"))
			continue // let the user try again without closing the socket
		}
		defer resp.Body.Close()

		// Stream response back to client
		decoder := json.NewDecoder(resp.Body)
		for decoder.More() {
			var chunk map[string]any
			if err := decoder.Decode(&chunk); err != nil {
				log.Println("Decode error:", err)
				break
			}
			if part, ok := chunk["response"].(string); ok {
				err := c.WriteMessage(websocket.TextMessage, []byte(part))
				if err != nil {
					log.Println("Send error:", err)
					break
				}
			}
		}

		// Mark end of streaming
		err = c.WriteMessage(websocket.TextMessage, []byte("[done]"))
		if err != nil {
			log.Println("Error sending done:", err)
			break
		}
	}

	log.Println("WebSocket connection closed")
	c.Close()
}
