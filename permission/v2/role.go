package v2

import (
	"errors"
	"slices"

	"github.com/oarkflow/pkg/maps"
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

type AttributeGroup struct {
	ID          string
	permissions maps.IMap[string, Attribute]
}

// Role represents a user role with its permissions
type Role struct {
	ID          string
	lock        bool
	permissions maps.IMap[string, *AttributeGroup]
	descendants maps.IMap[string, *Role]
}

func (r *Role) Lock() {
	r.lock = true
}

func (r *Role) Unlock() {
	r.lock = false
}

func (r *Role) Has(group, permissionName string, allowedDescendants ...string) bool {
	groupPermissions, ok := r.permissions.Get(group)
	if !ok {
		return false
	}
	if _, ok := groupPermissions.permissions.Get(permissionName); ok {
		return true
	}
	matched := false
	groupPermissions.permissions.ForEach(func(perm string, _ Attribute) bool {
		if MatchResource(permissionName, perm) {
			matched = true
			return false
		}
		return true
	})
	if matched {
		return true
	}
	totalD := len(allowedDescendants)
	// Check inherited permissions recursively
	for _, descendant := range r.GetDescendantRoles() {
		if totalD > 0 {
			if slices.Contains(allowedDescendants, descendant.ID) {
				if descendant.Has(group, permissionName, allowedDescendants...) {
					return true
				}
			}
		} else {
			if descendant.Has(group, permissionName, allowedDescendants...) {
				return true
			}
		}
	}
	return false
}

func (r *Role) GetDescendantRoles() []*Role {
	var descendants []*Role
	r.descendants.ForEach(func(_ string, child *Role) bool {
		descendants = append(descendants, child)
		descendants = append(descendants, child.GetDescendantRoles()...)
		return true
	})
	return descendants
}

// AddDescendent adds a new permission to the role
func (r *Role) AddDescendent(descendants ...*Role) error {
	if r.lock {
		return errors.New("changes not allowed")
	}
	for _, descendant := range descendants {
		r.descendants.Set(descendant.ID, descendant)
	}
	return nil
}

// AddPermission adds a new permission to the role
func (r *Role) AddPermission(group string, permissions ...Attribute) error {
	if r.lock {
		return errors.New("changes not allowed")
	}
	groupAttributes, exists := r.permissions.Get(group)
	if !exists || groupAttributes == nil {
		groupAttributes = &AttributeGroup{
			ID:          group,
			permissions: maps.New[string, Attribute](),
		}
	}
	for _, permission := range permissions {
		groupAttributes.permissions.Set(permission.String(), permission)
	}
	r.permissions.Set(group, groupAttributes)
	return nil
}

func (r *Role) AddPermissionGroup(group *AttributeGroup) error {
	if r.lock {
		return errors.New("changes not allowed")
	}
	r.permissions.Set(group.ID, group)
	return nil
}

func (r *Role) GetGroupPermissions(group string) (permissions []Attribute) {
	if grp, exists := r.permissions.Get(group); exists {
		grp.permissions.ForEach(func(_ string, attr Attribute) bool {
			permissions = append(permissions, attr)
			return true
		})
	}
	return
}

func (r *Role) GetAllImplicitPermissions(perm ...map[string][]Attribute) map[string][]Attribute {
	var grpPermissions map[string][]Attribute
	if len(perm) > 0 {
		grpPermissions = perm[0]
	} else {
		grpPermissions = make(map[string][]Attribute)
	}
	r.permissions.ForEach(func(group string, grp *AttributeGroup) bool {
		var permissions []Attribute
		grp.permissions.ForEach(func(_ string, attr Attribute) bool {
			permissions = append(permissions, attr)
			return true
		})
		grpPermissions[group] = append(grpPermissions[group], permissions...)
		return true
	})
	for _, descendant := range r.GetDescendantRoles() {
		descendant.GetAllImplicitPermissions(grpPermissions)
	}
	return grpPermissions
}

func (r *Role) GetPermissions() map[string][]Attribute {
	grpPermissions := make(map[string][]Attribute)
	r.permissions.ForEach(func(group string, grp *AttributeGroup) bool {
		var permissions []Attribute
		grp.permissions.ForEach(func(_ string, attr Attribute) bool {
			permissions = append(permissions, attr)
			return true
		})
		grpPermissions[group] = permissions
		return true
	})
	return grpPermissions
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
	role := &Role{
		ID:          id,
		permissions: maps.New[string, *AttributeGroup](),
		descendants: maps.New[string, *Role](),
		lock:        disable,
	}
	AddRole(role)
	return role
}
