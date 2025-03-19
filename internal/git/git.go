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
