package commit

import "testing"

func TestValidateWorkMode(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantOK  bool
	}{
		{"feat valid", "feat: add authentication", true},
		{"fix valid", "fix: resolve payment null error", true},
		{"docs valid", "docs: add api usage guide", true},
		{"ci valid", "ci: update release pipeline", true},
		{"refactor valid", "refactor: simplify user service", true},
		{"chore valid", "chore: improve test coverage", true},

		{"test valid", "test: update ci tests", true},
		{"perf valid", "perf: improve latency", true},
		{"build valid", "build: update docker image", true},
		{"style valid", "style: reformat auth package", true},

		{"breaking feat valid", "feat!: remove legacy auth flow", true},
		{"breaking fix valid", "fix!: change retry behavior", true},
		{"breaking footer valid", "feat: add auth\n\nBREAKING CHANGE: auth flow changed", true},
		{"breaking bang and footer valid", "feat!: add auth\n\nBREAKING CHANGE: auth flow changed", true},

		{"revert valid on work branch", "revert: undo previous change", true},

		{"unknown invalid", "banana: random message", false},
		{"malformed missing colon", "feat add authentication", false},
		{"malformed empty description", "feat: ", false},
		{"empty breaking footer invalid", "feat!: remove auth\n\nBREAKING CHANGE:", false},
		{"random invalid", "hello world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOK, _ := Validate(tt.message, ModeWork)
			if gotOK != tt.wantOK {
				t.Fatalf("Validate(%q, work) = %v, want %v", tt.message, gotOK, tt.wantOK)
			}
		})
	}
}

func TestValidateSquashMode(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantOK  bool
	}{
		{"feat valid", "feat: add authentication", true},
		{"test valid", "test: update ci tests", true},
		{"perf valid", "perf: improve latency", true},
		{"breaking footer valid", "feat!: remove auth\n\nBREAKING CHANGE: auth flow changed", true},

		{"revert invalid on squash", "revert: undo previous change", false},
		{"unknown invalid", "banana: random message", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOK, _ := Validate(tt.message, ModeSquash)
			if gotOK != tt.wantOK {
				t.Fatalf("Validate(%q, squash) = %v, want %v", tt.message, gotOK, tt.wantOK)
			}
		})
	}
}
