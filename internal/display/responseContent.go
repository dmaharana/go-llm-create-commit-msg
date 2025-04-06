package display

import (
	"strings"

	"github.com/gomarkdown/markdown"
)

func FormatContentToHTML(txt string) string {
	if txt == "" {
		return ""
	}

	md := markdown.ToHTML([]byte(txt), nil, nil)

	return strings.TrimSpace(string(md))
}
