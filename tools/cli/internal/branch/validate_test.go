package branch

import "testing"

func TestValidateName(t *testing.T) {
	tests := []struct {
		name   string
		branch string
		wantOK bool
	}{
		{"main valid", "main", true},
		{"develop valid", "develop", true},

		{"feat valid", "feat/add-authentication", true},
		{"fix valid", "fix/payment-null-error", true},
		{"docs valid", "docs/api-usage-guide", true},
		{"ci valid", "ci/update-release-pipeline", true},
		{"refactor valid", "refactor/simplify-user-service", true},
		{"chore valid", "chore/update-dependencies", true},
		{"wip valid", "wip/explore-oauth-options", true},

		{"release valid", "release/1.2.3", true},
		{"hotfix valid", "hotfix/1.2.4", true},
		{"sync valid", "sync/main-into-develop-1.2.4", true},

		{"feat empty suffix invalid", "feat/", false},
		{"fix empty suffix invalid", "fix/", false},
		{"release empty suffix invalid", "release/", false},
		{"hotfix empty suffix invalid", "hotfix/", false},
		{"sync empty suffix invalid", "sync/", false},
		{"wip empty suffix invalid", "wip/", false},

		{"unknown prefix invalid", "banana/foo", false},
		{"random invalid", "hello-world", false},
		{"empty invalid", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOK, _ := ValidateName(tt.branch)
			if gotOK != tt.wantOK {
				t.Fatalf("ValidateName(%q) = %v, want %v", tt.branch, gotOK, tt.wantOK)
			}
		})
	}
}
