package radix

type Attribute struct {
	Resource string
	Action   string
}

func (a Attribute) String() string {
	return a.Resource + " " + a.Action
}

// Role represents a user role with its permissions
type Role struct {
	name        string
	permissions map[string]Attribute
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
func (r *Role) AddPermission(permissions ...Attribute) {
	for _, permission := range permissions {
		r.permissions[permission.String()] = permission
	}
}
