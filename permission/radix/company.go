package radix

import (
	"errors"
	"slices"
)

type UserRole struct {
	User IUser
	Role IRole
}

type Company struct {
	name         string
	users        []*UserRole
	roles        map[string]IRole
	modules      map[string]*Module
	entities     map[string]*Entity
	userEntities map[string][]string
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

func (c *Company) Users() []*UserRole {
	return c.users
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

func (c *Company) AddModule(mod *Module, copyUserRoles, copyEntities bool) {
	module := &Module{
		Name: mod.Name,
	}
	if copyUserRoles {
		module.roles = c.roles
	}
	if copyEntities {
		for id := range c.entities {
			module.entities[id] = id
		}
	}
	c.modules[module.Name] = module
}

func (c *Company) AddEntity(entities ...*Entity) {
	for _, entity := range entities {
		c.entities[entity.ID] = entity
	}
}

func (c *Company) Entities() map[string]*Entity {
	return c.entities
}

func (c *Company) UserEntities() map[string][]string {
	return c.userEntities
}

func (c *Company) AddEntityToModule(module, entityID string) error {
	mod, ok := c.modules[module]
	if !ok {
		return errors.New("module not available for company")
	}
	entity, ok := c.entities[entityID]
	if !ok {
		return errors.New("entity not available for company")
	}
	mod.entities[entityID] = entity.ID
	return nil
}

func (c *Company) AssignEntityToUser(userID string, entityIDs []string) {
	var entities []string
	for id, entity := range c.entities {
		if slices.Contains(entityIDs, id) {
			entities = append(entities, entity.ID)
		}
	}
	if len(entities) > 0 {
		c.userEntities[userID] = append(c.userEntities[userID], entities...)
	}
}

func (c *Company) AssignEntityToUserInModules(userID string, entityIDs []string, modules []string) {
	var entities []string
	if len(modules) == 0 {
		return
	}
	for _, mod := range modules {
		if module, ok := c.modules[mod]; ok {
			for id, entity := range module.entities {
				if slices.Contains(entityIDs, id) {
					entities = append(entities, entity)
				}
			}
			if len(entities) > 0 {
				module.userEntities[userID] = append(module.userEntities[userID], entities...)
			}
		}
	}
}

func (c *Company) AddUserToModule(module string, user IUser, roles ...string) error {
	mod, ok := c.modules[module]
	if !ok {
		return errors.New("module not available for company")
	}
	rolesToAssign := len(roles)
	for name, role := range c.roles {
		if rolesToAssign > 0 && !slices.Contains(roles, name) {
			return errors.New("role not available for company")
		}
		if rolesToAssign > 0 && slices.Contains(roles, name) {
			mod.users = append(mod.users, &UserRole{
				User: user,
				Role: role,
			})
		}
		if rolesToAssign == 0 {
			mod.users = append(mod.users, &UserRole{
				User: user,
				Role: role,
			})
		}
	}
	return nil
}

type Module struct {
	Name         string
	users        []*UserRole
	roles        map[string]IRole
	entities     map[string]string
	userEntities map[string][]string
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
