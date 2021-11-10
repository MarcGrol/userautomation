package user

import (
	"context"
)

func GetUserFilterByName(ctx context.Context, name string) (FilterFunc, bool) {
	ff, found := userFilters[name]
	return ff, found
}
