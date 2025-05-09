package main

import (
	"flag"
	"time"

	graphql "github.com/dpranantha/go-llm-api/back-end/graphql"

	ginHandler "github.com/dpranantha/go-llm-api/back-end/rest/pooling/gin"
	sseGin "github.com/dpranantha/go-llm-api/back-end/rest/sse/gin"
	wsGin "github.com/dpranantha/go-llm-api/back-end/rest/ws/gin"
	gcors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	fiberHandler "github.com/dpranantha/go-llm-api/back-end/rest/pooling/fiber"
	sseFiber "github.com/dpranantha/go-llm-api/back-end/rest/sse/fiber"
	wsFiber "github.com/dpranantha/go-llm-api/back-end/rest/ws/fiber"
	"github.com/gofiber/fiber/v2"
	fcors "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

func main() {
	// Command line argument parsing
	useFiber := flag.Bool("fiber", false, "Set this flag to use Fiber instead of Gin")
	flag.Parse()
	if *useFiber {
		fiberServer()
	} else {
		ginServer()
	}
}

// Gin Server
func ginServer() {
	r := gin.Default()

	// ✅ CORS Middleware
	r.Use(gcors.New(gcors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // or "*"
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ✅ Logging middleware is included by default in gin.Default()
	// custom logging, use: r.Use(gin.Logger())

	// REST route
	r.POST("/prompt", ginHandler.HandlePrompt)

	// GraphQL routes
	graphql.RegisterGraphQLRoutesGin(r)

	// WS routes
	r.GET("/prompt/ws", func(c *gin.Context) {
		wsGin.HandleLLMStream(c.Writer, c.Request)
	})

	//SSE routes
	r.POST("/prompt/sse", sseGin.HandleLLMStreamSSE)

	r.Run(":8080")
}

func fiberServer() {
	// Initialize Fiber app
	app := fiber.New()

	// ✅ CORS Middleware
	app.Use(fcors.New(fcors.Config{
		AllowOrigins:     "http://localhost:5173", // Frontend origin
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour / time.Second),
	}))

	// ✅ REST API Route for Prompts
	app.Post("/prompt", fiberHandler.HandlePrompt)

	// ✅ Register GraphQL Routes
	graphql.RegisterGraphQLRoutesFiber(app)

	// WS routes
	app.Get("/prompt/ws", websocket.New(func(c *websocket.Conn) {
		// Handle LLM streaming logic in this handler
		wsFiber.HandleLLMStreamFiber(c)
	}))

	// SSE route
	app.Post("/prompt/sse", sseFiber.HandlePromptSSE)

	// Start the server on port 8080
	app.Listen(":8080")
}
