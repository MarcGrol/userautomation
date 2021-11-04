package userlookup

import (
	"context"
	"github.com/MarcGrol/userautomation/batch/batchcore"
)

//go:generate mockgen -source=user_lookup.go -destination=user_lookup_mock_test.go -package=userlookup UserLookuper
type UserLookuper interface {
	GetUserOnUid(c context.Context, uid string) (batchcore.User, error)
	GetUserOnQuery(c context.Context, whereClause string) ([]batchcore.User, error)
}
