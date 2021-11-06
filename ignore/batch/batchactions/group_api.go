package batchactions

import "context"

//go:generate mockgen -source=group_api.go -destination=group_api_mock.go -package=batchactions GroupApi
type GroupAPI interface {
	GroupExists(c context.Context, groupUID string) (bool, error)
	AddUserToGroup(c context.Context, groupName, userUID string) error
}
