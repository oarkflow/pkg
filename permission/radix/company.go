package radix

import (
	"errors"
)

type UserRole struct {
	User IUser
	Role IRole
}

type Company struct {
	name     string
	users    []*UserRole
	roles    map[string]IRole
	modules  map[string]*Module
	entities map[string]*Entity
}

func (c *Company) AddUser(user IUser, role string) error {
	if role, ok := c.roles[role]; ok {
		c.users = append(c.users, &UserRole{
			User: user,
			Role: role,
		})
		user.Assign(role)
		return nil
	}
	return errors.New("role not available for company")
}

func (c *Company) Roles() map[string]IRole {
	return c.roles
}

func (c *Company) AddRole(roles ...IRole) {
	for _, role := range roles {
		c.roles[role.Name()] = role
		for _, module := range c.modules {
			module.roles[role.Name()] = role
		}
	}
}

func (c *Company) GetModule(name string) (*Module, bool) {
	mod, ok := c.modules[name]
	return mod, ok
}

func (c *Company) AddModule(mod *Module, copyUserRoles, copyEntities bool) *Module {
	module := &Module{
		Name: mod.Name,
	}
	if copyUserRoles {
		module.roles = c.roles
	}
	if copyEntities {
		module.entities = c.entities
	}
	c.modules[module.Name] = module
	return module
}

func (c *Company) AddEntity(id string, entity *Entity) {
	c.entities[id] = entity
}

func (c *Company) AddUserToModule(module *Module, user IUser, role string) error {

	return errors.New("module not available for company")
}

type Module struct {
	Name     string
	users    []*UserRole
	roles    map[string]IRole
	entities map[string]*Entity
}

func (c *Module) Roles() map[string]IRole {
	return c.roles
}

func (c *Module) AddUser(user IUser, role string) error {
	if role, ok := c.roles[role]; ok {
		c.users = append(c.users, &UserRole{
			User: user,
			Role: role,
		})
		user.Assign(role)
		return nil
	}
	return errors.New("role not available for module")
}

type Entity struct {
	ID string
}

type CompanyAttribute struct {
	Company   Company
	Module    Module
	Attribute Attribute
	Entities  []Entity
}

type UserCompany struct {
	User    string
	Company Company
	Role    Role
}

type UserPermission struct {
	User     string
	Company  Company
	Module   Module
	Role     Role
	Entities []Entity
}
