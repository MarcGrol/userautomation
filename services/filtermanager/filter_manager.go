package filtermanager

import (
	"context"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/coredata/predefinedfilters"
)

type service struct{}

func New() user.FilterManager {
	return &service{}
}
func (s service) GetUserFilterByName(ctx context.Context, name string) (user.FilterFunc, bool) {
	ff, found := predefinedfilters.UserFilters[name]
	return ff, found
}
