package v2

// User represents a user with a role
type User struct {
	ID      string
	company string
	module  string
	entity  string
}

// Can check if a user is allowed to do an activity based on their role and inherited permissions
func (u *User) Can(activity string) bool {
	var allowed []string
	if u.company == "" {
		return false
	}
	companyUser := RoleManager.GetUserRoles(u.company, u.ID)
	if companyUser == nil {
		return false
	}
	var userRoles []*Role
	roles := RoleManager.GetAllowedRoles(companyUser, u.module, u.entity)
	for _, r := range companyUser.Company.Roles {
		for _, rt := range roles {
			if r.ID == rt {
				userRoles = append(userRoles, r)
			}
		}
		allowed = append(allowed, r.ID)
	}
	for _, role := range userRoles {
		if role.Has(activity, allowed...) {
			return true
		}
	}
	return false
}

func (u *User) WithCompany(company string) *User {
	return &User{
		ID:      u.ID,
		company: company,
		module:  u.module,
		entity:  u.entity,
	}
}

func (u *User) WithModule(module string) *User {
	return &User{
		ID:      u.ID,
		company: u.company,
		module:  module,
		entity:  u.entity,
	}
}

func (u *User) WithEntity(entity string) *User {
	return &User{
		ID:      u.ID,
		company: u.company,
		module:  u.module,
		entity:  entity,
	}
}
