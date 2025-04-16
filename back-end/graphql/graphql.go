package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dpranantha/go-llm-api/back-end/graphql/graph"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/dpranantha/go-llm-api/back-end/rest/services"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

// GIN version
func RegisterGraphQLRoutesGin(r *gin.Engine) {
	promptService := services.NewPromptService("http://localhost:8080/prompt")
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			PromptService: promptService,
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	r.GET("/graphql", gin.WrapH(playground.Handler("GraphQL Playground", "/query")))
	r.POST("/query", gin.WrapH(srv))
}

// Fiber version
func RegisterGraphQLRoutesFiber(app *fiber.App) {
	promptService := services.NewPromptService("http://localhost:8080/prompt")
	// Create the GraphQL server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			PromptService: promptService,
		},
	}))

	// Add the transports for different types of requests
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// Optional: Set up query cache and other extensions
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	app.Get("/graphql", adaptor.HTTPHandler(playground.Handler("GraphQL Playground", "/query")))
	app.Post("/query", adaptor.HTTPHandler(srv))
}
