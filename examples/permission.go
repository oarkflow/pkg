package main

import (
	"fmt"

	v2 "github.com/oarkflow/pkg/permission/v2"
)

func main() {
	company := v2.NewCompany("Edelberg")
	module := v2.NewModule("Coding")
	company.AddModule(module)
	// company.SetDefaultModule(module.ID)

	coder, qa, suspendManager, _ := myRoles()
	company.AddRole(coder, qa, suspendManager)

	e29 := v2.NewEntity("29")
	e30 := v2.NewEntity("30")
	e33 := v2.NewEntity("33")

	company.AddEntities(e29, e30, e33)

	sujit := v2.NewUser("sujit")
	alex := v2.NewUser("alex")
	josh := v2.NewUser("josh")

	company.AddUser(sujit.ID, coder.ID)
	company.AddUser(alex.ID, qa.ID)
	company.AddUser(josh.ID, suspendManager.ID)

	company.AssignEntitiesToUser(sujit.ID, e29.ID)
	company.AssignEntitiesToUser(alex.ID, e30.ID)
	company.AssignEntitiesToUser(josh.ID, e33.ID)
	fmt.Println("R:", v2.Can(sujit.ID, "Edelberg", "Coding", e29.ID, "route", "/coding/1/2/start-coding POST"), "E:", true)
	fmt.Println("R:", v2.Can(sujit.ID, "Edelberg", "Coding", e29.ID, "route", "/coding/1/open GET"), "E:", true)
	fmt.Println("R:", v2.Can(sujit.ID, "Edelberg", "Coding", e29.ID, "backend", "/coding/1/2/start-coding POST"), "E:", false)
}

func myRoles() (coder *v2.Role, qa *v2.Role, suspendManager *v2.Role, admin *v2.Role) {
	coder = v2.NewRole("coder")
	permission := []v2.Attribute{
		{"/coding/:wid/:eid/start-coding", "POST"},
		{"/coding/:wid/open", "GET"},
		{"/coding/:wid/in-progress", "GET"},
		{"/coding/:wid/:eid/review", "POST"},
	}
	coder.AddPermission("route", permission...)

	qa = v2.NewRole("qa")
	permission = []v2.Attribute{
		{"/coding/:wid/:eid/start-qa", "POST"},
		{"/coding/:wid/qa", "GET"},
		{"/coding/:wid/qa-in-progress", "GET"},
		{"/coding/:wid/:eid/qa-review", "POST"},
	}
	qa.AddPermission("route", permission...)

	suspendManager = v2.NewRole("suspend-manager")
	permission = []v2.Attribute{
		{"/coding/:wid/suspended", "GET"},
		{"/coding/:wid/:eid/release-suspend", "POST"},
		{"/coding/:wid/:eid/request-abandon", "POST"},
	}
	suspendManager.AddPermission("route", permission...)

	admin = v2.NewRole("admin")
	permission = []v2.Attribute{
		{"/admin/user/add", "POST"},
		{"/admin/user/edit", "PUT"},
	}
	admin.AddPermission("route", permission...)
	admin.AddDescendent(coder, qa, suspendManager)
	return
}