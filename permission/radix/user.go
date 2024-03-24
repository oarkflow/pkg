package radix

// User represents a user with a role
type User struct {
	name  string
	roles []IRole
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Roles() []IRole {
	return u.roles
}

// Can check if a user is allowed to do an activity based on their role and inherited permissions
func (u *User) Can(activity string) bool {
	for _, role := range u.roles {
		if role.Has(activity) {
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
