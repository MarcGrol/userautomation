package actions

import "context"

//go:generate mockgen -source=group_api.go -destination=group_api_mock.go -package=actions GroupApi
type GroupApi interface {
	GroupExists(c context.Context, groupUid string) (bool, error)
	AddUserToGroup(c context.Context, groupName, userUid string) error
}
