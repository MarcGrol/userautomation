package endtoend

import (
	"context"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/services/ondemandtriggerservice"
	"github.com/MarcGrol/userautomation/services/rulemanagement"
	"github.com/MarcGrol/userautomation/services/segmentmanagement"
	"github.com/MarcGrol/userautomation/services/usercchangevaluator"
	"github.com/MarcGrol/userautomation/services/usermanagement"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type EntireSystem interface {
	GetUserService() user.Management
	GetRuleService() rule.RuleService
	GetSegmentService() segmentmanagement.SegmentManagement
	GetOnDemandExecutionService() rule.TriggerRuleExecution
}

type entireSystemWiredTogether struct {
	userService     user.Management
	ruleService     rule.RuleService
	segmentService  segmentmanagement.SegmentManagement
	ondemandService rule.TriggerRuleExecution
}

func New(ctx context.Context) EntireSystem {
	pubsub := pubsub.NewPubSub()

	userStore := datastore.NewDatastore()
	userService := usermanagement.New(userStore, pubsub)

	ruleStore := datastore.NewDatastore()
	ruleService := rulemanagement.New(ruleStore)

	segmentStore := datastore.NewDatastore()
	segmentService := segmentmanagement.New(segmentStore, pubsub)

	userEventService := usercchangevaluator.New(pubsub, ruleService)
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

func (s *entireSystemWiredTogether) GetRuleService() rule.RuleService {
	return s.ruleService
}

func (s *entireSystemWiredTogether) GetSegmentService() segmentmanagement.SegmentManagement {
	return s.segmentService
}

func (s *entireSystemWiredTogether) GetOnDemandExecutionService() rule.TriggerRuleExecution {
	return s.ondemandService
}
