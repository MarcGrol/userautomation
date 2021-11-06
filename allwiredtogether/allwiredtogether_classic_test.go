package allwiredtogether

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUsingClassicSubTests(t *testing.T) {
	ctx := context.TODO()

	t.Run("create user, no rule exists", func(t *testing.T) {
		// setup
		_, userService := setupSut(ctx)
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given

		// when
		defer createUser(ctx, t, userService, 50)

		// then
	})

	t.Run("create user, no rule matched", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer createUser(ctx, t, userService, 50)

		// then
	})

	t.Run("create user, young age rule matched -> sms", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer createUser(ctx, t, userService, 12)

		// then
		mockSmser.EXPECT().Send(gomock.Any(), "+31611111111",
			"young rule fired for Marc: your age is 12").Return(nil)
	})

	t.Run("create user, old age rule matched -> email", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// when
		defer createUser(ctx, t, userService, 50)

		// then
		mockEmailer.EXPECT().Send(gomock.Any(), "marc@home.nl",
			"old rule fired", "Hoi Marc, your age is 50").Return(nil)

	})

	t.Run("modify user, no rule exist", func(t *testing.T) {
		// setup
		_, userService := setupSut(ctx)
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// expect

		// when
		defer modifyUser(ctx, t, userService, 12)

		// then
		createUser(ctx, t, userService, 50)

	})

	t.Run("modify user, no rule matched", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, userService, 12)
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// when
		defer modifyUser(ctx, t, userService, 14)

		// then

	})

	t.Run("modify user, young age rule matched -> sms", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, userService, 50)
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer modifyUser(ctx, t, userService, 12)

		// then
		mockSmser.EXPECT().Send(gomock.Any(), "+31611111111",
			"young rule fired for Marc: your age is 12").Return(nil)

	})

	t.Run("modify user, old age rule matched -> email", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)
		createUser(ctx, t, userService, 12)

		// when
		defer modifyUser(ctx, t, userService, 50)

		// then
		mockEmailer.EXPECT().Send(gomock.Any(), "marc@home.nl",
			"old rule fired", "Hoi Marc, your age is 50").Return(nil)

	})

	t.Run("modify user, remains young", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)
		createUser(ctx, t, userService, 12)

		// when
		defer modifyUser(ctx, t, userService, 14)

		// then
	})

	t.Run("delete user, no user exists", func(t *testing.T) {
		// setup
		_, userService := setupSut(ctx)
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given

		// when
		defer removeUser(ctx, t, userService)

		// then
	})

	t.Run("delete user, no rule exist", func(t *testing.T) {
		// setup
		_, userService := setupSut(ctx)
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, userService, 50)

		// when
		defer removeUser(ctx, t, userService)

		// then
	})

	t.Run("delete user, no rule matched", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, userService, 50)
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer removeUser(ctx, t, userService)

		// then
	})

	t.Run("delete user, young age rule matched", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, userService, 50)
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// when
		defer removeUser(ctx, t, userService)

		// then
	})
}
