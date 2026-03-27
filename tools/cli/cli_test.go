package main

import (
	"os/exec"
	"strings"
	"testing"
)

func runCLI(t *testing.T, args ...string) (int, string) {
	t.Helper()

	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err == nil {
		return 0, output
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("unexpected error: %v\noutput:\n%s", err, output)
	}

	return exitErr.ExitCode(), output
}

func assertContainsAll(t *testing.T, output string, want ...string) {
	t.Helper()

	for _, s := range want {
		if !strings.Contains(output, s) {
			t.Fatalf("expected output to contain %q\nfull output:\n%s", s, output)
		}
	}
}

func TestCLICommands(t *testing.T) {
	t.Run("validate merge feat to develop passes text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-merge",
			"--source", "feat/add-authentication",
			"--target", "develop",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"PASS: merge allowed: feat/add-authentication -> develop",
			"reason: regular work branches may merge only into develop",
		)
	})

	t.Run("validate merge feat to main fails text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-merge",
			"--source", "feat/add-authentication",
			"--target", "main",
		)

		if code != 1 {
			t.Fatalf("expected exit code 1, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"FAIL: merge not allowed: feat/add-authentication -> main",
			"reason: merge not allowed by Stability Flow: feat -> main",
		)
	})

	t.Run("validate merge release to main passes json", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-merge",
			"--source", "release/1.2.3",
			"--target", "main",
			"--format", "json",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			`"ok": true`,
			`"command": "validate-merge"`,
			`"source": "release/1.2.3"`,
			`"target": "main"`,
		)
	})

	t.Run("validate merge release to main passes jsonl", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-merge",
			"--source", "release/1.2.3",
			"--target", "main",
			"--format", "jsonl",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			`"ok":true`,
			`"command":"validate-merge"`,
			`"source":"release/1.2.3"`,
			`"target":"main"`,
		)
	})

	t.Run("validate origin hotfix from main passes text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-origin",
			"--branch", "hotfix/1.2.4",
			"--base", "main",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"PASS: branch origin allowed: hotfix/1.2.4 from main",
			"reason: hotfix/* must be created from main",
		)
	})

	t.Run("validate origin release from hotfix passes markdown", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-origin",
			"--branch", "release/1.2.4",
			"--base", "hotfix/1.2.4",
			"--format", "markdown",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"## Stability Flow Validation Result",
			"**Command:** `validate-origin`",
			"**Status:** ✅ Passed",
			"**Branch:** `release/1.2.4`",
			"**Base:** `hotfix/1.2.4`",
		)
	})

	t.Run("validate origin feat from main fails text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-origin",
			"--branch", "feat/add-authentication",
			"--base", "main",
		)

		if code != 1 {
			t.Fatalf("expected exit code 1, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"FAIL: branch origin not allowed: feat/add-authentication from main",
			"reason: regular work branches must branch from develop",
		)
	})

	t.Run("validate commit feat passes squash text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-commit",
			"--mode", "squash",
			"--message", "feat: add authentication",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"PASS: commit message allowed (squash): feat: add authentication",
			"reason: valid commit type: feat",
		)
	})

	t.Run("validate commit breaking feat with footer passes text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-commit",
			"--mode", "squash",
			"--message", "feat!: remove legacy auth flow\n\nBREAKING CHANGE: legacy auth removed",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"PASS: commit message allowed (squash): feat!: remove legacy auth flow",
			"reason: valid breaking change commit type: feat!",
		)
	})

	t.Run("validate commit revert passes work text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-commit",
			"--mode", "work",
			"--message", "revert: undo previous change",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"PASS: commit message allowed (work): revert: undo previous change",
			"reason: valid commit type: revert",
		)
	})

	t.Run("validate commit revert fails squash text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-commit",
			"--mode", "squash",
			"--message", "revert: undo previous change",
		)

		if code != 1 {
			t.Fatalf("expected exit code 1, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"FAIL: commit message not allowed (squash): revert: undo previous change",
			"reason: revert is allowed on work branches but not as the final squash commit",
		)
	})

	t.Run("validate commit test passes json", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-commit",
			"--mode", "squash",
			"--message", "test: update ci tests",
			"--format", "json",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			`"ok": true`,
			`"command": "validate-commit"`,
			`"message": "test: update ci tests"`,
			`"mode": "squash"`,
		)
	})

	t.Run("validate commit test passes jsonl", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-commit",
			"--mode", "squash",
			"--message", "test: update ci tests",
			"--format", "jsonl",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			`"ok":true`,
			`"command":"validate-commit"`,
			`"message":"test: update ci tests"`,
			`"mode":"squash"`,
		)
	})

	t.Run("validate commit perf passes markdown", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-commit",
			"--mode", "squash",
			"--message", "perf: improve query latency",
			"--format", "markdown",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"## Stability Flow Validation Result",
			"**Command:** `validate-commit`",
			"**Status:** ✅ Passed",
			"**Message:** `perf: improve query latency`",
			"**Mode:** `squash`",
		)
	})

	t.Run("validate commit unknown type fails text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-commit",
			"--mode", "squash",
			"--message", "banana: random message",
		)

		if code != 1 {
			t.Fatalf("expected exit code 1, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"FAIL: commit message not allowed (squash): banana: random message",
			"reason: unsupported commit type: banana",
		)
	})

	t.Run("validate branch name feat passes text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-branch-name",
			"--branch", "feat/add-authentication",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"PASS: branch name allowed: feat/add-authentication",
			"reason: valid branch type: feat",
		)
	})

	t.Run("validate branch name main passes markdown", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-branch-name",
			"--branch", "main",
			"--format", "markdown",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"## Stability Flow Validation Result",
			"**Command:** `validate-branch-name`",
			"**Status:** ✅ Passed",
			"**Branch:** `main`",
		)
	})

	t.Run("validate branch name release passes json", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-branch-name",
			"--branch", "release/1.2.3",
			"--format", "json",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			`"ok": true`,
			`"command": "validate-branch-name"`,
			`"branch": "release/1.2.3"`,
		)
	})

	t.Run("validate branch name release passes jsonl", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-branch-name",
			"--branch", "release/1.2.3",
			"--format", "jsonl",
		)

		if code != 0 {
			t.Fatalf("expected exit code 0, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			`"ok":true`,
			`"command":"validate-branch-name"`,
			`"branch":"release/1.2.3"`,
		)
	})

	t.Run("validate branch name empty suffix fails text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-branch-name",
			"--branch", "feat/",
		)

		if code != 1 {
			t.Fatalf("expected exit code 1, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"FAIL: branch name not allowed: feat/",
			"reason: branch name must include a non-empty suffix",
		)
	})

	t.Run("validate branch name unknown prefix fails text", func(t *testing.T) {
		code, output := runCLI(t,
			"validate-branch-name",
			"--branch", "banana/foo",
		)

		if code != 1 {
			t.Fatalf("expected exit code 1, got %d\noutput:\n%s", code, output)
		}

		assertContainsAll(t, output,
			"FAIL: branch name not allowed: banana/foo",
			"reason: invalid branch name: banana/foo",
		)
	})
}
