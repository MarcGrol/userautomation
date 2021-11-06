package allwiredtogether

import (
	"context"
	"github.com/MarcGrol/userautomation/triggers/userchanged"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/MarcGrol/userautomation/rules"
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

	ruleStore := datastore.NewDatastore()
	ruleService := rules.NewUserSegmentRuleService(ruleStore, pubsub)

	userStore := datastore.NewDatastore()
	userService := users.NewUserService(userStore, pubsub)

	userEventService := userchanged.NewUserEventService(pubsub, ruleService)
	userEventService.Subscribe(ctx)

	return &allComponentsWiredTogether{
		userService: userService,
		ruleService: ruleService,
	}
}

func (s *allComponentsWiredTogether) GetUserService() users.UserService {
	return s.userService
}

func (s *allComponentsWiredTogether) GetSegmentRuleService() rules.SegmentRuleService {
	return s.ruleService
}
