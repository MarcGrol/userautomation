package endtoend

import (
	"context"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/services/ondemandtriggerservice"
	"github.com/MarcGrol/userautomation/services/rulemanagementservice"
	"github.com/MarcGrol/userautomation/services/segmentmanagement"
	"github.com/MarcGrol/userautomation/services/usereventservice"
	"github.com/MarcGrol/userautomation/services/userservice"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type EntireSystem interface {
	GetUserService() user.Management
	GetRuleService() rule.SegmentRuleService
	GetSegmentService() segmentmanagement.SegmentManagement
	GetOnDemandExecutionService() rule.SegmentRuleExecutionTrigger
}

type entireSystemWiredTogether struct {
	userService     user.Management
	ruleService     rule.SegmentRuleService
	segmentService  segmentmanagement.SegmentManagement
	ondemandService rule.SegmentRuleExecutionTrigger
}

func New(ctx context.Context) EntireSystem {
	pubsub := pubsub.NewPubSub()

	userStore := datastore.NewDatastore()
	userService := userservice.NewUserService(userStore, pubsub)

	ruleStore := datastore.NewDatastore()
	ruleService := rulemanagementservice.NewUserSegmentRuleService(ruleStore)

	segmentStore := datastore.NewDatastore()
	segmentService := segmentmanagement.New(segmentStore, pubsub)

	userEventService := usereventservice.NewUserEventService(pubsub, ruleService)
	userEventService.Subscribe(ctx)

	ondemandService := ondemandtriggerservice.New(ruleService, pubsub)

	return &entireSystemWiredTogether{
		userService:     userService,
		ruleService:     ruleService,
		segmentService:  segmentService,
		ondemandService: ondemandService,
	}
}

func (s *entireSystemWiredTogether) GetUserService() user.Management {
	return s.userService
}

func (s *entireSystemWiredTogether) GetRuleService() rule.SegmentRuleService {
	return s.ruleService
}

func (s *entireSystemWiredTogether) GetSegmentService() segmentmanagement.SegmentManagement {
	return s.segmentService
}

func (s *entireSystemWiredTogether) GetOnDemandExecutionService() rule.SegmentRuleExecutionTrigger {
	return s.ondemandService
}
