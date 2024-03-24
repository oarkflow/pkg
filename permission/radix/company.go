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
		return nil
	}
	return errors.New("role not available for company")
}

func (c *Company) AddRole(role IRole) {
	c.roles[role.Name()] = role
}

func (c *Company) Can(role string) bool {
	_, ok := c.roles[role]
	return ok
}

func (c *Company) AddModule(name string, module *Module) {
	c.modules[name] = module
}

func (c *Company) AddModules(modules ...*Module) {
	for _, module := range modules {
		c.modules[module.Name] = module
	}
}

func (c *Company) AddEntity(id string, entity *Entity) {
	c.entities[id] = entity
}

type Module struct {
	Name  string
	users map[string]UserRole
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
