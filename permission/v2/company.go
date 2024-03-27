package v2

import (
	"github.com/oarkflow/pkg/maps"
)

type Company struct {
	ID            string
	defaultModule *Module
	Modules       maps.IMap[string, *Module]
	Roles         maps.IMap[string, *Role]
	Entities      maps.IMap[string, *Entity]
}

func NewCompany(id string) *Company {
	company := &Company{
		ID:       id,
		Modules:  maps.New[string, *Module](),
		Roles:    maps.New[string, *Role](),
		Entities: maps.New[string, *Entity](),
	}
	AddCompany(company)
	return company
}

func (c *Company) SetDefaultModule(module string) {
	if mod, ok := c.Modules.Get(module); ok {
		c.defaultModule = mod
	}
}

func (c *Company) AddModule(modules ...*Module) {
	for _, module := range modules {
		c.Modules.Set(module.ID, module)
	}
}

func (c *Company) AddRole(roles ...*Role) {
	for _, role := range roles {
		c.Roles.Set(role.ID, role)
	}
}

func (c *Company) AddEntities(entities ...*Entity) {
	for _, entity := range entities {
		c.Entities.Set(entity.ID, entity)
		if c.defaultModule != nil {
			c.defaultModule.Entities.Set(entity.ID, entity)
		}
	}
}

func (c *Company) AddEntitiesToModule(module string, entities ...string) {
	for _, id := range entities {
		entity, exists := c.Entities.Get(id)
		if !exists {
			return
		}
		if mod, ok := c.Modules.Get(module); ok {
			mod.Entities.Set(id, entity)
		} else {
			return
		}
	}
}

func (c *Company) AddRolesToModule(module string, roles ...string) {
	for _, id := range roles {
		role, exists := c.Roles.Get(id)
		if !exists {
			return
		}
		if mod, ok := c.Modules.Get(module); ok {
			mod.Roles.Set(id, role)
		} else {
			return
		}
	}
}

func (c *Company) AddUser(user, role string) {
	if _, ok := c.Roles.Get(role); ok {
		AddUserRole(user, role, c, nil, nil)
		if c.defaultModule != nil {
			AddUserRole(user, role, c, c.defaultModule, nil)
		}
	}
}

func (c *Company) AddUserInModule(user, module string, roles ...string) {
	mod, exists := c.Modules.Get(module)
	if !exists {
		return
	}
	if len(roles) > 0 {
		for _, r := range roles {
			if role, ok := c.Roles.Get(r); ok {
				AddUserRole(user, role.ID, c, mod, nil)
			}
		}
	} else {
		for _, ur := range GetUserRolesByCompany(c.ID) {
			if ur.UserID == user && ur.Module != nil && ur.Module.ID != module {
				AddUserRole(user, ur.RoleID, c, mod, nil)
			}
		}
	}
}

func (c *Company) AssignEntitiesToUser(userID string, entities ...string) {
	user := GetUserRoles(c.ID, userID)
	if user == nil {
		return
	}
	for _, role := range user.Roles {
		for _, id := range entities {
			if entity, ok := c.Entities.Get(id); ok {
				AddUserRole(userID, role.RoleID, c, nil, entity)
				if c.defaultModule != nil {
					AddUserRole(userID, role.RoleID, c, c.defaultModule, entity)
				}
			}
		}
	}
}

func (c *Company) AssignEntitiesWithRole(userID, roleId string, entities ...string) {
	if len(entities) == 0 {
		return
	}
	user := GetUserRoles(c.ID, userID)
	if user == nil {
		return
	}
	_, ok := c.Roles.Get(roleId)
	if !ok {
		return
	}
	for _, id := range entities {
		if entity, ok := c.Entities.Get(id); ok {
			AddUserRole(userID, roleId, c, nil, entity)
			if c.defaultModule != nil {
				AddUserRole(userID, roleId, c, c.defaultModule, entity)
			}
		}
	}
}

type Module struct {
	ID       string
	Roles    maps.IMap[string, *Role]
	Entities maps.IMap[string, *Entity]
}

func NewModule(id string) *Module {
	module := &Module{
		ID:       id,
		Roles:    maps.New[string, *Role](),
		Entities: maps.New[string, *Entity](),
	}
	AddModule(module)
	return module
}

type UserRole struct {
	UserID  string
	RoleID  string
	Company *Company
	Module  *Module
	Entity  *Entity
}
type Entity struct {
	ID string
}

func NewEntity(id string) *Entity {
	entity := &Entity{ID: id}
	AddEntity(entity)
	return entity
}
