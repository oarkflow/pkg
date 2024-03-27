package v2

import (
	"slices"

	"github.com/oarkflow/pkg/maps"
)

var RoleManager *UserRoleManager

func init() {
	RoleManager = NewUserRoleManager()
}

type CompanyUser struct {
	Company   *Company
	User      *User
	UserRoles []*UserRole
}

type UserRoleManager struct {
	userRoles maps.IMap[string, *CompanyUser]
}

func NewUserRoleManager() *UserRoleManager {
	return &UserRoleManager{userRoles: maps.New[string, *CompanyUser]()}
}

func (u *UserRoleManager) AddUserRole(userID string, roleID string, company *Company, module *Module, entity *Entity) {
	role := &UserRole{
		UserID:  userID,
		RoleID:  roleID,
		Company: company,
		Module:  module,
		Entity:  entity,
	}
	companyUser, ok := u.userRoles.Get(company.ID)
	if !ok {
		companyUser = &CompanyUser{
			Company: company,
			User:    &User{ID: userID},
		}
	}
	companyUser.UserRoles = append(companyUser.UserRoles, role)
	u.userRoles.Set(company.ID, companyUser)
}

func (u *UserRoleManager) GetCompanyUserRoles(company string) *CompanyUser {
	userRoles, ok := u.userRoles.Get(company)
	if !ok {
		return nil
	}
	return userRoles
}

func (u *UserRoleManager) GetUserRoles(company, userID string) *CompanyUser {
	userRoles, ok := u.userRoles.Get(company)
	if !ok {
		return nil
	}
	ur := &CompanyUser{
		Company: userRoles.Company,
	}
	for _, ut := range userRoles.UserRoles {
		if ut.UserID == userID {
			ur.UserRoles = append(ur.UserRoles, ut)
		}
	}
	return ur
}

func (u *UserRoleManager) GetUserRolesByCompany(company string) []*UserRole {
	userRoles, ok := u.userRoles.Get(company)
	if !ok {
		return nil
	}
	return userRoles.UserRoles
}

func (u *UserRoleManager) GetUserRoleByCompanyAndUser(company, userID string) (ut []*UserRole) {
	userRoles, ok := u.userRoles.Get(company)
	if !ok {
		return
	}
	for _, ur := range userRoles.UserRoles {
		if ur.UserID == userID {
			ut = append(ut, ur)
		}
	}
	return
}

type Params struct {
	U string
	C string
	M string
	E string
	R []*UserRole
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

	for _, userRole := range userRoles.UserRoles {
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
