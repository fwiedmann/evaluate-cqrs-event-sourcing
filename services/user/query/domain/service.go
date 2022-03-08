package domain

import (
	"fmt"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user/core"
)

type UserRepository interface {
	FindById(id string) (core.User, error)
	Find() ([]core.User, error)

	Create(user core.User) error
	CreateHistoryEntry(user core.User) error
	Update(user core.User) error
	Delete(id string) error
}

type Event struct {
	User core.User
}

func NewUserService(repo UserRepository) (*UserService, error) {

	us := &UserService{
		repo: repo,
	}

	return us, nil
}

type UserService struct {
	repo UserRepository
}

type CreateUser struct {
	Name    string
	Surname string
}

func (us *UserService) FindById(id string) (core.User, error) {
	return us.repo.FindById(id)
}
func (us *UserService) Find() ([]core.User, error) {
	return us.repo.Find()
}

func (us *UserService) OnUserCreated(event core.UserCreatedEvent) error {
	err := us.repo.Create(event.User)
	if err != nil {
		return err
	}
	err = us.repo.CreateHistoryEntry(event.User)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	return nil
}

func (us *UserService) OnUserUpdated(event core.UserUpdatedEvent) error {
	err := us.repo.Update(event.User)
	if err != nil {
		return err
	}
	err = us.repo.CreateHistoryEntry(event.User)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	return nil
}

func (us *UserService) OnUserDeleted(event core.UserDeletedEvent) error {
	err := us.repo.Delete(event.UserId)
	if err != nil {
		return err
	}
	return nil
}
