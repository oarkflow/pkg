package radix

type PermittedUser interface {
	Can(string) bool
}

type ICompany interface {
	ID() string
	AddUser(user IUser, role string) error
	Roles() map[string]IRole
	Users() []*UserRole
	UserEntities() map[string][]string
	AddRole(roles ...IRole)
	GetModule(name string) (*Module, bool)
	AddModule(module *Module, copyUserRoles, copyEntities bool)
	AddEntityToModule(module, entityID string) error
	AddUserToModule(module string, user IUser, roles ...string) error
	AssignEntityToUser(userID string, entityIDs []string)
	AssignEntityToUserInModules(userID string, entityIDs []string, modules []string)
	AddEntity(entity ...*Entity)
	Entities() map[string]*Entity
}

type IRole interface {
	ID() string
	Has(permission string, allowedDescendants ...string) bool
	Lock()
	Unlock()
	GetDescendantRoles() []IRole
	AddDescendent(descendants ...IRole) error
	AddPermission(permissions ...Attribute) error
}

type IUser interface {
	PermittedUser
	WithCompany(company string, module ...string) IUser
	WithEntity(entity ...string) IUser
	AssignTo(company ICompany)
	Assign(roles ...IRole)
	ID() string
	Roles() []IRole
}
