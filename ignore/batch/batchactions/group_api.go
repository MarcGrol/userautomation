package batchactions

import "context"

//go:generate mockgen -source=group_api.go -destination=group_api_mock.go -package=batchactions GroupApi
type GroupApi interface {
	GroupExists(c context.Context, groupUid string) (bool, error)
	AddUserToGroup(c context.Context, groupName, userUid string) error
}
