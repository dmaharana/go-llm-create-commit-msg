package git

import (
	"bytes"
	"log"
	"os/exec"
)

func GetChanges() string {
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

func GetCommitMessagesBetweenTags(tag1, tag2 string) ([]string, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%s", tag1+".."+tag2)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Println("Error getting commits between tags:", errb.String())
		return nil, err
	}
	output := outb.String()
	if output == "" {
		return []string{}, nil
	}
	lines := bytes.Split(outb.Bytes(), []byte("\n"))
	messages := make([]string, 0, len(lines))
	for _, line := range lines {
		messages = append(messages, string(line))
	}
	return messages, nil
}
