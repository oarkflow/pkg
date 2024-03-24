package radix

import (
	"slices"
)

// User represents a user with a role
type User struct {
	name     string
	roles    []IRole
	company  ICompany
	module   *Module
	entities []string
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Roles() []IRole {
	return u.roles
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
			if u.name == userRole.User.Name() {
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
	} else if u.company != nil {
		if len(u.company.Users()) == 0 {
			return false
		}

		// Check if the user's id matches any user in the company
		foundUser := false
		for _, userRole := range u.company.Users() {
			if u.name == userRole.User.Name() {
				foundUser = true
				break
			}
		}
		if !foundUser {
			return false
		}
		for role := range u.company.Roles() {
			allowed = append(allowed, role)
		}
	}
	if len(u.entities) > 0 {
		id := u.entities[0]
		var entityMap []string

		// Determine the entity map based on whether the user has a module or a company
		if u.module != nil {
			entityMap = u.module.userEntities[u.name]
			if entities := u.module.entities; len(entities) > 0 {
				// Check if the entity ID is present in the module's entities
				if _, found := entities[id]; !found {
					return false
				}
			}
		} else if u.company != nil {
			entityMap = u.company.UserEntities()[u.name]
			if entities := u.company.Entities(); len(entities) > 0 {
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

func (u *User) WithCompany(company ICompany, module ...string) IUser {
	user := &User{
		name:    u.name,
		roles:   u.roles,
		company: company,
	}
	if len(module) > 0 {
		mod, ok := company.GetModule(module[0])
		if ok {
			user.module = mod
		}
	}
	return user
}

func (u *User) WithEntity(entities ...string) IUser {
	return &User{
		name:     u.name,
		roles:    u.roles,
		company:  u.company,
		module:   u.module,
		entities: entities,
	}
}
