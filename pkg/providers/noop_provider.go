package providers

import (
	"context"
	"fmt"
)

// NoOpProvider is a placeholder provider used when no LLM API key is configured.
// It returns a user-friendly error message instead of crashing the application.
type NoOpProvider struct {
	reason string
}

func NewNoOpProvider(reason string) *NoOpProvider {
	return &NoOpProvider{reason: reason}
}

func (p *NoOpProvider) Chat(
	ctx context.Context,
	messages []Message,
	tools []ToolDefinition,
	model string,
	options map[string]any,
) (*LLMResponse, error) {
	return nil, fmt.Errorf("%s. Set OPENROUTER_API_KEY in your Railway environment variables (https://openrouter.ai/keys)", p.reason)
}

func (p *NoOpProvider) GetDefaultModel() string {
	return "none"
}
