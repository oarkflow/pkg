package v2

// User represents a user with a role
type User struct {
	ID string
}

// Can check if a user is allowed to do an activity based on their role and inherited permissions
func (u *User) Can(company, module, entity, group, activity string) bool {
	return Can(u.ID, company, module, entity, group, activity)
}

func Can(userID, company, module, entity, group, activity string) bool {
	var allowed []string
	if company == "" {
		return false
	}
	companyUser := RoleManager.GetUserRoles(company, userID)
	if companyUser == nil {
		return false
	}
	var userRoles []*Role
	roles := RoleManager.GetAllowedRoles(companyUser, module, entity)
	for _, r := range companyUser.Company.Roles {
		for _, rt := range roles {
			if r.ID == rt {
				userRoles = append(userRoles, r)
			}
		}
		allowed = append(allowed, r.ID)
	}
	for _, role := range userRoles {
		if role.Has(group, activity, allowed...) {
			return true
		}
	}
	return false
}
