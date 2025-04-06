package git

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type GitDiff struct {
	Filename        string   `json:"filename"`
	Change          string   `json:"change"`
	ModifiedLines   []int    `json:"modifiedLines"`
	ModifiedContent []string `json:"modifiedContent"`
}

type fileState struct {
	File            string
	Change          string
	ModifiedLines   []int
	ModifiedContent []string
}

func GetChanges() string {
	cmd := exec.Command("git", "diff", "--cached")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error creating stdout pipe:", err)
		return ""
	}

	var errb bytes.Buffer
	cmd.Stderr = &errb

	if err := cmd.Start(); err != nil {
		log.Println("Error starting git diff command:", errb.String(), err)
		return ""
	}

	var builder strings.Builder
	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		builder.WriteString(scanner.Text())
		builder.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading git diff output:", err)
		return ""
	}

	if err := cmd.Wait(); err != nil {
		log.Println("Error waiting for git diff command:", errb.String(), err)
		return ""
	}

	return builder.String()
}

func GetChangesInJSON() (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--unified=0")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Println("Error getting staged changes:", errb.String())
		return "", err
	}

	diffOutput := outb.String()
	changes := parseDiffOutput(diffOutput)

	jsonOutput, err := json.Marshal(changes)
	if err != nil {
		log.Println("Error marshaling changes to JSON:", err)
		return "", err
	}

	return string(jsonOutput), nil
}

func parseDiffOutput(diffOutput string) []GitDiff {
	diffLines := strings.Split(diffOutput, "\n")
	return parseDiffLines(diffLines)
}

func parseDiffLines(diffLines []string) []GitDiff {
	var changes []GitDiff
	var currentFileState = newFileState()

	for _, line := range diffLines {
		switch {
		case strings.HasPrefix(line, "diff --git"):
			handleNewFile(line, &currentFileState, &changes)
		case strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++"):
			continue
		case strings.HasPrefix(line, "@@"):
			handleModifiedLines(line, &currentFileState)
		case strings.HasPrefix(line, "+"):
			handleAddedLine(line, &currentFileState)
		}
	}

	if currentFileState.File != "" {
		changes = append(changes, currentFileState.ToGitDiff())
	}

	return changes
}

func newFileState() fileState {
	return fileState{
		ModifiedLines:   []int{},
		ModifiedContent: []string{},
	}
}

func (fs *fileState) ToGitDiff() GitDiff {
	return GitDiff{
		Filename:        fs.File,
		Change:          fs.Change,
		ModifiedLines:   fs.ModifiedLines,
		ModifiedContent: fs.ModifiedContent,
	}
}

func handleNewFile(line string, currentFileState *fileState, changes *[]GitDiff) {
	if currentFileState != nil && currentFileState.File != "" {
		*changes = append(*changes, currentFileState.ToGitDiff())
	}
	parts := strings.Split(line, " ")
	if len(parts) > 2 {
		currentFileState.File = strings.TrimPrefix(parts[len(parts)-1], "b/")
		currentFileState.Change = ""
		currentFileState.ModifiedLines = []int{}
		currentFileState.ModifiedContent = []string{}
	}
}

func handleModifiedLines(line string, currentFileState *fileState) {
	parts := strings.Split(line, " ")
	if len(parts) > 1 {
		lineInfo := parts[1]
		lineInfo = strings.TrimPrefix(lineInfo, "-")
		lineInfo = strings.TrimPrefix(lineInfo, "+")
		parts := strings.Split(lineInfo, ",")
		if len(parts) > 0 {
			startLine, err := strconv.Atoi(parts[0])
			if err == nil {
				currentFileState.ModifiedLines = append(currentFileState.ModifiedLines, startLine)
			} else {
				log.Println("Error parsing start line:", err)
			}
		}
		if len(parts) > 1 {
			endLine, err := strconv.Atoi(parts[1])
			if err == nil {
				currentFileState.ModifiedLines = append(currentFileState.ModifiedLines, endLine)
			} else {
				log.Println("Error parsing end line:", err)
			}
		}
	}
}

func handleAddedLine(line string, currentFileState *fileState) {
	currentFileState.Change = "modified"
	currentFileState.ModifiedContent = append(currentFileState.ModifiedContent, line)
}
