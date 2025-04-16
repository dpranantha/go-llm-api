package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type PromptService struct {
	Endpoint string
}

func NewPromptService(endpoint string) *PromptService {
	return &PromptService{Endpoint: endpoint}
}

// Helper function to make the REST call to /prompt
func (ps *PromptService) Generate(prompt string) (string, error) {
	payload := map[string]string{
		"prompt": prompt,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(ps.Endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["response"], nil
}
