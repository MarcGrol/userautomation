package predefinedfilters

import (
	"context"
	"strings"

	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/coredata/supportedattrs"
)

const (
	FilterYoungAgeName         = "young_age_filter"
	FilterOldAgeName           = "old_age_filter"
	FilterFirstnameStartsWithM = "first_name_starts_with_m"
)

var (
	FilterYoungAge = func(ctx context.Context, user user.User) (bool, error) {
		age, ok := user.Attributes[supportedattrs.Age].(int)
		if !ok {
			return false, nil
		}
		return age < 18, nil
	}

	FilterOldAge = func(ctx context.Context, user user.User) (bool, error) {
		age, ok := user.Attributes[supportedattrs.Age].(int)
		if !ok {
			return false, nil
		}
		return age > 40, nil
	}

	FilterFirstNameStart = func(ctx context.Context, user user.User) (bool, error) {
		name, ok := user.Attributes[supportedattrs.FirstName].(string)
		if !ok {
			return false, nil
		}
		return strings.HasSuffix(strings.ToLower(name), "m"), nil
	}

	UserFilters = map[string]user.FilterFunc{
		FilterYoungAgeName:         FilterYoungAge,
		FilterOldAgeName:           FilterOldAge,
		FilterFirstnameStartsWithM: FilterFirstNameStart,
	}
)
