package radix

func NewPermission(name string) Permission {
	return Permission{
		Name: name,
	}
}

func NewRole(name string) IRole {
	return &Role{
		name:        name,
		permissions: make(map[string]Permission),
		descendants: make(map[string]IRole),
	}
}

func NewUser(name string) IUser {
	return &User{
		name: name,
	}
}
