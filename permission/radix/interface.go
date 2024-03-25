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
	TopUpEntityToUser(userID, entityID, roleID string)
	GetModule(name string) (*Module, bool)
	AddModule(mod *Module, defaultModule, copyUserRoles, copyEntities bool)
	AddEntityToModule(module, entityID string) error
	AddUserToModule(module string, user IUser, roles ...string) error
	AssignEntityToUser(userID string, entityIDs []string, assignToModules bool)
	AssignEntityToUserInModules(userID string, entityIDs []string, modules []string)
	AddEntity(entity ...string)
	Entities() map[string]bool
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
