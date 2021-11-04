package main

import (
	"context"
	"github.com/MarcGrol/userautomation/realtime/realtimeactions"
	"github.com/MarcGrol/userautomation/realtime/realtimeservices"
	"log"

	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

func main() {
	ctx := context.TODO()

	pubsub := realtimeservices.NewPubSub()

	ruleService := realtimeservices.NewUserSegmentRuleService()
	userEventService := realtimeservices.NewUserEventHandler(pubsub, ruleService)
	userEventService.Subscribe(ctx)

	userService := realtimeservices.NewUserService(pubsub)

	// pre-provision users: no rules yet so nothing fires
	createMarc(ctx, userService)
	createEva(ctx, userService)

	// pre-provision segment rules
	createOldRule(ctx, ruleService, realtimeactions.NewEmailSender())
	createYoungRule(ctx, ruleService,  realtimeactions.NewSmsSender())

	// Start processing commands on users
	adjustMarc(ctx, userService) // young-rule fires, sms action
	deleteMarc(ctx, userService) // no rule fires

	createFreek(ctx, userService)      // young-rule fires, sms action
	adjustFreek(ctx, userService)      // still young-rule, no action
	adjustFreekAgain(ctx, userService) // old-rule fires, email action
}

func createMarc(ctx context.Context, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID:      "1",
		Attributes: map[string]interface{}{
			"firstname":    "Marc",
			"emailaddress": "marc@home.nl",
			"phonenumber":  "+31611111111",
			"age":          50, // old
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func adjustMarc(ctx context.Context, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID:      "1",
		Attributes: map[string]interface{}{
			"firstname":    "Marc",
			"emailaddress": "marc@home.nl",
			"phonenumber":  "+31611111111",
			"age":          10, // now young
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func deleteMarc(ctx context.Context, userService realtimecore.UserService) {
	err := userService.Remove(ctx, "1")
	if err != nil {
		log.Fatalln(err)
	}
}

func createEva(ctx context.Context, userService realtimecore.UserService) {

	err := userService.Put(ctx, realtimecore.User{
		UID:      "2",
		Attributes: map[string]interface{}{
			"firstname":    "Eva",
			"emailaddress": "eva@home.nl",
			"phonenumber":  "+31622222222",
			"age":          48, // old
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
}


func createFreek(ctx context.Context, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID:      "3",
		Attributes: map[string]interface{}{
			"firstname":    "Freek",
			"emailaddress": "freek@home.nl",
			"phonenumber":  "+31633333333",
			"age":          12, // young
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func adjustFreek(ctx context.Context, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID:      "3",
		Attributes: map[string]interface{}{
			"firstname":    "Freek",
			"emailaddress": "freek@home.nl",
			"phonenumber":  "+31633333333",
			"age":          13, // slightly older but still young
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func adjustFreekAgain(ctx context.Context, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID:      "3",
		Attributes: map[string]interface{}{
			"firstname":    "Freek",
			"emailaddress": "freek@home.nl",
			"phonenumber":  "+31633333333",
			"age":          41, // big increase in age, now old
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func createOldRule(ctx context.Context, segmentService realtimecore.SegmentRuleService,
	emailSender realtimeactions.Emailer) {
	err := segmentService.Put(ctx, realtimecore.UserSegmentRule{
		Name: "OldRule",
		IsApplicableForUser: func(ctx context.Context, user realtimecore.User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age > 40, nil
		},
		PerformAction: realtimeactions.EmailerAction("old rule fired", "Hoi {{.firstname}}, your age is {{.age}}", emailSender),
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func createYoungRule(ctx context.Context, segmentService realtimecore.SegmentRuleService, smsSender realtimeactions.SmsSender){

	err := segmentService.Put(ctx, realtimecore.UserSegmentRule{
		Name: "YoungRule",
		IsApplicableForUser: func(ctx context.Context, user realtimecore.User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age < 18, nil
		},
		PerformAction: realtimeactions.SmsAction("young rule fired for {{.firstname}}: your age is {{.age}}", smsSender),
	})
	if err != nil {
		log.Fatalln(err)
	}
}
