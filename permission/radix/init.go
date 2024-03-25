package radix

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
		id:          id,
		permissions: make(map[string]Attribute),
		descendants: make(map[string]*Role),
		lock:        disable,
	}
}

func NewUser(id string) *User {
	return &User{
		id:        id,
		companies: make(map[string]*Company),
	}
}

func NewCompany(id string) *Company {
	return &Company{
		id:           id,
		roles:        make(map[string]*Role),
		modules:      make(map[string]*Module),
		entities:     make(map[string]bool),
		userEntities: make(map[string][]string),
	}
}

func NewModule(id string) *Module {
	return &Module{
		id:       id,
		roles:    make(map[string]*Role),
		entities: make(map[string]bool),
	}
}
