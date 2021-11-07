package endtoend

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestOnDemand(t *testing.T) {
	ctx := context.TODO()

	t.Run("execute rule, no rule exists", func(t *testing.T) {
		// setup
		sut := New(ctx)

		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given

		// when
		err := executeYoungAgeRuleReturnError(ctx, t, sut.GetOnDemandExecutionService())
		assert.Error(t, err)

		// then
	})

	t.Run("execute rule, no rule matched", func(t *testing.T) {
		// setup
		sut := New(ctx)

		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createOldAgeRule(ctx, t, sut.GetRuleService(), mockEmailer)

		// when
		err := executeYoungAgeRuleReturnError(ctx, t, sut.GetOnDemandExecutionService())
		assert.Error(t, err)

		// then
	})

	t.Run("execute rule, young age rule matched with no users", func(t *testing.T) {
		// setup
		sut := New(ctx)

		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createYoungAgeRule(ctx, t, sut.GetRuleService(), mockSmser)

		// when
		defer executeYoungAgeRule(ctx, t, sut.GetOnDemandExecutionService())

		// then
	})

	t.Run("execute rule, young age rule matched -> sms", func(t *testing.T) {
		// setup
		sut := New(ctx)

		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, sut.GetUserService(), 12)
		createYoungAgeRule(ctx, t, sut.GetRuleService(), mockSmser)

		// when
		defer executeYoungAgeRule(ctx, t, sut.GetOnDemandExecutionService())

		// then
		mockSmser.EXPECT().Send(gomock.Any(), "+31611111111",
			"young rule fired for Marc: your age is 12").Return(nil)
	})

	t.Run("execute rule, old age rule matched -> no users", func(t *testing.T) {
		// setup
		sut := New(ctx)

		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createYoungAgeRule(ctx, t, sut.GetRuleService(), mockSmser)

		// when
		defer executeYoungAgeRule(ctx, t, sut.GetOnDemandExecutionService())

		// then
	})

	t.Run("execute rule, old age rule matched -> email", func(t *testing.T) {
		// setup
		sut := New(ctx)

		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, sut.GetUserService(), 50)
		createOldAgeRule(ctx, t, sut.GetRuleService(), mockEmailer)

		// when
		defer executeOldAgeRule(ctx, t, sut.GetOnDemandExecutionService())

		// then
		mockEmailer.EXPECT().Send(gomock.Any(), "marc@home.nl",
			"old rule fired", "Hoi Marc, your age is 50").Return(nil)

	})

	t.Run("execute rule, old age rule matched multiple users -> 2 emails", func(t *testing.T) {
		// setup
		sut := New(ctx)

		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, sut.GetUserService(), 50)
		createOtherUser(ctx, t, sut.GetUserService(), 50)
		createOldAgeRule(ctx, t, sut.GetRuleService(), mockEmailer)

		// when
		defer executeOldAgeRule(ctx, t, sut.GetOnDemandExecutionService())

		// then
		mockEmailer.EXPECT().Send(gomock.Any(), gomock.Any(),
			"old rule fired", gomock.Any()).Return(nil).Times(2)

	})

}
