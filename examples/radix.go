package main

import (
	"fmt"
)

type Allowed interface {
	Can(string) bool
}

// Permission defines a single permission for an activity
type Permission struct {
	Name string
}

func NewRole(name string) *Role {
	return &Role{
		Name:        name,
		permissions: make(map[string]Permission),
		descendants: make(map[string]*Role),
	}
}

// Role represents a user role with its permissions
type Role struct {
	Name        string
	permissions map[string]Permission
	descendants map[string]*Role
}

func (r *Role) Has(permissionName string) bool {
	if _, ok := r.permissions[permissionName]; ok {
		return true
	}
	// Check inherited permissions recursively
	for _, descendant := range r.getDescendantRoles() {
		if descendant.Has(permissionName) {
			return true
		}
	}
	return false
}

func (r *Role) getDescendantRoles() []*Role {
	var descendants []*Role
	for _, child := range r.descendants { // Simulate descendent-child relationship through permissions
		descendants = append(descendants, child)
		descendants = append(descendants, child.getDescendantRoles()...)
	}
	return descendants
}

// AddDescendent adds a new permission to the role
func (r *Role) AddDescendent(descendants ...*Role) {
	for _, descendant := range descendants {
		r.descendants[descendant.Name] = descendant
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
	Name      string
	Roles     []*Role
	WorkItems map[string][]WorkItem
}

// Can check if a user is allowed to do an activity based on their role and inherited permissions
func (u *User) Can(activity string) bool {
	for _, role := range u.Roles {
		if role.Has(activity) {
			return true
		}
	}
	return false
}

func (u *User) Assign(roles ...*Role) {
	if len(roles) == 0 {
		return
	}
	u.Roles = append(u.Roles, roles...)
}

type WorkItem struct {
	ID string
}

type Organization struct {
	Name      string
	WorkItems []WorkItem
}

func NewUser(name string) *User {
	return &User{
		Name:      name,
		Roles:     nil,
		WorkItems: make(map[string][]WorkItem),
	}
}

func main() {
	coderRole := NewRole("Coder")
	coderRole.AddPermission(Permission{Name: "code"})

	qaRole := NewRole("QA")
	qaRole.AddPermission(Permission{Name: "qa"})

	suspendManagerRole := NewRole("SuspendManager")
	suspendManagerRole.AddPermission(Permission{Name: "suspend"})

	adminRole := NewRole("Admin")
	adminRole.AddPermission(Permission{Name: "add-user"})

	accountManagerRole := NewRole("AccountManager")
	accountManagerRole.AddPermission(Permission{"add-company"})

	adminRole.AddDescendent(coderRole, qaRole, suspendManagerRole)
	accountManagerRole.AddDescendent(adminRole)

	userA := NewUser("userA")
	userA.Assign(coderRole)
	userB := NewUser("userB")
	userB.Assign(qaRole)
	userC := NewUser("userC")
	userC.Assign(adminRole)
	userD := NewUser("userD")
	userD.Assign(accountManagerRole)

	// Check permissions
	fmt.Println(userA.Name, "can code:", userA.Can("code"))
	fmt.Println(userB.Name, "can suspend:", userB.Can("suspend"))
	fmt.Println(userC.Name, "can create user:", userC.Can("add-user")) // Inherited from AccountManager
	fmt.Println(userD.Name, "can qa:", userD.Can("qa"))
	fmt.Println(userD.Name, "can qa:", userD.Can("delete-user"))

	// Add a new permission dynamically (inherited by Admin)
	newPermission := Permission{Name: "delete-user"}

	adminRole.AddPermission(newPermission)

	fmt.Println(userC.Name, "can delete user (after adding permission to Admin):", userC.Can("delete-user"))
	fmt.Println(userD.Name, "can qa user (after adding permission to Admin):", userD.Can("delete-user"))
}
