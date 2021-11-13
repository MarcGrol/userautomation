package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/core/userrule"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/MarcGrol/userautomation/integrations/emailsending"
	"github.com/MarcGrol/userautomation/integrations/smssending"
	"github.com/MarcGrol/userautomation/services/actionmanager"
	"github.com/MarcGrol/userautomation/services/filtermanagement"
	"github.com/MarcGrol/userautomation/services/rulemanagement"
	"github.com/MarcGrol/userautomation/services/segmentchangeevaluator"
	"github.com/MarcGrol/userautomation/services/segmentmanagement"
	"github.com/MarcGrol/userautomation/services/segmentqueryservice"
	"github.com/MarcGrol/userautomation/services/segmentruleevaluator"
	"github.com/MarcGrol/userautomation/services/segmentruletrigger"
	"github.com/MarcGrol/userautomation/services/segmentusermanagement"
	"github.com/MarcGrol/userautomation/services/usermanagement"
	userruleevaluator "github.com/MarcGrol/userautomation/services/userruleevaluation"
	"github.com/MarcGrol/userautomation/services/userruletriggerservice"
	"github.com/MarcGrol/userautomation/services/usertaskexecutor"
	"github.com/gorilla/mux"
)

type system struct {
	ps             pubsub.Pubsub
	userManager    user.Management
	segmentManager segment.Management
	ruleManager    segmentrule.Management
	actionManager  action.Management

	userRuleTrigger    userrule.TriggerRuleExecution
	segmentRuleTrigger segmentruletrigger.Service

	segmentUserManager      segmentusermanagement.SegmentUserManager
	segmentUserQueryManager segment.Querier

	userRuleEvaluator      userruleevaluator.Service
	segmentRuleEvaluator   segmentruleevaluator.Service
	segmentChangeEvaluator segmentchangeevaluator.Service

	userTaskExecutor usertaskexecutor.Service
}

func new(ctx context.Context) *system {
	ps := pubsub.NewPubSub()

	fm := filtermanagement.New()

	userManager := usermanagement.New(datastore.NewDatastoreStub(), fm, ps)

	segmentManager := segmentmanagement.New(datastore.NewDatastoreStub(), ps)

	ruleManager := rulemanagement.New(datastore.NewDatastoreStub(), ps)

	actionManager := actionmanager.New(datastore.NewDatastoreStub())

	userRuleTriggerService := userruletriggerservice.New(ps)

	segmentRuleTrigger := segmentruletrigger.New(ruleManager, ps)

	segmentUserManager := segmentusermanagement.New(datastore.NewDatastoreStub(), userManager, fm, ps)

	segmentQueryService := segmentqueryservice.New(segmentUserManager)

	userRuleEvaluator := userruleevaluator.New(ps)

	segmentRuleEvaluator := segmentruleevaluator.New(ps, ruleManager, userManager)

	segmentcChangeEvaluator := segmentchangeevaluator.New(ps, ruleManager)

	userTaskExecutor := usertaskexecutor.New(ps, usertask.NewExecutionReporterStub(),
		smssending.NewSmsSender(), emailsending.NewEmailSender())

	return &system{
		ps:                      ps,
		userManager:             userManager,
		segmentManager:          segmentManager,
		ruleManager:             ruleManager,
		actionManager:           actionManager,
		segmentUserManager:      segmentUserManager,
		segmentUserQueryManager: segmentQueryService,
		userRuleTrigger:         userRuleTriggerService,
		segmentRuleTrigger:      segmentRuleTrigger,
		userRuleEvaluator:       userRuleEvaluator,
		segmentRuleEvaluator:    segmentRuleEvaluator,
		segmentChangeEvaluator:  segmentcChangeEvaluator,
		userTaskExecutor:        userTaskExecutor,
	}
}

func (s *system) Subscribe(ctx context.Context, router *mux.Router) error {
	s.segmentUserManager.Subscribe(ctx, router)
	s.userRuleEvaluator.Subscribe(ctx, router)
	s.segmentRuleEvaluator.Subscribe(ctx, router)
	s.segmentChangeEvaluator.Subscribe(ctx, router)
	s.userTaskExecutor.Subscribe(ctx, router)

	return nil
}

func (s *system) Register(ctx context.Context, router *mux.Router) error {
	s.userManager.RegisterEndpoints(ctx, router)
	s.segmentManager.RegisterEndpoints(ctx, router)
	s.ruleManager.RegisterEndpoints(ctx, router)
	s.actionManager.RegisterEndpoints(ctx, router)
	s.userRuleTrigger.RegisterEndpoints(ctx, router)
	s.segmentRuleTrigger.RegisterEndpoints(ctx, router)
	s.segmentUserQueryManager.RegisterEndpoints(ctx, router)

	return nil
}

func (s *system) Preprov(ctx context.Context) error {
	s.userManager.Preprov(ctx)
	s.segmentManager.Preprov(ctx)
	s.ruleManager.Preprov(ctx)
	s.actionManager.Preprov(ctx)

	return nil
}

func (s *system) Start(ctx context.Context, router *mux.Router) error {
	http.Handle("/", router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

	return nil
}

func main() {
	ctx := context.Background()
	router := mux.NewRouter()

	sys := new(ctx)

	sys.Subscribe(ctx, router)

	sys.Preprov(ctx)

	sys.Register(ctx, router)

	sys.Start(ctx, router)

}
