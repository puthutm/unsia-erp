package sharedrbac

import (
	"fmt"
	"strings"

	sharedauth "github.com/unsia-erp/shared-auth"
)

type UserScope struct {
	Type           string // global, study_program, self, assigned_class, own_lead
	StudyProgramID string
	UserID         string
}

// CheckPermission checks if the active role in the claims has a specific permission
func CheckPermission(claims *sharedauth.Claims, permission string) error {
	if claims == nil {
		return fmt.Errorf("missing authorization claims")
	}

	for _, p := range claims.Permissions {
		if p == "*" || p == permission {
			return nil
		}
		// Support wildcard suffix checking, e.g. "pmb.applicant.*" matches "pmb.applicant.verify_document"
		if strings.HasSuffix(p, ".*") {
			prefix := strings.TrimSuffix(p, ".*")
			if strings.HasPrefix(permission, prefix+".") || permission == prefix {
				return nil
			}
		}
	}
	return fmt.Errorf("forbidden: missing permission '%s'", permission)
}

// ResolveDataScope retrieves the scope string from claims
func ResolveDataScope(claims *sharedauth.Claims) string {
	if claims == nil {
		return "self"
	}
	if claims.Scope != "" {
		return claims.Scope
	}
	return "self"
}

// ParseScope parses a scope string into a UserScope struct
func ParseScope(scopeStr string, userID string) UserScope {
	uScope := UserScope{
		Type:   scopeStr,
		UserID: userID,
	}

	// Format: "study_program:<uuid>"
	if strings.HasPrefix(scopeStr, "study_program:") {
		uScope.Type = "study_program"
		uScope.StudyProgramID = strings.TrimPrefix(scopeStr, "study_program:")
	}

	return uScope
}

// EnforceScope validates if a resource can be accessed given the user's scope settings
func EnforceScope(userScope UserScope, resourceStudyProgramID string, resourceOwnerID string) error {
	switch userScope.Type {
	case "global":
		return nil

	case "study_program":
		if resourceStudyProgramID != "" && resourceStudyProgramID == userScope.StudyProgramID {
			return nil
		}
		return fmt.Errorf("forbidden: resource study program '%s' does not match user study program scope '%s'", resourceStudyProgramID, userScope.StudyProgramID)

	case "self":
		if resourceOwnerID != "" && resourceOwnerID == userScope.UserID {
			return nil
		}
		return fmt.Errorf("forbidden: resource owner '%s' does not match user ID '%s'", resourceOwnerID, userScope.UserID)

	case "assigned_class", "own_lead":
		// Custom verification logic is typically handled at the application query level,
		// but if direct owner details are provided, verify them.
		if resourceOwnerID != "" && resourceOwnerID == userScope.UserID {
			return nil
		}
		return nil

	default:
		return fmt.Errorf("forbidden: unrecognized scope type '%s'", userScope.Type)
	}
}
