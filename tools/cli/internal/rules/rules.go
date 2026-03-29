package rules

import (
	"fmt"

	"stability-flow/internal/branch"
)

func ValidateMerge(source, target string, allowNonPrefixedBranchesToDevelop bool) (bool, string) {
	sourceType := branch.Classify(source)
	targetType := branch.Classify(target)

	// special compatibility mode:
	// allow non-prefixed contributor branches only when targeting develop
	if allowNonPrefixedBranchesToDevelop &&
		sourceType == branch.TypeInvalid &&
		targetType == branch.TypeDevelop {
		return true, "non-prefixed contributor branches may merge into develop when compatibility mode is enabled, using squash merge"
	}

	if sourceType == branch.TypeInvalid {
		return false, fmt.Sprintf("invalid source branch type: %s", source)
	}
	if targetType == branch.TypeInvalid {
		return false, fmt.Sprintf("invalid target branch type: %s", target)
	}

	if sourceType == branch.TypeWIP {
		return false, "wip/* branches are exploratory only and must never be merged"
	}

	// regular work -> develop
	if branch.IsRegularWork(sourceType) && targetType == branch.TypeDevelop {
		return true, "regular work branches may merge only into develop, using squash merge"
	}

	// release/* -> main
	if sourceType == branch.TypeRelease && targetType == branch.TypeMain {
		return true, "only release/* may merge into main, using fast-forward only"
	}

	// main -> sync/*
	if sourceType == branch.TypeMain && targetType == branch.TypeSync {
		return true, "main may merge only into sync/* for reconciliation, using a regular merge commit"
	}

	// sync/* -> develop
	if sourceType == branch.TypeSync && targetType == branch.TypeDevelop {
		return true, "sync/* may merge only into develop, using a regular merge commit"
	}

	// block direct main -> develop
	if sourceType == branch.TypeMain && targetType == branch.TypeDevelop {
		return false, "direct reconciliation from main into develop is prohibited; use sync/*"
	}

	// hotfix branches do not merge directly in this model
	if sourceType == branch.TypeHotfix {
		return false, "hotfix/* branches are not merge targets in Stability Flow; create release/* from hotfix/* instead"
	}

	return false, fmt.Sprintf("merge not allowed by Stability Flow: %s -> %s", sourceType, targetType)
}

func ValidateOrigin(branchName, base string) (bool, string) {
	branchType := branch.Classify(branchName)
	baseType := branch.Classify(base)

	if branchType == branch.TypeInvalid {
		return false, fmt.Sprintf("invalid branch type: %s", branchName)
	}
	if baseType == branch.TypeInvalid {
		return false, fmt.Sprintf("invalid base branch type: %s", base)
	}

	// regular work must branch from develop
	if branch.IsRegularWork(branchType) {
		if baseType == branch.TypeDevelop {
			return true, "regular work branches must be created from develop"
		}
		return false, "regular work branches must branch from develop"
	}

	// hotfix/* must branch from main
	if branchType == branch.TypeHotfix {
		if baseType == branch.TypeMain {
			return true, "hotfix/* must be created from main"
		}
		return false, "hotfix/* must branch from main"
	}

	// release/* must branch from develop or hotfix/*
	if branchType == branch.TypeRelease {
		if baseType == branch.TypeDevelop || baseType == branch.TypeHotfix {
			return true, "release/* must be created from develop or hotfix/*"
		}
		return false, "release/* must branch from develop or hotfix/*"
	}

	// sync/* must branch from develop
	if branchType == branch.TypeSync {
		if baseType == branch.TypeDevelop {
			return true, "sync/* must be created from develop"
		}
		return false, "sync/* must branch from develop"
	}

	// wip/* may branch from develop or main for temporary exploration, but they are not integration branches and are not mergeable under Stability Flow
	if branchType == branch.TypeWIP {
		if baseType == branch.TypeDevelop {
			return true, "wip/* may be created from develop for temporary exploration"
		}
		if baseType == branch.TypeMain {
			return true, "wip/* may be created from main for hotfix troubleshooting"
		}
		return false, "wip/* must branch from develop, or from main when used for hotfix troubleshooting"
	}

	// main/develop should not be created
	if branchType == branch.TypeMain || branchType == branch.TypeDevelop {
		return false, fmt.Sprintf("%s is a long-lived branch and should not be created from another branch", branchType)
	}

	return false, fmt.Sprintf("branch origin not allowed by Stability Flow: %s from %s", branchType, baseType)
}
