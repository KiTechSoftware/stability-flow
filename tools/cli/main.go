package main

import (
	"flag"
	"fmt"
	"os"

	"stability-flow/internal/branch"
	"stability-flow/internal/commit"
	"stability-flow/internal/output"
	"stability-flow/internal/rules"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "validate-merge":
		validateMerge(os.Args[2:])
	case "validate-origin":
		validateOrigin(os.Args[2:])
	case "validate-commit":
		validateCommit(os.Args[2:])
	case "validate-branch-name":
		validateBranchName(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func validateMerge(args []string) {
	fs := flag.NewFlagSet("validate-merge", flag.ExitOnError)
	source := fs.String("source", "", "source branch")
	target := fs.String("target", "", "target branch")
	formatValue := fs.String("format", "text", "output format: text, json, jsonl, markdown")
	fs.Parse(args)

	if *source == "" || *target == "" {
		failUsage("--source and --target are required")
	}

	format := mustParseFormat(*formatValue)

	ok, reason := rules.ValidateMerge(*source, *target)
	renderAndExit(format, output.ValidationResult{
		OK:      ok,
		Command: "validate-merge",
		Reason:  reason,
		Fields: map[string]string{
			"source": *source,
			"target": *target,
		},
	})
}

func validateOrigin(args []string) {
	fs := flag.NewFlagSet("validate-origin", flag.ExitOnError)
	branchName := fs.String("branch", "", "branch being created")
	base := fs.String("base", "", "base branch")
	formatValue := fs.String("format", "text", "output format: text, json, jsonl, markdown")
	fs.Parse(args)

	if *branchName == "" || *base == "" {
		failUsage("--branch and --base are required")
	}

	format := mustParseFormat(*formatValue)

	ok, reason := rules.ValidateOrigin(*branchName, *base)
	renderAndExit(format, output.ValidationResult{
		OK:      ok,
		Command: "validate-origin",
		Reason:  reason,
		Fields: map[string]string{
			"branch": *branchName,
			"base":   *base,
		},
	})
}

func validateCommit(args []string) {
	fs := flag.NewFlagSet("validate-commit", flag.ExitOnError)
	message := fs.String("message", "", "commit message")
	mode := fs.String("mode", "squash", "validation mode: work or squash")
	formatValue := fs.String("format", "text", "output format: text, json, jsonl, markdown")
	fs.Parse(args)

	if *message == "" {
		failUsage("--message is required")
	}

	var commitMode commit.Mode
	switch *mode {
	case "work":
		commitMode = commit.ModeWork
	case "squash":
		commitMode = commit.ModeSquash
	default:
		failUsage("--mode must be one of: work, squash")
	}

	format := mustParseFormat(*formatValue)

	ok, reason := commit.Validate(*message, commitMode)
	renderAndExit(format, output.ValidationResult{
		OK:      ok,
		Command: "validate-commit",
		Reason:  reason,
		Fields: map[string]string{
			"message": *message,
			"mode":    *mode,
		},
	})
}

func validateBranchName(args []string) {
	fs := flag.NewFlagSet("validate-branch-name", flag.ExitOnError)
	name := fs.String("branch", "", "branch name")
	formatValue := fs.String("format", "text", "output format: text, json, jsonl, markdown")
	fs.Parse(args)

	if *name == "" {
		failUsage("--branch is required")
	}

	format := mustParseFormat(*formatValue)

	ok, reason := branch.ValidateName(*name)
	renderAndExit(format, output.ValidationResult{
		OK:      ok,
		Command: "validate-branch-name",
		Reason:  reason,
		Fields: map[string]string{
			"branch": *name,
		},
	})
}

func mustParseFormat(value string) output.Format {
	format, err := output.ParseFormat(value)
	if err != nil {
		failUsage(err.Error())
	}
	return format
}

func renderAndExit(format output.Format, result output.ValidationResult) {
	rendered, err := output.Render(format, result)
	if err != nil {
		fmt.Printf("FAIL: %s\n", err.Error())
		os.Exit(2)
	}

	fmt.Print(rendered)

	if result.OK {
		os.Exit(0)
	}
	os.Exit(1)
}

func failUsage(message string) {
	fmt.Printf("FAIL: %s\n", message)
	os.Exit(2)
}

func usage() {
	fmt.Println(`stability-flow

Usage:
  stability-flow validate-merge --source <branch> --target <branch> [--format text|json|jsonl|markdown]
  stability-flow validate-origin --branch <branch> --base <branch> [--format text|json|jsonl|markdown]
  stability-flow validate-commit --mode <work|squash> --message "<type>: <description>" [--format text|json|jsonl|markdown]
  stability-flow validate-branch-name --branch <branch> [--format text|json|jsonl|markdown]`)
}
