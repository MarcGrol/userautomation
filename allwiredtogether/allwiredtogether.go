package allwiredtogether

import (
	"context"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/triggers/usereventhandler"
	"github.com/MarcGrol/userautomation/users"
)

type EntireSystem interface {
	GetUserService() users.UserService
	GetSegmentRuleService() rules.SegmentRuleService
}

type allComponentsWiredTogether struct {
	userService users.UserService
	ruleService rules.SegmentRuleService
}

func New(ctx context.Context) EntireSystem {
	pubsub := pubsub.NewPubSub()

	ruleService := rules.NewUserSegmentRuleService()

	userEventService := usereventhandler.NewUserEventService(pubsub, ruleService)
	userEventService.Subscribe(ctx)

	userService := users.NewUserService(pubsub)

	return &allComponentsWiredTogether{
		userService: userService,
		ruleService: ruleService,
	}
}

func (s *allComponentsWiredTogether)GetUserService() users.UserService {
	return s.userService
}

func (s *allComponentsWiredTogether)GetSegmentRuleService() rules.SegmentRuleService {
	return s.ruleService
}
