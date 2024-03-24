package radix

type PermittedUser interface {
	Can(string) bool
}

type IRole interface {
	Name() string
	Has(string) bool
	GetDescendantRoles() []IRole
	AddDescendent(descendants ...IRole)
	AddPermission(permissions ...Permission)
}

type IUser interface {
	PermittedUser
	Assign(roles ...IRole)
	Name() string
	Roles() []IRole
}
