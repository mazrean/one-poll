package domain

import "github.com/cs-sysimpl/suzukake/domain/values"

type User struct {
	id             values.UserID
	name           values.UserName
	hashedPassword values.UserHashedPassword
}

func NewUser(
	id values.UserID,
	name values.UserName,
	hashedPassword values.UserHashedPassword,
) *User {
	return &User{
		id:             id,
		name:           name,
		hashedPassword: hashedPassword,
	}
}

func (u *User) GetID() values.UserID {
	return u.id
}

func (u *User) GetName() values.UserName {
	return u.name
}

func (u *User) SetName(name values.UserName) {
	u.name = name
}

func (u *User) GetHashedPassword() values.UserHashedPassword {
	return u.hashedPassword
}

func (u *User) SetHashedPassword(hashedPassword values.UserHashedPassword) {
	u.hashedPassword = hashedPassword
}
