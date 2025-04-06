package main

import (
	"encoding/json"
	"fmt"
	"log"

	"createCommitMsg/internal/action"
	"createCommitMsg/internal/changelog"
	"createCommitMsg/internal/git"
	"createCommitMsg/internal/llm"
	"flag"

	"github.com/pkg/errors"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	// Define the mode flag
	var mode string
	flag.StringVar(&mode, "mode", "b", "Mode: 'r' (review), 'c' (comment), or 'b' (both)")

	// Define the output flag
	var output string
	flag.StringVar(&output, "output", "b", "Output: 'r' (review), 'c' (comment), or 'b' (both)")

	// Define the format flag
	var format string
	flag.StringVar(&format, "format", "r", "Format: 'r' (raw) or 'j' (json)")

	var tag1 string
	var tag2 string
	flag.StringVar(&tag1, "tag1", "", "Older git tag for changelog generation")
	flag.StringVar(&tag2, "tag2", "", "Newer git tag for changelog generation")

	flag.Parse()

	if tag1 != "" && tag2 != "" {
		changelogStr, err := changelog.GenerateFormattedChangelog(tag1, tag2)
		if err != nil {
			log.Fatalf("Error generating changelog: %v", err)
		}
		fmt.Println(changelogStr)
		return
	}

	// Normalize mode flag
	switch mode {
	case "r":
		mode = "review"
	case "c":
		mode = "comment"
	case "b":
		mode = "both"
	}

	// Normalize output flag
	switch output {
	case "r":
		output = "review"
	case "c":
		output = "comment"
	case "b":
		output = "both"
	}

	// Normalize format flag
	switch format {
	case "r":
		format = "raw"
	case "j":
		format = "json"
	}

	// Set up roles based on mode
	var roles []string
	switch mode {
	case "review":
		roles = []string{action.CODE_REVIEW}
	case "comment":
		roles = []string{action.COMMIT_MESSAGE}
	case "both":
		roles = []string{action.CODE_REVIEW, action.COMMIT_MESSAGE}
	default:
		log.Printf("Invalid mode: %s. Defaulting to 'both'", mode)
		roles = []string{action.CODE_REVIEW, action.COMMIT_MESSAGE}
	}

	changes := git.GetChanges()
	log.Println("Staged changes length:", len(changes))
	if changes == "" {
		log.Println("No staged changes found.")
		return
	}

	response, err := llm.CallLLM(changes, roles)
	if err != nil {
		log.Printf("Error calling LLM: %v\n", errors.WithStack(err))
		return
	}

	printLLMResponse(response, output, format)
}

func printLLMResponse(response map[string]string, output, format string) {
	log.Print("Response: =============================>\n\n")

	if format == "json" {
		filtered := make(map[string]string)

		if output == "review" || output == "both" {
			if val, ok := response[action.CODE_REVIEW]; ok {
				filtered["review"] = val
			}
		}

		if output == "comment" || output == "both" {
			if val, ok := response[action.COMMIT_MESSAGE]; ok {
				filtered["comment"] = val
			}
		}

		jsonBytes, err := json.MarshalIndent(filtered, "", "  ")
		if err != nil {
			log.Printf("Error marshaling JSON: %v\n", err)
			return
		}
		log.Print(string(jsonBytes))
	} else {
		if output == "review" || output == "both" {
			if _, ok := response[action.CODE_REVIEW]; ok {
				log.Print("Code review comments: =============================>\n\n")
				log.Print(response[action.CODE_REVIEW])
			}
		}

		if output == "comment" || output == "both" {
			if _, ok := response[action.COMMIT_MESSAGE]; ok {
				log.Print("Commit message: =============================>\n\n")
				log.Print(response[action.COMMIT_MESSAGE])
			}
		}
	}
}
