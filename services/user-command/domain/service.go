package domain

type UserRepository interface {
	Append(command string, user User) error
}

type EventBus interface {
	PublishEvent(kind string, user User) error
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

const (
	CreateEvent = "userCreated"
	UpdateEvent = "userUpdated"
	DeleteEvent = "userDeleted"
)

func (us *UserService) Create(user CreateUser) (User, error) {
	u, err := NewUser(user.Name, user.Surname)
	if err != nil {
		return User{}, err
	}
	return us.storeAndPublish(CreateCommand, CreateEvent, u)
}

type UpdateUser struct {
	Id      string
	Name    string
	Surname string
}

func (us *UserService) Update(user UpdateUser) (User, error) {
	u := User{
		Id:      user.Id,
		Name:    user.Name,
		Surname: user.Surname,
	}

	if err := u.Validate(); err != nil {
		return User{}, err
	}

	return us.storeAndPublish(UpdateCommand, UpdateEvent, u)
}

func (us *UserService) Delete(id string) error {
	u := User{
		Id: id,
	}
	_, err := us.storeAndPublish(DeleteCommand, DeleteEvent, u)
	return err
}

func (us *UserService) storeAndPublish(command, eventKind string, user User) (User, error) {

	if err := us.repo.Append(command, user); err != nil {
		return User{}, err
	}

	if err := us.eventBus.PublishEvent(eventKind, user); err != nil {
		return User{}, err
	}
	return user, nil
}
