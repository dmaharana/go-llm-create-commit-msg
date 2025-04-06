package llm

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"createCommitMsg/internal/action"
	"createCommitMsg/internal/constant"
	"createCommitMsg/internal/prompt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

const (
	openRouterApi = "https://openrouter.ai/api/v1"
	defaultModel  = "google/gemini-2.0-flash-lite-preview-02-05:free"

	// maxTokens is the maximum number of tokens to generate.
	maxTokens = 4096
)

var model string

func init() {
	model = os.Getenv("MODEL_NAME")
	if model == "" {
		model = defaultModel
	}
}

func CallLLM(changes string, roles []string) (map[string]string, error) {
	if changes == "" {
		log.Println("No staged changes found.")
		return nil, fmt.Errorf("no staged changes found")
	}

	if len(roles) == 0 {
		log.Println("No roles found.")
		return nil, fmt.Errorf("no roles found")
	}

	// sanitize changes
	changes = SanitizeInput(changes)

	llmBaseURL := openRouterApi

	client, err := openai.New(
		openai.WithBaseURL(llmBaseURL),                 // Set OpenRouter base URL
		openai.WithToken(os.Getenv("OPEN_ROUTER_KEY")), // Use OpenRouter API key
		openai.WithModel(model),                        // Specify the model
	)
	if err != nil {
		log.Fatal("Error initializing LLM:", err)
		return nil, err
	}

	// Define the context
	ctx := context.Background()

	// Set the system message and user message
	systemMessages := []string{}
	for _, role := range roles {
		systemMessages = append(systemMessages, getSystemMessage(role))
	}

	systemMessages = RemoveDuplicates(systemMessages)

	content := []llms.MessageContent{}
	for _, systemMessage := range systemMessages {
		content = append(content, llms.TextParts(llms.ChatMessageTypeSystem, systemMessage))
	}
	content = append(content, llms.TextParts(llms.ChatMessageTypeHuman, changes))

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
		return nil, err
	}

	// Print the token count
	PrintTokenCounts(resp)

	// Print the response
	if len(resp.Choices) > 0 {
		content := map[string]string{}
		for i, role := range roles {
			if i < len(resp.Choices) {
				content[role] = resp.Choices[i].Content
			}
		}
		return content, nil
	} else {
		log.Println("No response choices received")
		return nil, fmt.Errorf("no response choices received")
	}
}

// select system message based on action
func getSystemMessage(role string) string {
	switch role {
	case action.COMMIT_MESSAGE:
		return prompt.CommitSummarizer
	case action.CODE_REVIEW:
		return prompt.CodeReviewer
	default:
		return prompt.CommitSummarizer
	}
}

// print number of tokens used
// func PrintTokenCount(resp *llms.ContentResponse) {
func PrintTokenCount(choice *llms.ContentChoice) {
	log.Println("Printing token count...")

	if choice == nil {
		log.Println("No choice found in response")
		return
	}

	content := choice.Content
	if content == "" {
		log.Println("No content found in response")
		return
	}

	log.Printf("Stop reason: %s\n", choice.StopReason)
	if choice.ReasoningContent != "" {
		log.Printf("Reasoning content: %s\n", choice.ReasoningContent)
	}

	if choice.FuncCall != nil {
		log.Printf("Function call: %+v\n", choice.FuncCall)
	}

	if len(choice.ToolCalls) > 0 {
		log.Printf("Tool calls: %+v\n", choice.ToolCalls)
	}

	// Check if usage stats are available (OpenAI-specific)
	if choice.GenerationInfo != nil {
		if promptTokens, ok := choice.GenerationInfo[constant.PromptTokens].(int); ok {
			fmt.Printf("Tokens sent (exact)              : %4d\n", promptTokens)
		}
		if reasoningTokens, ok := choice.GenerationInfo[constant.ReasoningTokens].(int); ok {
			fmt.Printf("Tokens used for reasoning (exact): %4d\n", reasoningTokens)
		}
		if completionTokens, ok := choice.GenerationInfo[constant.CompletionTokens].(int); ok {
			fmt.Printf("Tokens received (exact)          : %4d\n", completionTokens)
		}
	}
}

// print token for all the choices
func PrintTokenCounts(resp *llms.ContentResponse) {
	for _, choice := range resp.Choices {
		PrintTokenCount(choice)
	}
}

// remove duplicate strings
func RemoveDuplicates(strSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// sanitize input to remove potentially harmful content
func SanitizeInput(input string) string {
	tmp := regexp.MustCompile(`[^a-zA-Z0-9 ]`).ReplaceAllString(input, " ")
	return regexp.MustCompile(`[<>]`).ReplaceAllString(tmp, "")
}
