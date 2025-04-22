package ws

import (
	"log"
	"net/http"

	"bytes"
	"encoding/json"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins — update this in production
		return true
	},
}

func HandleLLMStream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	// Don't defer conn.Close() here

	log.Println("WebSocket connected")

	for {
		// Wait for message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Client disconnected or error reading:", err)
			break // exit loop and close the connection
		}

		prompt := string(message)
		log.Println("Received prompt:", prompt)

		// Send request to LLM
		body, _ := json.Marshal(map[string]any{
			"model":  "llama2",
			"prompt": prompt,
			"stream": true,
		})

		resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Println("Error contacting Ollama:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("[Error contacting LLM]"))
			continue
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for decoder.More() {
			var chunk map[string]any
			if err := decoder.Decode(&chunk); err != nil {
				log.Println("Error decoding chunk:", err)
				break
			}
			if part, ok := chunk["response"].(string); ok {
				conn.WriteMessage(websocket.TextMessage, []byte(part))
			}
		}

		// Signal the end of the response
		conn.WriteMessage(websocket.TextMessage, []byte("[done]"))
	}

	log.Println("Closing WebSocket connection")
	conn.Close()
}
