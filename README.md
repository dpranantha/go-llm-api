# Go LLM API: Gin and Fiber versions

This project provides two versions of a simple Go-based API using [Ollama](https://ollama.com) and LLaMA 2 for local LLM inference on macOS:

- ✅ Gin version
- ✅ Fiber version
- Use Rest API and GraphQL accessing Rest API
- 🧠 Local inference with LLaMA 2 via Ollama
- ⚠️ Docker and `.env` support not implemented yet

## Prerequisites

- Go 1.21+
- macOS (tested with Apple Silicon)
- [Ollama](https://ollama.com/download) installed and running LLaMA 2 model pulled locally:
  - On MacOS
  ```bash
  brew install ollama
  brew services start ollama
  ollama run llama2

## Run (Gin version)
  ```bash
  make run
  ```
The Gin API will start on http://localhost:8080

## Run (Fiber version)
  ```bash
  make run FRAMEWORK=fiber
  ```
The Fiber API will start on http://localhost:8080

## API Usage
Via curl POST:

```bash
curl -X POST http://localhost:8080/prompt \
  -H "Content-Type: application/json" \
  -d '{"prompt": "What is the capital of France?"}'
```

Via Graphql UI in http://localhost:8080/graphql:
```
query {
  promptResponse(prompt: "Tell me a joke about gophers")
}
```

Via ChatUI
  1. Go to front-end/gochat folder
  2. Install dependencies and run locally - assuming you have npm
  ```
  npm install
  npm run dev
  ```
  3. Two flavours

    - For Rest API non stream: go to http://localhost:5173/chat - Uses GraphQL interface
    - For Websocket stream style: go to http://localhost:5173/chatws - Uses REST + Websocket
    
  4. Ask your question in the chat box
  
## Roadmap
 - Add .env support
 - Extracting LLM call to internal utilities
 - Add Docker setup for cross-platform usage
 - Add streaming support via SSE / WebSockets
 - Add Chat UI

## License
MIT
