package user

import (
	"context"
)

const (
	FilterYoungAgeName = "young_age"
	FilterOldAgeName   = "old_age"
)

var (
	FilterYoungAge = func(ctx context.Context, user User) (bool, error) {
		age, ok := user.Attributes["age"].(int)
		if !ok {
			return false, nil
		}
		return age < 18, nil
	}

	FilterOldAge = func(ctx context.Context, user User) (bool, error) {
		age, ok := user.Attributes["age"].(int)
		if !ok {
			return false, nil
		}
		return age > 40, nil
	}

	userFilters = map[string]FilterFunc{
		FilterYoungAgeName: FilterYoungAge,
		FilterOldAgeName:   FilterOldAge,
	}
)
