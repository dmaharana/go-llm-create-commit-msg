package llm

import (
	"context"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

const (
	openRouterApi  = "https://openrouter.ai/api/v1"
	defaultModel   = "google/gemini-2.0-flash-lite-preview-02-05:free"
	system         = "You are a senior technical lead and will review code to provide conventional commit message based on the changes staged for commit. Provide key changes or features to be included."
	systemShortMsg = "I want you to act as a commit message generator. I will provide you with information about the task and the prefix for the task code, and I would like you to generate an appropriate commit message using the conventional commit format. Do not write any explanations or other words, just reply with the commit message."
	maxTokens      = 4096
)

var model string

func init() {
	model = os.Getenv("MODEL_NAME")
	if model == "" {
		model = defaultModel
	}
}

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
		return resp.Choices[0].Content
	} else {
		log.Println("No response choices received")
		return ""
	}
}
