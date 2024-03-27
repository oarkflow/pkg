package v2

import (
	"testing"
)

func myRoles() (coder *Role, qa *Role, suspendManager *Role) {
	coder = NewRole("coder")
	permission := []Attribute{
		{"/coding/:wid/:eid/start-coding", "POST"},
		{"/coding/:wid/open", "GET"},
		{"/coding/:wid/in-progress", "GET"},
		{"/coding/:wid/:eid/review", "POST"},
	}
	coder.AddPermission("route", permission...)

	qa = NewRole("qa")
	permission = []Attribute{
		{"/coding/:wid/:eid/start-qa", "POST"},
		{"/coding/:wid/qa", "GET"},
		{"/coding/:wid/qa-in-progress", "GET"},
		{"/coding/:wid/:eid/qa-review", "POST"},
	}
	qa.AddPermission("route", permission...)

	suspendManager = NewRole("suspend-manager")
	permission = []Attribute{
		{"/coding/:wid/suspended", "GET"},
		{"/coding/:wid/:eid/release-suspend", "POST"},
		{"/coding/:wid/:eid/request-abandon", "POST"},
	}
	suspendManager.AddPermission("route", permission...)
	return
}

func BenchmarkMatchResource(b *testing.B) {
	company := NewCompany("Edelberg")
	module := NewModule("Coding")
	company.AddModule(module)
	// company.SetDefaultModule(module.ID)

	coder, qa, suspendManager := myRoles()
	company.AddRole(coder, qa, suspendManager)

	e29 := &Entity{ID: "29"}
	e30 := &Entity{ID: "30"}
	e33 := &Entity{ID: "33"}

	company.AddEntities(e29, e30, e33)

	sujit := &User{ID: "sujit"}
	alex := &User{ID: "alex"}
	josh := &User{ID: "josh"}

	company.AddUser(sujit.ID, coder.ID)
	company.AddUser(alex.ID, qa.ID)
	company.AddUser(josh.ID, suspendManager.ID)

	company.AssignEntitiesToUser(sujit.ID, e29.ID)
	company.AssignEntitiesToUser(alex.ID, e30.ID)
	company.AssignEntitiesToUser(josh.ID, e33.ID)
	for i := 0; i < b.N; i++ {
		sujit.Can("Edelberg", "Coding", e29.ID, "route", "/coding/1/2/start-coding POST")
	}
}
