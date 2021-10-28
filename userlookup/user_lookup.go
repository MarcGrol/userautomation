package userlookup

import "github.com/MarcGrol/userautomation/core"

//go:generate mockgen -source=user_lookup.go -destination=user_lookup_mock.go -package=userlookup UserLookuper
type UserLookuper interface {
	GetUserOnUid(uid string) (core.User, error)
	GetUserOnQuery(whereClause string) ([]core.User, error)
}
