package main

import (
	"fmt"

	"github.com/oarkflow/pkg/permission/radix"
)

func main() {
	companyRolePermission()
}

func companyRolePermission() {
	entity1 := &radix.Entity{ID: "1"}
	entity2 := &radix.Entity{ID: "2"}
	entity3 := &radix.Entity{ID: "3"}
	entity4 := &radix.Entity{ID: "4"}

	company := radix.NewCompany("Edelberg")

	module := radix.NewModule("Coding")
	company.AddModule(module, true, true)

	coderRole, _, _, adminRole, _ := addRoles()
	company.AddRole(adminRole)

	userA := radix.NewUser("userA")
	err := company.AddUser(userA, adminRole.Name())
	if err != nil {
		panic(err)
	}
	pattern := `[Roles: %s, Feature: %s, Company: %s, Module: %s, Assigned To Module: %s, Entities: %s, Valid Entities: %s] ===> Expected %s, Actual %v`
	company.AddEntity(entity1, entity2, entity3, entity4)
	company.AddEntityToModule(module.Name, entity2.ID)
	fmt.Println(fmt.Sprintf(pattern, "Admin", "code", "no", "no", "no", "no", "n/a", "true", userA.Can("code add")))
	fmt.Println(fmt.Sprintf(pattern, "Admin", "code", "yes", "no", "no", "no", "n/a", "false", userA.WithCompany(company).Can("code add")))
	fmt.Println(fmt.Sprintf(pattern, "Admin", "code", "yes", "yes", "no", "no", "n/a", "false", userA.WithCompany(company, module.Name).Can("code add")))
	company.AddRole(coderRole)
	fmt.Println(fmt.Sprintf(pattern, "Admin&Coder", "code", "yes", "no", "no", "no", "n/a", "true", userA.WithCompany(company).Can("code add")))
	fmt.Println(fmt.Sprintf(pattern, "Admin&Coder", "code", "yes", "yes", "no", "no", "n/a", "false", userA.WithCompany(company, module.Name).Can("code add")))
	fmt.Println(fmt.Sprintf(pattern, "Admin&Coder", "code", "yes", "no", "no", "yes", "yes", "true", userA.WithCompany(company).WithEntity(entity1.ID).Can("code add")))
	fmt.Println(fmt.Sprintf(pattern, "Admin&Coder", "code", "yes", "no", "no", "yes", "no", "false", userA.WithCompany(company).WithEntity("5").Can("code add")))
	fmt.Println(fmt.Sprintf(pattern, "Admin&Coder", "code", "yes", "yes", "no", "yes", "no", "false", userA.WithCompany(company, module.Name).Can("code add")))
	fmt.Println(fmt.Sprintf(pattern, "Admin&Coder", "code", "yes", "yes", "no", "yes", "yes", "false", userA.WithCompany(company, module.Name).WithEntity(entity1.ID).Can("code add")))
	fmt.Println(fmt.Sprintf(pattern, "Admin&Coder", "code", "yes", "yes", "no", "yes", "no", "false", userA.WithCompany(company, module.Name).WithEntity("6").Can("code add")))
	company.AddUserToModule(module.Name, userA)
	fmt.Println(fmt.Sprintf(pattern, "Admin&Coder", "code", "yes", "yes", "yes", "no", "n/a", "true", userA.WithCompany(company, module.Name).Can("code add")))
	fmt.Println(fmt.Sprintf(pattern, "Admin&Coder", "code", "yes", "yes", "yes", "yes", "yes", "false", userA.WithCompany(company, module.Name).WithEntity(entity1.ID).Can("code add")))
	fmt.Println(fmt.Sprintf(pattern, "Admin&Coder", "code", "yes", "yes", "yes", "yes", "yes", "true", userA.WithCompany(company, module.Name).WithEntity(entity2.ID).Can("code add")))
}

func addRoles() (radix.IRole, radix.IRole, radix.IRole, radix.IRole, radix.IRole) {
	coderRole := radix.NewRole("Coder")
	coderRole.AddPermission(radix.NewAttribute("code", "add"))

	qaRole := radix.NewRole("QA")
	qaRole.AddPermission(radix.NewAttribute("qa", "add"))

	suspendManagerRole := radix.NewRole("SuspendManager")
	suspendManagerRole.AddPermission(radix.NewAttribute("suspend", "release"))

	adminRole := radix.NewRole("Admin")
	adminRole.AddPermission(radix.NewAttribute("user", "add"))

	accountManagerRole := radix.NewRole("AccountManager")
	accountManagerRole.AddPermission(radix.NewAttribute("company", "add"))

	adminRole.AddDescendent(coderRole, qaRole, suspendManagerRole)
	accountManagerRole.AddDescendent(adminRole)
	return coderRole, qaRole, suspendManagerRole, adminRole, accountManagerRole
}

func userRolePermission() {
	coderRole, qaRole, _, adminRole, accountManagerRole := addRoles()
	userA := radix.NewUser("userA")
	userA.Assign(coderRole)
	userB := radix.NewUser("userB")
	userB.Assign(qaRole)
	userC := radix.NewUser("userC")
	userC.Assign(adminRole)
	userD := radix.NewUser("userD")
	userD.Assign(accountManagerRole)

	// Check permissions
	fmt.Println(userA.Name(), "can code:", userA.Can("code add"))
	fmt.Println(userB.Name(), "can suspend:", userB.Can("suspend add"))
	fmt.Println(userC.Name(), "can create user:", userC.Can("user add")) // Inherited from AccountManager
	fmt.Println(userD.Name(), "can qa:", userD.Can("qa add"))
	fmt.Println(userD.Name(), "can qa:", userD.Can("user delete"))

	// Add a new permission dynamically (inherited by Admin)
	newPermission := radix.NewAttribute("user", "delete")

	adminRole.AddPermission(newPermission)

	fmt.Println(userC.Name(), "can delete user (after adding permission to Admin):", userC.Can("user delete"))
	fmt.Println(userD.Name(), "can qa user (after adding permission to Admin):", userD.Can("user delete"))
}
