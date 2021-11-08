package segment

import (
	"context"
	"github.com/MarcGrol/userautomation/core/user"
)

const (
	FilterYoungAge = "young_age"
	FilterOldAge   = "old_age"
)

var (
	userFilters = map[string]user.FilterFunc{
		FilterYoungAge: func(ctx context.Context, user user.User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age < 18, nil
		},
		FilterOldAge: func(ctx context.Context, user user.User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age > 40, nil
		},
	}
)

func GetUserFilterByName(name string) (user.FilterFunc, bool) {
	ff, found := userFilters[name]
	return ff, found
}
