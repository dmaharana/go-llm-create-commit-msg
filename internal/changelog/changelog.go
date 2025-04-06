package changelog

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

// GenerateChangelog returns commit messages between two git tags.
func GenerateChangelog(tag1, tag2 string) ([]string, error) {
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

// GenerateFormattedChangelog returns a categorized changelog string between two tags.
func GenerateFormattedChangelog(tag1, tag2 string) (string, error) {
	messages, err := GenerateChangelog(tag1, tag2)
	if err != nil {
		return "", err
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

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "# Changelog\n\n")
	fmt.Fprintf(&buf, "Changes between `%s` and `%s`\n\n", tag1, tag2)

	fmt.Fprintf(&buf, "## New Features\n")
	if len(newFeatures) == 0 {
		fmt.Fprintf(&buf, "_None_\n")
	} else {
		for _, feat := range newFeatures {
			fmt.Fprintf(&buf, "- %s\n", feat)
		}
	}

	fmt.Fprintf(&buf, "\n## Bug Fixes\n")
	if len(bugFixes) == 0 {
		fmt.Fprintf(&buf, "_None_\n")
	} else {
		for _, fix := range bugFixes {
			fmt.Fprintf(&buf, "- %s\n", fix)
		}
	}

	fmt.Fprintf(&buf, "\n## Other Stuff\n")
	if len(otherStuff) == 0 {
		fmt.Fprintf(&buf, "_None_\n")
	} else {
		for _, other := range otherStuff {
			fmt.Fprintf(&buf, "- %s\n", other)
		}
	}

	return buf.String(), nil
}
