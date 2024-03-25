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
	id            string
	users         []*UserRole
	roles         map[string]IRole
	modules       map[string]*Module
	entities      map[string]bool
	userEntities  map[string][]string
	defaultModule *Module
}

func (c *Company) ID() string {
	return c.id
}

func (c *Company) AddUser(user IUser, roleID string) error {
	if role, ok := c.roles[roleID]; ok {
		c.users = append(c.users, &UserRole{
			User: user,
			Role: role,
		})
		user.AssignTo(c)
		user.Assign(role)
		if c.defaultModule != nil {
			c.AddUserToModule(c.defaultModule.id, user, roleID)
		}
		return nil
	}
	return errors.New("role not available for company")
}

func (c *Company) Roles() map[string]IRole {
	return c.roles
}

func (c *Company) SetDefaultModule(mod string) {
	if module, ok := c.modules[mod]; ok {
		c.defaultModule = module
	}
}

func (c *Company) Users() []*UserRole {
	return c.users
}

func (c *Company) AddRole(roles ...IRole) {
	for _, role := range roles {
		c.roles[role.ID()] = role
		for _, module := range c.modules {
			module.roles[role.ID()] = role
		}
	}
}

func (c *Company) GetModule(name string) (*Module, bool) {
	mod, ok := c.modules[name]
	return mod, ok
}

func (c *Company) AddModule(mod *Module, defaultModule, copyUserRoles, copyEntities bool) {
	module := NewModule(mod.id)
	if copyUserRoles {
		module.roles = c.roles
	}
	if copyEntities {
		for id := range c.entities {
			module.entities[id] = true
		}
	}
	c.modules[module.id] = module
	if defaultModule {
		c.SetDefaultModule(module.id)
	}
}

func (c *Company) AddEntity(entities ...string) {
	for _, entity := range entities {
		c.entities[entity] = true
		if c.defaultModule != nil {
			c.defaultModule.entities[entity] = true
		}
	}
}

func (c *Company) Entities() map[string]bool {
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
	_, ok = c.entities[entityID]
	if !ok {
		return errors.New("entity not available for company")
	}
	mod.entities[entityID] = true
	return nil
}

func (c *Company) AssignEntityToUser(userID string, entityIDs []string) {
	var entities []string
	for id := range c.entities {
		if slices.Contains(entityIDs, id) {
			entities = append(entities, id)
		}
	}
	if len(entities) > 0 {
		c.userEntities[userID] = append(c.userEntities[userID], entities...)
	}
	if c.defaultModule != nil {
		c.defaultModule.userEntities = c.userEntities
	}
}

func (c *Company) AssignEntityToUserInModules(userID string, entityIDs []string, modules []string) {
	var entities []string
	if len(modules) == 0 {
		return
	}
	for _, mod := range modules {
		if module, ok := c.modules[mod]; ok {
			for id := range module.entities {
				if slices.Contains(entityIDs, id) {
					entities = append(entities, id)
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
	id           string
	users        []*UserRole
	roles        map[string]IRole
	entities     map[string]bool
	userEntities map[string][]string
}

func (m *Module) ID() string {
	return m.id
}
