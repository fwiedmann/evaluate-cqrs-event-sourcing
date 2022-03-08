package core

import "github.com/google/uuid"

func NewUser(name, surname string) (User, error) {
	user := User{
		Id:      uuid.NewString(),
		Name:    name,
		Surname: surname,
	}
	return user, user.Validate()
}

type User struct {
	Id      string
	Name    string
	Surname string
}

func (u User) Validate() error {
	return nil
}
