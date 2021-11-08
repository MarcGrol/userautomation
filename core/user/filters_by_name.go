package user

import (
	"context"
)

const (
	FilterYoungAge = "young_age"
	FilterOldAge   = "old_age"
)

var (
	userFilters = map[string]FilterFunc{
		FilterYoungAge: func(ctx context.Context, user User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age < 18, nil
		},
		FilterOldAge: func(ctx context.Context, user User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age > 40, nil
		},
	}
)

func GetUserFilterByName(ctx context.Context, name string) (FilterFunc, bool) {
	ff, found := userFilters[name]
	return ff, found
}
