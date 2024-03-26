package v2

import (
	"fmt"
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
	if entity != "" && module != "" {

	}
	var otherRole, userModuleRole, userCompanyEntityRole, userCompanyRole, userModuleEntityRole []*UserRole
	for _, userRole := range userRoles.UserRoles {
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
	if module != "" && entity != "" {
		if len(userModuleEntityRole) > 0 {
			found := false
			for _, r := range userModuleEntityRole {
				if r.Module.ID == module && r.Entity.ID == entity {
					found = true
					ut[r.RoleID] = r.RoleID
				}
			}
			if found {

			}
		}
	}
	for _, r := range otherRole {
		fmt.Println(r.RoleID, r.Company, r.Module, r.Entity)
	}
	/*fmt.Println(userCompanyRole)
	fmt.Println(userCompanyEntityRole)
	fmt.Println(userModuleEntityRole)
	fmt.Println(otherRole)*/
	if entity != "" && module != "" {
		entityModFound := false
		for _, ur := range userRoles.UserRoles {
			if ur.Module != nil && ur.Module.ID == module && ur.Entity != nil && ur.Entity.ID == entity {
				entityModFound = true
				ut[ur.RoleID] = ur.RoleID
			}
		}
		if !entityModFound {
			/*for _, ur := range userRoles.UserRoles {
				if ur.Module == nil && ur.Entity != nil && ur.Entity.ID == entity {
					ut[ur.RoleID] = ur.RoleID
				}
			}*/
		}
	}
	/*for _, ur := range userRoles.UserRoles {
		if ur.Entity != nil && ur.Entity.ID == entity {
			// Entity role found, return roles
			for _, ur := range userRoles {
				if ur.Entity != nil && ur.Entity.ID == entity {
					ut[ur.RoleID] = ur.RoleID
				}
			}
			return ut
		} else if ur.Module != nil && ur.Module.ID == module && ur.Entity == nil {
			// Module role found, return roles (including entity roles within module)
			for _, ur := range userRoles {
				if (ur.Module != nil && ur.Module.ID == module) || (ur.Entity != nil && ur.Entity.ID == entity) {
					ut[ur.RoleID] = ur.RoleID
				}
			}
			return ut
		} else if ur.Company != nil {
			// Company role found, return roles
			for _, ur := range userRoles {
				if ur.Company != nil {
					ut[ur.RoleID] = ur.RoleID
				}
			}
			return ut
		}
	}*/
	return ut
}
