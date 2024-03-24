package main

import (
	"fmt"

	"github.com/oarkflow/pkg/permission/radix"
)

func main() {
	coderRole := radix.NewRole("Coder")
	coderRole.AddPermission(radix.NewPermission("code"))

	qaRole := radix.NewRole("QA")
	qaRole.AddPermission(radix.NewPermission("qa"))

	suspendManagerRole := radix.NewRole("SuspendManager")
	suspendManagerRole.AddPermission(radix.NewPermission("suspend"))

	adminRole := radix.NewRole("Admin")
	adminRole.AddPermission(radix.NewPermission("add-user"))

	accountManagerRole := radix.NewRole("AccountManager")
	accountManagerRole.AddPermission(radix.NewPermission("add-company"))

	adminRole.AddDescendent(coderRole, qaRole, suspendManagerRole)
	accountManagerRole.AddDescendent(adminRole)

	userA := radix.NewUser("userA")
	userA.Assign(coderRole)
	userB := radix.NewUser("userB")
	userB.Assign(qaRole)
	userC := radix.NewUser("userC")
	userC.Assign(adminRole)
	userD := radix.NewUser("userD")
	userD.Assign(accountManagerRole)

	// Check permissions
	fmt.Println(userA.Name(), "can code:", userA.Can("code"))
	fmt.Println(userB.Name(), "can suspend:", userB.Can("suspend"))
	fmt.Println(userC.Name(), "can create user:", userC.Can("add-user")) // Inherited from AccountManager
	fmt.Println(userD.Name(), "can qa:", userD.Can("qa"))
	fmt.Println(userD.Name(), "can qa:", userD.Can("delete-user"))

	// Add a new permission dynamically (inherited by Admin)
	newPermission := radix.NewPermission("delete-user")

	adminRole.AddPermission(newPermission)

	fmt.Println(userC.Name(), "can delete user (after adding permission to Admin):", userC.Can("delete-user"))
	fmt.Println(userD.Name(), "can qa user (after adding permission to Admin):", userD.Can("delete-user"))
}
