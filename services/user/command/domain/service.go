package domain

import "github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user/core"

type UserRepository interface {
	Append(command string, user core.User) error
}

type EventBus interface {
	PublishCreateEvent(event core.UserCreatedEvent) error
	PublishUpdated(event core.UserUpdatedEvent) error
	PublishDeleteEvent(event core.UserDeletedEvent) error
}

func NewUserService(repo UserRepository, bus EventBus) *UserService {
	return &UserService{
		repo:     repo,
		eventBus: bus,
	}
}

type UserService struct {
	repo     UserRepository
	eventBus EventBus
}

type CreateUser struct {
	Name    string
	Surname string
}

const (
	CreateCommand = "userCreate"
	UpdateCommand = "userUpdate"
	DeleteCommand = "userDelete"
)

func (us *UserService) Create(user CreateUser) (core.User, error) {
	u, err := core.NewUser(user.Name, user.Surname)
	if err != nil {
		return core.User{}, err
	}
	if err := us.repo.Append(CreateCommand, u); err != nil {
		return core.User{}, err
	}
	if err := us.eventBus.PublishUpdated(core.UserUpdatedEvent{User: u}); err != nil {
		return core.User{}, err
	}
	return u, nil
}

type UpdateUser struct {
	Id      string
	Name    string
	Surname string
}

func (us *UserService) Update(user UpdateUser) (core.User, error) {
	u := core.User{
		Id:      user.Id,
		Name:    user.Name,
		Surname: user.Surname,
	}

	if err := us.repo.Append(UpdateCommand, u); err != nil {
		return core.User{}, err
	}
	if err := us.eventBus.PublishUpdated(core.UserUpdatedEvent{User: u}); err != nil {
		return core.User{}, err
	}

	return u, nil
}

func (us *UserService) Delete(id string) error {
	u := core.User{
		Id: id,
	}

	if err := us.repo.Append(DeleteCommand, u); err != nil {
		return err
	}

	return us.eventBus.PublishDeleteEvent(core.UserDeletedEvent{UserId: u.Id})
}
