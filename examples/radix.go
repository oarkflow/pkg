package main

import (
	"fmt"

	"github.com/oarkflow/pkg/permission/radix"
)

func main() {
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
