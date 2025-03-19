package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

const (
	openRouterApi = "https://openrouter.ai/api/v1"
	model         = "google/gemini-2.0-flash-lite-preview-02-05:free"
	system        = "You are a senior technical lead and will review code to provide conventional commit message based on the changes staged for commit. Provide key changes or features to be included."
	maxTokens     = 4096
)

func getChanges() string {
	cmd := exec.Command("git", "diff", "--cached")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Println("Error getting staged changes:", errb.String())
		return ""
	}
	return outb.String()
}

func getCommitMessageFromLLM(changes string) string {
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

// example with prompt template
func generateCommitMessage() {
	// API details (consider using environment variables)
	apiKey := os.Getenv("OPEN_ROUTER_KEY") // Changed to OPENAI_API_KEY
	if apiKey == "" {
		log.Fatal("OPEN_ROUTER_KEY environment variable is not set") // Updated error message
	}

	// 1. Get staged changes
	cmd := exec.Command("git", "diff", "--cached")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Println("Error getting staged changes:", errb.String())
		return
	}
	stagedChanges := outb.String()

	if stagedChanges == "" {
		log.Println("No staged changes found.")
		return
	}

	// 2. Initialize the LLM
	llmBaseURL := openRouterApi

	// llm, err := openai.New(openai.WithToken(apiKey))
	// Create a new OpenAI client configured for OpenRouter
	llm, err := openai.New(
		openai.WithBaseURL(llmBaseURL), // Set OpenRouter base URL
		openai.WithToken(apiKey),       // Use OpenRouter API key
	)
	if err != nil {
		log.Fatal("Error initializing LLM:", err)
		return
	}

	// 3. Create a PromptTemplate
	promptTemplate := prompts.NewPromptTemplate(
		"You are a senior technical lead and will review code to provide CONVENTIONAL commit message based on the changes staged for commit.\n\nChanges:\n{changes}\n\nCommit message:",
		[]string{"changes"},
	)

	// 4. Create a Chain
	chain := chains.NewLLMChain(llm, promptTemplate)

	// 5. Execute the Chain
	ctx := context.Background() // Add context
	vars := map[string]any{
		"changes": stagedChanges,
	}
	output, err := chain.Call(ctx, vars) // Pass context to Call
	if err != nil {
		log.Fatal("Error executing chain:", err)
		return
	}

	// 6. Print the result
	log.Printf("Staged changes:\n%s\n", stagedChanges)
	log.Println(output["text"])

	log.Printf("Commit message:\n%+v\n", output)
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	// generateCommitMessage()
	changes := getChanges()
	log.Println("Staged changes length:", len(changes))
	if changes == "" {
		log.Println("No staged changes found.")
		return
	}
	commitMsg := getCommitMessageFromLLM(changes)
	log.Println("Commit message:", commitMsg)
}
