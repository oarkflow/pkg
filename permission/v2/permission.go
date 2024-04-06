package v2

import (
	"github.com/oarkflow/maps"
)

func Can(userID, company, module, entity, group, activity string) bool {
	var allowed []string
	if company == "" {
		return false
	}
	companyUser := GetUserRoles(company, userID)
	if companyUser == nil {
		return false
	}
	var userRoles []*Role
	roles := GetAllowedRoles(companyUser, module, entity)
	companyUser.Company.Roles.ForEach(func(_ string, r *Role) bool {
		for _, rt := range roles {
			if r.ID == rt {
				userRoles = append(userRoles, r)
			}
		}
		allowed = append(allowed, r.ID)
		return true
	})
	for _, role := range userRoles {
		if role.Has(group, activity, allowed...) {
			return true
		}
	}
	return false
}
func AddRole(role *Role) {
	roleManager.roles.Set(role.ID, role)
}
func GetRole(role string) (*Role, bool) {
	return roleManager.roles.Get(role)
}
func Roles() map[string]*Role {
	return roleManager.roles.AsMap()
}
func AddUserRole(userID string, roleID string, company *Company, module *Module, entity *Entity) {
	roleManager.AddUserRole(userID, roleID, company, module, entity)
}
func GetCompanyUserRoles(company string) *CompanyUser {
	return roleManager.GetCompanyUserRoles(company)
}
func GetUserRoles(company, userID string) *CompanyUser {
	return roleManager.GetUserRoles(company, userID)
}
func GetUserRolesByCompany(company string) []*UserRole {
	return roleManager.GetUserRolesByCompany(company)
}
func GetUserRoleByCompanyAndUser(company, userID string) (ut []*UserRole) {
	return roleManager.GetUserRoleByCompanyAndUser(company, userID)
}
func GetAllowedRoles(userRoles *CompanyUser, module, entity string) []string {
	return roleManager.GetAllowedRoles(userRoles, module, entity)
}
func AddCompany(data *Company) {
	roleManager.AddCompany(data)
}
func GetCompany(id string) (*Company, bool) {
	return roleManager.GetCompany(id)
}
func Companies() map[string]*Company {
	return roleManager.Companies()
}
func AddModule(data *Module) {
	roleManager.AddModule(data)
}
func GetModule(id string) (*Module, bool) {
	return roleManager.GetModule(id)
}
func Modules() map[string]*Module {
	return roleManager.Modules()
}
func AddUser(data *User) {
	roleManager.AddUser(data)
}
func GetUser(id string) (*User, bool) {
	return roleManager.GetUser(id)
}
func Users() map[string]*User {
	return roleManager.Users()
}
func AddEntity(data *Entity) {
	roleManager.AddEntity(data)
}
func GetEntity(id string) (*Entity, bool) {
	return roleManager.GetEntity(id)
}
func Entities() map[string]*Entity {
	return roleManager.Entities()
}
func NewCompany(id string) *Company {
	company := &Company{
		ID:          id,
		Modules:     maps.New[string, *Module](),
		Roles:       maps.New[string, *Role](),
		Entities:    maps.New[string, *Entity](),
		descendants: maps.New[string, *Company](),
	}
	AddCompany(company)
	return company
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
func NewEntity(id string) *Entity {
	entity := &Entity{ID: id}
	AddEntity(entity)
	return entity
}
func NewRole(id string, lock ...bool) *Role {
	var disable bool
	if len(lock) > 0 {
		disable = lock[0]
	}
	role := &Role{
		ID:          id,
		permissions: maps.New[string, *AttributeGroup](),
		descendants: maps.New[string, *Role](),
		lock:        disable,
	}
	AddRole(role)
	return role
}
func NewAttribute(resource, action string) Attribute {
	return Attribute{
		Resource: resource,
		Action:   action,
	}
}
func NewUser(id string) *User {
	user := &User{ID: id}
	AddUser(user)
	return user
}
