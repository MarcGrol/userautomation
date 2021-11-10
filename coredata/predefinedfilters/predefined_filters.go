package predefinedfilters

import (
	"context"
	"github.com/MarcGrol/userautomation/core/user"
)

const (
	FilterYoungAgeName = "young_age_filter"
	FilterOldAgeName   = "old_age_filter"
)

var (
	FilterYoungAge = func(ctx context.Context, user user.User) (bool, error) {
		age, ok := user.Attributes["age"].(int)
		if !ok {
			return false, nil
		}
		return age < 18, nil
	}

	FilterOldAge = func(ctx context.Context, user user.User) (bool, error) {
		age, ok := user.Attributes["age"].(int)
		if !ok {
			return false, nil
		}
		return age > 40, nil
	}

	UserFilters = map[string]user.FilterFunc{
		FilterYoungAgeName: FilterYoungAge,
		FilterOldAgeName:   FilterOldAge,
	}
)
