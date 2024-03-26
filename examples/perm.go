package main

import (
	"fmt"

	v2 "github.com/oarkflow/pkg/permission/v2"
)

func main() {
	company, module := v2.NewCompany("Edelberg"), v2.NewModule("Coding")
	company.AddModule(module)
	// company.SetDefaultModule(module.ID)
	entity1, entity2, entity3, entity4 := &v2.Entity{ID: "entity1"}, &v2.Entity{ID: "entity2"}, &v2.Entity{ID: "entity3"}, &v2.Entity{ID: "entity4"}
	company.AddEntities(entity1, entity2, entity3, entity4)
	coder, qa, suspend, admin, _ := roles()
	company.AddRole(coder, qa, suspend, admin)
	user := v2.User{ID: "sujit"}
	company.AddUser(user.ID, admin.ID)
	company.AssignEntitiesWithRole("sujit", coder.ID, entity1.ID, entity2.ID)
	fmt.Println("R:", user.Can("Edelberg", "Coding", entity1.ID, "backend", "qa add"), "E:", true)
	fmt.Println("R:", user.Can("Edelberg", "Coding", entity2.ID, "backend", "qa add"), "E:", true)
	fmt.Println("R:", user.Can("Edelberg", "Coding", entity3.ID, "backend", "user add"), "E:", true)
	fmt.Println("R:", user.Can("Edelberg", "Coding", entity4.ID, "backend", "suspend release"), "E:", true)
	company.AddEntitiesToModule("Coding", entity1.ID, entity2.ID)
	company.AddRolesToModule("Coding", admin.ID, coder.ID, qa.ID)
	fmt.Println("After adding entities to module")
	fmt.Println("R:", user.Can("Edelberg", "Coding", entity1.ID, "backend", "code add"), "E:", true)
	fmt.Println("R:", user.Can("Edelberg", "Coding", entity1.ID, "backend", "qa add"), "E:", true)
	fmt.Println("R:", user.Can("Edelberg", "Coding", entity2.ID, "backend", "qa add"), "E:", true)
	fmt.Println("R:", user.Can("Edelberg", "Coding", entity3.ID, "backend", "user add"), "E:", false)
	fmt.Println("R:", user.Can("Edelberg", "Coding", entity4.ID, "backend", "suspend release"), "E:", false)
}

func roles() (*v2.Role, *v2.Role, *v2.Role, *v2.Role, *v2.Role) {
	coderRole := v2.NewRole("Coder")
	coderRole.AddPermission("backend", v2.NewAttribute("code", "add"))

	qaRole := v2.NewRole("QA")
	qaRole.AddPermission("backend", v2.NewAttribute("qa", "add"))

	suspendManagerRole := v2.NewRole("SuspendManager")
	suspendManagerRole.AddPermission("backend", v2.NewAttribute("suspend", "release"))

	adminRole := v2.NewRole("Admin")
	adminRole.AddPermission("backend", v2.NewAttribute("user", "add"))

	accountManagerRole := v2.NewRole("AccountManager")
	accountManagerRole.AddPermission("backend", v2.NewAttribute("company", "add"))

	adminRole.AddDescendent(coderRole, qaRole, suspendManagerRole)
	accountManagerRole.AddDescendent(adminRole)
	return coderRole, qaRole, suspendManagerRole, adminRole, accountManagerRole
}
