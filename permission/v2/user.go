package v2

// User represents a user with a role
type User struct {
	ID string
}

// Can check if a user is allowed to do an activity based on their role and inherited permissions
func (u *User) Can(company, module, entity, group, activity string) bool {
	return Can(u.ID, company, module, entity, group, activity)
}

func NewUser(id string) *User {
	user := &User{ID: id}
	AddUser(user)
	return user
}
