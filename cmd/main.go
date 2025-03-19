package main

import (
	"log"

	"createCommitMsg/internal/git"
	"createCommitMsg/internal/llm"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	// generateCommitMessage()
	changes := git.GetChanges()
	log.Println("Staged changes length:", len(changes))
	if changes == "" {
		log.Println("No staged changes found.")
		return
	}

	commitMsg := llm.GetCommitMessageFromLLM(changes)
	log.Println("Commit message:", commitMsg)
}
