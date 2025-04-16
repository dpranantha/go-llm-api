package graph

import (
	"github.com/dpranantha/go-llm-api/back-end/rest/services"
)

// Resolver struct is the entry point for all resolvers
type Resolver struct {
	PromptService *services.PromptService
}
