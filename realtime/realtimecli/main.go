package main

import (
	"context"
	"github.com/MarcGrol/userautomation/realtime/realtimeactions"
	"github.com/MarcGrol/userautomation/realtime/realtimeservices"
	"log"

	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

func main() {
	ruleService := realtimeservices.NewUserSegmentRuleService()
	userEventService := realtimeservices.NewUserEventHandler(ruleService)
	userService := realtimeservices.NewUserService(userEventService)
	emailSender := realtimeactions.NewEmailSender()
	smsSender := realtimeactions.NewSmsSender()

	ctx := context.TODO()

	preprovisionUsers(ctx, userService) // no rules present, nothing fires
	preprovisionUserSegmentRules(ctx, ruleService, emailSender, smsSender)

	adjustMarc(ctx, userService) // young-rule fires, sms action
	deleteMarc(ctx, userService) // no rule fires

	createFreek(ctx, userService)      // young-rule fires, sms action
	adjustFreek(ctx, userService)      // still young-rule, no action
	adjustFreekAgain(ctx, userService) // old-rule fires, email actio
}

func adjustMarc(ctx context.Context, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID:      "1",
		FullName: "Marc",
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
	err := userService.Delete(ctx, "1")
	if err != nil {
		log.Fatalln(err)
	}
}

func createFreek(ctx context.Context, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID:      "3",
		FullName: "Freek",
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
		FullName: "Freek",
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
		FullName: "Freek",
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

func preprovisionUsers(ctx context.Context, userService realtimecore.UserService) {
	createMarc(ctx, userService)
	createEva(ctx, userService)
}

func createMarc(ctx context.Context, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID:      "1",
		FullName: "Marc",
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

func createEva(ctx context.Context, userService realtimecore.UserService) {

	err := userService.Put(ctx, realtimecore.User{
		UID:      "2",
		FullName: "Eva",
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

func preprovisionUserSegmentRules(ctx context.Context, segmentService realtimecore.SegmentRuleService,
	emailSender realtimeactions.Emailer,
	smsSender realtimeactions.SmsSender) {
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

	err = segmentService.Put(ctx, realtimecore.UserSegmentRule{
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
