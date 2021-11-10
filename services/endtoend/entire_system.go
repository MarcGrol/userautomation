package endtoend

import (
	"context"
	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/services/ondemandtriggerservice"
	"github.com/MarcGrol/userautomation/services/rulemanagement"
	"github.com/MarcGrol/userautomation/services/segmentmanagement"
	"github.com/MarcGrol/userautomation/services/usermanagement"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type EntireSystem interface {
	GetUserService() user.Management
	GetRuleService() segmentrule.Service
	GetSegmentService() segmentmanagement.SegmentManagement
	GetOnDemandExecutionService() segmentrule.TriggerRuleExecution
}

type entireSystemWiredTogether struct {
	userService     user.Management
	ruleService     segmentrule.Service
	segmentService  segmentmanagement.SegmentManagement
	ondemandService segmentrule.TriggerRuleExecution
}

func New(ctx context.Context) EntireSystem {
	pubsub := pubsub.NewPubSub()

	userStore := datastore.NewDatastore()
	userService := usermanagement.New(userStore, pubsub)

	ruleStore := datastore.NewDatastore()
	ruleService := rulemanagement.New(ruleStore)

	segmentStore := datastore.NewDatastore()
	segmentService := segmentmanagement.New(segmentStore, pubsub)

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

func (s *entireSystemWiredTogether) GetRuleService() segmentrule.Service {
	return s.ruleService
}

func (s *entireSystemWiredTogether) GetSegmentService() segmentmanagement.SegmentManagement {
	return s.segmentService
}

func (s *entireSystemWiredTogether) GetOnDemandExecutionService() segmentrule.TriggerRuleExecution {
	return s.ondemandService
}
