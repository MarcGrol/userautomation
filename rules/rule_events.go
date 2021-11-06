package rules

type RuleCreatedEvent struct {
	State UserSegmentRule
}

type RuleModifiedEvent struct {
	OldState UserSegmentRule
	NewState UserSegmentRule
}

type RuleRemovedEvent struct {
	State UserSegmentRule
}
