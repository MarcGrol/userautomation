package realtimecore


type UserCreatedEvent struct{
	State User
}

type UserModifiedEvent struct {
	OldState User
	NewState User
}

type UserRemovedEvent struct {
	State User
}

