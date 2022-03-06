package domain

import "fmt"

type UserRepository interface {
	FindById(id string) (User, error)
	Find() ([]User, error)

	Create(user User) error
	Update(user User) error
	Delete(id string) error
}

type EventKind string

type Event struct {
	Kind EventKind
	User User
}

type EventBus interface {
	Subscribe(handler func(e Event) error) error
}

func NewUserService(repo UserRepository, bus EventBus) (*UserService, error) {

	us := &UserService{
		repo:     repo,
		eventBus: bus,
	}

	if err := bus.Subscribe(us.eventHandler); err != nil {
		return nil, err
	}

	return us, nil
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
	CreateEvent = "userCreated"
	UpdateEvent = "userUpdated"
	DeleteEvent = "userDeleted"
)

func (us *UserService) FindById(id string) (User, error) {
	return us.repo.FindById(id)
}
func (us *UserService) Find() ([]User, error) {
	return us.repo.Find()
}

func (us *UserService) eventHandler(e Event) error {
	fmt.Printf("%+v\n", e)
	switch e.Kind {
	case CreateEvent:
		fmt.Println(e.Kind)
		err := us.repo.Create(e.User)
		if err != nil {
			return err
		}
	case UpdateEvent:
		fmt.Println(e.Kind)
		err := us.repo.Update(e.User)
		if err != nil {
			fmt.Printf(err.Error())
			return err
		}
	case DeleteEvent:
		fmt.Println(e.Kind)
		err := us.repo.Delete(e.User.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
