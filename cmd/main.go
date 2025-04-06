package main

import (
	"encoding/json"
	"fmt"
	"log"

	"createCommitMsg/internal/action"
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
		messages, err := git.GetCommitMessagesBetweenTags(tag1, tag2)
		if err != nil {
			log.Fatalf("Error fetching commits between tags: %v", err)
		}

		var newFeatures []string
		var bugFixes []string
		var otherStuff []string

		for _, msg := range messages {
			lmsg := msg
			if len(lmsg) >= 5 && (lmsg[:5] == "feat:" || lmsg[:5] == "Feat:" || lmsg[:8] == "feature:" || lmsg[:8] == "Feature:") {
				newFeatures = append(newFeatures, msg)
			} else if len(lmsg) >= 4 && (lmsg[:4] == "fix:" || lmsg[:4] == "Fix:") {
				bugFixes = append(bugFixes, msg)
			} else {
				otherStuff = append(otherStuff, msg)
			}
		}

		fmt.Println("# Changelog")
		fmt.Printf("\nChanges between `%s` and `%s`\n\n", tag1, tag2)

		fmt.Println("## New Features")
		if len(newFeatures) == 0 {
			fmt.Println("_None_")
		} else {
			for _, feat := range newFeatures {
				fmt.Println("- " + feat)
			}
		}

		fmt.Println("\n## Bug Fixes")
		if len(bugFixes) == 0 {
			fmt.Println("_None_")
		} else {
			for _, fix := range bugFixes {
				fmt.Println("- " + fix)
			}
		}

		fmt.Println("\n## Other Stuff")
		if len(otherStuff) == 0 {
			fmt.Println("_None_")
		} else {
			for _, other := range otherStuff {
				fmt.Println("- " + other)
			}
		}
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
