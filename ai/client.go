package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// Client wraps the OpenAI client
type Client struct {
	client *openai.Client
}

// NewClient creates a new AI client with custom proxy
func NewClient() (*Client, error) {
	apiKey := os.Getenv("AI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("AI_KEY environment variable not set")
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://proxy.fuelix.ai/v1"

	return &Client{
		client: openai.NewClientWithConfig(config),
	}, nil
}

// GenerateMRDescription sends the git context to AI and returns MR description
func (c *Client) GenerateMRDescription(commits, diff, model string) (string, error) {
	prompt := buildPrompt(commits, diff)

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("OpenAI API call failed: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

func buildPrompt(commits, diff string) string {
	template := getMRTemplate()

	return fmt.Sprintf(`Based on the following git commits and code changes, please fill out this merge request template with a clear, concise description.

%s

## Commits:
%s

## Code Changes:
%s

Please provide a well-structured merge request description following the template above.`,
		template, commits, diff)
}

func getMRTemplate() string {
	return `# Template for Merge Request Description
## Description
[Provide a brief overview of what this MR accomplishes]

## Changes Made
[List the key changes]

## Type of Change
- [ ] Feature (new functionality)
- [ ] Bug Fix
- [ ] Documentation Update
- [ ] Code Refactor
- [ ] Performance Improvement
- [ ] Other (please specify)

## Testing Done
[Describe the testing approach and results]

## Checklist

- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have updated the documentation accordingly

## Related Issues
[Link any related tickets or issues]`
}
