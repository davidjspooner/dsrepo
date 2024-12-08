package access

type PolicyName string

type Policy struct {
	Name      PolicyName
	Actions   PatternList
	Deny      bool
	Resources PatternList
}

type PolicyList []*Policy

func (pl *PolicyList) Allow(action, resource string) (allowed bool, reason PolicyName) {
	var allowBy, denyBy PolicyName
	for _, policy := range *pl {
		if policy.Actions.MatchString(action) && policy.Resources.MatchString(resource) {
			if policy.Deny {
				denyBy = policy.Name
			} else {
				allowBy = policy.Name
			}
		}
	}
	if denyBy != "" {
		return false, denyBy
	}
	if allowBy != "" {
		return true, allowBy
	}
	return false, "no matching policy"
}
