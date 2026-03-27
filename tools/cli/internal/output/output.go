package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Format string

const (
	FormatText     Format = "text"
	FormatJSON     Format = "json"
	FormatJSONL    Format = "jsonl"
	FormatMarkdown Format = "markdown"
)

type ValidationResult struct {
	OK      bool              `json:"ok"`
	Command string            `json:"command"`
	Reason  string            `json:"reason"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func ParseFormat(value string) (Format, error) {
	switch Format(value) {
	case FormatText, FormatJSON, FormatJSONL, FormatMarkdown:
		return Format(value), nil
	default:
		return "", fmt.Errorf("invalid format %q; must be one of: text, json, jsonl, markdown", value)
	}
}

func Render(format Format, result ValidationResult) (string, error) {
	switch format {
	case FormatText:
		return renderText(result), nil
	case FormatJSON:
		return renderJSON(result)
	case FormatJSONL:
		return renderJSONL(result)
	case FormatMarkdown:
		return renderMarkdown(result), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func renderText(result ValidationResult) string {
	status := "FAIL"
	if result.OK {
		status = "PASS"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s: %s\n", status, summaryLine(result))
	fmt.Fprintf(&b, "reason: %s\n", result.Reason)
	return b.String()
}

func renderJSON(result ValidationResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func renderJSONL(result ValidationResult) (string, error) {
	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func renderMarkdown(result ValidationResult) string {
	status := "❌ Failed"
	if result.OK {
		status = "✅ Passed"
	}

	var b strings.Builder
	b.WriteString("## Stability Flow Validation Result\n\n")
	fmt.Fprintf(&b, "- **Command:** `%s`\n", result.Command)
	fmt.Fprintf(&b, "- **Status:** %s\n", status)

	keys := sortedKeys(result.Fields)
	for _, k := range keys {
		fmt.Fprintf(&b, "- **%s:** `%s`\n", titleize(k), result.Fields[k])
	}

	fmt.Fprintf(&b, "- **Reason:** %s\n", result.Reason)
	return b.String()
}

func summaryLine(result ValidationResult) string {
	switch result.Command {
	case "validate-merge":
		return fmt.Sprintf("merge %s: %s -> %s",
			statusWord(result.OK),
			result.Fields["source"],
			result.Fields["target"],
		)
	case "validate-origin":
		return fmt.Sprintf("branch origin %s: %s from %s",
			statusWord(result.OK),
			result.Fields["branch"],
			result.Fields["base"],
		)
	case "validate-commit":
		return fmt.Sprintf("commit message %s (%s): %s",
			statusWord(result.OK),
			result.Fields["mode"],
			result.Fields["message"],
		)
	case "validate-branch-name":
		return fmt.Sprintf("branch name %s: %s",
			statusWord(result.OK),
			result.Fields["branch"],
		)
	default:
		return fmt.Sprintf("validation %s", statusWord(result.OK))
	}
}

func statusWord(ok bool) string {
	if ok {
		return "allowed"
	}
	return "not allowed"
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func titleize(s string) string {
	switch s {
	case "ok":
		return "OK"
	case "base":
		return "Base"
	case "branch":
		return "Branch"
	case "command":
		return "Command"
	case "message":
		return "Message"
	case "mode":
		return "Mode"
	case "reason":
		return "Reason"
	case "source":
		return "Source"
	case "target":
		return "Target"
	default:
		if s == "" {
			return s
		}
		return strings.ToUpper(s[:1]) + s[1:]
	}
}
