package commit

import (
	"fmt"
	"regexp"
	"strings"
)

type Mode string

const (
	ModeWork   Mode = "work"
	ModeSquash Mode = "squash"
)

type Result struct {
	Type               string
	Description        string
	IsBreaking         bool
	HasBreakingFooter  bool
	BreakingFooterText string
}

var headerRE = regexp.MustCompile(`^([a-z]+)(!)?: (.+)$`)

var allowedTypes = map[string]struct{}{
	"feat":     {},
	"fix":      {},
	"docs":     {},
	"ci":       {},
	"refactor": {},
	"chore":    {},

	// supported Angular-style commit types
	"perf":  {},
	"test":  {},
	"build": {},
	"style": {},

	// allowed on work branches, blocked for final squash
	"revert": {},
}

func Validate(message string, mode Mode) (bool, string) {
	_, reason, err := Parse(message, mode)
	if err != nil {
		return false, err.Error()
	}
	return true, reason
}

func Parse(message string, mode Mode) (Result, string, error) {
	var result Result

	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return result, "", fmt.Errorf("commit message must not be empty")
	}

	lines := strings.Split(trimmed, "\n")
	header := strings.TrimSpace(lines[0])

	matches := headerRE.FindStringSubmatch(header)
	if matches == nil {
		return result, "", fmt.Errorf("commit message must match '<type>: <description>' or '<type>!: <description>'")
	}

	commitType := matches[1]
	description := strings.TrimSpace(matches[3])

	if description == "" {
		return result, "", fmt.Errorf("commit description must not be empty")
	}

	if _, ok := allowedTypes[commitType]; !ok {
		return result, "", fmt.Errorf("unsupported commit type: %s", commitType)
	}

	if commitType == "revert" && mode == ModeSquash {
		return result, "", fmt.Errorf("revert is allowed on work branches but not as the final squash commit; categorize the squash commit by impact")
	}

	result.Type = commitType
	result.Description = description
	result.IsBreaking = matches[2] == "!"

	footerText, hasBreakingFooter, err := extractBreakingFooter(lines[1:])
	if err != nil {
		return result, "", err
	}
	result.HasBreakingFooter = hasBreakingFooter
	result.BreakingFooterText = footerText

	switch {
	case result.IsBreaking:
		return result, fmt.Sprintf("valid breaking change commit type: %s!", commitType), nil
	case result.HasBreakingFooter:
		return result, fmt.Sprintf("valid commit type: %s", commitType), nil
	default:
		return result, fmt.Sprintf("valid commit type: %s", commitType), nil
	}
}

func extractBreakingFooter(lines []string) (string, bool, error) {
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "BREAKING CHANGE:") {
			text := strings.TrimSpace(strings.TrimPrefix(trimmed, "BREAKING CHANGE:"))
			if text == "" {
				return "", false, fmt.Errorf("BREAKING CHANGE footer must include a description")
			}

			extra := []string{text}
			for _, next := range lines[i+1:] {
				nextTrimmed := strings.TrimSpace(next)
				if nextTrimmed == "" {
					break
				}
				extra = append(extra, nextTrimmed)
			}

			return strings.Join(extra, " "), true, nil
		}
	}
	return "", false, nil
}
