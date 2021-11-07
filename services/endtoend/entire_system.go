package endtoend

import (
	"context"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/services/ondemandservice"
	"github.com/MarcGrol/userautomation/services/ruleservice"
	"github.com/MarcGrol/userautomation/services/segmentservice"
	"github.com/MarcGrol/userautomation/services/usereventservice"
	"github.com/MarcGrol/userautomation/services/userservice"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type EntireSystem interface {
	GetUserService() user.Service
	GetRuleService() rule.SegmentRuleService
	GetSegmentService() segmentservice.SegmentService
	GetOnDemandExecutionService() rule.SegmentRuleExecutionService
}

type entireSystemWiredTogether struct {
	userService     user.Service
	ruleService     rule.SegmentRuleService
	segmentService  segmentservice.SegmentService
	ondemandService rule.SegmentRuleExecutionService
}

func New(ctx context.Context) EntireSystem {
	pubsub := pubsub.NewPubSub()

	userStore := datastore.NewDatastore()
	userService := userservice.NewUserService(userStore, pubsub)

	ruleStore := datastore.NewDatastore()
	ruleService := ruleservice.NewUserSegmentRuleService(ruleStore, pubsub)

	segmentStore := datastore.NewDatastore()
	segmentService := segmentservice.New(segmentStore, userService, pubsub)

	userEventService := usereventservice.NewUserEventService(pubsub, ruleService)
	userEventService.Subscribe(ctx)

	ondemandService := ondemandservice.New(ruleService, userService)

	return &entireSystemWiredTogether{
		userService:     userService,
		ruleService:     ruleService,
		segmentService:  segmentService,
		ondemandService: ondemandService,
	}
}

func (s *entireSystemWiredTogether) GetUserService() user.Service {
	return s.userService
}

func (s *entireSystemWiredTogether) GetRuleService() rule.SegmentRuleService {
	return s.ruleService
}

func (s *entireSystemWiredTogether) GetSegmentService() segmentservice.SegmentService {
	return s.segmentService
}

func (s *entireSystemWiredTogether) GetOnDemandExecutionService() rule.SegmentRuleExecutionService {
	return s.ondemandService
}
