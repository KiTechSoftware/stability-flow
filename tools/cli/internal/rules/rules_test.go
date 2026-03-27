package rules

import "testing"

func TestValidateMerge(t *testing.T) {
	tests := []struct {
		name   string
		source string
		target string
		wantOK bool
	}{
		{"feat to develop", "feat/add-auth", "develop", true},
		{"fix to develop", "fix/payment-null", "develop", true},
		{"release to main", "release/1.2.3", "main", true},
		{"sync to develop", "sync/main-into-develop-1.2.3", "develop", true},

		{"feat to main blocked", "feat/add-auth", "main", false},
		{"wip to develop blocked", "wip/explore-auth", "develop", false},
		{"main to develop blocked", "main", "develop", false},
		{"hotfix to main blocked", "hotfix/1.2.4", "main", false},
		{"hotfix to develop blocked", "hotfix/1.2.4", "develop", false},
		{"release to develop blocked", "release/1.2.3", "develop", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOK, _ := ValidateMerge(tt.source, tt.target)
			if gotOK != tt.wantOK {
				t.Fatalf("ValidateMerge(%q, %q) = %v, want %v", tt.source, tt.target, gotOK, tt.wantOK)
			}
		})
	}
}

func TestValidateOrigin(t *testing.T) {
	tests := []struct {
		name   string
		branch string
		base   string
		wantOK bool
	}{
		{"feat from develop", "feat/add-auth", "develop", true},
		{"fix from develop", "fix/payment-null", "develop", true},
		{"hotfix from main", "hotfix/1.2.4", "main", true},
		{"release from develop", "release/1.2.5", "develop", true},
		{"release from hotfix", "release/1.2.4", "hotfix/1.2.4", true},
		{"sync from develop", "sync/main-into-develop-1.2.4", "develop", true},
		{"wip from develop", "wip/explore-auth", "develop", true},

		{"feat from main blocked", "feat/add-auth", "main", false},
		{"hotfix from develop blocked", "hotfix/1.2.4", "develop", false},
		{"release from main blocked", "release/1.2.4", "main", false},
		{"sync from main blocked", "sync/main-into-develop-1.2.4", "main", false},
		{"main from develop blocked", "main", "develop", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOK, _ := ValidateOrigin(tt.branch, tt.base)
			if gotOK != tt.wantOK {
				t.Fatalf("ValidateOrigin(%q, %q) = %v, want %v", tt.branch, tt.base, gotOK, tt.wantOK)
			}
		})
	}
}
