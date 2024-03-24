package radix

type PermittedUser interface {
	Can(string) bool
}

type ICompany interface {
	AddUser(user IUser, role string) error
	AddRole(roles ...IRole)
	Roles() map[string]IRole
	Users() []*UserRole
	GetModule(name string) (*Module, bool)
	AddModule(module *Module, copyUserRoles, copyEntities bool) *Module
	AddEntity(id string, entity *Entity)
}

type IModule interface {
	AddUser(user IUser, role string) error
	AddRole(roles ...IRole)
	Roles() map[string]IRole
	AddEntity(id string, entity *Entity)
}

type IRole interface {
	Name() string
	Has(permission string, allowedDescendants ...string) bool
	Lock()
	Unlock()
	GetDescendantRoles() []IRole
	AddDescendent(descendants ...IRole) error
	AddPermission(permissions ...Attribute) error
}

type IUser interface {
	PermittedUser
	WithCompany(company ICompany, module ...string) IUser
	Assign(roles ...IRole)
	Name() string
	Roles() []IRole
}
