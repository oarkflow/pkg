package radix

func NewAttribute(resource, action string) Attribute {
	return Attribute{
		Resource: resource,
		Action:   action,
	}
}

func NewRole(name string, lock ...bool) IRole {
	var disable bool
	if len(lock) > 0 {
		disable = lock[0]
	}
	return &Role{
		name:        name,
		permissions: make(map[string]Attribute),
		descendants: make(map[string]IRole),
		lock:        disable,
	}
}

func NewUser(name string) IUser {
	return &User{
		name: name,
	}
}

func NewCompany(name string) *Company {
	return &Company{
		id:           name,
		roles:        make(map[string]IRole),
		modules:      make(map[string]*Module),
		entities:     make(map[string]*Entity),
		userEntities: make(map[string][]string),
	}
}

func NewModule(name string) *Module {
	return &Module{
		Name:     name,
		roles:    make(map[string]IRole),
		entities: make(map[string]string),
	}
}
