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
	company.AddModule(module, false, true, true)

	coderRole, _, suspendManagerRole, adminRole, _ := addRoles()
	company.AddRole(adminRole)
	company.AddRole(suspendManagerRole)

	userA := radix.NewUser("userA")
	err := company.AddUser(userA, adminRole.ID())
	if err != nil {
		panic(err)
	}
	pattern := `%d) [Roles: %s, Feature: %s, Company: %s, Module: %s, Assigned To Module: %s, Entities: %s, Valid Entities: %s] ===> Expected %s, Actual %v`
	company.AddEntity(entity1, entity2, entity3, entity4)
	company.AddEntityToModule(module.ID(), entity2.ID)
	fmt.Println(fmt.Sprintf(pattern, 1, "Admin", "code", "no", "no", "no", "no", "n/a", "true", panicIfNotExpected(userA.Can("code add"), "true")))
	fmt.Println(fmt.Sprintf(pattern, 2, "Admin", "code", "yes", "no", "no", "no", "n/a", "false", panicIfNotExpected(userA.WithCompany(company.ID()).Can("code add"), "false")))
	fmt.Println(fmt.Sprintf(pattern, 3, "Admin", "code", "yes", "yes", "no", "no", "n/a", "false", panicIfNotExpected(userA.WithCompany(company.ID(), module.ID()).Can("code add"), "false")))
	fmt.Println(fmt.Sprintf(pattern, 3, "Admin", "code", "yes", "yes", "no", "no", "n/a", "false", panicIfNotExpected(userA.WithCompany(company.ID(), module.ID()).Can("suspend release"), "false")))
	company.AddRole(coderRole)
	fmt.Println(fmt.Sprintf(pattern, 4, "Admin&Coder", "code", "yes", "no", "no", "no", "n/a", "true", panicIfNotExpected(userA.WithCompany(company.ID()).Can("code add"), "true")))
	fmt.Println(fmt.Sprintf(pattern, 5, "Admin&Coder", "code", "yes", "yes", "no", "no", "n/a", "false", panicIfNotExpected(userA.WithCompany(company.ID(), module.ID()).Can("code add"), "false")))
	fmt.Println(fmt.Sprintf(pattern, 6, "Admin&Coder", "code", "yes", "no", "no", "yes", "yes", "true", panicIfNotExpected(userA.WithCompany(company.ID()).WithEntity(entity1.ID).Can("code add"), "true")))
	fmt.Println(fmt.Sprintf(pattern, 7, "Admin&Coder", "code", "yes", "no", "no", "yes", "no", "false", panicIfNotExpected(userA.WithCompany(company.ID()).WithEntity("5").Can("code add"), "false")))
	fmt.Println(fmt.Sprintf(pattern, 8, "Admin&Coder", "code", "yes", "yes", "no", "yes", "no", "false", panicIfNotExpected(userA.WithCompany(company.ID(), module.ID()).Can("code add"), "false")))
	fmt.Println(fmt.Sprintf(pattern, 9, "Admin&Coder", "code", "yes", "yes", "no", "yes", "yes", "false", panicIfNotExpected(userA.WithCompany(company.ID(), module.ID()).WithEntity(entity1.ID).Can("code add"), "false")))
	fmt.Println(fmt.Sprintf(pattern, 10, "Admin&Coder", "code", "yes", "yes", "no", "yes", "no", "false", panicIfNotExpected(userA.WithCompany(company.ID(), module.ID()).WithEntity("6").Can("code add"), "false")))
	company.AddUserToModule(module.ID(), userA)
	fmt.Println(fmt.Sprintf(pattern, 3, "Admin", "code", "yes", "yes", "no", "no", "n/a", "false", panicIfNotExpected(userA.WithCompany(company.ID(), module.ID()).Can("suspend release"), "true")))
	fmt.Println(fmt.Sprintf(pattern, 11, "Admin&Coder", "code", "yes", "yes", "yes", "no", "n/a", "true", panicIfNotExpected(userA.WithCompany(company.ID(), module.ID()).Can("code add"), "true")))
	fmt.Println(fmt.Sprintf(pattern, 12, "Admin&Coder", "code", "yes", "yes", "yes", "yes", "yes", "false", panicIfNotExpected(userA.WithCompany(company.ID(), module.ID()).WithEntity(entity1.ID).Can("code add"), "false")))
	fmt.Println(fmt.Sprintf(pattern, 12, "Admin&Coder", "code", "yes", "yes", "yes", "yes", "yes", "true", panicIfNotExpected(userA.WithCompany(company.ID(), module.ID()).WithEntity(entity2.ID).Can("code add"), "true")))
}

func panicIfNotExpected(condition bool, expected string) bool {
	if fmt.Sprintf("%v", condition) != expected {
		panic(fmt.Sprintf("not a match: expected %v, actual: %v", expected, condition))
	}
	return condition
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
	fmt.Println(userA.ID(), "can code:", userA.Can("code add"))
	fmt.Println(userB.ID(), "can suspend:", userB.Can("suspend add"))
	fmt.Println(userC.ID(), "can create user:", userC.Can("user add")) // Inherited from AccountManager
	fmt.Println(userD.ID(), "can qa:", userD.Can("qa add"))
	fmt.Println(userD.ID(), "can qa:", userD.Can("user delete"))

	// Add a new permission dynamically (inherited by Admin)
	newPermission := radix.NewAttribute("user", "delete")

	adminRole.AddPermission(newPermission)

	fmt.Println(userC.ID(), "can delete user (after adding permission to Admin):", userC.Can("user delete"))
	fmt.Println(userD.ID(), "can qa user (after adding permission to Admin):", userD.Can("user delete"))
}
