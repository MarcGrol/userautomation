package actions

//go:generate mockgen -source=group_api.go -destination=group_api_mock.go -package=actions GroupApi
type GroupApi interface {
	GroupExists(groupUid string) (bool, error)
	AddUserToGroup(groupName, userUid string) error
}
