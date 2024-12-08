package access

import "fmt"

type RoleName string

type Role struct {
	Name     RoleName
	Policies []PolicyName
	policies PolicyList
}

func CrossLink(allRoles []*Role, allPolicies PolicyList) error {
	errors := ErrorList{}
	policiesMap := make(map[PolicyName]*Policy, len(allPolicies))
	for _, policy := range allPolicies {
		if _, exists := policiesMap[policy.Name]; exists {
			errors = append(errors, fmt.Errorf("duplicate policy %q", policy.Name))
			continue
		}
		policiesMap[policy.Name] = policy
	}
	for _, role := range allRoles {
		for _, policyName := range role.Policies {
			policy, exists := policiesMap[policyName]
			if !exists {
				errors = append(errors, fmt.Errorf("role %q references unknown policy %q", role.Name, policyName))
				continue
			}
			role.policies = append(role.policies, policy)
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (r *Role) Allow(action, resource string) (allowed bool, reason PolicyName) {
	return r.policies.Allow(action, resource)
}

type RoleList []*Role
