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

type globalArgs struct {
	command                           string
	branch                            string
	base                              string
	source                            string
	target                            string
	message                           string
	mode                              string
	format                            string
	allowNonPrefixedBranchesToDevelop bool
}

func main() {
	args := parseArgs()

	switch args.command {
	case "validate-merge":
		validateMerge(args)
	case "validate-origin":
		validateOrigin(args)
	case "validate-commit":
		validateCommit(args)
	case "validate-branch-name":
		validateBranchName(args)
	default:
		usage()
		os.Exit(2)
	}
}

func parseArgs() globalArgs {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	command := os.Args[1]

	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	branchName := fs.String("branch", "", "branch name for validate-branch-name or branch being created for validate-origin")
	base := fs.String("base", "", "base branch for validate-origin")
	source := fs.String("source", "", "source branch for validate-merge")
	target := fs.String("target", "", "target branch for validate-merge or contextual validation for validate-branch-name")
	message := fs.String("message", "", "commit message for validate-commit")
	mode := fs.String("mode", "squash", "validation mode: work or squash")
	formatValue := fs.String("format", "text", "output format: text, json, jsonl, markdown")
	allowNonPrefixedBranchesToDevelop := fs.Bool("allow-non-prefixed-branches-to-develop", false, "allow branches without specific prefixes to merge into develop")

	if err := fs.Parse(os.Args[2:]); err != nil {
		failUsage(err.Error())
	}

	return globalArgs{
		command:                           command,
		branch:                            *branchName,
		base:                              *base,
		source:                            *source,
		target:                            *target,
		message:                           *message,
		mode:                              *mode,
		format:                            *formatValue,
		allowNonPrefixedBranchesToDevelop: *allowNonPrefixedBranchesToDevelop,
	}
}

func validateMerge(args globalArgs) {
	if args.source == "" || args.target == "" {
		failUsage("--source and --target are required")
	}

	format := mustParseFormat(args.format)

	ok, reason := rules.ValidateMerge(
		args.source,
		args.target,
		args.allowNonPrefixedBranchesToDevelop,
	)

	renderAndExit(format, output.ValidationResult{
		OK:      ok,
		Command: "validate-merge",
		Reason:  reason,
		Fields: map[string]string{
			"source": args.source,
			"target": args.target,
		},
	})
}

func validateOrigin(args globalArgs) {
	if args.branch == "" || args.base == "" {
		failUsage("--branch and --base are required")
	}

	format := mustParseFormat(args.format)

	ok, reason := rules.ValidateOrigin(args.branch, args.base)
	renderAndExit(format, output.ValidationResult{
		OK:      ok,
		Command: "validate-origin",
		Reason:  reason,
		Fields: map[string]string{
			"branch": args.branch,
			"base":   args.base,
		},
	})
}

func validateCommit(args globalArgs) {
	if args.message == "" {
		failUsage("--message is required")
	}

	var commitMode commit.Mode
	switch args.mode {
	case "work":
		commitMode = commit.ModeWork
	case "squash":
		commitMode = commit.ModeSquash
	default:
		failUsage("--mode must be one of: work, squash")
	}

	format := mustParseFormat(args.format)

	ok, reason := commit.Validate(args.message, commitMode)
	renderAndExit(format, output.ValidationResult{
		OK:      ok,
		Command: "validate-commit",
		Reason:  reason,
		Fields: map[string]string{
			"message": args.message,
			"mode":    args.mode,
		},
	})
}

func validateBranchName(args globalArgs) {
	if args.branch == "" {
		failUsage("--branch is required")
	}

	format := mustParseFormat(args.format)

	ok, reason := branch.ValidateName(
		args.branch,
		args.target,
		args.allowNonPrefixedBranchesToDevelop,
	)

	fields := map[string]string{
		"branch": args.branch,
	}
	if args.target != "" {
		fields["target"] = args.target
	}

	renderAndExit(format, output.ValidationResult{
		OK:      ok,
		Command: "validate-branch-name",
		Reason:  reason,
		Fields:  fields,
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
  stability-flow validate-branch-name --branch <branch> [--target <branch>] [--allow-non-prefixed-branches-to-develop] [--format text|json|jsonl|markdown]`)
}
