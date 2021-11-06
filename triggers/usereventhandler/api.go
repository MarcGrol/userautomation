package usereventhandler

import "context"

type UserEventService interface {
	Subscribe(c context.Context) error
	OnEvent(c context.Context, topic string, event interface{}) error
}
