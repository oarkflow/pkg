package v2

import (
	"slices"
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
	userRoles map[string]*CompanyUser
}

func NewUserRoleManager() *UserRoleManager {
	return &UserRoleManager{userRoles: make(map[string]*CompanyUser)}
}

func (u *UserRoleManager) AddUserRole(userID string, roleID string, company *Company, module *Module, entity *Entity) {
	role := &UserRole{
		UserID:  userID,
		RoleID:  roleID,
		Company: company,
		Module:  module,
		Entity:  entity,
	}
	if _, ok := u.userRoles[company.ID]; !ok {
		u.userRoles[company.ID] = &CompanyUser{
			Company: company,
			User:    &User{ID: userID},
		}
	}
	u.userRoles[company.ID].UserRoles = append(u.userRoles[company.ID].UserRoles, role)
}

func (u *UserRoleManager) GetCompanyUserRoles(company string) *CompanyUser {
	userRoles, ok := u.userRoles[company]
	if !ok {
		return nil
	}
	return userRoles
}

func (u *UserRoleManager) GetUserRoles(company, userID string) *CompanyUser {
	userRoles, ok := u.userRoles[company]
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
	userRoles, ok := u.userRoles[company]
	if !ok {
		return nil
	}
	return userRoles.UserRoles
}

func (u *UserRoleManager) GetUserRoleByCompanyAndUser(company, userID string) (ut []*UserRole) {
	userRoles, ok := u.userRoles[company]
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

func (u *UserRoleManager) GetAllowedRoles(userRoles *CompanyUser, module, entity string) map[string]string {
	ut := make(map[string]string)
	if userRoles == nil {
		return ut
	}
	if entity != "" {
		if _, exists := userRoles.Company.Entities[entity]; !exists {
			return ut
		}
	}
	if module != "" {
		if _, exists := userRoles.Company.Modules[module]; !exists {
			return ut
		}
	}
	var moduleEntities, entities []string
	if mod, ok := userRoles.Company.Modules[module]; ok {
		for id := range mod.Entities {
			moduleEntities = append(moduleEntities, id)
		}
	}
	var otherRole, userModuleRole, userCompanyEntityRole, userCompanyRole, userModuleEntityRole []*UserRole
	for _, userRole := range userRoles.UserRoles {
		if userRole.Entity != nil {
			entities = append(entities, userRole.Entity.ID)
		}
		if userRole.Module != nil && userRole.Entity != nil { // if role for module and entity
			userModuleEntityRole = append(userModuleEntityRole, userRole)
		} else if userRole.Module == nil && userRole.Entity == nil { // if role for company
			userCompanyRole = append(userCompanyRole, userRole)
		} else if userRole.Module == nil && userRole.Entity != nil {
			userCompanyEntityRole = append(userCompanyEntityRole, userRole)
		} else if userRole.Module != nil && userRole.Entity == nil {
			userModuleRole = append(userModuleRole, userRole)
		} else {
			otherRole = append(otherRole, userRole)
		}
	}
	for _, r := range userCompanyRole {
		_, ex := r.Company.Entities[entity]
		if !ex {
			return nil
		}
		ut[r.RoleID] = r.RoleID
	}
	if !slices.Contains(entities, entity) && len(userCompanyRole) == 0 {
		return nil
	}
	if len(moduleEntities) > 0 && !slices.Contains(moduleEntities, entity) {
		return nil
	}
	if module != "" && entity != "" {
		if len(userModuleEntityRole) > 0 {
			for _, r := range userModuleEntityRole {
				if r.Module.ID == module && r.Entity.ID == entity {
					ut[r.RoleID] = r.RoleID
				}
			}
		}
	}
	for role := range ut {
		if _, ok := userRoles.Company.Roles[role]; !ok {
			return nil
		}
	}
	return ut
}
