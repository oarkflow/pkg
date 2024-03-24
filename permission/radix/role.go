package radix

import (
	"errors"
	"slices"
)

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
	lock        bool
	permissions map[string]Attribute
	descendants map[string]IRole
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) Lock() {
	r.lock = true
}

func (r *Role) Unlock() {
	r.lock = false
}

func (r *Role) Has(permissionName string, allowedDescendants ...string) bool {
	if _, ok := r.permissions[permissionName]; ok {
		return true
	}
	totalD := len(allowedDescendants)
	// Check inherited permissions recursively
	for _, descendant := range r.GetDescendantRoles() {
		if totalD > 0 {
			if slices.Contains(allowedDescendants, descendant.Name()) {
				if descendant.Has(permissionName, allowedDescendants...) {
					return true
				}
			}
		} else {
			if descendant.Has(permissionName, allowedDescendants...) {
				return true
			}
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
func (r *Role) AddDescendent(descendants ...IRole) error {
	if r.lock {
		return errors.New("changes not allowed")
	}
	for _, descendant := range descendants {
		r.descendants[descendant.Name()] = descendant
	}
	return nil
}

// AddPermission adds a new permission to the role
func (r *Role) AddPermission(permissions ...Attribute) error {
	if r.lock {
		return errors.New("changes not allowed")
	}
	for _, permission := range permissions {
		r.permissions[permission.String()] = permission
	}
	return nil
}
