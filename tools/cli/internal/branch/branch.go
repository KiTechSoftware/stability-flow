package branch

import "strings"

type Type string

const (
	TypeInvalid Type = "invalid"

	TypeMain    Type = "main"
	TypeDevelop Type = "develop"

	TypeRelease Type = "release"
	TypeHotfix  Type = "hotfix"
	TypeSync    Type = "sync"

	TypeFeat     Type = "feat"
	TypeFix      Type = "fix"
	TypeDocs     Type = "docs"
	TypeCI       Type = "ci"
	TypeRefactor Type = "refactor"
	TypeChore    Type = "chore"
	TypeWIP      Type = "wip"
)

func Classify(name string) Type {
	switch {
	case name == "main":
		return TypeMain
	case name == "develop":
		return TypeDevelop
	case strings.HasPrefix(name, "release/"):
		return TypeRelease
	case strings.HasPrefix(name, "hotfix/"):
		return TypeHotfix
	case strings.HasPrefix(name, "sync/"):
		return TypeSync
	case strings.HasPrefix(name, "feat/"):
		return TypeFeat
	case strings.HasPrefix(name, "fix/"):
		return TypeFix
	case strings.HasPrefix(name, "docs/"):
		return TypeDocs
	case strings.HasPrefix(name, "ci/"):
		return TypeCI
	case strings.HasPrefix(name, "refactor/"):
		return TypeRefactor
	case strings.HasPrefix(name, "chore/"):
		return TypeChore
	case strings.HasPrefix(name, "wip/"):
		return TypeWIP
	default:
		return TypeInvalid
	}
}

func IsRegularWork(t Type) bool {
	switch t {
	case TypeFeat, TypeFix, TypeDocs, TypeCI, TypeRefactor, TypeChore:
		return true
	default:
		return false
	}
}
