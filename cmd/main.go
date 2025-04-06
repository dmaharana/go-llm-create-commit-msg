package main

import (
	"log"

	"createCommitMsg/internal/action"
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

	// reviewCode(changes)
	codeReviewComments := llm.CallLLM(changes, action.CODE_REVIEW)
	log.Println("Code review comments: \n\n", codeReviewComments)

	// generate commit message
	commitMsg := llm.CallLLM(changes, action.COMMIT_MESSAGE)
	log.Println("Commit message:", commitMsg)
}
