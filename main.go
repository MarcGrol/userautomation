package main

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/integrations/emailsending"
	"github.com/MarcGrol/userautomation/integrations/smssending"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
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
	"github.com/MarcGrol/userautomation/services/userruleevaluation"
	"github.com/MarcGrol/userautomation/services/userruletriggerservice"
	"github.com/MarcGrol/userautomation/services/usertaskexecutor"
)

func main() {
	ctx := context.Background()

	var router = mux.NewRouter()

	ps := pubsub.NewPubSub()

	fm := filtermanagement.New()

	userManager := usermanagement.New(datastore.NewDatastoreStub(), fm, ps)
	{
		userManager.Preprov(ctx)
		userManager.RegisterEndpoints(ctx, router)
	}

	segmentManager := segmentmanagement.New(datastore.NewDatastoreStub(), ps)
	{
		segmentManager.Preprov(ctx)
		segmentManager.RegisterEndpoints(ctx, router)
	}

	ruleManager := rulemanagement.New(datastore.NewDatastoreStub(), ps)
	{
		ruleManager.Preprov(ctx)
		ruleManager.RegisterEndpoints(ctx, router)
	}

	{
		actionManager := actionmanager.New(datastore.NewDatastoreStub())
		actionManager.Preprov(ctx)
		actionManager.RegisterEndpoints(ctx, router)
	}

	{
		service := userruletriggerservice.New(ps)
		service.RegisterEndpoints(ctx, router)
	}

	{
		service := segmentruletrigger.New(ruleManager, ps)
		service.RegisterEndpoints(ctx, router)
	}

	segmentUserManager := segmentusermanagement.New(datastore.NewDatastoreStub(), userManager, fm, ps)
	{
		segmentUserManager.Subscribe(ctx, router)
	}

	{
		segmentQueryService := segmentqueryservice.New(segmentUserManager)
		segmentQueryService.RegisterEndpoints(ctx, router)
	}

	{
		evaluator := userruleevaluator.New(ps)
		evaluator.Subscribe(ctx, router)
	}

	{
		evaluator := segmentruleevaluator.New(ps, ruleManager, userManager)
		evaluator.Subscribe(ctx, router)
	}
	{
		evaluator := segmentchangeevaluator.New(ps, ruleManager)
		evaluator.Subscribe(ctx, router)
	}

	{
		evaluator := usertaskexecutor.New(ps, usertask.NewExecutionReporterStub(),
			smssending.NewSmsSender(), emailsending.NewEmailSender())
		evaluator.Subscribe(ctx, router)
	}

	http.Handle("/", router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
