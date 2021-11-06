package allwiredtogether

import (
	"context"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/services/ruleservice"
	"github.com/MarcGrol/userautomation/services/usereventservice"
	"github.com/MarcGrol/userautomation/services/userservice"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type EntireSystem interface {
	GetUserService() user.UserService
	GetSegmentRuleService() rule.SegmentRuleService
}

type entireSystemWiredTogether struct {
	userService user.UserService
	ruleService rule.SegmentRuleService
}

func New(ctx context.Context) EntireSystem {
	pubsub := pubsub.NewPubSub()

	userStore := datastore.NewDatastore()
	userService := userservice.NewUserService(userStore, pubsub)

	ruleStore := datastore.NewDatastore()
	ruleService := ruleservice.NewUserSegmentRuleService(ruleStore, pubsub)

	userEventService := usereventservice.NewUserEventService(pubsub, ruleService)
	userEventService.Subscribe(ctx)

	return &entireSystemWiredTogether{
		userService: userService,
		ruleService: ruleService,
	}
}

func (s *entireSystemWiredTogether) GetUserService() user.UserService {
	return s.userService
}

func (s *entireSystemWiredTogether) GetSegmentRuleService() rule.SegmentRuleService {
	return s.ruleService
}
