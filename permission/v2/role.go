package v2

import (
	"errors"
	"slices"
)

type Attribute struct {
	Resource string
	Action   string
}

func (a Attribute) String(delimiter ...string) string {
	delim := " "
	if len(delimiter) > 0 {
		delim = delimiter[0]
	}
	return a.Resource + delim + a.Action
}

// Role represents a user role with its permissions
type Role struct {
	ID          string
	lock        bool
	permissions map[string]Attribute
	descendants map[string]*Role
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
			if slices.Contains(allowedDescendants, descendant.ID) {
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

func (r *Role) GetDescendantRoles() []*Role {
	var descendants []*Role
	for _, child := range r.descendants { // Simulate descendent-child relationship through permissions
		descendants = append(descendants, child)
		descendants = append(descendants, child.GetDescendantRoles()...)
	}
	return descendants
}

// AddDescendent adds a new permission to the role
func (r *Role) AddDescendent(descendants ...*Role) error {
	if r.lock {
		return errors.New("changes not allowed")
	}
	for _, descendant := range descendants {
		r.descendants[descendant.ID] = descendant
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

func NewAttribute(resource, action string) Attribute {
	return Attribute{
		Resource: resource,
		Action:   action,
	}
}

func NewRole(id string, lock ...bool) *Role {
	var disable bool
	if len(lock) > 0 {
		disable = lock[0]
	}
	return &Role{
		ID:          id,
		permissions: make(map[string]Attribute),
		descendants: make(map[string]*Role),
		lock:        disable,
	}
}
