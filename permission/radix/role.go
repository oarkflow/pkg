package radix

// Permission defines a single permission for an activity
type Permission struct {
	Name string
}

// Role represents a user role with its permissions
type Role struct {
	name        string
	permissions map[string]Permission
	descendants map[string]IRole
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) Has(permissionName string) bool {
	if _, ok := r.permissions[permissionName]; ok {
		return true
	}
	// Check inherited permissions recursively
	for _, descendant := range r.GetDescendantRoles() {
		if descendant.Has(permissionName) {
			return true
		}
	}
	return false
}

func (r *Role) GetDescendantRoles() []IRole {
	var descendants []IRole
	for _, child := range r.descendants { // Simulate descendent-child relationship through permissions
		descendants = append(descendants, child)
		descendants = append(descendants, child.GetDescendantRoles()...)
	}
	return descendants
}

// AddDescendent adds a new permission to the role
func (r *Role) AddDescendent(descendants ...IRole) {
	for _, descendant := range descendants {
		r.descendants[descendant.Name()] = descendant
	}
}

// AddPermission adds a new permission to the role
func (r *Role) AddPermission(permissions ...Permission) {
	for _, permission := range permissions {
		r.permissions[permission.Name] = permission
	}
}

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
