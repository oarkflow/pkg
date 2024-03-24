package radix

func NewAttribute(resource, action string) Attribute {
	return Attribute{
		Resource: resource,
		Action:   action,
	}
}

func NewRole(name string) IRole {
	return &Role{
		name:        name,
		permissions: make(map[string]Attribute),
		descendants: make(map[string]IRole),
	}
}

func NewUser(name string) IUser {
	return &User{
		name: name,
	}
}

func NewCompany(name string) *Company {
	return &Company{
		name:     name,
		roles:    make(map[string]IRole),
		modules:  make(map[string]*Module),
		entities: make(map[string]*Entity),
	}
}
