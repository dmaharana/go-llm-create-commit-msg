package main

import (
	"fmt"
	"log"

	"createCommitMsg/internal/git"
)

func main() {
	gitDiffs, err := git.GetChangesInJSON()
	if err != nil {
		log.Printf("Error getting staged changes in JSON: %v\n", err)
		return
	}
	if gitDiffs == "" {
		log.Println("No staged changes found in JSON")
		return
	}
	log.Println("Staged changes in JSON:")
	fmt.Println(string(gitDiffs))
}
