package v2

import (
	"slices"

	"github.com/oarkflow/pkg/maps"
)

var roleManager *UserRoleManager

func init() {
	roleManager = NewUserRoleManager()
}

// User represents a user with a role
type User struct {
	ID string
}

// Can check if a user is allowed to do an activity based on their role and inherited permissions
func (u *User) Can(company, module, entity, group, activity string) bool {
	return Can(u.ID, company, module, entity, group, activity)
}

type CompanyUser struct {
	Company *Company
	User    *User
	Roles   []*UserRole
}

type UserRoleManager struct {
	companies    maps.IMap[string, *Company]
	modules      maps.IMap[string, *Module]
	entities     maps.IMap[string, *Entity]
	users        maps.IMap[string, *User]
	roles        maps.IMap[string, *Role]
	companyUsers maps.IMap[string, *CompanyUser]
}

func NewUserRoleManager() *UserRoleManager {
	return &UserRoleManager{
		companies:    maps.New[string, *Company](),
		modules:      maps.New[string, *Module](),
		entities:     maps.New[string, *Entity](),
		users:        maps.New[string, *User](),
		roles:        maps.New[string, *Role](),
		companyUsers: maps.New[string, *CompanyUser](),
	}
}

func (u *UserRoleManager) AddRole(role *Role) {
	u.roles.Set(role.ID, role)
}

func (u *UserRoleManager) GetRole(role string) (*Role, bool) {
	return u.roles.Get(role)
}

func (u *UserRoleManager) Roles() map[string]*Role {
	return u.roles.AsMap()
}

func (u *UserRoleManager) AddCompany(data *Company) {
	u.companies.Set(data.ID, data)
}

func (u *UserRoleManager) GetCompany(id string) (*Company, bool) {
	return u.companies.Get(id)
}

func (u *UserRoleManager) Companies() map[string]*Company {
	return u.companies.AsMap()
}

func (u *UserRoleManager) AddModule(data *Module) {
	u.modules.Set(data.ID, data)
}

func (u *UserRoleManager) GetModule(id string) (*Module, bool) {
	return u.modules.Get(id)
}

func (u *UserRoleManager) Modules() map[string]*Module {
	return u.modules.AsMap()
}

func (u *UserRoleManager) AddUser(data *User) {
	u.users.Set(data.ID, data)
}

func (u *UserRoleManager) GetUser(id string) (*User, bool) {
	return u.users.Get(id)
}

func (u *UserRoleManager) Users() map[string]*User {
	return u.users.AsMap()
}

func (u *UserRoleManager) AddEntity(data *Entity) {
	u.entities.Set(data.ID, data)
}

func (u *UserRoleManager) GetEntity(id string) (*Entity, bool) {
	return u.entities.Get(id)
}

func (u *UserRoleManager) Entities() map[string]*Entity {
	return u.entities.AsMap()
}

func (u *UserRoleManager) AddUserRole(userID string, roleID string, company *Company, module *Module, entity *Entity) {
	role := &UserRole{
		UserID:  userID,
		RoleID:  roleID,
		Company: company,
		Module:  module,
		Entity:  entity,
	}
	companyUser, ok := u.companyUsers.Get(company.ID)
	if !ok {
		companyUser = &CompanyUser{
			Company: company,
			User:    &User{ID: userID},
		}
	}
	companyUser.Roles = append(companyUser.Roles, role)
	u.companyUsers.Set(company.ID, companyUser)
}

func (u *UserRoleManager) GetCompanyUserRoles(company string) *CompanyUser {
	userRoles, ok := u.companyUsers.Get(company)
	if !ok {
		return nil
	}
	return userRoles
}

func (u *UserRoleManager) GetUserRoles(company, userID string) *CompanyUser {
	userRoles, ok := u.companyUsers.Get(company)
	if !ok {
		return nil
	}
	ur := &CompanyUser{
		Company: userRoles.Company,
	}
	for _, ut := range userRoles.Roles {
		if ut.UserID == userID {
			ur.Roles = append(ur.Roles, ut)
		}
	}
	return ur
}

func (u *UserRoleManager) GetUserRolesByCompany(company string) []*UserRole {
	userRoles, ok := u.companyUsers.Get(company)
	if !ok {
		return nil
	}
	return userRoles.Roles
}

func (u *UserRoleManager) GetUserRoleByCompanyAndUser(company, userID string) (ut []*UserRole) {
	userRoles, ok := u.companyUsers.Get(company)
	if !ok {
		return
	}
	for _, ur := range userRoles.Roles {
		if ur.UserID == userID {
			ut = append(ut, ur)
		}
	}
	return
}

func (u *UserRoleManager) GetAllowedRoles(userRoles *CompanyUser, module, entity string) []string {
	if userRoles == nil {
		return nil
	}
	// Reusable slices
	moduleEntities := stringSlice.Get()
	moduleRoles := stringSlice.Get()
	entities := stringSlice.Get()
	allowedRoles := stringSlice.Get()
	userCompanyRole := userRoleSlice.Get()
	userModuleEntityRole := userRoleSlice.Get()
	defer func() {
		stringSlice.Put(moduleEntities)
		stringSlice.Put(moduleRoles)
		stringSlice.Put(entities)
		stringSlice.Put(allowedRoles)
		userRoleSlice.Put(userCompanyRole)
		userRoleSlice.Put(userModuleEntityRole)
	}()

	mod, modExists := userRoles.Company.Modules.Get(module)
	_, entExists := userRoles.Company.Entities.Get(entity)
	if (entity != "" && !entExists) || (module != "" && !modExists) {
		return nil
	}

	if modExists {
		mod.Entities.ForEach(func(id string, _ *Entity) bool {
			moduleEntities = append(moduleEntities, id)
			return true
		})
		mod.Roles.ForEach(func(id string, _ *Role) bool {
			moduleRoles = append(moduleRoles, id)
			return true
		})
	}

	for _, userRole := range userRoles.Roles {
		if userRole.Entity != nil {
			entities = append(entities, userRole.Entity.ID)
		}
		if userRole.Module != nil && userRole.Entity != nil { // if role for module and entity
			userModuleEntityRole = append(userModuleEntityRole, userRole)
		} else if userRole.Module == nil && userRole.Entity == nil { // if role for company
			userCompanyRole = append(userCompanyRole, userRole)
		}
	}

	if len(moduleRoles) > 0 {
		for _, modRole := range moduleRoles {
			allowedRoles = append(allowedRoles, modRole)
		}
	} else {
		for _, r := range userCompanyRole {
			allowedRoles = append(allowedRoles, r.RoleID)
		}
	}

	noCompanyEntities := !slices.Contains(entities, entity) && len(userCompanyRole) == 0
	noModuleEntities := len(moduleEntities) > 0 && !slices.Contains(moduleEntities, entity)
	if noCompanyEntities || noModuleEntities {
		return nil
	}

	if module != "" && entity != "" && len(userModuleEntityRole) > 0 {
		for _, r := range userModuleEntityRole {
			if r.Module.ID == module && r.Entity.ID == entity {
				allowedRoles = append(allowedRoles, r.RoleID)
			}
		}
	}

	for _, role := range allowedRoles {
		if _, ok := userRoles.Company.Roles.Get(role); !ok {
			return nil
		}
	}
	return slices.Compact(allowedRoles)
}
