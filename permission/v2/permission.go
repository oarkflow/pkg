package v2

type Company struct {
	ID            string
	defaultModule *Module
	Modules       map[string]*Module
	Roles         map[string]*Role
	Entities      map[string]*Entity
}

func NewCompany(id string) *Company {
	return &Company{
		ID:       id,
		Modules:  make(map[string]*Module),
		Roles:    make(map[string]*Role),
		Entities: make(map[string]*Entity),
	}
}

func (c *Company) SetDefaultModule(module string) {
	if mod, ok := c.Modules[module]; ok {
		c.defaultModule = mod
	}
}

func (c *Company) AddModule(modules ...*Module) {
	for _, module := range modules {
		c.Modules[module.ID] = module
	}
}

func (c *Company) AddRole(roles ...*Role) {
	for _, role := range roles {
		c.Roles[role.ID] = role
	}
}

func (c *Company) AddEntities(entities ...*Entity) {
	for _, entity := range entities {
		c.Entities[entity.ID] = entity
		if c.defaultModule != nil {
			c.defaultModule.Entities[entity.ID] = entity
		}
	}
}

func (c *Company) AddEntitiesToModule(module string, entities ...string) {
	for _, id := range entities {
		entity, exists := c.Entities[id]
		if !exists {
			return
		}
		if mod, ok := c.Modules[module]; ok {
			mod.Entities[id] = entity
		} else {
			return
		}
	}
}

func (c *Company) AddRolesToModule(module string, roles ...string) {
	for _, id := range roles {
		role, exists := c.Roles[id]
		if !exists {
			return
		}
		if mod, ok := c.Modules[module]; ok {
			mod.Roles[id] = role
		} else {
			return
		}
	}
}

func (c *Company) AddUser(user, role string) {
	if _, ok := c.Roles[role]; ok {
		RoleManager.AddUserRole(user, role, c, nil, nil)
		if c.defaultModule != nil {
			RoleManager.AddUserRole(user, role, c, c.defaultModule, nil)
		}
	}
}

func (c *Company) AddUserInModule(user, module string, roles ...string) {
	mod, exists := c.Modules[module]
	if !exists {
		return
	}
	if len(roles) > 0 {
		for _, r := range roles {
			if role, ok := c.Roles[r]; ok {
				RoleManager.AddUserRole(user, role.ID, c, mod, nil)
			}
		}
	} else {
		for _, ur := range RoleManager.GetUserRolesByCompany(c.ID) {
			if ur.UserID == user && ur.Module != nil && ur.Module.ID != module {
				RoleManager.AddUserRole(user, ur.RoleID, c, mod, nil)
			}
		}
	}
}

func (c *Company) AssignEntitiesToUser(userID string, entities ...string) {
	user := RoleManager.GetUserRoles(c.ID, userID)
	if user == nil {
		return
	}
	for _, role := range user.UserRoles {
		for _, id := range entities {
			if entity, ok := c.Entities[id]; ok {
				RoleManager.AddUserRole(userID, role.RoleID, c, nil, entity)
				if c.defaultModule != nil {
					RoleManager.AddUserRole(userID, role.RoleID, c, c.defaultModule, entity)
				}
			}
		}
	}
}

func (c *Company) AssignEntitiesWithRole(userID, roleId string, entities ...string) {
	if len(entities) == 0 {
		return
	}
	user := RoleManager.GetUserRoles(c.ID, userID)
	if user == nil {
		return
	}
	_, ok := c.Roles[roleId]
	if !ok {
		return
	}
	for _, id := range entities {
		if entity, ok := c.Entities[id]; ok {
			RoleManager.AddUserRole(userID, roleId, c, nil, entity)
			if c.defaultModule != nil {
				RoleManager.AddUserRole(userID, roleId, c, c.defaultModule, entity)
			}
		}
	}
}

type Module struct {
	ID       string
	Roles    map[string]*Role   // Subset or all roles from company
	Entities map[string]*Entity // Subset or all entities from company
}

func NewModule(id string) *Module {
	return &Module{
		ID:       id,
		Roles:    make(map[string]*Role),
		Entities: make(map[string]*Entity),
	}
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
