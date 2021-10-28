package userlookup

import "github.com/MarcGrol/userautomation/api"

//go:generate mockgen -source=user_lookup.go -destination=user_lookup_mock.go -package=userlookup UserLookuper
type UserLookuper interface {
	GetUserOnUid(uid string) (api.User, error)
	GetUserOnQuery(whereClause string) ([]api.User, error)
}
