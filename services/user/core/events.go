package core

const UserCreatedEventKey = "userCreated"

type UserCreatedEvent struct {
	User User
}

const UserUpdatedEventKey = "userUpdated"

type UserUpdatedEvent struct {
	User User
}

const UserDeletedEventKey = "userDeleted"

type UserDeletedEvent struct {
	UserId string
}
