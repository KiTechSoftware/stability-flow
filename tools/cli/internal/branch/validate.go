package branch

import "fmt"

func ValidateName(name string, target string, allowNonPrefixed bool) (bool, string) {
	t := Classify(name)

	switch t {
	case TypeMain:
		return true, "valid long-lived branch: main"
	case TypeDevelop:
		return true, "valid long-lived branch: develop"
	case TypeRelease:
		return validateSuffixed(name, "release/")
	case TypeHotfix:
		return validateSuffixed(name, "hotfix/")
	case TypeSync:
		return validateSuffixed(name, "sync/")
	case TypeFeat:
		return validateSuffixed(name, "feat/")
	case TypeFix:
		return validateSuffixed(name, "fix/")
	case TypeDocs:
		return validateSuffixed(name, "docs/")
	case TypeCI:
		return validateSuffixed(name, "ci/")
	case TypeRefactor:
		return validateSuffixed(name, "refactor/")
	case TypeChore:
		return validateSuffixed(name, "chore/")
	case TypeWIP:
		return validateSuffixed(name, "wip/")
	default:
		if name == "" {
			return false, "branch name must not be empty"
		}

		if allowNonPrefixed && target == "develop" {
			return true, "valid non-prefixed contributor branch for develop compatibility mode"
		}

		return false, fmt.Sprintf("invalid branch name: %s", name)
	}
}

func validateSuffixed(name, prefix string) (bool, string) {
	if len(name) <= len(prefix) {
		return false, fmt.Sprintf("branch name must include a non-empty suffix after %q", prefix)
	}
	return true, fmt.Sprintf("valid branch type: %s", Classify(name))
}
