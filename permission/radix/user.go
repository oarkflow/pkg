package radix

import (
	"slices"
)

// User represents a user with a role
type User struct {
	id        string
	roles     []IRole
	company   string
	companies map[string]ICompany
	module    *Module
	entities  []string
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Roles() []IRole {
	return u.roles
}

func (u *User) AssignTo(company ICompany) {
	u.companies[company.ID()] = company
}

// Can check if a user is allowed to do an activity based on their role and inherited permissions
func (u *User) Can(activity string) bool {
	var allowed []string
	if u.module != nil {
		if len(u.module.users) == 0 {
			return false
		}
		// Check if the user's id matches any user in the module
		foundUser := false
		for _, userRole := range u.module.users {
			if u.id == userRole.User.ID() {
				foundUser = true
				break
			}
		}
		if !foundUser {
			return false
		}
		for role := range u.module.roles {
			allowed = append(allowed, role)
		}
	} else if u.company != "" {
		company, ok := u.companies[u.company]
		if !ok {
			return false
		}
		if len(company.Users()) == 0 {
			return false
		}

		// Check if the user's id matches any user in the company
		foundUser := false
		for _, userRole := range company.Users() {
			if u.id == userRole.User.ID() {
				foundUser = true
				break
			}
		}
		if !foundUser {
			return false
		}
		for role := range company.Roles() {
			allowed = append(allowed, role)
		}
	}
	if len(u.entities) > 0 {
		id := u.entities[0]
		var entityMap []string

		// Determine the entity map based on whether the user has a module or a company
		if u.module != nil {
			entityMap = u.module.userEntities[u.id]
			if entities := u.module.entities; len(entities) > 0 {
				// Check if the entity ID is present in the module's entities
				if _, found := entities[id]; !found {
					return false
				}
			}
		} else if u.company != "" {
			company, ok := u.companies[u.company]
			if !ok {
				return false
			}
			entityMap = company.UserEntities()[u.id]
			if entities := company.Entities(); len(entities) > 0 {
				// Check if the entity ID is present in the company's entities
				if _, found := entities[id]; !found {
					return false
				}
			}
		}
		if len(entityMap) > 0 && !slices.Contains(entityMap, id) {
			return false
		}
	}
	for _, role := range u.roles {
		if role.Has(activity, allowed...) {
			return true
		}
	}
	return false
}

func (u *User) Assign(roles ...IRole) {
	if len(roles) == 0 {
		return
	}
	u.roles = append(u.roles, roles...)
}

func (u *User) WithCompany(company string, module ...string) IUser {
	user := &User{
		id:        u.id,
		roles:     u.roles,
		companies: u.companies,
		company:   company,
	}
	if len(module) > 0 {
		company, ok := u.companies[company]
		if ok {
			mod, o := company.GetModule(module[0])
			if o {
				user.module = mod
			}
		}
	}
	return user
}

func (u *User) WithEntity(entities ...string) IUser {
	return &User{
		id:        u.id,
		roles:     u.roles,
		company:   u.company,
		companies: u.companies,
		module:    u.module,
		entities:  entities,
	}
}
