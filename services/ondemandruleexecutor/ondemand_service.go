package ondemandruleexecutor

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/user"
)

type OnDemandRuleExecutor struct {
	ruleService rule.SegmentRuleService
	userService user.Service
}

func New(ruleService rule.SegmentRuleService, userService user.Service) rule.SegmentRuleExecutionService {
	return &OnDemandRuleExecutor{
		ruleService: ruleService,
		userService: userService,
	}
}

func (s *OnDemandRuleExecutor) Trigger(ctx context.Context, ruleUID string) error {
	r, exists, err := s.ruleService.Get(ctx, ruleUID)
	if err != nil {
		return fmt.Errorf("Error getting rule with uid %s: %s", ruleUID, err)
	}
	if !exists {
		return fmt.Errorf("Rule with uid %s does not exist: %s", ruleUID, err)
	}

	// TODO fix this
	//if (r.AllowedTriggers & rule.TriggerOnDemand) == 0 {
	//	return fmt.Errorf("Rule with uid %s cannot be executed on demand", ruleUID)
	//}

	users, err := s.userService.Query(ctx, r.UserSegment.IsApplicableForUser)
	if err != nil {
		return fmt.Errorf("Error getting all users while executing rule %s: %s", r.UID, err)
	}

	for _, user := range users {
		_, err = executeRuleForUser(ctx, r, user)
		if err != nil {
			return fmt.Errorf("Error executing rule %s: %s", r.UID, err)
		}
	}

	return nil
}

func executeRuleForUser(ctx context.Context, r rule.UserSegmentRule, user user.User) (bool, error) {
	applicable, err := r.UserSegment.IsApplicableForUser(ctx, user)
	if err != nil {
		return false, fmt.Errorf("Error determining if rule %s is applicable for u %s: %s", r.UID, user.UID, err)
	}

	if !applicable {
		return false, nil
	}

	act := action.UserAction{
		RuleUID:  r.UID,
		Reason:   action.ReasonIsOnDemand,
		OldState: nil,
		NewState: &user,
	}

	err = r.Action.Perform(ctx, act)
	if err != nil {
		return false, fmt.Errorf("Error performing action for rule %s and useer %s: %s", r.UID, user.UID, err)
	}

	return true, nil
}
