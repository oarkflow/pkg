package radix

// User represents a user with a role
type User struct {
	name    string
	roles   []IRole
	company ICompany
	module  *Module
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
		for _, userRole := range u.module.users {
			if u.name != userRole.User.Name() {
				return false
			}
		}
		for role := range u.module.roles {
			allowed = append(allowed, role)
		}
	} else if u.company != nil {
		if len(u.company.Users()) == 0 {
			return false
		}
		for _, userRole := range u.company.Users() {
			if u.name != userRole.User.Name() {
				return false
			}
		}
		for role := range u.company.Roles() {
			allowed = append(allowed, role)
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
