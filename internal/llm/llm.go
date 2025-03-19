package llm

import (
	"context"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

const (
	openRouterApi = "https://openrouter.ai/api/v1"
	model         = "google/gemini-2.0-flash-lite-preview-02-05:free"
	system        = "You are a senior technical lead and will review code to provide conventional commit message based on the changes staged for commit. Provide key changes or features to be included."
	maxTokens     = 4096
)

func GetCommitMessageFromLLM(changes string) string {
	if changes == "" {
		log.Println("No staged changes found.")
		return ""
	}

	llmBaseURL := openRouterApi

	client, err := openai.New(
		openai.WithBaseURL(llmBaseURL),                 // Set OpenRouter base URL
		openai.WithToken(os.Getenv("OPEN_ROUTER_KEY")), // Use OpenRouter API key
		openai.WithModel(model),                        // Specify the model
	)
	if err != nil {
		log.Fatal("Error initializing LLM:", err)
		return ""
	}

	// Define the context
	ctx := context.Background()

	// Specify the model and messages
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, system),
		llms.TextParts(llms.ChatMessageTypeHuman, changes),
	}

	// Call the OpenRouter API via the client if streaming is needed
	// resp, err := client.GenerateContent(ctx, content,
	// 	llms.WithMaxTokens(1024),
	// 	llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	// 		log.Print(string(chunk))
	// 		return nil
	// 	}))

	// create chat request
	resp, err := client.GenerateContent(ctx, content, llms.WithMaxTokens(maxTokens))
	if err != nil {
		// log.Fatal(err)
		log.Printf("Error calling API: %v\n", err)
		return ""
	}

	// Print the response
	if len(resp.Choices) > 0 {
		// log.Printf("Response: %+v", resp.Choices[0].Content)
		return resp.Choices[0].Content
	} else {
		log.Println("No response choices received")
		return ""
	}
}
