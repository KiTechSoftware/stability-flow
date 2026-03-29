package branch

import "testing"

func TestValidateName(t *testing.T) {
	tests := []struct {
		name             string
		branch           string
		target           string
		allowNonPrefixed bool
		wantOK           bool
	}{
		{"main valid", "main", "", false, true},
		{"develop valid", "develop", "", false, true},

		{"feat valid", "feat/add-authentication", "", false, true},
		{"fix valid", "fix/payment-null-error", "", false, true},
		{"docs valid", "docs/api-usage-guide", "", false, true},
		{"ci valid", "ci/update-release-pipeline", "", false, true},
		{"refactor valid", "refactor/simplify-user-service", "", false, true},
		{"chore valid", "chore/update-dependencies", "", false, true},
		{"wip valid", "wip/explore-oauth-options", "", false, true},

		{"release valid", "release/1.2.3", "", false, true},
		{"hotfix valid", "hotfix/1.2.4", "", false, true},
		{"sync valid", "sync/main-into-develop-1.2.4", "", false, true},

		{"feat empty suffix invalid", "feat/", "", false, false},
		{"fix empty suffix invalid", "fix/", "", false, false},
		{"release empty suffix invalid", "release/", "", false, false},
		{"hotfix empty suffix invalid", "hotfix/", "", false, false},
		{"sync empty suffix invalid", "sync/", "", false, false},
		{"wip empty suffix invalid", "wip/", "", false, false},

		{"unknown prefix invalid by default", "banana/foo", "", false, false},
		{"random invalid by default", "hello-world", "", false, false},

		{"unknown prefix valid in compatibility mode to develop", "banana/foo", "develop", true, true},
		{"random valid in compatibility mode to develop", "hello-world", "develop", true, true},
		{"unknown prefix invalid in compatibility mode to main", "banana/foo", "main", true, false},
		{"unknown prefix invalid in compatibility mode without target", "banana/foo", "", true, false},
		{"empty still invalid in compatibility mode", "", "develop", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOK, _ := ValidateName(tt.branch, tt.target, tt.allowNonPrefixed)
			if gotOK != tt.wantOK {
				t.Fatalf("ValidateName(%q, %q, %v) = %v, want %v", tt.branch, tt.target, tt.allowNonPrefixed, gotOK, tt.wantOK)
			}
		})
	}
}
